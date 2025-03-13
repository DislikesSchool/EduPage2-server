package main

import (
	"net/http"
	"time"

	"github.com/DislikesSchool/EduPage2-server/cmd/server/apimodel"
	"github.com/DislikesSchool/EduPage2-server/edupage"
	"github.com/DislikesSchool/EduPage2-server/edupage/model"
	"github.com/gin-gonic/gin"
)

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

	var cacheKey string
	var err error
	if shouldCache {
		cacheKey, err = CacheKeyFromEPClient(client, "timetable")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		cached, err := IsCached(cacheKey)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if cached {
			var timetable apimodel.CompleteTimetable
			read, err := ReadCache(cacheKey, &timetable)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			if read {
				c.JSON(http.StatusOK, timetable)

				go func() {
					timetable, err := client.GetRecentTimetable()
					if err != nil {
						return
					}

					user, err := client.GetUser(false)
					if err != nil {
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

					_ = CacheData(cacheKey, completeTimetable, TTLFromType("timetable"))
				}()

				return
			}
		}
	}

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

	if shouldCache {
		_ = CacheData(cacheKey, completeTimetable, TTLFromType("timetable"))
	}
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
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

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
