package edupage

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/publicsuffix"
)

var (
	ErrAuthorization   = errors.New("failed to authorize")
	ErrRedirect        = errors.New("redirect")
	ErrBadCredentials  = errors.New("bad credentials (wrong username or password)")
	ErrCaptcha         = errors.New("captcha required - log in via browser once, then retry")
	ErrTwoFactor       = errors.New("two-factor authentication required - approve in EduPage app or use email code")
	ErrLoginToken      = errors.New("EduPage did not provide a login token")
	ErrInvalidResponse = errors.New("invalid response from EduPage")
	ErrInvalidServer   = errors.New("invalid server: must be an edupage.org subdomain")
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
	Username     string
	Server       string
	LoginServer  string
	PasswordHash string
	httpClient   *http.Client
}

// validSubdomain matches a single DNS label (no dots, ports, userinfo, etc.)
// so that "https://<sub>.edupage.org" can only ever resolve to edupage.org.
var validSubdomain = regexp.MustCompile(`^[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?$`)

// isEdupageHost reports whether host is edupage.org or one of its subdomains.
func isEdupageHost(host string) bool {
	host = strings.ToLower(strings.TrimSuffix(host, "."))
	return host == edupageDomain || strings.HasSuffix(host, "."+edupageDomain)
}

func normalizeSubdomain(server string) string {
	server = strings.ToLower(strings.TrimSpace(server))
	server = strings.TrimPrefix(server, "http://")
	server = strings.TrimPrefix(server, "https://")
	// strip path
	if i := strings.Index(server, "/"); i != -1 {
		server = server[:i]
	}
	server = strings.TrimSuffix(server, ".edupage.org")
	server = strings.Trim(server, ".")
	if server == "" {
		server = "login1"
	}
	return server
}

func newEduClient(jar http.CookieJar) *http.Client {
	if jar == nil {
		var err error
		jar, err = cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
		if err != nil {
			jar = nil
		}
	}
	return &http.Client{
		Jar:     jar,
		Timeout: 20 * time.Second,
		// Follow redirects (default behavior). We inspect resp.Request.URL
		// to detect bad=1 / cap=1 / twofactor, like requests.Session does.
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return errors.New("too many redirects")
			}
			if req.URL.Scheme != "https" || !isEdupageHost(req.URL.Hostname()) {
				return fmt.Errorf("refusing redirect to non-edupage host %q", req.URL.Hostname())
			}
			// Preserve headers on redirect
			if len(via) > 0 {
				req.Header.Set("User-Agent", EduPageUserAgent)
				req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
			}
			return nil
		},
	}
}

func setCommonHeaders(req *http.Request, referer string) {
	req.Header.Set("User-Agent", EduPageUserAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "sk-SK,sk;q=0.9,cs;q=0.8,en;q=0.7")
	req.Header.Set("Cache-Control", "max-age=0")
	if referer != "" {
		req.Header.Set("Referer", referer)
	}
}

func httpGet(client *http.Client, rawurl string) (*http.Response, []byte, error) {
	req, err := http.NewRequest("GET", rawurl, nil)
	if err != nil {
		return nil, nil, err
	}
	setCommonHeaders(req, "")
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 10<<20))
	if err != nil {
		return resp, nil, err
	}
	return resp, body, nil
}

func httpPostFormRaw(client *http.Client, rawurl, body string, referer string) (*http.Response, []byte, error) {
	req, err := http.NewRequest("POST", rawurl, strings.NewReader(body))
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("User-Agent", EduPageUserAgent)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	if referer != "" {
		req.Header.Set("Referer", referer)
		req.Header.Set("Origin", referer)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(io.LimitReader(resp.Body, 10<<20))
	if err != nil {
		return resp, nil, err
	}
	return resp, b, nil
}

// Login creates EdupageClient credentials you can use to interact with the edupage api.
// It mirrors the current EduPage web login (mainlogin.js):
//  1. RPC: getToken -> login -> GET redirectUrl
//  2. Fallback: GET csrftoken -> POST edubarLogin.php
//
// Returns Credentials or error (ErrBadCredentials / ErrCaptcha / ErrTwoFactor).
func Login(username, password, server, loginserver string) (Credentials, error) {
	sub := normalizeSubdomain(server)
	if !validSubdomain.MatchString(sub) {
		return Credentials{}, ErrInvalidServer
	}
	if loginserver == "" {
		loginserver = sub
	} else {
		loginserver = normalizeSubdomain(loginserver)
		if !validSubdomain.MatchString(loginserver) {
			return Credentials{}, ErrInvalidServer
		}
	}
	fqdn := sub + "." + edupageDomain
	// Keep global for backwards compat (used by old callers).
	Server = fqdn

	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return Credentials{}, err
	}
	client := newEduClient(jar)

	// 1) Try modern RPC login first.
	if ok, rpcServer, err := loginWithRPC(client, username, password, sub); err == nil && ok {
		return buildCredentials(client, username, rpcServer, loginserver, password)
	} else {
		// Only fall through on recoverable errors. Hard failures
		// (bad creds, captcha, 2FA) must be returned directly.
		if errors.Is(err, ErrBadCredentials) || errors.Is(err, ErrCaptcha) || errors.Is(err, ErrTwoFactor) {
			return Credentials{}, err
		}
		// otherwise try legacy flow
		_ = rpcServer
	}

	// 2) Fallback to edubarLogin.php with csrftoken.
	actualServer, err := loginWithEdubar(client, username, password, sub)
	if err != nil {
		return Credentials{}, err
	}
	return buildCredentials(client, username, actualServer, loginserver, password)
}

func buildCredentials(client *http.Client, username, actualServerFQDN, loginserver, password string) (Credentials, error) {
	var credentials Credentials
	credentials.Username = username
	credentials.Server = actualServerFQDN
	credentials.LoginServer = loginserver
	hash, err := HashPassword(password)
	if err != nil {
		return Credentials{}, err
	}
	credentials.PasswordHash = hash
	credentials.httpClient = client
	// Update global for compat.
	Server = actualServerFQDN
	return credentials, nil
}

// loginWithRPC mirrors edupage-api's Login.__login_with_rpc.
func loginWithRPC(client *http.Client, username, password, subdomain string) (bool, string, error) {
	baseURL := fmt.Sprintf("https://%s.%s", subdomain, edupageDomain)

	// GET login page (sets cookies, validates host)
	resp, _, err := httpGet(client, baseURL+"/login/?cmd=MainLogin")
	if err != nil {
		return false, "", fmt.Errorf("get login page: %w", err)
	}
	if resp.StatusCode != 200 {
		return false, "", fmt.Errorf("get login page: status %d", resp.StatusCode)
	}

	// POST getToken
	tokenParams, _ := json.Marshal(map[string]string{
		"username": username,
		"edupage":  "",
	})
	tokenBody, err := EncodeRequestBody(map[string]string{
		"rpcparams": string(tokenParams),
	})
	if err != nil {
		return false, "", err
	}
	_, tokenRaw, err := httpPostFormRaw(client, baseURL+"/login/?cmd=MainLogin&akcia=getToken", tokenBody, baseURL+"/login/?cmd=MainLogin")
	if err != nil {
		return false, "", fmt.Errorf("getToken: %w", err)
	}
	tokenResp, err := ParseRPCResponse(string(tokenRaw))
	if err != nil || tokenResp == nil {
		return false, "", fmt.Errorf("getToken: bad response")
	}
	token, _ := tokenResp["token"].(string)
	if token == "" {
		return false, "", ErrLoginToken
	}

	// POST login
	loginMap := map[string]interface{}{
		"username":  username,
		"password":  password,
		"userToken": token,
		"edupage":   "",
		"ctxt":      "",
		"tu":        nil,
		"gu":        nil,
		"au":        nil,
	}
	loginJSON, _ := json.Marshal(loginMap)
	loginBody, err := EncodeRequestBody(map[string]string{
		"rpcparams": string(loginJSON),
	})
	if err != nil {
		return false, "", err
	}
	_, loginRaw, err := httpPostFormRaw(client, baseURL+"/login/?cmd=MainLogin&akcia=login", loginBody, baseURL+"/login/?cmd=MainLogin")
	if err != nil {
		return false, "", fmt.Errorf("login rpc: %w", err)
	}
	loginResp, err := ParseRPCResponse(string(loginRaw))
	if err != nil || loginResp == nil {
		return false, "", fmt.Errorf("login rpc: bad response")
	}
	var errorID string
	if errobj, ok := loginResp["err"].(map[string]interface{}); ok {
		if id, ok := errobj["error_id"].(string); ok {
			errorID = id
		}
	}
	redirectURL, _ := loginResp["redirectUrl"].(string)
	if errorID == "invalid_token" || redirectURL == "" {
		// Let caller try fallback (or treat as bad creds if explicit).
		if errorID != "" && errorID != "invalid_token" {
			// e.g. bad credentials come back without redirect
			return false, "", ErrBadCredentials
		}
		return false, "", fmt.Errorf("rpc failed, falling back")
	}

	// GET redirectUrl (may be relative)
	fullRedirect := redirectURL
	if strings.HasPrefix(redirectURL, "/") {
		fullRedirect = baseURL + redirectURL
	} else if !strings.HasPrefix(redirectURL, "http") {
		fullRedirect = baseURL + "/" + strings.TrimPrefix(redirectURL, "/")
	}
	if ru, err := url.Parse(fullRedirect); err != nil || ru.Scheme != "https" || !isEdupageHost(ru.Hostname()) {
		return false, "", ErrInvalidResponse
	}
	resp2, body2, err := httpGet(client, fullRedirect)
	if err != nil {
		return false, "", fmt.Errorf("get redirect: %w", err)
	}
	finalURL := ""
	if resp2.Request != nil && resp2.Request.URL != nil {
		finalURL = resp2.Request.URL.String()
	}
	// Detect captcha / bad / twofactor in final URL
	if strings.Contains(finalURL, "cap=1") || strings.Contains(finalURL, "lerr=b43b43") {
		return false, "", ErrCaptcha
	}
	if strings.Contains(finalURL, "bad=1") {
		return false, "", ErrBadCredentials
	}
	if strings.Contains(finalURL, "twofactor") {
		return false, "", ErrTwoFactor
	}
	// login1 portal: extract real school from final host
	actualSub := subdomain
	if subdomain == "login1" {
		if u, err := url.Parse(finalURL); err == nil && u.Hostname() != "" {
			parts := strings.Split(u.Hostname(), ".")
			if len(parts) >= 3 {
				actualSub = parts[0]
			}
		}
		// Also try to parse from body if redirect landed on generic page
		if actualSub == "login1" {
			if m := regexp.MustCompile(`https?://([a-z0-9\-]+)\.edupage\.org`).FindStringSubmatch(string(body2)); len(m) == 2 {
				actualSub = m[1]
			}
		}
	}
	actualFQDN := actualSub + "." + edupageDomain

	// Validate session: body should contain userhome or gsechash.
	// If not, fetch /user/ explicitly.
	bodyStr := string(body2)
	if !strings.Contains(bodyStr, "userhome(") && !strings.Contains(bodyStr, "gsechash") {
		u := fmt.Sprintf("https://%s/user/", actualFQDN)
		_, ubody, err := httpGet(client, u)
		if err == nil && (strings.Contains(string(ubody), "userhome(") || strings.Contains(string(ubody), "gsechash")) {
			return true, actualFQDN, nil
		}
		// If still no marker, consider login failed.
		// Distinguish bad creds: EduPage redirects to login?bad=1 which we
		// already checked; otherwise generic auth error.
		return false, "", ErrAuthorization
	}
	return true, actualFQDN, nil
}

var (
	reCsrftokenJSON = regexp.MustCompile(`"csrftoken"\s*:\s*"([^"]+)"`)
	reCsrfauthHTML  = regexp.MustCompile(`name="csrfauth"\s+value="([^"]+)"`)
	reCsrfauthHTML2 = regexp.MustCompile(`name='csrfauth'\s+value='([^']+)'`)
)

func extractCSRF(body string) string {
	if m := reCsrftokenJSON.FindStringSubmatch(body); len(m) == 2 {
		return m[1]
	}
	if m := reCsrfauthHTML.FindStringSubmatch(body); len(m) == 2 {
		return m[1]
	}
	if m := reCsrfauthHTML2.FindStringSubmatch(body); len(m) == 2 {
		return m[1]
	}
	return ""
}

// loginWithEdubar is the legacy fallback: GET csrftoken, POST edubarLogin.php.
func loginWithEdubar(client *http.Client, username, password, subdomain string) (string, error) {
	baseURL := fmt.Sprintf("https://%s.%s", subdomain, edupageDomain)

	resp, body, err := httpGet(client, baseURL+"/login/?cmd=MainLogin")
	if err != nil {
		return "", fmt.Errorf("get login page: %w", err)
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("get login page: status %d", resp.StatusCode)
	}
	csrf := extractCSRF(string(body))
	if csrf == "" {
		// Try older /login/index.php page
		_, body2, err := httpGet(client, baseURL+"/login/index.php")
		if err == nil {
			csrf = extractCSRF(string(body2))
		}
	}
	if csrf == "" {
		return "", ErrLoginToken
	}

	form := url.Values{
		"username": []string{username},
		"password": []string{password},
		"csrfauth": []string{csrf},
	}
	// POST with plain form (NOT eqap) - matches edupage-api fallback.
	req, err := http.NewRequest("POST", baseURL+"/login/edubarLogin.php", strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("User-Agent", EduPageUserAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Referer", baseURL+"/login/?cmd=MainLogin")
	req.Header.Set("Origin", baseURL)

	rsp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer rsp.Body.Close()
	finalBody, _ := io.ReadAll(io.LimitReader(rsp.Body, 10<<20))
	finalURL := ""
	if rsp.Request != nil && rsp.Request.URL != nil {
		finalURL = rsp.Request.URL.String()
	}
	// Also check Location header for non-followed redirects (defensive).
	if loc := rsp.Header.Get("Location"); loc != "" && finalURL == "" {
		finalURL = loc
	}

	if strings.Contains(finalURL, "cap=1") || strings.Contains(finalURL, "lerr=b43b43") {
		return "", ErrCaptcha
	}
	if strings.Contains(finalURL, "bad=1") {
		return "", ErrBadCredentials
	}
	if strings.Contains(finalURL, "twofactor") {
		return "", ErrTwoFactor
	}
	// Success is a redirect to /user/ (or dashboard). Old code expected
	// exactly "/user/"; be lenient: any page containing userhome/gsechash.
	actualSub := subdomain
	if subdomain == "login1" {
		if u, err := url.Parse(finalURL); err == nil && u.Hostname() != "" {
			parts := strings.Split(u.Hostname(), ".")
			if len(parts) >= 3 && parts[0] != "" && parts[0] != "login1" {
				actualSub = parts[0]
			}
		}
	}
	actualFQDN := actualSub + "." + edupageDomain

	// Validate: final page or fresh /user/ must contain login markers.
	combined := string(finalBody)
	if !strings.Contains(combined, "userhome(") && !strings.Contains(combined, "gsechash") {
		_, ubody, err := httpGet(client, fmt.Sprintf("https://%s/user/", actualFQDN))
		if err != nil {
			return "", ErrAuthorization
		}
		combined = string(ubody)
		if strings.Contains(finalURL, "bad=1") {
			return "", ErrBadCredentials
		}
		if !strings.Contains(combined, "userhome(") && !strings.Contains(combined, "gsechash") {
			// No session markers -> auth failed.
			return "", ErrAuthorization
		}
	}
	return actualFQDN, nil
}
