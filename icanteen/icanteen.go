package icanteen

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type ICanteenLunch struct {
	Name      string `json:"name"`
	Ordered   bool   `json:"ordered"`
	CanOrder  bool   `json:"can_order"`
	ChangeURL string `json:"change_url"`
}

type ICanteenDay struct {
	Day     string          `json:"day"`
	Lunches []ICanteenLunch `json:"lunches"`
}

type ICanteenData struct {
	Days   []ICanteenDay `json:"days"`
	Credit string        `json:"credit"`
}

func NormalizeServerURL(server string) (string, error) {
	// Add protocol if missing
	if !strings.HasPrefix(server, "http://") && !strings.HasPrefix(server, "https://") {
		server = "https://" + server
	}

	// Remove /login if present
	if strings.HasSuffix(server, "/login") {
		server = strings.TrimSuffix(server, "/login")
	}

	// Remove trailing slash if present
	if strings.HasSuffix(server, "/") {
		server = strings.TrimSuffix(server, "/")
	}

	// Validate the URL
	_, err := url.ParseRequestURI(server)
	if err != nil {
		return "", err
	}

	return server, nil
}

func login(username, password, server string) (*cookiejar.Jar, *url.URL, *http.Client, error) {
	cookieJar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: cookieJar}

	server, err := NormalizeServerURL(server)
	if err != nil {
		return nil, nil, nil, err
	}

	loginURL := server + "/login"
	parsedURL, err := url.Parse(loginURL)
	if err != nil {
		return nil, nil, nil, err
	}
	resp, err := client.Get(loginURL)
	if err != nil {
		return nil, nil, nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, nil, nil, err
	}
	csrfToken := doc.Find("input[name='_csrf']").AttrOr("value", "")

	loginData := url.Values{}
	loginData.Set("j_username", username)
	loginData.Set("j_password", password)
	loginData.Set("terminal", "false")
	loginData.Set("_csrf", csrfToken)
	loginData.Set("targetUrl", "/faces/secured/main.jsp?terminal=false&status=true&printer=&keyboard=")

	loginHeaders := map[string]string{
		"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8",
		"Accept-Encoding": "text/html",
		"Cache-Control":   "max-age=0",
		"Connection":      "keep-alive",
		"Content-Type":    "application/x-www-form-urlencoded",
		"Host":            parsedURL.Host,
		"Origin":          server,
		"Referer":         server + "/login",
	}
	loginURL = server + "/j_spring_security_check"
	req, err := http.NewRequest(http.MethodPost, loginURL, strings.NewReader(loginData.Encode()))
	if err != nil {
		return nil, nil, nil, err
	}
	for key, value := range loginHeaders {
		req.Header.Set(key, value)
	}
	resp, err = client.Do(req)
	if err != nil {
		return nil, nil, nil, err
	}
	defer resp.Body.Close()

	return cookieJar, parsedURL, client, nil
}

func Login(username, password, server string) ([]*http.Cookie, error) {
	cookieJar, url, _, err := login(username, password, server)
	if err != nil {
		return nil, err
	}

	return cookieJar.Cookies(url), nil
}

func LoadLunchesWithCookies(cookies []*http.Cookie, server string) (ICanteenData, error) {
	server, err := NormalizeServerURL(server)
	if err != nil {
		return ICanteenData{}, err
	}

	url, err := url.Parse(server)
	if err != nil {
		return ICanteenData{}, err
	}

	cookieJar, _ := cookiejar.New(nil)
	cookieJar.SetCookies(url, cookies)

	client := &http.Client{Jar: cookieJar}

	monthURL := server + "/faces/secured/month.jsp"
	req, err := http.NewRequest(http.MethodGet, monthURL, nil)
	if err != nil {
		return ICanteenData{}, err
	}
	loginHeaders := map[string]string{
		"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8",
		"Accept-Encoding": "text/html",
		"Cache-Control":   "max-age=0",
		"Connection":      "keep-alive",
		"Content-Type":    "application/x-www-form-urlencoded",
		"Host":            url.Host,
		"Origin":          server,
		"Referer":         server + "/login",
	}
	for key, value := range loginHeaders {
		req.Header.Set(key, value)
	}
	resp, err := client.Do(req)
	if err != nil {
		return ICanteenData{}, err
	}
	defer resp.Body.Close()

	// Cookies no longer valid
	if resp.Header.Get("Cache-Control") == "max-age=0" {
		return ICanteenData{}, http.ErrNoCookie
	}

	// Step 6: Parse lunch data
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return ICanteenData{}, err
	}

	icanteenData := ICanteenData{
		Credit: doc.Find("#Kredit").Eq(0).Text(),
	}
	doc.Find(".orderContent").Each(func(index int, td *goquery.Selection) {
		day := strings.TrimPrefix(td.AttrOr("id", ""), "orderContent")
		lunches := []ICanteenLunch{}
		td.Find(".jidelnicekItemWrapper").Each(func(index int, lunch *goquery.Selection) {
			lunchEntry := ICanteenLunch{
				Name:      strings.TrimSpace(strings.Split(lunch.Children().Eq(1).Text(), "\n")[1]),
				Ordered:   lunch.Children().Eq(0).Children().Eq(0).Children().Length() == 5,
				CanOrder:  lunch.Children().Eq(2).Children().Length() != 3,
				ChangeURL: strings.Split(strings.TrimPrefix(strings.TrimSpace(lunch.Children().Eq(0).Children().Eq(0).AttrOr("onclick", "not found")), "ajaxOrder(this, '"), "'")[0],
			}
			lunches = append(lunches, lunchEntry)
		})
		dayEntry := ICanteenDay{
			Day:     day,
			Lunches: lunches,
		}
		icanteenData.Days = append(icanteenData.Days, dayEntry)
	})

	return icanteenData, nil
}

func ChangeOrder(cookies []*http.Cookie, server, changeURL string) (ICanteenData, error) {
	server, err := NormalizeServerURL(server)
	if err != nil {
		return ICanteenData{}, err
	}

	url, err := url.Parse(server)
	if err != nil {
		return ICanteenData{}, err
	}

	cookieJar, _ := cookiejar.New(nil)
	cookieJar.SetCookies(url, cookies)

	client := &http.Client{Jar: cookieJar}

	changeURL = server + "/faces/secured/" + changeURL
	req, err := http.NewRequest(http.MethodGet, changeURL, nil)
	if err != nil {
		return ICanteenData{}, err
	}
	loginHeaders := map[string]string{
		"Accept":          "text/html, */*; q=0.01",
		"Accept-Encoding": "text/html",
		"Cache-Control":   "max-age=0",
		"Connection":      "keep-alive",
		"Content-Type":    "application/x-www-form-urlencoded",
		"Host":            url.Host,
		"Origin":          server,
		"Referer":         server + "/faces/secured/month.jsp",
	}
	for key, value := range loginHeaders {
		req.Header.Set(key, value)
	}
	resp, err := client.Do(req)
	if err != nil {
		return ICanteenData{}, err
	}
	defer resp.Body.Close()

	// Cookies no longer valid
	if resp.Header.Get("Cache-Control") == "max-age=0" {
		return ICanteenData{}, http.ErrNoCookie
	}

	// Step 6: Parse lunch data
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return ICanteenData{}, err
	}

	icanteenData := ICanteenData{
		Credit: doc.Find("#Kredit").Eq(0).Text(),
	}

	req, err = http.NewRequest(http.MethodGet, server+"/faces/secured/db/dbJidelnicekOnDayView.jsp?day="+req.URL.Query().Get("day"), nil)
	for key, value := range loginHeaders {
		req.Header.Set(key, value)
	}
	resp, err = client.Do(req)
	if err != nil {
		return ICanteenData{}, err
	}
	defer resp.Body.Close()

	// Cookies no longer valid
	if resp.Header.Get("Cache-Control") == "max-age=0" {
		return ICanteenData{}, http.ErrNoCookie
	}

	// Step 6: Parse lunch data
	doc, err = goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return ICanteenData{}, err
	}

	doc.Find(".orderContent").Each(func(index int, td *goquery.Selection) {
		day := strings.TrimPrefix(td.AttrOr("id", ""), "orderContent")
		lunches := []ICanteenLunch{}
		td.Find(".jidelnicekItemWrapper").Each(func(index int, lunch *goquery.Selection) {
			lunchEntry := ICanteenLunch{
				Name:      strings.TrimSpace(strings.Split(lunch.Children().Eq(1).Text(), "\n")[1]),
				Ordered:   lunch.Children().Eq(0).Children().Eq(0).Children().Length() == 5,
				CanOrder:  lunch.Children().Eq(2).Children().Length() != 3,
				ChangeURL: strings.Split(strings.TrimPrefix(strings.TrimSpace(lunch.Children().Eq(0).Children().Eq(0).AttrOr("onclick", "not found")), "ajaxOrder(this, '"), "'")[0],
			}
			lunches = append(lunches, lunchEntry)
		})
		dayEntry := ICanteenDay{
			Day:     day,
			Lunches: lunches,
		}
		icanteenData.Days = append(icanteenData.Days, dayEntry)
	})

	return icanteenData, nil
}

func TryLogin(username, password, server string) error {
	_, _, _, err := login(username, password, server)
	return err
}

func LoadLunches(username, password, server string) ([]ICanteenDay, error) {
	_, url, client, err := login(username, password, server)

	server, err = NormalizeServerURL(server)
	if err != nil {
		return nil, err
	}

	monthURL := server + "/faces/secured/month.jsp"
	req, err := http.NewRequest(http.MethodGet, monthURL, nil)
	if err != nil {
		return nil, err
	}
	loginHeaders := map[string]string{
		"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8",
		"Accept-Encoding": "text/html",
		"Cache-Control":   "max-age=0",
		"Connection":      "keep-alive",
		"Content-Type":    "application/x-www-form-urlencoded",
		"Host":            url.Host,
		"Origin":          server,
		"Referer":         server + "/login",
	}
	for key, value := range loginHeaders {
		req.Header.Set(key, value)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Step 6: Parse lunch data
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	lunchData := []ICanteenDay{}
	doc.Find(".orderContent").Each(func(index int, td *goquery.Selection) {
		day := strings.TrimPrefix(td.AttrOr("id", ""), "orderContent")
		lunches := []ICanteenLunch{}
		td.Find(".jidelnicekItemWrapper").Each(func(index int, lunch *goquery.Selection) {
			lunchEntry := ICanteenLunch{
				Name:     strings.TrimSpace(strings.Split(lunch.Children().Eq(1).Text(), "\n")[1]),
				Ordered:  lunch.Children().Eq(0).Children().Eq(0).Children().Length() == 5,
				CanOrder: lunch.Children().Eq(2).Children().Length() != 3,
			}
			lunches = append(lunches, lunchEntry)
		})
		dayEntry := ICanteenDay{
			Day:     day,
			Lunches: lunches,
		}
		lunchData = append(lunchData, dayEntry)
	})

	return lunchData, nil
}
