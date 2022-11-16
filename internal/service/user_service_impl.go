package service

import (
	"context"
	"database/sql"

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
	Jwt            *helper.Jwt
}

func NewUserService(userRepository repository.UserRepository, DB *sql.DB, validate *validator.Validate, jwtClient *helper.Jwt) UserService {
	return &UserServiceImpl{
		UserRepository: userRepository,
		DB:             DB,
		Validate:       validate,
		Jwt:            jwtClient,
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

	user, repoErr := service.UserRepository.Login(ctx, tx, user)
	token, valid_until := service.Jwt.NewToken(user.Username, user.Email)
	tokenResponse := web.TokenResponse{
		AccessToken: token,
		ValidUntil:  valid_until,
	}

	return tokenResponse, repoErr
}

func (service *UserServiceImpl) Verify(ctx context.Context, request web.UserVerifyRequest) error {
	err := service.Validate.Struct(request)
	helper.PanicIfError(err)

	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	repoErr := service.UserRepository.Verify(ctx, tx, request.Email, request.Code)
	return repoErr
}

func (service *UserServiceImpl) ChangePassword(ctx context.Context, request web.UserChangePasswordRequest, token string) error {
	err := service.Validate.Struct(request)
	helper.PanicIfError(err)

	claims := service.Jwt.GetClaims(token)
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	repoErr := service.UserRepository.UpdatePassword(ctx, tx, claims.Username, request.CurrentPassword, request.NewPassword)
	return repoErr
}

func (service *UserServiceImpl) UpdateInformation(ctx context.Context, request web.UserUpdateInfoRequest, token string) (web.UserResponse, error) {
	err := service.Validate.Struct(request)
	helper.PanicIfError(err)

	claims := service.Jwt.GetClaims(token)
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	userRequest := domain.User{
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Bio:       request.Bio,
	}

	userData, err := service.UserRepository.FindByUsername(ctx, tx, claims.Username)
	helper.PanicIfError(err)

	helper.UserSetDefaultValue(&userRequest, &userData)
	userResult, repoErr := service.UserRepository.Update(ctx, tx, userRequest)
	userResponse := helper.ToUserResponse(userResult)

	return userResponse, repoErr
}
