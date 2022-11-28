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
}

type SocialMediaLinkService interface {
	GetAllTypes(ctx context.Context) ([]web.SocialMediaTypeResponse, error)
	CreateLink(ctx context.Context, request web.SocialMediaLinkCreateRequest, host string, jwtToken string) (web.SocialMediaLinkResponse, error)
	UpdateLink(ctx context.Context, request web.SocialMediaLinkUpdateRequest, host string, jwtToken string) (web.SocialMediaLinkResponse, error)
	GetAllLink(ctx context.Context, host string, jwtToken string) ([]web.SocialMediaLinkResponse, error)
	RedirectLink(ctx context.Context, request web.SocialMediaLinkRedirectRequest) (string, error)
}
