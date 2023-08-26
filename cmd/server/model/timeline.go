package main

import "github.com/DislikesSchool/EduPage2-server/edupage/model"

type RecentTimelineSuccessResponse struct {
	Success  bool             `json:"success" example:"true"`
	Error    string           `json:"error" example:""`
	Timeline []model.Timeline `json:"timeline"`
}

type RecentTimelineUnauthorizedResponse struct {
	Success bool   `json:"success" example:"false"`
	Error   string `json:"error" example:"Unauthorized"`
}

type RecentTimelineInternalErrorResponse struct {
	Success bool   `json:"success" example:"false"`
	Error   string `json:"error" example:"failed to create payload"`
}
