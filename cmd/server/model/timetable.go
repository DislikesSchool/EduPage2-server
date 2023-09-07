package main

type TimetableRequest struct {
	From string `json:"from" example:"2022-01-01T00:00:00Z2022-01-01T00:00:00Z"`
	To   string `json:"to" example:"2022-01-01T00:00:00Z" default:"time.Now()"`
}

type TimetableUnauthorizedResponse struct {
	Error string `json:"error" example:"Unauthorized"`
}

type TimetableInternalErrorResponse struct {
	Error string `json:"error" example:"failed to create payload"`
}
