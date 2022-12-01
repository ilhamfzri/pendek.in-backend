package web

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
	LinkID string `uri:"link_id" binding:"required"`
}

type CustomLinkCheckShortCodeRequest struct {
	LinkID string `uri:"code" binding:"required"`
}
