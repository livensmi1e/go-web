package models

type LoginRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
	Type  string `json:"type"`
}

type LoginResponseBody struct {
	Data       *LoginResponse `json:"data"`
	StatusCode int            `json:"status_code"`
}

type RegisterRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	Id    string `json:"id"`
	Email string `json:"email"`
}

type RegisterResponseBody struct {
	Data       *RegisterResponse `json:"data"`
	StatusCode int               `json:"status_code"`
}

type GetMeResponseBody struct {
	Data       string `json:"data"`
	StatusCode int    `json:"status_code"`
}
