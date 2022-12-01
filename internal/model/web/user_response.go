package web

type UserResponse struct {
	ID         string `json:"id,omitempty"`
	Username   string `json:"username,omitempty"`
	FullName   string `json:"full_name,omitempty"`
	Bio        string `json:"bio,omitempty"`
	Email      string `json:"email,omitempty"`
	ProfilePic string `json:"profile_pic"`
}

type UserProfileResponse struct {
	Username    string                           `json:"username"`
	FullName    string                           `json:"full_name,omitempty"`
	Bio         string                           `json:"bio,omitempty"`
	ProfilePic  string                           `json:"profile_pic,omitempty"`
	SocialMedia []UserProfileSocialMediaResponse `json:"social_media"`
	Link        []UserProfileCustomLinkResponse  `json:"link"`
}

type UserProfileSocialMediaResponse struct {
	Name    string `json:"name"`
	Link    string `json:"link"`
	IconUrl string `json:"icon_url"`
}

type UserProfileCustomLinkResponse struct {
	Title        string `json:"title"`
	Link         string `json:"link"`
	ThumbnailUrl string `json:"thumbnail_url"`
}
