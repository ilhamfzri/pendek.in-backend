package service

import (
	"context"

	"github.com/ilhamfzri/pendek.in/app/logger"
	"github.com/ilhamfzri/pendek.in/helper"
	"github.com/ilhamfzri/pendek.in/internal/model/web"
	"github.com/ilhamfzri/pendek.in/internal/repository"
	"gorm.io/gorm"
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

func NewUserService(repository repository.UserRepository, DB *gorm.DB, logger *logger.Logger, jwt *helper.Jwt) UserService {
	return &UserServiceImpl{
		Repository: repository,
		DB:         DB,
		Logger:     logger,
		Jwt:        jwt,
	}
}
