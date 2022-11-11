package service

import (
	"context"

	"github.com/ilhamfzri/pendek.in/internal/model/web"
)

type UserService interface {
	Register(ctx context.Context, request web.UserRegisterRequest) (web.UserResponse, error)
	FindByUsername(ctx context.Context) (web.UserResponse, error)
}
