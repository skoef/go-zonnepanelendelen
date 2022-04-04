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
	Metrics   Metrics `json:"metrics"`
}

// Metrics is a container for all sorts of project metrics
type Metrics struct {
	Parts     int         `json:"solar_parts"`
	Interest  float64     `json:"solar_interest"`
	Value     int         `json:"net_value"`
	Today     TimeMetrics `json:"production_today"`
	LastMonth TimeMetrics `json:"production_last_month"`
	LastYear  TimeMetrics `json:"production_last_year"`
	All       TimeMetrics `json:"production_all"`
}

// TimeMetrics combine certain metrics over a specific period in time
type TimeMetrics struct {
	TotalPower    float64       `json:"total_power_kWh"`
	TotalExpected float64       `json:"total_power_expected_kWh"`
	TotalProfit   float64       `json:"total_profit"`
	ROI           float64       `json:"return_on_investment"`
	Measurements  []Measurement `json:"data"`
}

// Measurement holds measurements for a specific moment in time
type Measurement struct {
	Production         float64 `json:"production_kWh"`
	Expected           float64 `json:"expected_production_kWh"`
	Cumulative         float64 `json:"cumulative_production_kWh"`
	ExpectedCumulative float64 `json:"expected_cumulative_production_kWh"`
	Timestamp          string  `json:"timestamp"`
}

// HTTPClient is the interface that should be implemented by API HTTP clients
// while this usually is just the default HTTP client, but it allows to override
// the HTTP client for additional control over API calls or testing purposes
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}
