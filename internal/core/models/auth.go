package models

type AuthTokens struct {
	AccessToken  string
	RefreshToken string
}

type RefreshUser struct {
	Id    string
	Email string
	// Add other fields as necessary like Roles, Permissions, etc.
}
