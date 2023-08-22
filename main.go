package main

import (
	"DislikesSchool/EduPage2-server/edupage"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.GET("/login", func(c *gin.Context) {
		username, password, server := c.Query("username"), c.Query("password"), c.Query("server")
		var handle edupage.Handle
		var err error
		if server == "" {
			handle, err = edupage.LoginAuto(username, password)
		} else {
			handle, err = edupage.Login(server, username, password)
		}
		if err != nil {
			c.JSON(200, gin.H{
				"message": "failed",
			})
		} else {
			data := handle.RefreshUser()
			c.JSON(200, gin.H{
				"message": "success",
				"user":    data.UserRow.Firstname + " " + data.UserRow.Lastname,
			})
		}
	})

	router.Run()
}
