package helper

import (
	"fmt"

	"github.com/ilhamfzri/pendek.in/internal/model/domain"
	"github.com/ilhamfzri/pendek.in/internal/model/web"
)

var thumbnailResourceEndpointPath = "v1/resources/thumbnail"

func UserDomainToResponse(user *domain.User) web.UserResponse {
	return web.UserResponse{
		ID:         user.ID,
		Username:   user.Username,
		FullName:   user.FullName,
		Bio:        user.Bio,
		Email:      user.Email,
		ProfilePic: user.ProfilePic,
	}
}

func SocialMediaLinkDomainToResponse(smld *domain.SocialMediaLink, host string, username string) web.SocialMediaLinkResponse {
	return web.SocialMediaLinkResponse{
		TypeID:          smld.TypeID,
		SocialMediaName: smld.SocialMediaType.Name,
		LinkOrUsername:  smld.LinkOrUsername,
		Activate:        smld.Activate,
		RedirectLink:    GenerateRedirectLink(host, username, smld.SocialMediaType.Name),
	}
}

func DeviceAnalyticDomainToResponse(da *domain.DeviceAnalytic) web.DeviceAnalyticResponse {
	return web.DeviceAnalyticResponse{
		Mobile:  da.Mobile,
		Tablet:  da.Tablet,
		Desktop: da.Desktop,
		Other:   da.Other,
	}
}

func SocialMediaAnalyticDomainToResponse(sma *domain.SocialMediaAnalytic, socialMediaName string) web.SocialMediaAnalyticResponse {
	return web.SocialMediaAnalyticResponse{
		SocialMediaLinkID: sma.SocialMediaLinkID,
		SocialMediaName:   socialMediaName,
		ClickCount:        sma.ClickCount,
		ViewCount:         sma.ViewCount,
		DeviceAnalytic:    DeviceAnalyticDomainToResponse(&sma.DeviceAnalytic),
		Datetime:          sma.Date.Format("2006-01-02"),
		LastUpdated:       sma.UpdatedAt,
	}
}

func CustomThumbnailDomainToResponse(ct *domain.CustomThumbnail, domain string) web.ThumbnailResponse {
	thumbnailUrl := fmt.Sprintf("%s/%s/%s.jpg", domain, thumbnailResourceEndpointPath, ct.ImageID)
	return web.ThumbnailResponse{
		ID:           ct.ID,
		ThumbnailUrl: thumbnailUrl,
	}
}

func ThumbnailDomainToResponse(t *domain.Thumbnail) web.ThumbnailResponse {
	return web.ThumbnailResponse{
		ID:           t.ID,
		Name:         t.Name,
		ThumbnailUrl: t.IconUrl,
	}
}

func CustomLinkDomainToResponse(l *domain.CustomLink) web.CustomLinkResponse {
	customLinkResponse := web.CustomLinkResponse{
		ID:            l.ID,
		Title:         l.Title,
		ShortLinkCode: l.ShortLinkCode,
		LongLink:      l.LongLink,
		ShowOnProfile: l.ShowOnProfile,
		Activate:      l.Activate,
	}

	if l.ThumbnailID != nil {
		customLinkResponse.ThumbnailID = *l.ThumbnailID
	}

	if l.CustomThumbnailID != nil {
		customLinkResponse.CustomThumbnailID = *l.CustomThumbnailID
	}

	return customLinkResponse
}
