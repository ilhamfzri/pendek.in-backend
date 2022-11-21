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
