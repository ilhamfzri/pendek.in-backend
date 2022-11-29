package web

import "time"

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

type SocialMediaAnalyticInteractionRequest struct {
	ClientIP          string
	UserAgent         string
	SocialMediaLinkID uint
}

type SocialMediaAnalyticGetRequest struct {
	TypeID    int       `form:"type_id" binding:"required"`
	StartDate time.Time `form:"start_date" binding:"required" time_format:"2006-01-02"`
	EndDate   time.Time `form:"end_date" binding:"required" time_format:"2006-01-02"`
}
