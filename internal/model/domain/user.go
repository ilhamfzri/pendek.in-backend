package domain

import "time"

type User struct {
	Id        int
	Username  string
	FirstName string
	LastName  string
	Bio       string
	Email     string
	Password  string
	Verified  bool
	LastLogin time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}
