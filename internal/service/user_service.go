package service

import (
	"context"

	"github.com/ilhamfzri/pendek.in/internal/model/web"
)

type UserService interface {
	Register(ctx context.Context, request web.UserRegisterRequest) (web.UserResponse, error)
	Login(ctx context.Context, request web.UserLoginRequest) (web.TokenResponse, error)
	Verify(ctx context.Context, request web.UserVerifyRequest) error
	ChangePassword(ctx context.Context, request web.UserChangePasswordRequest, token string) error
	UpdateInformation(ctx context.Context, request web.UserUpdateInfoRequest, token string) (web.UserResponse, error)
}
