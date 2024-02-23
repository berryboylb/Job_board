package auth

import "time"

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

type GoogleResponse struct {
	Aud           string    `json:"aud"`
	Email         string    `json:"email"`
	EmailVerified bool      `json:"email_verified"`
	Exp           int       `json:"exp"`
	FamilyName    string    `json:"family_name"`
	GivenName     string    `json:"given_name"`
	Iat           int       `json:"iat"`
	Iss           string    `json:"iss"`
	Locale        string    `json:"locale"`
	Name          string    `json:"name"`
	Nickname      string    `json:"nickname"`
	Picture       string    `json:"picture"`
	Sid           string    `json:"sid"`
	Sub           string    `json:"sub"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type EmailResponse struct {
	Aud           string    `json:"aud"`
	Email         string    `json:"email"`
	EmailVerified bool      `json:"email_verified"`
	Exp           int       `json:"exp"`
	Iat           int       `json:"iat"`
	Iss           string    `json:"iss"`
	Name          string    `json:"name"`
	Nickname      string    `json:"nickname"`
	Picture       string    `json:"picture"`
	Sid           string    `json:"sid"`
	Sub           string    `json:"sub"`
	UpdatedAt     time.Time `json:"updated_at"`
}
