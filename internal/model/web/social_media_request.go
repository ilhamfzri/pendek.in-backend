package web

type SocialMediaLinkCreateRequest struct {
	TypeID         int    `validate:"required" json:"type_id"`
	LinkOrUsername string `validate:"required" json:"link_or_username"`
}

type SocialMediaLinkUpdateRequest struct {
	TypeID            int    `uri:"type_id" binding:"required" `
	NewLinkOrUsername string `json:"new_link_or_username"`
	Activate          *bool  `json:"activate"`
}

type SocialMediaLinkRedirectRequest struct {
	Username        string `uri:"username" binding:"required"`
	SocialMediaName string `uri:"social-media" binding:"required"`
}
