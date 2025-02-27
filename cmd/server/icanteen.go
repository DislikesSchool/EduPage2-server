package main

import (
	"encoding/json"
	"net/http"

	"github.com/DislikesSchool/EduPage2-server/icanteen"
	"github.com/gin-gonic/gin"
)

// ICanteenLoginHandler godoc
// @Summary iCanteen login
// @Schemes
// @Description Logs in to iCanteen and returns the cookies
// @Tags lunches
// @Accept multipart/form-data
// @Accept x-www-form-urlencoded
// @Consumes application/x-www-form-urlencoded
// @Param username formData string true "Username"
// @Param password formData string true "Password"
// @Param server formData string true "Server"
// @Produce json
// @Success 200 {object} []http.Cookie
// @Failure 400 {object} apimodel.ICanteenBadRequestResponse
// @Failure 500 {object} apimodel.ICanteenInternalErrorResponse
// @Router /icanteen/login [post]
func ICanteenLoginHandler(ctx *gin.Context) {
	var username string
	var password string
	var server string

	if username = ctx.PostForm("username"); username == "" {
		ctx.AbortWithStatusJSON(400, gin.H{"error": "username is missing"})
		return
	}

	if password = ctx.PostForm("password"); password == "" {
		ctx.AbortWithStatusJSON(400, gin.H{"error": "password is missing"})
		return
	}

	if server = ctx.PostForm("server"); server == "" {
		ctx.AbortWithStatusJSON(400, gin.H{"error": "server is missing"})
		return
	}

	cookies, err := icanteen.Login(username, password, server)
	if err != nil {
		ctx.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, cookies)
}

// ICanteenMonthHandler godoc
// @Summary Fetch lunches for month
// @Schemes
// @Description Fetches the lunches for the month using the provided cookies
// @Tags lunches
// @Accept multipart/form-data
// @Accept x-www-form-urlencoded
// @Consumes application/x-www-form-urlencoded
// @Param cookies formData string true "JSON object with cookies from login"
// @Param server formData string true "Server"
// @Produce json
// @Success 200 {object} icanteen.ICanteenData
// @Failure 400 {object} apimodel.ICanteenBadRequestResponse
// @Failure 500 {object} apimodel.ICanteenInternalErrorResponse
// @Router /icanteen/month [post]
func ICanteenMonthHandler(ctx *gin.Context) {
	var cookies string
	var server string

	if cookies = ctx.PostForm("cookies"); cookies == "" {
		ctx.AbortWithStatusJSON(400, gin.H{"error": "cookies are missing"})
		return
	}

	if server = ctx.PostForm("server"); server == "" {
		ctx.AbortWithStatusJSON(400, gin.H{"error": "server is missing"})
		return
	}

	parsedCookies, err := parseCookies(cookies)
	if err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{"error": "failed to parse cookies"})
		return
	}

	lunches, err := icanteen.LoadLunchesWithCookies(parsedCookies, server)
	if err != nil {
		if err == http.ErrNoCookie {
			ctx.AbortWithStatusJSON(400, gin.H{"error": "cookies are no longer valid"})
			return
		} else {
			ctx.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
			return
		}
	}

	ctx.JSON(200, lunches)
}

// ICanteenChangeOrderHandler godoc
// @Summary Change lunch order
// @Schemes
// @Description Changes the lunch order using the provided cookies and changeURL
// @Tags lunches
// @Accept multipart/form-data
// @Accept x-www-form-urlencoded
// @Consumes application/x-www-form-urlencoded
// @Param cookies formData string true "JSON object with cookies from login"
// @Param server formData string true "Server"
// @Param changeURL formData string true "URL to change the order"
// @Produce json
// @Success 200 {object} icanteen.ICanteenData
// @Failure 400 {object} apimodel.ICanteenBadRequestResponse
// @Failure 500 {object} apimodel.ICanteenInternalErrorResponse
// @Router /icanteen/change [post]
func ICanteenChangeOrderHandler(ctx *gin.Context) {
	var cookies string
	var server string
	var changeURL string

	if cookies = ctx.PostForm("cookies"); cookies == "" {
		ctx.AbortWithStatusJSON(400, gin.H{"error": "cookies are missing"})
		return
	}

	if server = ctx.PostForm("server"); server == "" {
		ctx.AbortWithStatusJSON(400, gin.H{"error": "server is missing"})
		return
	}

	if changeURL = ctx.PostForm("changeURL"); changeURL == "" {
		ctx.AbortWithStatusJSON(400, gin.H{"error": "changeURL is missing"})
		return
	}

	parsedCookies, err := parseCookies(cookies)
	if err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{"error": "failed to parse cookies"})
		return
	}

	lunches, err := icanteen.ChangeOrder(parsedCookies, server, changeURL)
	if err != nil {
		if err == http.ErrNoCookie {
			ctx.AbortWithStatusJSON(400, gin.H{"error": "cookies are no longer valid"})
			return
		} else {
			ctx.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
			return
		}
	}

	ctx.JSON(200, lunches)
}

func parseCookies(cookies string) ([]*http.Cookie, error) {
	var parsedCookies []*http.Cookie
	err := json.Unmarshal([]byte(cookies), &parsedCookies)
	if err != nil {
		return nil, err
	}

	return parsedCookies, nil
}
