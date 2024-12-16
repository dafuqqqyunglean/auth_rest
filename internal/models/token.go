package models

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
	User         User   `json:"user"`
}
