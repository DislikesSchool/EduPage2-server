package edupage

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
)

func (h Handle) RefreshUser() {
	urlStr := fmt.Sprintf("https://%s/user/?", h.server)
	url, _ := url.Parse(urlStr)
	fmt.Println(url)
	rs, err := h.hc.Do(&http.Request{
		Method: "GET",
		URL:    url,
		Header: http.Header{
			"Referer":          []string{fmt.Sprintf("https://%s/", h.server)},
			"Accept":           []string{"application/json, text/javascript, */*; q=0.01"},
			"Content-Type":     []string{"application/x-www-form-urlencoded; charset=UTF-8"},
			"X-Requested-With": []string{"XMLHttpRequest"},
		},
	})
	if err != nil {
		fmt.Println(err)
	}

	defer rs.Body.Close()

	b, err := io.ReadAll(rs.Body)
	if err != nil {
		log.Fatalln(err)
	}

	html := string(b)

	// Parse raw JSON data from html
	_ = parse(html)
}

func parse(html string) RawDataObject {
	data := RawDataObject{
		Edubar: make(map[string]interface{}),
	}

	re := regexp.MustCompile(`\.userhome\((.+?)\);`)
	match := re.FindStringSubmatch(html)
	if len(match) < 2 {
		log.Fatalf("Failed to parse Edupage data from html: %s", html)
	}

	err := json.Unmarshal([]byte(match[1]), &data)
	if err != nil {
		log.Fatalf("Failed to parse JSON from Edupage html: %s, %s, %s", html, match[1], err)
	}

	// Parse additional edubar data
	re = regexp.MustCompile(`edubar\(([\s\S]*?)\);`)
	match = re.FindStringSubmatch(html)
	if len(match) < 2 {
		log.Fatalf("Failed to parse edubar data from html: %s", html)
	}

	err = json.Unmarshal([]byte(match[1]), &data.Edubar)
	if err != nil {
		log.Fatalf("Failed to parse JSON from edubar html: %s, %s, %s", html, match[1], err)
	}

	return data
}
