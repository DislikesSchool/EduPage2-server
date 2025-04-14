package routes

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/DislikesSchool/EduPage2-server/cmd/server/util"
	"github.com/DislikesSchool/EduPage2-server/edupage"
	"github.com/gin-gonic/gin"
	"github.com/meilisearch/meilisearch-go"
)

// SearchMessagesRequest represents the request parameters for message search
type SearchMessagesRequest struct {
	Query    string `form:"query" binding:"required"`
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"pageSize,default=20"`
	SenderID string `form:"senderId"`
}

// SearchMessagesResponse represents the response structure for message searches
type SearchMessagesResponse struct {
	TotalHits      int64                    `json:"totalHits"`
	Page           int                      `json:"page"`
	PageSize       int                      `json:"pageSize"`
	TotalPages     int                      `json:"totalPages"`
	ProcessingTime string                   `json:"processingTime"`
	Messages       []map[string]interface{} `json:"messages"`
}

// SearchMessagesHandler godoc
// @Summary Search through messages
// @Schemes
// @Description Searches through user's messages using Meilisearch
// @Tags messages,search
// @Param Authorization header string true "JWT token"
// @Param query query string true "Search query"
// @Param page query int false "Page number (default: 1)"
// @Param pageSize query int false "Page size (default: 20, max: 100)"
// @Param senderId query string false "Filter by sender ID"
// @Produce json
// @Security Bearer
// @Success 200 {object} SearchMessagesResponse
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Search functionality disabled"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/search/messages [get]
func SearchMessagesHandler(c *gin.Context) {
	// Check if search is enabled
	if !util.ShouldSearch {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Message search functionality is not enabled on this server",
		})
		return
	}

	// Get client
	client := c.MustGet("client").(*edupage.EdupageClient)

	// Create owner ID that will be used for filtering (username@server format)
	ownerID := fmt.Sprintf("%s@%s", client.Credentials.Username, client.Credentials.Server)

	// Parse request
	var req SearchMessagesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate and limit page size
	if req.PageSize <= 0 {
		req.PageSize = 20
	} else if req.PageSize > 100 {
		req.PageSize = 100 // Cap at 100 for safety
	}

	// Validate page number
	if req.Page <= 0 {
		req.Page = 1
	}

	// Prepare search request
	searchReq := &meilisearch.SearchRequest{
		// Always filter by owner ID to ensure privacy
		Filter: fmt.Sprintf("ownerid = \"%s\"", ownerID),

		// Sort by timestamp descending (newest first)
		Sort: []string{"timestamp:desc"},

		// Pagination parameters
		Limit:  int64(req.PageSize),
		Offset: int64((req.Page - 1) * req.PageSize),
	}

	// Add sender filter if provided
	if req.SenderID != "" {
		searchReq.Filter = fmt.Sprintf("%s AND sender_id = \"%s\"", searchReq.Filter, req.SenderID)
	}

	// Perform search
	searchResult, err := util.MeiliIndex.Search(req.Query, searchReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Search failed: " + err.Error(),
		})
		return
	}

	// Calculate total pages
	totalPages := int(searchResult.EstimatedTotalHits) / req.PageSize
	if int(searchResult.EstimatedTotalHits)%req.PageSize > 0 {
		totalPages++
	}

	// Format processing time
	processingTime := fmt.Sprintf("%d ms", searchResult.ProcessingTimeMs)

	// Extract messages while ensuring no sensitive data is included
	messages := make([]map[string]interface{}, 0, len(searchResult.Hits))
	for _, hit := range searchResult.Hits {
		if doc, ok := hit.(map[string]interface{}); ok {
			// Keep only necessary fields, removing anything that might contain
			// sensitive data or internal flags like the ownerID itself
			safeDoc := map[string]interface{}{
				"id":        doc["id"],
				"text":      doc["text"],
				"timestamp": doc["timestamp"],
				"sender_id": doc["sender_id"],
			}

			// Include receiver_id if present
			if receiverId, exists := doc["receiver_id"]; exists {
				safeDoc["receiver_id"] = receiverId
			}

			messages = append(messages, safeDoc)
		}
	}

	// Return response
	c.JSON(http.StatusOK, SearchMessagesResponse{
		TotalHits:      searchResult.EstimatedTotalHits,
		Page:           req.Page,
		PageSize:       req.PageSize,
		TotalPages:     totalPages,
		ProcessingTime: processingTime,
		Messages:       messages,
	})
}

// ConversationSearchHandler godoc
// @Summary Search messages between two users
// @Schemes
// @Description Search through messages between the current user and a specific other user
// @Tags messages,search
// @Param Authorization header string true "JWT token"
// @Param userId path string true "The ID of the other user in the conversation"
// @Param query query string true "Search query"
// @Param page query int false "Page number (default: 1)"
// @Param pageSize query int false "Page size (default: 20, max: 100)"
// @Produce json
// @Security Bearer
// @Success 200 {object} SearchMessagesResponse
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Search functionality disabled"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/search/conversation/{userId} [get]
func ConversationSearchHandler(c *gin.Context) {
	// Check if search is enabled
	if !util.ShouldSearch {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Message search functionality is not enabled on this server",
		})
		return
	}

	// Get client
	client := c.MustGet("client").(*edupage.EdupageClient)

	// Get user info
	user, err := client.GetUser(false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get user information: " + err.Error(),
		})
		return
	}

	// Create owner ID for filtering
	ownerID := fmt.Sprintf("%s@%s", client.Credentials.Username, client.Credentials.Server)

	// Get conversation partner ID
	otherUserId := c.Param("userId")
	if otherUserId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	// Parse query parameters
	query := c.Query("query")
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	if err != nil || pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100 // Max page size
	}

	// Get current user ID
	currentUserId := user.UserRow.UserID

	// Prepare search request with complex filter:
	// 1. Messages must belong to the current user (ownerID)
	// 2. Messages must be either:
	//    a. Sent by currentUser to otherUser OR
	//    b. Sent by otherUser to currentUser
	filter := fmt.Sprintf(
		"ownerid = \"%s\" AND ((sender_id = \"%s\" AND receiver_id = \"%s\") OR (sender_id = \"%s\" AND receiver_id = \"%s\"))",
		ownerID, currentUserId, otherUserId, otherUserId, currentUserId,
	)

	searchReq := &meilisearch.SearchRequest{
		Filter: filter,
		Sort:   []string{"timestamp:desc"},
		Limit:  int64(pageSize),
		Offset: int64((page - 1) * pageSize),
	}

	// Perform search
	searchResult, err := util.MeiliIndex.Search(query, searchReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Search failed: " + err.Error(),
		})
		return
	}

	// Calculate total pages
	totalPages := int(searchResult.EstimatedTotalHits) / pageSize
	if int(searchResult.EstimatedTotalHits)%pageSize > 0 {
		totalPages++
	}

	// Format processing time
	processingTime := fmt.Sprintf("%d ms", searchResult.ProcessingTimeMs)

	// Extract messages
	messages := make([]map[string]interface{}, 0, len(searchResult.Hits))
	for _, hit := range searchResult.Hits {
		if doc, ok := hit.(map[string]interface{}); ok {
			// Create safe document with only necessary fields
			safeDoc := map[string]interface{}{
				"id":          doc["id"],
				"text":        doc["text"],
				"timestamp":   doc["timestamp"],
				"sender_id":   doc["sender_id"],
				"receiver_id": doc["receiver_id"],
				// Add a direction flag to easily identify message direction in UI
				"outgoing": doc["sender_id"] == currentUserId,
			}

			messages = append(messages, safeDoc)
		}
	}

	// Return response
	c.JSON(http.StatusOK, SearchMessagesResponse{
		TotalHits:      searchResult.EstimatedTotalHits,
		Page:           page,
		PageSize:       pageSize,
		TotalPages:     totalPages,
		ProcessingTime: processingTime,
		Messages:       messages,
	})
}

// MessageFulltextSearchHandler godoc
// @Summary Advanced message search with multiple filters
// @Schemes
// @Description Search through messages with advanced filtering options
// @Tags messages,search
// @Param Authorization header string true "JWT token"
// @Param query query string false "Search query text"
// @Param senderIds query string false "Comma-separated list of sender IDs to include"
// @Param receiverIds query string false "Comma-separated list of receiver IDs to include"
// @Param startDate query string false "Search messages after this date (RFC3339 format)"
// @Param endDate query string false "Search messages before this date (RFC3339 format)"
// @Param page query int false "Page number (default: 1)"
// @Param pageSize query int false "Page size (default: 20, max: 100)"
// @Param sortDir query string false "Sort direction: asc or desc (default: desc)"
// @Produce json
// @Security Bearer
// @Success 200 {object} SearchMessagesResponse
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Search functionality disabled"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/search/advanced [get]
func MessageFulltextSearchHandler(c *gin.Context) {
	// Check if search is enabled
	if !util.ShouldSearch {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Message search functionality is not enabled on this server",
		})
		return
	}

	// Get client
	client := c.MustGet("client").(*edupage.EdupageClient)

	// Create owner ID for filtering
	ownerID := fmt.Sprintf("%s@%s", client.Credentials.Username, client.Credentials.Server)

	// Extract all query parameters
	query := c.Query("query")
	senderIdsRaw := c.Query("senderIds")
	receiverIdsRaw := c.Query("receiverIds")
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")
	sortDir := strings.ToLower(c.DefaultQuery("sortDir", "desc"))

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	if err != nil || pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100 // Max page size
	}

	// Validate sort direction
	if sortDir != "asc" && sortDir != "desc" {
		sortDir = "desc" // Default to newest first
	}

	// Start building the filter with owner ID (always required)
	filterParts := []string{fmt.Sprintf("ownerid = \"%s\"", ownerID)}

	// Add sender IDs filter if provided
	if senderIdsRaw != "" {
		senderIds := strings.Split(senderIdsRaw, ",")
		if len(senderIds) > 0 {
			senderFilter := []string{}
			for _, id := range senderIds {
				if id = strings.TrimSpace(id); id != "" {
					senderFilter = append(senderFilter, fmt.Sprintf("sender_id = \"%s\"", id))
				}
			}
			if len(senderFilter) > 0 {
				filterParts = append(filterParts, "("+strings.Join(senderFilter, " OR ")+")")
			}
		}
	}

	// Add receiver IDs filter if provided
	if receiverIdsRaw != "" {
		receiverIds := strings.Split(receiverIdsRaw, ",")
		if len(receiverIds) > 0 {
			receiverFilter := []string{}
			for _, id := range receiverIds {
				if id = strings.TrimSpace(id); id != "" {
					receiverFilter = append(receiverFilter, fmt.Sprintf("receiver_id = \"%s\"", id))
				}
			}
			if len(receiverFilter) > 0 {
				filterParts = append(filterParts, "("+strings.Join(receiverFilter, " OR ")+")")
			}
		}
	}

	// Add date range filters if provided
	// Note: timestamp is stored as Unix timestamp (number) in Meilisearch
	if startDate != "" {
		filterParts = append(filterParts, fmt.Sprintf("timestamp >= %s", startDate))
	}

	if endDate != "" {
		filterParts = append(filterParts, fmt.Sprintf("timestamp <= %s", endDate))
	}

	// Combine all filter parts with AND
	filter := strings.Join(filterParts, " AND ")

	// Prepare search request
	searchReq := &meilisearch.SearchRequest{
		Filter: filter,
		Sort:   []string{fmt.Sprintf("timestamp:%s", sortDir)},
		Limit:  int64(pageSize),
		Offset: int64((page - 1) * pageSize),
	}

	// Perform search
	searchResult, err := util.MeiliIndex.Search(query, searchReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Search failed: " + err.Error(),
		})
		return
	}

	// Calculate total pages
	totalPages := int(searchResult.EstimatedTotalHits) / pageSize
	if int(searchResult.EstimatedTotalHits)%pageSize > 0 {
		totalPages++
	}

	// Format processing time
	processingTime := fmt.Sprintf("%d ms", searchResult.ProcessingTimeMs)

	// Extract messages
	messages := make([]map[string]interface{}, 0, len(searchResult.Hits))
	for _, hit := range searchResult.Hits {
		if doc, ok := hit.(map[string]interface{}); ok {
			// Create safe document with only necessary fields
			safeDoc := map[string]interface{}{
				"id":          doc["id"],
				"text":        doc["text"],
				"timestamp":   doc["timestamp"],
				"sender_id":   doc["sender_id"],
				"receiver_id": doc["receiver_id"],
			}

			messages = append(messages, safeDoc)
		}
	}

	// Return response
	c.JSON(http.StatusOK, SearchMessagesResponse{
		TotalHits:      searchResult.EstimatedTotalHits,
		Page:           page,
		PageSize:       pageSize,
		TotalPages:     totalPages,
		ProcessingTime: processingTime,
		Messages:       messages,
	})
}
