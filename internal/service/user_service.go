package service

import (
	"context"
	"errors"

	"github.com/ilhamfzri/pendek.in/app/logger"
	"github.com/ilhamfzri/pendek.in/helper"
	"github.com/ilhamfzri/pendek.in/internal/model/domain"
	"github.com/ilhamfzri/pendek.in/internal/model/web"
	"github.com/ilhamfzri/pendek.in/internal/repository"
	"gorm.io/gorm"
)

type UserServiceImpl struct {
	Repository repository.UserRepository
	DB         *gorm.DB
	Logger     *logger.Logger
	Jwt        *helper.Jwt
}

var (
	ErrUserService              = "[User Service] Failed Execute User Service"
	ErrUsernameFound            = errors.New("username is already used")
	ErrEmailFound               = errors.New("email is already registered")
	ErrEmailNotFound            = errors.New("email isn't registered")
	ErrEmailNotVerified         = errors.New("email isn't verified")
	ErrPasswordIncorrect        = errors.New("password incorrect")
	ErrCurrentPasswordIncorrect = errors.New("current password incorrect")
)

func (service *UserServiceImpl) Register(ctx context.Context, request web.UserRegisterRequest) (web.UserResponse, error) {
	// It's a transaction.
	tx := service.DB.Begin()
	defer helper.CommitOrRollback(tx)

	user := domain.User{
		Username: request.Username,
		Email:    request.Email,
		Password: request.Password,
	}

	// It's checking if the username is already used or not.
	userData, repoErr := service.Repository.FindByUsername(ctx, tx, user.Username)
	if (userData != domain.User{}) {
		return web.UserResponse{}, ErrUsernameFound
	}

	if !errors.Is(repoErr, gorm.ErrRecordNotFound) && repoErr != nil {
		service.Logger.PanicIfErr(repoErr, ErrUserService)
	}

	// It's checking if the email is already registered or not.
	userData, repoErr = service.Repository.FindByEmail(ctx, tx, user.Email)
	if (userData != domain.User{}) {
		return web.UserResponse{}, ErrEmailFound
	}

	if !errors.Is(repoErr, gorm.ErrRecordNotFound) && repoErr != nil {
		service.Logger.PanicIfErr(repoErr, ErrUserService)
	}

	// It's generating a 6 digit OTP and assign it to user.VerificationCode
	verificationCode, err := helper.GenerateOTP(6)
	service.Logger.PanicIfErr(err, ErrUserService)
	user.VerificationCode = verificationCode

	// It's hashing the password before saving it to the database.
	hashPassword, err := helper.HashPassword(user.Password)
	service.Logger.PanicIfErr(err, ErrUserService)
	user.Password = hashPassword

	// It's creating a new user and assign it to user variable.
	user, repoErr = service.Repository.Create(ctx, tx, user)
	service.Logger.PanicIfErr(repoErr, ErrUserService)

	webResponse := helper.UserDomainToResponse(&user)
	return webResponse, nil
}

func (service *UserServiceImpl) Login(ctx context.Context, request web.UserLoginRequest) (web.TokenResponse, error) {
	// It's a transaction.
	tx := service.DB.Begin()
	defer helper.CommitOrRollback(tx)

	user := domain.User{
		Email:    request.Email,
		Password: request.Password,
	}

	// It's checking if the email is already registered or not.
	userData, repoErr := service.Repository.FindByEmail(ctx, tx, user.Email)
	if repoErr != nil {
		if errors.Is(repoErr, gorm.ErrRecordNotFound) {
			return web.TokenResponse{}, ErrEmailNotFound
		} else {
			service.Logger.PanicIfErr(repoErr, ErrUserService)
		}
	}

	// It's checking if the password is correct or not.
	valid := helper.CheckPasswordHash(user.Password, userData.Password)
	if !valid {
		return web.TokenResponse{}, ErrPasswordIncorrect
	}

	// It's checking if the email is verified or not.
	if !userData.Verified {
		return web.TokenResponse{}, ErrEmailNotVerified
	}

	// It's creating a new token and assign it to accessToken variable.
	accessToken, validUntil, err := service.Jwt.NewToken(userData.ID, userData.Username, userData.Email)
	service.Logger.PanicIfErr(err, ErrUserService)

	webResponse := web.TokenResponse{
		AccessToken: accessToken,
		ValidUntil:  validUntil,
	}

	return webResponse, nil
}

func (service *UserServiceImpl) ChangePassword(ctx context.Context, request web.UserChangePasswordRequest, jwtToken string) error {

	// It's getting the claims from the token.
	claims := service.Jwt.GetClaims(jwtToken)

	// It's a transaction.
	tx := service.DB.Begin()
	defer helper.CommitOrRollback(tx)

	user, errRepo := service.Repository.FindByEmail(ctx, tx, claims.Email)
	service.Logger.PanicIfErr(errRepo, ErrUserService)

	valid := helper.CheckPasswordHash(request.CurrentPassword, user.Password)
	if !valid {
		return ErrCurrentPasswordIncorrect
	}

	newHashPassword, err := helper.HashPassword(request.NewPassword)
	service.Logger.PanicIfErr(err, ErrUserService)

	errRepo = service.Repository.UpdatePassword(ctx, tx, user.ID, newHashPassword)
	service.Logger.PanicIfErr(errRepo, ErrUserService)

	return nil
}
