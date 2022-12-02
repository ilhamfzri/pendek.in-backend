package service

import (
	"context"

	"github.com/ilhamfzri/pendek.in/internal/model/web"
)

type UserService interface {
	Register(ctx context.Context, request web.UserRegisterRequest) (web.UserResponse, error)
	Login(ctx context.Context, request web.UserLoginRequest) (web.TokenResponse, error)
	ChangePassword(ctx context.Context, request web.UserChangePasswordRequest, jwtToken string) error
	Update(ctx context.Context, request web.UserUpdateRequest, jwtToken string) (web.UserResponse, error)
	EmailVerification(ctx context.Context, request web.UserEmailVerificationRequest) (web.UserResponse, error)
	GenerateToken(ctx context.Context, jwtToken string) (web.TokenResponse, error)
	ChangeProfilePicture(ctx context.Context, imgByte []byte, jwtToken string) error
	GetProfileData(ctx context.Context, request web.UserProfileRequest) web.UserResponse
	GetCurrentProfile(ctx context.Context, jwtToken string) (web.UserResponse, error)
}

type SocialMediaLinkService interface {
	GetAllTypes(ctx context.Context) ([]web.SocialMediaTypeResponse, error)
	CreateLink(ctx context.Context, request web.SocialMediaLinkCreateRequest, host string, jwtToken string) (web.SocialMediaLinkResponse, error)
	UpdateLink(ctx context.Context, request web.SocialMediaLinkUpdateRequest, host string, jwtToken string) (web.SocialMediaLinkResponse, error)
	GetAllLink(ctx context.Context, host string, jwtToken string) ([]web.SocialMediaLinkResponse, error)
	RedirectLink(ctx context.Context, request web.SocialMediaLinkRedirectRequest) (string, uint, error)
	GetAllLinkProfile(ctx context.Context, domainName string, userID string, username string) []web.UserProfileSocialMediaResponse
}

type SocialMediaAnalytic interface {
	SaveInteraction(ctx context.Context, request web.SocialMediaAnalyticInteractionRequest) error
	GetLinkAnalytic(ctx context.Context, request web.SocialMediaAnalyticGetRequest, jwtToken string) ([]web.SocialMediaAnalyticResponse, error)
	GetSummaryLinkAnalytic(ctx context.Context, jwtToken string) (web.SocialMediaAnalyticSummaryResponse, error)
}

type CustomLinkService interface {
	CreateLink(ctx context.Context, request web.CustomLinkCreateRequest, domainName string, jwtToken string) (web.CustomLinkResponse, error)
	UpdateLink(ctx context.Context, request web.CustomLinkUpdateRequest, domainName string, jwtToken string) (web.CustomLinkResponse, error)
	GetLink(ctx context.Context, request web.CustomLinkGetRequest, domainName string, jwtToken string) (web.CustomLinkResponse, error)
	GetAllLink(ctx context.Context, domainName string, jwtToken string) ([]web.CustomLinkResponse, error)
	GetAllThumbnail(ctx context.Context) ([]web.ThumbnailResponse, error)
	GetUserThumbnail(ctx context.Context, domainName string, jwtToken string) ([]web.ThumbnailResponse, error)
	UploadCustomThumbnail(ctx context.Context, imgData []byte, domainName string, jwtToken string) (web.ThumbnailResponse, error)
	CheckShortLinkAvaibility(ctx context.Context, request web.CustomLinkCheckShortCodeAvaibilityRequest) error
	RedirectLink(ctx context.Context, request web.CustomLinkRedirectRequest) (string, uint, error)
	GetAllLinkProfile(ctx context.Context, domainName string, userID string, username string) []web.UserProfileCustomLinkResponse
}

type CustomLinkAnalyticService interface {
	SaveInteraction(ctx context.Context, request web.CustomLinkAnalyticInteractionRequest) error
	GetLinkAnalytic(ctx context.Context, request web.CustomLinkAnalyticGetRequest, jwtToken string) ([]web.CustomLinkAnalyticResponse, error)
	GetSummaryLinkAnalytic(ctx context.Context, jwtToken string) (web.CustomLinkAnalyticSummaryResponse, error)
}
