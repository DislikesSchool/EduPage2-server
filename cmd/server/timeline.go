package main

import (
	"encoding/json"
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
// @Success 200 {object} apimodel.Timeline
// @Failure 401 {object} apimodel.UnauthorizedResponse
// @Failure 500 {object} apimodel.InternalErrorResponse
// @Router /api/timeline/recent [get]
func RecentTimelineHandler(c *gin.Context) {
	client := c.MustGet("client").(*edupage.EdupageClient)

	var cacheKey string
	if shouldCache {
		cacheKey, err := CacheKeyFromEPClient(client, "timeline")
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
			var timeline apimodel.Timeline
			read, err := ReadCache(cacheKey, &timeline)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			if read {
				c.JSON(http.StatusOK, timeline)
				return
			}
		}
	}

	timeline, err := client.GetRecentTimeline()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, timeline)

	if shouldCache {
		_ = CacheData(cacheKey, timeline, TTLFromType("timeline"))
	}
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
// @Success 200 {object} apimodel.Timeline
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
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	timeline, err := client.GetTimeline(dateFrom, dateTo)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, timeline)
}

// SendMessageHandler godoc
// @Summary Send a message
// @Schemes
// @Description Sends a message to a recipient.
// @Tags messages
// @Param Authorization header string true "JWT token"
// @Param message body apimodel.SendMessageRequest true "Message"
// @Produce json
// @Security Bearer
// @Success 200
// @Failure 401 {object} apimodel.UnauthorizedResponse
// @Failure 500 {object} apimodel.InternalErrorResponse
// @Router /api/message [post]
func SendMessageHandler(c *gin.Context) {
	client := c.MustGet("client").(*edupage.EdupageClient)

	recipient := c.PostForm("recipient")
	optsJson := c.PostForm("message")

	var opts edupage.MessageOptions
	if err := json.Unmarshal([]byte(optsJson), &opts); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := client.SendMessage(recipient, opts); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// TimelineItemHandler godoc
// @Summary Get the timeline item by ID
// @Schemes
// @Description Returns the timeline item by ID.
// @Tags timeline
// @Param Authorization header string true "JWT token"
// @Param id path string true "Timeline item ID"
// @Produce json
// @Security Bearer
// @Success 200 {object} apimodel.TimelineItemWithOrigin
// @Failure 401 {object} apimodel.UnauthorizedResponse
// @Failure 500 {object} apimodel.InternalErrorResponse
// @Router /api/timelineitem/{id} [get]
func TimelineItemHandler(c *gin.Context) {
	client := c.MustGet("client").(*edupage.EdupageClient)

	id := c.Param("id")
	date := c.Query("date")

	var timeline model.Timeline
	var err error

	if date != "" {
		dateTime, err := time.Parse(time.RFC3339, date)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		timeline, err = client.GetTimeline(dateTime, time.Now())
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		timeline, err = client.GetRecentTimeline()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	timelineItem := timeline.Items[id]

	var replies []apimodel.TimelineItemWithOrigin
	for _, msg := range timeline.Items {
		if msg.ReactionTo == timelineItem.ID {
			replies = append(replies, apimodel.TimelineItemWithOrigin{
				ID:              msg.ID,
				Timestamp:       msg.Timestamp,
				ReactionTo:      msg.ReactionTo,
				Type:            msg.Type,
				User:            msg.User,
				TargetUser:      msg.TargetUser,
				UserName:        msg.UserName,
				OtherID:         msg.OtherID,
				Text:            msg.Text,
				TimeAdded:       msg.TimeAdded,
				TimeEvent:       msg.TimeEvent,
				Data:            msg.Data,
				Owner:           msg.Owner,
				OwnerName:       msg.OwnerName,
				ReactionCount:   msg.ReactionCount,
				LastReaction:    msg.LastReaction,
				PomocnyZaznam:   msg.PomocnyZaznam,
				Removed:         msg.Removed,
				TimeAddedBTC:    msg.TimeAddedBTC,
				LastReactionBTC: msg.LastReactionBTC,
				OriginServer:    client.Credentials.Server,
			})
		}

	}

	c.JSON(http.StatusOK, apimodel.TimelineItemWithOrigin{
		ID:              timelineItem.ID,
		Timestamp:       timelineItem.Timestamp,
		ReactionTo:      timelineItem.ReactionTo,
		Type:            timelineItem.Type,
		User:            timelineItem.User,
		TargetUser:      timelineItem.TargetUser,
		UserName:        timelineItem.UserName,
		OtherID:         timelineItem.OtherID,
		Text:            timelineItem.Text,
		TimeAdded:       timelineItem.TimeAdded,
		TimeEvent:       timelineItem.TimeEvent,
		Data:            timelineItem.Data,
		Owner:           timelineItem.Owner,
		OwnerName:       timelineItem.OwnerName,
		ReactionCount:   timelineItem.ReactionCount,
		LastReaction:    timelineItem.LastReaction,
		PomocnyZaznam:   timelineItem.PomocnyZaznam,
		Removed:         timelineItem.Removed,
		TimeAddedBTC:    timelineItem.TimeAddedBTC,
		LastReactionBTC: timelineItem.LastReactionBTC,
		OriginServer:    client.Credentials.Server,
		Replies:         replies,
	})
}
