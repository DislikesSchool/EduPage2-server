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
