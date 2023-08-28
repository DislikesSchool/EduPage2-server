package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RecentTimelineHandler godoc
// @Summary Get the user's recent timeline
// @Schemes
// @Description Returns the user's timeline from today to 30 days in the past.
// @Tags timeline
// @Param token header string true "JWT token"
// @Produce json
// @Success 200 {object} RecentTimelineSuccessResponse
// @Failure 401 {object} RecentTimelineUnauthorizedResponse
// @Failure 500 {object} RecentTimelineInternalErrorResponse
// @Router /api/timeline/recent [get]
func RecentTimelineHandler(c *gin.Context) {
	claims, err := getClaims(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	userID := claims["userID"].(string)
	username := claims["username"].(string)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	client, ok := clients[userID+username]
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "client not found"})
		return
	}

	err = client.LoadRecentTimeline()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, client.Timeline)
}
