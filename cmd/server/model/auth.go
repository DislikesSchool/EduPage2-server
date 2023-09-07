package main

type LoginSuccessResponse struct {
	Success   bool   `json:"success" example:"true"`
	Error     string `json:"error" example:""`
	Token     string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM"`
	Firstname string `json:"firstname" example:"John"`
	Lastname  string `json:"lastname" example:"Doe"`
}

type LoginBadRequestResponse struct {
	Success bool   `json:"success" example:"false"`
	Error   string `json:"error" example:"Username and Password are required"`
}

type LoginUnauthorizedResponse struct {
	Success bool   `json:"success" example:"false"`
	Error   string `json:"error" example:"Unexpected response from server, make sure credentials are specified correctly"`
}

type LoginInternalErrorResponse struct {
	Success bool   `json:"success" example:"false"`
	Error   string `json:"error" example:"failed to login: Post https://example.edupage.org/login/edubarLogin.php: dial tcp: lookup example.edupage.org: no such host"`
}

type ValidateTokenSuccessResponse struct {
	Success bool   `json:"success" example:"true"`
	Error   string `json:"error" example:""`
	Expires string `json:"expires" example:"1620000000"`
}

type ValidateTokenUnauthorizedResponse struct {
	Success bool   `json:"success" example:"false"`
	Error   string `json:"error" example:"Unauthorized"`
}

type UnauthorizedResponse struct {
	Error string `json:"error" example:"Unauthorized"`
}

type InternalErrorResponse struct {
	Error string `json:"error" example:"failed to create payload"`
}
