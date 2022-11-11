package web

type UserResponse struct {
	Id        int    `json:"id,omitempty"`
	Username  string `json:"username,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Bio       string `json:"bio,omitempty"`
	Email     string `json:"email,omitempty"`
}
