package web

import "time"

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

type CustomLinkAnalyticResponse struct {
	LinkID         uint                   `json:"link_id"`
	ClickCount     int                    `json:"click_count"`
	ViewCount      int                    `json:"view_count"`
	DeviceAnalytic DeviceAnalyticResponse `json:"device_analytic"`
	Datetime       string                 `json:"datetime"`
	LastUpdated    time.Time              `json:"last_updated"`
}

type TotalCustomLinkAnalyticResponse struct {
	LinkID          int `json:"link_id"`
	TotalClickCount int `json:"total_click_count"`
	TotalViewCount  int `json:"total_view_count"`
}

type CustomLinkAnalyticSummaryResponse struct {
	CustomLink     []TotalCustomLinkAnalyticResponse `json:"link"`
	DeviceAnalytic DeviceAnalyticResponse            `json:"device_analytic"`
	LastUpdated    time.Time                         `json:"last_updated"`
}
