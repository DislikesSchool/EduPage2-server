package edupage

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"path"

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

type Credentials struct {
	Server     string
	httpClient *http.Client
}

type mAppData struct {
	LoggedUser         string              `json:"loggedUser"`
	LoggedChild        int                 `json:"loggedChild"`
	LoggedUserName     string              `json:"loggedUserName"`
	Lang               string              `json:"lang"`
	Edupage            string              `json:"edupage"`
	SchoolType         string              `json:"school_type"`
	TimezoneDiff       int                 `json:"timezonediff"`
	SchoolCountry      string              `json:"school_country"`
	SchoolyearTurnover string              `json:"schoolyear_turnover"`
	FirstDayOfWeek     int                 `json:"firstDayOfWeek"`
	SortNameCol        string              `json:"sort_name_col"`
	SelectedYear       int                 `json:"selectedYear"`
	AutoYear           int                 `json:"autoYear"`
	YearTurnover       string              `json:"year_turnover"`
	VyucovacieDni      []bool              `json:"vyucovacieDni"`
	Server             string              `json:"server"`
	SyncIntervalMult   int                 `json:"syncIntervalMultiplier"`
	Ascspl             interface{}         `json:"ascspl"`
	JePro              bool                `json:"jePro"`
	JeZUS              bool                `json:"jeZUS"`
	Rtl                bool                `json:"rtl"`
	RtlAvailable       bool                `json:"rtlAvailable"`
	Uidsgn             string              `json:"uidsgn"`
	Webpageadmin       bool                `json:"webpageadmin"`
	EduRequestProps    edupageRequestProps `json:"edurequestProps"`
	Gsechash           string              `json:"gsechash"`
	Email              string              `json:"email"`
	Userrights         []interface{}       `json:"userrights"`
	IsAdult            bool                `json:"isAdult"`
}

type edupageRequestProps struct {
	Edupage            string        `json:"edupage"`
	Lang               string        `json:"lang"`
	SchoolName         string        `json:"school_name"`
	SchoolCountry      string        `json:"school_country"`
	SchoolState        string        `json:"school_state"`
	SchoolyearTurnover string        `json:"schoolyear_turnover"`
	CustomTurnover     []interface{} `json:"custom_turnover"`
	FirstDayOfWeek     int           `json:"firstDayOfWeek"`
	WeekendDays        []int         `json:"weekendDays"`
	Timezone           string        `json:"timezone"`
	SortNameCol        string        `json:"sort_name_col"`
	DtFormats          struct {
		Date string `json:"date"`
		Time string `json:"time"`
	} `json:"dtFormats"`
	Jsmodulemode     string        `json:"jsmodulemode"`
	LoggedUser       string        `json:"loggedUser"`
	LoggedUserRights []interface{} `json:"loggedUserRights"`
	IsAsc            bool          `json:"isAsc"`
	IsAgenda         bool          `json:"isAgenda"`
}

type mAuthUser struct {
	Userid       string   `json:"userid"`
	Typ          string   `json:"typ"`
	Edupage      string   `json:"edupage"`
	Edumeno      string   `json:"edumeno"`
	Eduheslo     string   `json:"eduheslo"`
	Firstname    string   `json:"firstname"`
	Lastname     string   `json:"lastname"`
	Esid         string   `json:"esid"`
	Appdata      mAppData `json:"appdata"`
	PortalUserid string   `json:"portal_userid"`
	PortalEmail  string   `json:"portal_email"`
	Need2fa      *string  `json:"need2fa,omitempty"`
}

type mAuthResponse struct {
	Users       []mAuthUser `json:"users"`
	NeedEdupage bool        `json:"needEdupage"`
	Edid        string      `json:"edid"`
	T2FASec     interface{} `json:"t2fasec,omitempty"`
}

// Login creates EdupageClient you can use to interact the edupage api with.
// Returns EdupageClient or error.
func Login(username, password, server string) (Credentials, error) {
	Server = server + "." + edupageDomain

	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return Credentials{}, err
	}

	client := &http.Client{
		CheckRedirect: noRedirect,
		Jar:           jar,
	}

	u := fmt.Sprintf("https://%s", path.Join(Server, loginPath))
	d := url.Values{
		"username": []string{username},
		"password": []string{password},
	}

	rs, err := client.PostForm(u, d)

	if err != nil && rs != nil {
		if rs.StatusCode == 302 {
			if rs.Header.Get("Location") != "/user/" {
				return Credentials{}, ErrAuthorization
			} else if rs.Header.Get("Location") == "/user/" {
				var credentials Credentials
				credentials.Server = Server
				credentials.httpClient = client

				return credentials, nil
			}
		} else {
			return Credentials{}, fmt.Errorf("failed to login: %s", err)
		}
	}

	return Credentials{}, errors.New("unexpected response from server, make sure credentials are specified correctly")
}

func LoginAuto(username, password string) (Credentials, error) {
	if len(username) == 0 || len(password) == 0 {
		return Credentials{}, errors.New("invalid credentials")
	}

	var credentials Credentials

	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return Credentials{}, err
	}

	client := &http.Client{
		CheckRedirect: noRedirect,
		Jar:           jar,
	}

	payload := url.Values{
		"m":             {username},
		"h":             {password},
		"edupage":       {""},
		"plgc":          {""},
		"ajheslo":       {"1"},
		"hasujheslo":    {"1"},
		"ajportal":      {"1"},
		"ajportallogin": {"1"},
		"mobileLogin":   {"1"},
		"version":       {"2020.0.18"},
		"fromEdupage":   {""},
		"device_name":   {"EduPage2 Public backend server"},
		"device_id":     {""},
		"device_key":    {""},
		"os":            {"Linux"},
		"murl":          {""},
		"edid":          {""},
	}

	loginServer := "login1"
	skip2Fa := true

	resp, err := http.PostForm(fmt.Sprintf("https://%s.edupage.org/login/mauth", loginServer), payload)
	if err != nil {
		return Credentials{}, err
	}
	defer resp.Body.Close()

	var authResponse mAuthResponse
	err = json.NewDecoder(resp.Body).Decode(&authResponse)
	if err != nil {
		return Credentials{}, err
	}

	//Error handling
	if len(authResponse.Users) == 0 {
		return Credentials{}, errors.New("failed to login: Incorrect password. (If you are sure that the password is correct, try providing 'edupage' option)")
	}

	//Process response
	if len(authResponse.Users) == 1 {
		if authResponse.Users[0].Need2fa != nil && !skip2Fa {
			if authResponse.T2FASec != "" {
				log.Printf("[Login] 2FA code is invalid\n")
				return Credentials{}, errors.New("invalid 2FA code")
			} else {
				log.Printf("[Login] 2FA was requested by the Edupage\n")
				return Credentials{}, nil
			}
		}
	} else {
		return Credentials{}, errors.New("multiple users found. Please, pass the selected user as 'user' option to login options")
	}

	origin := authResponse.Users[0].Edupage
	credentials.Server = origin + "." + edupageDomain
	u, err := url.Parse(fmt.Sprintf("https://%s/login/edubarLogin.php", credentials.Server))
	if err != nil {
		return Credentials{}, err
	}
	cookies := resp.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "PHPSESSID" {
			cookie.Value = authResponse.Users[0].Esid
		}
	}

	client.Jar.SetCookies(u, cookies)
	credentials.httpClient = client
	return credentials, nil
}

func noRedirect(_ *http.Request, _ []*http.Request) error {
	return ErrRedirect
}
