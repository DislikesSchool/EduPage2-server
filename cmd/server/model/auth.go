package main

type LoginRequestUsernamePassword struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginRequestToken struct {
	Token string `json:"token" binding:"required"`
}

type LoginRequestUsernamePasswordServer struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Server   string `json:"server" binding:"required"`
}

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
