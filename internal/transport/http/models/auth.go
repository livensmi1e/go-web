package models

type LoginRequestBody struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=3"`
}

type LoginResponse struct {
	Token string `json:"token"`
	Type  string `json:"type"`
}

type LoginResponseBody struct {
	Data       *LoginResponse `json:"data"`
	StatusCode int            `json:"statusCode"`
}

type RegisterRequestBody struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=3"`
}

type RegisterResponse struct {
	Id    string `json:"id"`
	Email string `json:"email"`
}

type RegisterResponseBody struct {
	Data       *RegisterResponse `json:"data"`
	StatusCode int               `json:"statusCode"`
}

type GetMeResponseBody struct {
	Data       string `json:"data"`
	StatusCode int    `json:"statusCode"`
}
