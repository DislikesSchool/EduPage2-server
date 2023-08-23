package edupage

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"path"
	"strings"

	"golang.org/x/net/publicsuffix"
)

var (
	ErrAuthorization = errors.New("failed to authorize")
)

var (
	edupageDomain = "edupage.org"
	Server        = ""
	loginPath     = "login/edubarLogin.php"
)

// Login creates a Handle which you can interact the Edupage API with.
// Returns Handle or error.
func Login(server, username, password string) (Handle, error) {
	Server = server + "." + edupageDomain
	var h Handle
	h.hc = http.DefaultClient
	h.hc.CheckRedirect = noRedirect
	h.hc.Jar, _ = cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})

	u := fmt.Sprintf("https://%s", path.Join(Server, loginPath))
	d := url.Values{
		"username": []string{username},
		"password": []string{password},
	}

	rs, err := http.PostForm(u, d)

	if err != nil && rs != nil {
		if rs.StatusCode == 302 {
			if rs.Header.Get("Location") != "/user/" {
				return Handle{}, ErrAuthorization
			} else if rs.Header.Get("Location") == "/user/" {
				h.hc.Jar.SetCookies(rs.Request.URL, rs.Cookies())
				h.server = Server
				return h, nil
			}
		} else {
			return Handle{}, fmt.Errorf("failed to login: %s", err)
		}
	}
	return Handle{}, errors.New("unexpected response from server, make sure credentials are specified correctly")
}

func LoginAuto(username string, password string) (Handle, error) {
	var h Handle
	h.hc = http.DefaultClient
	h.hc.CheckRedirect = noRedirect
	h.hc.Jar, _ = cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})

	u := "https://portal.edupage.org/index.php?jwid=jw3&module=Login&lang=sk"
	d := url.Values{
		"meno":  []string{username},
		"heslo": []string{password},
		"akcia": []string{"login"},
	}

	rs, err := http.PostForm(u, d)
	if err != nil && rs != nil {
		if rs.StatusCode == 302 {
			if !strings.Contains(rs.Header.Get("Location"), "edupage.org/user/") {
				return Handle{}, ErrAuthorization
			} else if strings.Contains(rs.Header.Get("Location"), "edupage.org/user/") {
				domain := strings.Split(rs.Header.Get("Location"), "/")[2]
				h2, err := Login(strings.Split(domain, ".")[0], username, password)
				return h2, err
			}
		} else {
			return Handle{}, fmt.Errorf("failed to login: %s", err)
		}
	}
	return Handle{}, errors.New("unexpected response from server, make sure server is specified correctly")
}

func noRedirect(_ *http.Request, _ []*http.Request) error {
	return errors.New("redirect")
}
