package main

import (
	"net/http"
	"time"

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

	timeline, err := client.GetRecentTimeline()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, timeline)
}

// TimelineHandler godoc
// @Summary Get the user's timeline
// @Schemes
// @Description Returns the user's timeline from any date to any other date or today.
// @Tags timeline
// @Param token header string true "JWT token"
// @Param range query TimelineRequest true "Date range"
// @Produce json
// @Success 200 {object} TimelineSuccessResponse
// @Failure 401 {object} TimelineUnauthorizedResponse
// @Failure 500 {object} TimelineInternalErrorResponse
// @Router /api/timeline [get]
func TimelineHandler(c *gin.Context) {
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

	dateFromString := c.Query("from")
	dateToString := c.Query("to")

	var dateTo time.Time
	var dateFrom time.Time

	if dateToString == "" {
		dateTo = time.Now()
	} else {
		dateTo, err = time.Parse(time.RFC3339, dateToString)
	}
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	dateFrom, err = time.Parse(time.RFC3339, dateFromString)

	timeline, err := client.GetTimeline(dateFrom, dateTo)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, timeline)
}
