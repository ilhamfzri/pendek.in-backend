package web

type ThumbnailResponse struct {
	ID           uint   `json:"id"`
	Name         string `json:"name,omitempty"`
	ThumbnailUrl string `json:"thumbnail_url"`
}

type CustomLinkResponse struct {
	ID                uint   `json:"id"`
	Title             string `json:"title"`
	ShortLinkCode     string `json:"short_link_code"`
	LongLink          string `json:"long_link"`
	RedirectLink      string `json:"redirect_link"`
	ShowOnProfile     bool   `json:"show_on_profile"`
	Activate          bool   `json:"activate"`
	ThumbnailID       uint   `json:"thumbnail_id,omitempty"`
	CustomThumbnailID uint   `json:"custom_thumbnail_id,omitempty"`
	ThumbnailUrl      string `json:"thumbnail_url,omitempty"`
}
