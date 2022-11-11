package web

import "time"

type UserResponse struct {
	Id          int       `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	ActivatedAt time.Time `json:"activated_at"`
}
