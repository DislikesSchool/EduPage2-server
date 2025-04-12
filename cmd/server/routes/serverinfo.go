package routes

import (
	"net/http"

	"github.com/DislikesSchool/EduPage2-server/config"
	"github.com/gin-gonic/gin"
)

// ServerVersion godoc
// @Summary Get server version information
// @Description Returns the server version details from the version.json file
// @Tags server,info
// @Produce json
// @Success 200 {object} object "Server version information"
// @Router /api/version [get]
func ServerVersion(c *gin.Context) {
	c.Status(http.StatusOK)
	c.File("./cmd/server/web/version.json")
}

// ServerCapabilities godoc
// @Summary Get server capabilities
// @Description Returns information about enabled server features like cache and storage
// @Tags server,info
// @Produce json
// @Success 200 {object} object "Server capabilities with cache and storage status"
// @Router /api/capabilities [get]
func ServerCapabilities(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"cache":      config.AppConfig.Redis.Enabled,
		"storage":    config.AppConfig.Database.Enabled,
		"encryption": config.AppConfig.Encryption.Enabled,
	})
}
