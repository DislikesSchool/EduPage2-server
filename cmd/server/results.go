package main

import (
	"net/http"
	"time"

	"github.com/DislikesSchool/EduPage2-server/edupage"
	"github.com/gin-gonic/gin"
)

// ResultsHandler godoc
// @Summary Get the user's grades
// @Schemes
// @Description Returns the user's grades.
// @Tags grades
// @Param Authorization header string true "JWT token"
// @Param year query string false "Year"
// @Param half query string false "Half"
// @Produce json
// @Security Bearer
// @Success 200 {object} []model.Results
// @Failure 401 {object} apimodel.UnauthorizedResponse
// @Failure 500 {object} apimodel.InternalErrorResponse
// @Router /api/grades [get]
func ResultsHandler(c *gin.Context) {
	client := c.MustGet("client").(*edupage.EdupageClient)

	year := c.Query("year")
	half := c.Query("half")

	if year == "" {
		month := time.Now().Month()
		if month >= time.January && month <= time.August {
			year = time.Now().AddDate(-1, 0, 0).Format("2006")
		}
		if month >= time.September && month <= time.December {
			year = time.Now().Format("2006")
		}
	}
	if half == "" {
		month := time.Now().Month()
		if month == time.January {
			half = "P1"
		}
		if month >= time.February && month <= time.August {
			half = "P2"
		}
		if month >= time.September && month <= time.December {
			half = "P1"
		}
	}

	results, err := client.GetResults(year, half)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}
