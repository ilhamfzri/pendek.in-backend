package web

import "time"

type CustomLinkCreateRequest struct {
	Title           string `json:"title" binding:"required,min=1,max=20"`
	ShortLinkCode   string `json:"short_link_code" binding:"required,min=5,max=20,alphanum"`
	LongLink        string `json:"long_link" binding:"required,url"`
	UserThumbnailID *uint  `json:"user_thumbnail_id"`
	ThumbnailID     *uint  `json:"thumbnail_id"`
}

type CustomLinkUpdateRequest struct {
	CustomLinkID    uint   `uri:"link_id" binding:"required"`
	Title           string `json:"title" binding:"omitempty,min=1,max=20"`
	ShortLinkCode   string `json:"short_link_code" binding:"omitempty,min=5,max=20,alphanum"`
	LongLink        string `json:"long_link" binding:"omitempty,url"`
	UserThumbnailID *uint  `json:"user_thumbnail_id"`
	ThumbnailID     *uint  `json:"thumbnail_id"`
	ShowOnProfile   *bool  `json:"show_on_profile"`
	Activate        *bool  `json:"activate"`
}

type CustomLinkGetRequest struct {
	LinkID uint `uri:"link_id" binding:"required"`
}

type CustomLinkRedirectRequest struct {
	ShortLinkCode string `uri:"short_link_code" binding:"required"`
}

type CustomLinkCheckShortCodeAvaibilityRequest struct {
	Code string `form:"code" binding:"required,min=5,max=20"`
}

type CustomLinkAnalyticInteractionRequest struct {
	ClientIP     string
	UserAgent    string
	CustomLinkID uint
}

type CustomLinkAnalyticGetRequest struct {
	LinkID    int       `form:"link_id" binding:"required"`
	StartDate time.Time `form:"start_date" binding:"required" time_format:"2006-01-02"`
	EndDate   time.Time `form:"end_date" binding:"required" time_format:"2006-01-02"`
}
