package helper

import (
	"github.com/ilhamfzri/pendek.in/internal/model/domain"
	"github.com/ilhamfzri/pendek.in/internal/model/web"
)

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
