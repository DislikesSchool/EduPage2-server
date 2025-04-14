package main

import (
	"net/http"
	"time"

	"github.com/DislikesSchool/EduPage2-server/cmd/server/dbmodel"
	"github.com/DislikesSchool/EduPage2-server/cmd/server/util"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetSecurityStatusHandler godoc
// @Summary Get user's data storage status
// @Schemes
// @Description Returns the current data storage preferences for the user
// @Tags security
// @Security Bearer
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /security/status [get]
func GetSecurityStatusHandler(c *gin.Context) {
	claims, err := getClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	username := claims["username"].(string)
	server := claims["server"].(string)

	clientData, ok := util.Clients[server+username]
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "client not found"})
		return
	}

	dataStorage := clientData.DataStorage

	// Initialize response with current in-memory settings
	response := gin.H{
		"dataStorage": gin.H{
			"enabled":     dataStorage.Enabled,
			"credentials": dataStorage.Credentials,
			"messages":    dataStorage.Messages,
			"timeline":    dataStorage.Timeline,
		},
	}

	// If database storage is enabled, add user record info
	if util.ShouldStore && dataStorage.Enabled {
		userModel := &dbmodel.User{}
		result := util.Db.First(userModel, "username = ?", username)
		if result.Error == nil {
			response["database"] = gin.H{
				"lastOnline":    userModel.LastOnline,
				"storeMessages": userModel.StoreMessages,
				"storeTimeline": userModel.StoreTimeline,
			}
		}
	}

	c.JSON(http.StatusOK, response)
}

// UpdateDataStoragePrefsHandler godoc
// @Summary Update user's data storage preferences
// @Schemes
// @Description Updates the user's data storage preferences and removes any data that is no longer permitted
// @Tags security
// @Security Bearer
// @Accept json
// @Param preferences body util.DataStorageConfig true "Data storage preferences"
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /security/preferences [post]
func UpdateDataStoragePrefsHandler(c *gin.Context) {
	var prefs util.DataStorageConfig
	if err := c.ShouldBindJSON(&prefs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	claims, err := getClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	username := claims["username"].(string)
	server := claims["server"].(string)

	clientData, ok := util.Clients[server+username]
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "client not found"})
		return
	}

	// Update the in-memory data storage preferences
	clientData.DataStorage = &prefs

	changes := gin.H{
		"dataStorageUpdated": true,
	}

	// If database storage is enabled, update the database
	if util.ShouldStore {
		userModel := &dbmodel.User{}
		result := util.Db.First(userModel, "username = ?", username)

		if result.Error == nil {
			// Clean up data if permissions have been revoked
			dataCleanup := gin.H{}

			// Handle credential storage change - delete entire record if credentials storage is disabled
			if !prefs.Credentials || !prefs.Enabled {
				// Delete the entire user record
				if err := util.Db.Delete(userModel).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user record: " + err.Error()})
					return
				}
				dataCleanup["userRecordDeleted"] = true
				changes["dataCleanup"] = dataCleanup
			} else {
				// User wants to keep credentials but maybe change other preferences

				// Handle messages storage change
				if userModel.StoreMessages && !prefs.Messages {
					// TODO: Clean up message data when we implement message storage
					dataCleanup["messagesRemoved"] = true
				}

				// Handle timeline storage change
				if userModel.StoreTimeline && !prefs.Timeline {
					// TODO: Clean up timeline data when we implement timeline storage
					dataCleanup["timelineRemoved"] = true
				}

				// Update user model with new preferences
				userModel.LastOnline = time.Now()
				userModel.StoreMessages = prefs.Messages
				userModel.StoreTimeline = prefs.Timeline

				if err := util.Db.Save(userModel).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user preferences: " + err.Error()})
					return
				}

				changes["dataCleanup"] = dataCleanup
			}
		} else if !prefs.Enabled || !prefs.Credentials {
			// User does not exist in DB and doesn't want to be stored - nothing to do
			changes["databaseAction"] = "none"
		} else if result.Error == gorm.ErrRecordNotFound {
			userModel := &dbmodel.User{
				Username:      username,
				LastOnline:    time.Now(),
				StoreMessages: prefs.Messages,
				StoreTimeline: prefs.Timeline,
			}

			if err := util.Db.Create(userModel).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user: " + err.Error()})
				return
			}

			changes["databaseAction"] = "created"
		} else {
			// Would have created the user but encountered an error
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}
	}

	// Clean up any cached data if needed
	if util.ShouldCache {
		client := clientData.Client
		// If timeline permission revoked, clear timeline cache
		if !prefs.Timeline {
			cacheKey, err := util.CacheKeyFromEPClient(client, "timeline")
			if err == nil {
				util.Rdb.Del(util.Ctx, cacheKey)
				changes["timelineCacheCleared"] = true
			}
		}

		// Clean any other cache types here when needed
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"changes": changes,
	})
}

// DeleteUserDataHandler godoc
// @Summary Delete all user data
// @Schemes
// @Description Deletes all stored data for the current user
// @Tags security
// @Security Bearer
// @Param confirm query bool true "Confirmation flag set to true"
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /security/delete-data [delete]
func DeleteUserDataHandler(c *gin.Context) {
	// Require confirmation to prevent accidental deletion
	confirm := c.Query("confirm")
	if confirm != "true" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Confirmation required. Set 'confirm=true' query parameter"})
		return
	}

	claims, err := getClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	username := claims["username"].(string)
	server := claims["server"].(string)

	clientData, ok := util.Clients[server+username]
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "client not found"})
		return
	}

	result := gin.H{
		"databaseDeletion": "not applicable",
		"cacheDeletion":    "not applicable",
	}

	// Remove data from database if storage is enabled
	if util.ShouldStore {
		dbResult := util.Db.Where("username = ?", username).Delete(&dbmodel.User{})
		if dbResult.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user from database: " + dbResult.Error.Error()})
			return
		}
		result["databaseDeletion"] = "success"
		result["recordsDeleted"] = dbResult.RowsAffected
	}

	// Clear any cached data
	if util.ShouldCache {
		client := clientData.Client
		user, err := client.GetUser(false)
		if err == nil {
			schoolId := server
			userId := user.UserRow.UserID

			// Use pattern matching to delete all keys related to this user
			pattern := schoolId + ":" + userId + ":*"
			keys, err := util.Rdb.Keys(util.Ctx, pattern).Result()
			if err == nil {
				if len(keys) > 0 {
					util.Rdb.Del(util.Ctx, keys...)
				}
				result["cacheDeletion"] = "success"
				result["keysDeleted"] = len(keys)
			} else {
				result["cacheDeletion"] = "error: " + err.Error()
			}
		} else {
			result["cacheDeletion"] = "error: " + err.Error()
		}
	}

	// Update in-memory preferences to disable storage
	clientData.DataStorage = &util.DataStorageConfig{
		Enabled:     false,
		Credentials: false,
		Messages:    false,
		Timeline:    false,
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"result":  result,
	})
}

// RequestDataExportHandler godoc
// @Summary Request export of user data
// @Schemes
// @Description Returns all data stored for the current user in a JSON format
// @Tags security
// @Security Bearer
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /security/export-data [get]
func RequestDataExportHandler(c *gin.Context) {
	claims, err := getClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	username := claims["username"].(string)
	server := claims["server"].(string)

	_, ok := util.Clients[server+username]
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "client not found"})
		return
	}

	export := gin.H{
		"metadata": gin.H{
			"exportDate": time.Now(),
			"username":   username,
			"server":     server,
		},
	}

	if util.ShouldStore {
		userModel := &dbmodel.User{}
		result := util.Db.First(userModel, "username = ?", username)
		if result.Error == nil {
			// Exclude sensitive fields like Password
			export["databaseRecord"] = gin.H{
				"id":            userModel.ID,
				"username":      userModel.Username,
				"server":        userModel.Server,
				"lastOnline":    userModel.LastOnline,
				"storeMessages": userModel.StoreMessages,
				"storeTimeline": userModel.StoreTimeline,
			}

			// TODO: When implemented, include other stored data like messages, timeline events, etc.
		}
	}

	c.JSON(http.StatusOK, export)
}
