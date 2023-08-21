package edupage

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"path"

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
				return h, nil
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
