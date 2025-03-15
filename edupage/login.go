package edupage

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"path"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/publicsuffix"
)

var (
	ErrAuthorization = errors.New("failed to authorize")
	ErrRedirect      = errors.New("redirect")
)

var (
	edupageDomain = "edupage.org"
	Server        = ""
	loginPath     = "login/edubarLogin.php"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

type Credentials struct {
	Server       string
	PasswordHash string
	httpClient   *http.Client
}

// Login creates EdupageClient you can use to interact the edupage api with.
// Returns EdupageClient or error.
func Login(username, password, server string) (Credentials, error) {
	server = strings.TrimPrefix(server, "http://")
	server = strings.TrimPrefix(server, "https://")
	server = strings.TrimSuffix(server, ".edupage.org")
	Server = server + "." + edupageDomain

	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return Credentials{}, err
	}

	client := &http.Client{
		Jar: jar,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return ErrRedirect
		},
	}

	u := fmt.Sprintf("https://%s", path.Join(Server, loginPath))
	d := url.Values{
		"username": []string{username},
		"password": []string{password},
	}

	rs, err := client.PostForm(u, d)

	if rs != nil && err != nil {
		if rs.StatusCode == 302 {
			if rs.Header.Get("Location") != "/user/" {
				loc := rs.Header.Get("Location")
				parsed, err := url.Parse(loc)
				if err != nil {
					return Credentials{}, err
				}

				sp := strings.Split(parsed.Hostname(), ".")
				sub := sp[0]

				return Login(username, password, sub)
			} else {
				var credentials Credentials
				credentials.Server = Server
				credentials.PasswordHash, err = HashPassword(password)
				if err != nil {
					return Credentials{}, err
				}
				credentials.httpClient = client

				return credentials, nil
			}
		} else {
			return Credentials{}, err
		}
	}
	return Credentials{}, err
}
