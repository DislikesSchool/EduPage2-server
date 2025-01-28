package main

import (
	"encoding/json"
	"flag"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var (
	username    string
	password    string
	server      string
	name        string
	ic_username string
	ic_password string
	ic_server   string
)

func init() {
	flag.StringVar(&username, "username", "", "Edupage username")
	flag.StringVar(&password, "password", "", "Edupage password")
	flag.StringVar(&server, "server", "", "Edupage server")
	flag.StringVar(&name, "name", "", "Name of the user (firstname lastname)")
	flag.StringVar(&ic_username, "ic_username", "", "iCanteen username")
	flag.StringVar(&ic_password, "ic_password", "", "iCanteen password")
	flag.StringVar(&ic_server, "ic_server", "", "iCanteen server")
}

func TestLoginAuto(t *testing.T) {
	if len(username) == 0 {
		t.Log("Username parameter missing, (-username=?)")
		return
	}

	if len(password) == 0 {
		t.Log("Password parameter missing, (-password=?)")
		return
	}

	gin.SetMode(gin.TestMode)

	// Test case 1: successful login
	data := url.Values{}
	data.Set("username", username)
	data.Set("password", password)
	req, _ := http.NewRequest("POST", "/login", strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	router := gin.Default()
	router.POST("/login", LoginHandler)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Error   string `json:"error"`
		Success bool   `json:"success"`
		Name    string `json:"name"`
		Token   string `json:"token"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "", response.Error)
	assert.True(t, response.Success)
	assert.Equal(t, name, response.Name)
	assert.NotEmpty(t, response.Token)
}

func TestLogin(t *testing.T) {
	if len(username) == 0 {
		t.Log("Username parameter missing, (-username=?)")
		return
	}

	if len(password) == 0 {
		t.Log("Password parameter missing, (-password=?)")
		return
	}

	if len(server) == 0 {
		t.Log("Server parameter missing, (-server=?)")
		return
	}

	gin.SetMode(gin.TestMode)

	// Test case 1: successful login
	data := url.Values{}
	data.Set("server", server)
	data.Set("username", username)
	data.Set("password", password)
	req, _ := http.NewRequest("POST", "/login", strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	router := gin.Default()
	router.POST("/login", LoginHandler)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Error   string `json:"error"`
		Success bool   `json:"success"`
		Name    string `json:"name"`
		Token   string `json:"token"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "", response.Error)
	assert.True(t, response.Success)
	assert.Equal(t, name, response.Name)
	assert.NotEmpty(t, response.Token)
}

func getAuthToken(t *testing.T) string {
	// Login to get auth token first
	data := url.Values{}
	data.Set("username", username)
	data.Set("password", password)
	req, _ := http.NewRequest("POST", "/login", strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	router := gin.Default()
	router.POST("/login", LoginHandler)
	router.ServeHTTP(w, req)

	var response struct {
		Token string `json:"token"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	return response.Token
}

func TestTimelineHandler(t *testing.T) {
	if len(username) == 0 || len(password) == 0 {
		t.Skip("Skipping test: credentials not provided")
	}

	gin.SetMode(gin.TestMode)
	token := getAuthToken(t)

	router := gin.Default()
	api := router.Group("/api")
	api.Use(authMiddleware())
	api.GET("/timeline/recent", RecentTimelineHandler)

	req, _ := http.NewRequest("GET", "/api/timeline/recent", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestTimetableHandler(t *testing.T) {
	if len(username) == 0 || len(password) == 0 {
		t.Skip("Skipping test: credentials not provided")
	}

	gin.SetMode(gin.TestMode)
	token := getAuthToken(t)

	router := gin.Default()
	api := router.Group("/api")
	api.Use(authMiddleware())
	api.GET("/timetable/recent", RecentTimetableHangler)

	req, _ := http.NewRequest("GET", "/api/timetable/recent", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGradesHandler(t *testing.T) {
	if len(username) == 0 || len(password) == 0 {
		t.Skip("Skipping test: credentials not provided")
	}

	gin.SetMode(gin.TestMode)
	token := getAuthToken(t)

	router := gin.Default()
	api := router.Group("/api")
	api.Use(authMiddleware())
	api.GET("/grades", ResultsHandler)

	req, _ := http.NewRequest("GET", "/api/grades", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPeriodsHandler(t *testing.T) {
	if len(username) == 0 || len(password) == 0 {
		t.Skip("Skipping test: credentials not provided")
	}

	gin.SetMode(gin.TestMode)
	token := getAuthToken(t)

	router := gin.Default()
	api := router.Group("/api")
	api.Use(authMiddleware())
	api.GET("/periods", PeriodsHandler)

	req, _ := http.NewRequest("GET", "/api/periods", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRecipientsHandler(t *testing.T) {
	if len(username) == 0 || len(password) == 0 {
		t.Skip("Skipping test: credentials not provided")
	}

	gin.SetMode(gin.TestMode)
	token := getAuthToken(t)

	router := gin.Default()
	api := router.Group("/api")
	api.Use(authMiddleware())
	api.GET("/recipients", RecipientsHandler)

	req, _ := http.NewRequest("GET", "/api/recipients", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
