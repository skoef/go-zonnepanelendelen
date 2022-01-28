package zonnepanelendelen

import "errors"

var (
	ErrNotAuthenticated = errors.New("failed to authenticate")
)

// API is a container for holding authentication state for API interfacing
type API struct {
	token *AuthToken

	username string
	password string
}

// AuthToken is the data structure as returned by /obtain-auth-token
type AuthToken struct {
	Name  string `json:"name"`
	Token string `json:"token"`
}

type Project struct {
	Name      string  `json:"name"`
	ID        int     `json:"id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	IsHidden  bool    `json:"is_hidden"`
}
