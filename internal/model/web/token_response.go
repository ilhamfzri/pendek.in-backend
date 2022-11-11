package web

import "time"

type TokenResponse struct {
	AccessToken string    `json:"access_token,omitempty"`
	ValidUntil  time.Time `json:"valid_until,omitempty"`
}
