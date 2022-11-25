package web

type SocialMediaTypeResponse struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Example string `json:"example"`
	IconUrl string `json:"icon_url,omitempty"`
}

type SocialMediaLinkResponse struct {
	TypeID          uint   `json:"type_id"`
	SocialMediaName string `json:"social_media"`
	LinkOrUsername  string `json:"link_or_username"`
	Activate        bool   `json:"activate"`
	RedirectLink    string `json:"redirect_link,omitempty"`
}
