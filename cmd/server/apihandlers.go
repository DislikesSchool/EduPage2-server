package main

import (
	"net/http"
	"time"

	"github.com/DislikesSchool/EduPage2-server/edupage"
	"github.com/gin-gonic/gin"
)

// RecentTimelineHandler godoc
// @Summary Get the user's recent timeline
// @Schemes
// @Description Returns the user's timeline from today to 30 days in the past.
// @Tags timeline
// @Param token header string true "JWT token"
// @Produce json
// @Success 200 {object} Timeline
// @Failure 401 {object} RecentTimelineUnauthorizedResponse
// @Failure 500 {object} RecentTimelineInternalErrorResponse
// @Router /api/timeline/recent [get]
func RecentTimelineHandler(c *gin.Context) {
	client := c.MustGet("client").(*edupage.EdupageClient)

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
// @Success 200 {object} Timeline
// @Failure 401 {object} TimelineUnauthorizedResponse
// @Failure 500 {object} TimelineInternalErrorResponse
// @Router /api/timeline [get]
func TimelineHandler(c *gin.Context) {
	client := c.MustGet("client").(*edupage.EdupageClient)

	dateFromString := c.Query("from")
	dateToString := c.Query("to")

	var dateTo time.Time
	var dateFrom time.Time
	var err error

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

// RecentTimetableHandler godoc
// @Summary Get the user's recent timetable
// @Schemes
// @Description Returns the user's timetable from before yesterday to 7 days in the future.
// @Tags timetable
// @Param token header string true "JWT token"
// @Produce json
// @Success 200 {object} model.Timetable
// @Failure 401 {object} TimetableUnauthorizedResponse
// @Failure 500 {object} TimetableInternalErrorResponse
// @Router /api/timetable/recent [get]
func RecentTimetableHangler(c *gin.Context) {
	client := c.MustGet("client").(*edupage.EdupageClient)

	timetable, err := client.GetRecentTimetable()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, timetable)
}

// RecentTimetableHandler godoc
// @Summary Get the user's  timetable
// @Schemes
// @Description Returns the user's timetable from date specified to date specified or today.
// @Tags timetable
// @Param token header string true "JWT token"
// @Param range query TimetableRequest true "Date range"
// @Produce json
// @Success 200 {object} model.Timetable
// @Failure 401 {object} TimetableUnauthorizedResponse
// @Failure 500 {object} TimetableInternalErrorResponse
// @Router /api/timetable [get]
func TimetableHandler(c *gin.Context) {
	client := c.MustGet("client").(*edupage.EdupageClient)

	dateFromString := c.Query("from")
	dateToString := c.Query("to")

	var dateTo time.Time
	var dateFrom time.Time
	var err error

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

	timetable, err := client.GetTimetable(dateFrom, dateTo)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, timetable)
}

// SubjectHandler godoc
// @Summary Get the subject by ID
// @Schemes
// @Description Returns the subject by ID.
// @Tags DBI
// @Param token header string true "JWT token"
// @Param id param string true "Subject ID"
// @Produce json
// @Success 200 {object} Subject
// @Failure 401 {object} SubjectUnauthorizedResponse
// @Failure 500 {object} SubjectInternalErrorResponse
// @Router /api/subject/:id [get]
func SubjectHandler(c *gin.Context) {
	client := c.MustGet("client").(*edupage.EdupageClient)

	id := c.Param("id")

	subject, err := client.GetSubjectByID(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, subject)
}
