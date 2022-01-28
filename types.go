package zonnepanelendelen

import (
	"net/http"
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

// Project is a single zonnepanelendelen project
type Project struct {
	Name      string  `json:"name"`
	ID        int     `json:"id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	IsHidden  bool    `json:"is_hidden"`
}

// HTTPClient is the interface that should be implemented by API HTTP clients
// while this usually is just the default HTTP client, but it allows to override
// the HTTP client for additional control over API calls or testing purposes
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}
