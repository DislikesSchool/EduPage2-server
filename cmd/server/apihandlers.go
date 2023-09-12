package main

import (
	"net/http"
	"time"

	"github.com/DislikesSchool/EduPage2-server/cmd/server/apimodel"
	"github.com/DislikesSchool/EduPage2-server/edupage"
	"github.com/DislikesSchool/EduPage2-server/edupage/model"
	"github.com/gin-gonic/gin"
)

// RecentTimelineHandler godoc
// @Summary Get the user's recent timeline
// @Schemes
// @Description Returns the user's timeline from today to 30 days in the past.
// @Tags timeline
// @Param Authorization header string true "JWT token"
// @Produce json
// @Security Bearer
// @Success 200 {object} model.Timeline
// @Failure 401 {object} apimodel.UnauthorizedResponse
// @Failure 500 {object} apimodel.InternalErrorResponse
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
// @Param Authorization header string true "JWT token"
// @Param range query apimodel.TimelineRequest true "Date range"
// @Produce json
// @Security Bearer
// @Success 200 {object} model.Timeline
// @Failure 401 {object} apimodel.UnauthorizedResponse
// @Failure 500 {object} apimodel.InternalErrorResponse
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
// @Param Authorization header string true "JWT token"
// @Produce json
// @Security Bearer
// @Success 200 {object} apimodel.CompleteTimetable
// @Failure 401 {object} apimodel.UnauthorizedResponse
// @Failure 500 {object} apimodel.InternalErrorResponse
// @Router /api/timetable/recent [get]
func RecentTimetableHangler(c *gin.Context) {
	client := c.MustGet("client").(*edupage.EdupageClient)

	timetable, err := client.GetRecentTimetable()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user, err := client.GetUser(false)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	completeTimetable := apimodel.CompleteTimetable{
		Days: make(map[string][]apimodel.CompleteTimetableItem, len(timetable.Days)),
	}

	for date, items := range timetable.Days {
		for _, item := range items {
			completeItem := apimodel.CompleteTimetableItem{
				Type:       item.Type,
				Date:       item.Date,
				Period:     item.Period,
				StartTime:  item.StartTime,
				EndTime:    item.EndTime,
				Subject:    user.DBI.Subjects[item.SubjectID],
				Classes:    make([]model.Class, len(item.ClassIDs)),
				GroupNames: item.GroupNames,
				IGroupID:   item.IGroupID,
				Teachers:   make([]model.Teacher, len(item.TeacherIDs)),
				Classrooms: make([]model.Classroom, len(item.ClassroomIDs)),
				StudentIDs: item.StudentIDs,
				Colors:     item.Colors,
			}

			for i, classID := range item.ClassIDs {
				completeItem.Classes[i] = user.DBI.Classes[classID]
			}

			for i, teacherID := range item.TeacherIDs {
				completeItem.Teachers[i] = user.DBI.Teachers[teacherID]
			}

			for i, classroomID := range item.ClassroomIDs {
				completeItem.Classrooms[i] = user.DBI.Classrooms[classroomID]
			}

			if original, ok := completeTimetable.Days[date]; ok {
				completeTimetable.Days[date] = append(original, completeItem)
			} else {
				completeTimetable.Days[date] = []apimodel.CompleteTimetableItem{completeItem}
			}
		}
	}

	c.JSON(http.StatusOK, completeTimetable)
}

// RecentTimetableHandler godoc
// @Summary Get the user's  timetable
// @Schemes
// @Description Returns the user's timetable from date specified to date specified or today.
// @Tags timetable
// @Param Authorization header string true "JWT token"
// @Param range query apimodel.TimetableRequest true "Date range"
// @Produce json
// @Security Bearer
// @Success 200 {object} apimodel.CompleteTimetable
// @Failure 401 {object} apimodel.UnauthorizedResponse
// @Failure 500 {object} apimodel.InternalErrorResponse
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

	user, err := client.GetUser(false)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	completeTimetable := apimodel.CompleteTimetable{
		Days: make(map[string][]apimodel.CompleteTimetableItem, len(timetable.Days)),
	}

	for date, items := range timetable.Days {
		for _, item := range items {
			completeItem := apimodel.CompleteTimetableItem{
				Type:       item.Type,
				Date:       item.Date,
				Period:     item.Period,
				StartTime:  item.StartTime,
				EndTime:    item.EndTime,
				Subject:    user.DBI.Subjects[item.SubjectID],
				Classes:    make([]model.Class, len(item.ClassIDs)),
				GroupNames: item.GroupNames,
				IGroupID:   item.IGroupID,
				Teachers:   make([]model.Teacher, len(item.TeacherIDs)),
				Classrooms: make([]model.Classroom, len(item.ClassroomIDs)),
				StudentIDs: item.StudentIDs,
				Colors:     item.Colors,
			}

			for i, classID := range item.ClassIDs {
				completeItem.Classes[i] = user.DBI.Classes[classID]
			}

			for i, teacherID := range item.TeacherIDs {
				completeItem.Teachers[i] = user.DBI.Teachers[teacherID]
			}

			for i, classroomID := range item.ClassroomIDs {
				completeItem.Classrooms[i] = user.DBI.Classrooms[classroomID]
			}

			if original, ok := completeTimetable.Days[date]; ok {
				completeTimetable.Days[date] = append(original, completeItem)
			} else {
				completeTimetable.Days[date] = []apimodel.CompleteTimetableItem{completeItem}
			}
		}
	}

	c.JSON(http.StatusOK, completeTimetable)
}

// SubjectHandler godoc
// @Summary Get the subject by ID
// @Schemes
// @Description Returns the subject by ID.
// @Tags DBI
// @Param Authorization header string true "JWT token"
// @Param id path string true "Subject ID"
// @Produce json
// @Security Bearer
// @Success 200 {object} model.Subject
// @Failure 401 {object} apimodel.UnauthorizedResponse
// @Failure 500 {object} apimodel.InternalErrorResponse
// @Router /api/subject/{id} [get]
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

// TeacherHandler godoc
// @Summary Get the teacher by ID
// @Schemes
// @Description Returns the teacher by ID.
// @Tags DBI
// @Param Authorization header string true "JWT token"
// @Param id path string true "Teacher ID"
// @Produce json
// @Security Bearer
// @Success 200 {object} model.Teacher
// @Failure 401 {object} apimodel.UnauthorizedResponse
// @Failure 500 {object} apimodel.InternalErrorResponse
// @Router /api/teacher/{id} [get]
func TeacherHandler(c *gin.Context) {
	client := c.MustGet("client").(*edupage.EdupageClient)

	id := c.Param("id")

	subject, err := client.GetTeacherByID(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, subject)
}

// ClassroomHandler godoc
// @Summary Get the classroom by ID
// @Schemes
// @Description Returns the classroom by ID.
// @Tags DBI
// @Param Authorization header string true "JWT token"
// @Param id path string true "Classroom ID"
// @Produce json
// @Security Bearer
// @Success 200 {object} model.Classroom
// @Failure 401 {object} apimodel.UnauthorizedResponse
// @Failure 500 {object} apimodel.InternalErrorResponse
// @Router /api/classroom/{id} [get]
func ClassroomHandler(c *gin.Context) {
	client := c.MustGet("client").(*edupage.EdupageClient)

	id := c.Param("id")

	subject, err := client.GetClassroomByID(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, subject)
}

// ICanteenHandler godoc
// @Summary Load lunches from iCanteen
// @Schemes
// @Description Loads the lunches from iCanteen for the next month.
// @Tags lunches
// @Accept multipart/form-data
// @Accept x-www-form-urlencoded
// @Consumes application/x-www-form-urlencoded
// @Param username formData string true "Username"
// @Param password formData string true "Password"
// @Param server formData string true "Server"
// @Produce json
// @Success 200 {object} []edupage.ICanteenDay
// @Failure 400 {object} apimodel.ICanteenBadRequestResponse
// @Failure 500 {object} apimodel.ICanteenInternalErrorResponse
// @Router /icanteen [post]
func ICanteenHandler(ctx *gin.Context) {
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

	lunches, err := edupage.LoadLunches(username, password, server)
	if err != nil {
		ctx.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, lunches)
}
