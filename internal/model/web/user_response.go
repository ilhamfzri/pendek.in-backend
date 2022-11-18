package web

type UserResponse struct {
	ID         string `json:"id,omitempty"`
	Username   string `json:"username,omitempty"`
	FullName   string `json:"full_name,omitempty"`
	Bio        string `json:"bio,omitempty"`
	Email      string `json:"email,omitempty"`
	ProfilePic string `json:"profile_pic"`
}
