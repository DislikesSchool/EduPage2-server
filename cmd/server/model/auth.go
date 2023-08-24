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
	Token     string `json:"token"`
	TempToken bool   `json:"temp_token"`
	Success   bool   `json:"success"`
	Error     string `json:"error"`
}

type LoginBadRequestResponse struct {
	Error string `json:"error" example:"Username and Password are required"`
}

type LoginUnauthorizedResponse struct {
	Error string `json:"error" example:"Invalid username or password"`
}
