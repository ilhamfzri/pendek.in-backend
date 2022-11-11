package service

import (
	"context"
	"database/sql"

	"github.com/go-playground/validator"
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

	user := domain.User{
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Username:  request.Username,
		Email:     request.Email,
		Password:  request.Password,
	}

	user, repoErr := service.UserRepository.Create(ctx, tx, user)
	return helper.ToUserResponse(user), repoErr
}
