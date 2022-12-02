package web

import "time"

type SocialMediaTypeResponse struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Example string `json:"example"`
	IconUrl string `json:"icon_url"`
}

type SocialMediaLinkResponse struct {
	TypeID          uint   `json:"type_id"`
	SocialMediaName string `json:"social_media"`
	LinkOrUsername  string `json:"link_or_username"`
	Activate        bool   `json:"activate"`
	RedirectLink    string `json:"redirect_link,omitempty"`
}

type SocialMediaAnalyticResponse struct {
	SocialMediaLinkID uint                   `json:"social_media_link_id"`
	SocialMediaName   string                 `json:"social_media_name"`
	ClickCount        int                    `json:"click_count"`
	ViewCount         int                    `json:"view_count"`
	DeviceAnalytic    DeviceAnalyticResponse `json:"device_analytic"`
	Datetime          string                 `json:"datetime"`
	LastUpdated       time.Time              `json:"last_updated"`
}

type TotalSocialMediaAnalyticResponse struct {
	SocialMediaName string `json:"name"`
	TotalClickCount int    `json:"total_click_count"`
	TotalViewCount  int    `json:"total_view_count"`
}

type SocialMediaAnalyticSummaryResponse struct {
	SocialMedia    []TotalSocialMediaAnalyticResponse `json:"social_media"`
	DeviceAnalytic DeviceAnalyticResponse             `json:"device_analytic"`
	LastUpdated    time.Time                          `json:"last_updated"`
}

type DeviceAnalyticResponse struct {
	Mobile  int `json:"mobile"`
	Tablet  int `json:"tablet"`
	Desktop int `json:"desktop"`
	Other   int `json:"other"`
}
