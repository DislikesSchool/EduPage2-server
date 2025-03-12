package main

import (
	"net/http"

	"github.com/DislikesSchool/EduPage2-server/cmd/server/apimodel"
	"github.com/DislikesSchool/EduPage2-server/edupage"
	"github.com/gin-gonic/gin"
)

// RecipientsHandler godoc
// @Summary Get recipients
// @Schemes
// @Description Returns the possible recipients for messages.
// @Tags messages
// @Param Authorization header string true "JWT token"
// @Produce json
// @Security Bearer
// @Success 200 {object} []apimodel.Recipient
// @Failure 401 {object} apimodel.UnauthorizedResponse
// @Failure 500 {object} apimodel.InternalErrorResponse
// @Router /api/recipients [get]
func RecipientsHandler(c *gin.Context) {
	client := c.MustGet("client").(*edupage.EdupageClient)

	user, err := client.GetUser(false)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	students := user.DBI.Students
	teachers := user.DBI.Teachers

	recipients := make([]apimodel.Recipient, len(students)+len(teachers))

	i := 0
	for _, student := range students {
		recipients[i] = apimodel.Recipient{
			ID:   student.ID,
			Type: "student",
			Name: student.Firstname + " " + student.Lastname,
		}
		i++
	}

	for _, teacher := range teachers {
		recipients[i] = apimodel.Recipient{
			ID:   teacher.ID,
			Type: "teacher",
			Name: teacher.Firstname + " " + teacher.Lastname,
		}
		i++
	}

	c.JSON(http.StatusOK, recipients)
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

// PeriodsHandler godoc
// @Summary Get the school's periods
// @Schemes
// @Description Returns the school's periods.
// @Tags DBI
// @Param Authorization header string true "JWT token"
// @Produce json
// @Security Bearer
// @Success 200 {object} map[string]model.Period
// @Failure 401 {object} apimodel.UnauthorizedResponse
// @Failure 500 {object} apimodel.InternalErrorResponse
// @Router /api/periods [get]
func PeriodsHandler(c *gin.Context) {
	client := c.MustGet("client").(*edupage.EdupageClient)

	user, err := client.GetUser(false)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user.DBI.Periods)
}
