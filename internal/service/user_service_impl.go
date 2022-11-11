package service

import (
	"context"
	"database/sql"
	"time"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/ilhamfzri/pendek.in/internal/helper"
	"github.com/ilhamfzri/pendek.in/internal/model/domain"
	"github.com/ilhamfzri/pendek.in/internal/model/web"
	"github.com/ilhamfzri/pendek.in/internal/repository"
)

type UserServiceImpl struct {
	UserRepository repository.UserRepository
	DB             *sql.DB
	Validate       *validator.Validate
}

func NewUserService(userRepository repository.UserRepository, DB *sql.DB, validate *validator.Validate) UserService {
	return &UserServiceImpl{
		UserRepository: userRepository,
		DB:             DB,
		Validate:       validate,
	}
}

func (service *UserServiceImpl) Register(ctx context.Context, request web.UserRegisterRequest) (web.UserResponse, error) {
	err := service.Validate.Struct(request)
	helper.PanicIfError(err)

	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	// TODO : Hash Password
	user := domain.User{
		Username: request.Username,
		Email:    request.Email,
		Password: request.Password,
	}

	user, repoErr := service.UserRepository.Create(ctx, tx, user)
	if repoErr != nil {
		return helper.ToUserResponse(user), repoErr
	}

	code := uuid.New().String()[:8]
	err = service.UserRepository.CreateVerifyCode(ctx, tx, user.Id, code)
	helper.PanicIfError(err)

	user.Id = 0
	return helper.ToUserResponse(user), nil
}

func (service *UserServiceImpl) Login(ctx context.Context, request web.UserLoginRequest) (web.TokenResponse, error) {
	err := service.Validate.Struct(request)
	helper.PanicIfError(err)

	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	user := domain.User{
		Email:    request.Email,
		Password: request.Password,
	}

	repoErr := service.UserRepository.Login(ctx, tx, user)

	// TODO : GENERATE JWT TOKEN
	token := uuid.New().String()
	valid_until := time.Now()

	tokenResponse := web.TokenResponse{
		AccessToken: token,
		ValidUntil:  valid_until,
	}

	return tokenResponse, repoErr

}
