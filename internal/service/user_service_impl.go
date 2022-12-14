package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"path"
	"time"

	"github.com/google/uuid"
	"github.com/ilhamfzri/pendek.in/app/logger"
	"github.com/ilhamfzri/pendek.in/app/mail"
	"github.com/ilhamfzri/pendek.in/helper"
	"github.com/ilhamfzri/pendek.in/internal/model/domain"
	"github.com/ilhamfzri/pendek.in/internal/model/web"
	"github.com/ilhamfzri/pendek.in/internal/repository"
	"github.com/nfnt/resize"
	"gorm.io/gorm"
)

type UserServiceImpl struct {
	Repository repository.UserRepository
	MailClient *mail.MailClient
	DB         *gorm.DB
	Logger     *logger.Logger
	Jwt        helper.IJwt
}

func NewUserService(repository repository.UserRepository, mailClient *mail.MailClient, DB *gorm.DB, logger *logger.Logger, jwt helper.IJwt) UserService {
	return &UserServiceImpl{
		Repository: repository,
		MailClient: mailClient,
		DB:         DB,
		Logger:     logger,
		Jwt:        jwt,
	}
}

var (
	ErrUserService              = "[User Service] Failed Execute User Service"
	ErrUsernameFound            = errors.New("username is already used")
	ErrEmailFound               = errors.New("email is already registered")
	ErrEmailNotFound            = errors.New("email isn't registered")
	ErrEmailNotVerified         = errors.New("email isn't verified")
	ErrPasswordIncorrect        = errors.New("password incorrect")
	ErrCurrentPasswordIncorrect = errors.New("current password incorrect")
	ErrVerificationCodeInvalid  = errors.New("verification code expired or invalid")
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
	if repoErr == nil && userData.Username != "" {
		return web.UserResponse{}, ErrUsernameFound
	}

	if !errors.Is(repoErr, gorm.ErrRecordNotFound) && repoErr != nil {
		service.Logger.PanicIfErr(repoErr, ErrUserService)
	}

	// It's checking if the email is already registered or not.
	userData, repoErr = service.Repository.FindByEmail(ctx, tx, user.Email)
	if repoErr == nil && userData.Email != "" {
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

	errMail := service.MailClient.SendVerificationEmail(user.Email, verificationCode)
	service.Logger.PanicIfErr(errMail, ErrUserService)

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

	// It's updating the last login time of the user.
	userData.LastLogin = time.Now()
	_, errRepo := service.Repository.Update(ctx, tx, userData)
	service.Logger.PanicIfErr(errRepo, ErrUserService)

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

func (service *UserServiceImpl) Update(ctx context.Context, request web.UserUpdateRequest,
	jwtToken string) (web.UserResponse, error) {

	// It's getting the claims from the token.
	claims := service.Jwt.GetClaims(jwtToken)

	// It's a transaction.
	tx := service.DB.Begin()
	defer helper.CommitOrRollback(tx)

	// It's getting the user data from the database.
	user, errRepo := service.Repository.FindByEmail(ctx, tx, claims.Email)
	service.Logger.PanicIfErr(errRepo, ErrUserService)

	if request.FullName != "" {
		user.FullName = request.FullName
	}

	if request.Bio != "" {
		user.Bio = request.Bio
	}

	// It's updating the user data in the database.
	user, errRepo = service.Repository.Update(ctx, tx, user)
	service.Logger.PanicIfErr(errRepo, ErrUserService)

	webResponse := helper.UserDomainToResponse(&user)
	return webResponse, errRepo
}

func (service *UserServiceImpl) EmailVerification(ctx context.Context,
	request web.UserEmailVerificationRequest) (web.UserResponse, error) {

	// It's a transaction.
	tx := service.DB.Begin()
	defer helper.CommitOrRollback(tx)

	// It's checking if the email is already registered or not.
	userDomain, errRepo := service.Repository.FindByEmail(ctx, tx, request.Email)
	if errors.Is(errRepo, gorm.ErrRecordNotFound) {
		return web.UserResponse{}, ErrEmailNotFound
	}

	service.Logger.PanicIfErr(errRepo, ErrUserService)

	// It's checking if the verification code is correct or not.
	if request.VerificationCode != userDomain.VerificationCode {
		return web.UserResponse{}, ErrVerificationCodeInvalid
	}

	// It's updating the user verification in the database.
	userDomain.Verified = true
	userDomain, errRepo = service.Repository.Update(ctx, tx, userDomain)
	service.Logger.PanicIfErr(errRepo, ErrUserService)

	userResponse := helper.UserDomainToResponse(&userDomain)
	return userResponse, nil
}

func (service *UserServiceImpl) GenerateToken(ctx context.Context, jwtToken string) (web.TokenResponse, error) {
	// It's getting the claims from the token.
	claims := service.Jwt.GetClaims(jwtToken)

	// It's creating a new token and assign it to accessToken variable.
	accessToken, validUntil, err := service.Jwt.NewToken(claims.Id, claims.Username, claims.Email)
	service.Logger.PanicIfErr(err, ErrUserService)

	webResponse := web.TokenResponse{
		AccessToken: accessToken,
		ValidUntil:  validUntil,
	}

	return webResponse, nil
}

func (service *UserServiceImpl) ChangeProfilePicture(ctx context.Context, imgByte []byte, jwtToken string) error {
	// It's getting the claims from the token.
	claims := service.Jwt.GetClaims(jwtToken)

	// It's a transaction.
	tx := service.DB.Begin()
	defer helper.CommitOrRollback(tx)

	// It's decoding the image from the byte array.
	reader := bytes.NewReader(imgByte)
	img, _, err := image.Decode(reader)
	service.Logger.PanicIfErr(err, ErrUserService)

	// It's resizing the image to 180x180 pixels.
	resizeImg := resize.Resize(300, 300, img, resize.Lanczos3)

	uuid := uuid.New().String()
	fileName := fmt.Sprintf("%s.jpg", uuid)
	userResourcePath := os.Getenv("PROFILE_IMG_DIR")
	filePath := path.Join(userResourcePath, fileName)

	out, err := os.Create(filePath)
	service.Logger.PanicIfErr(err, ErrUserService)
	defer out.Close()

	jpeg.Encode(out, resizeImg, nil)

	// It's checking if the email is already registered or not.
	userDomain, errRepo := service.Repository.FindByEmail(ctx, tx, claims.Email)
	service.Logger.PanicIfErr(errRepo, ErrUserService)

	userDomain.ProfilePic = uuid
	_, errRepo = service.Repository.Update(ctx, tx, userDomain)
	service.Logger.PanicIfErr(errRepo, ErrUserService)

	return nil
}

func (service *UserServiceImpl) GetProfileData(ctx context.Context, request web.UserProfileRequest) web.UserResponse {
	// It's a transaction.
	tx := service.DB.Begin()
	defer helper.CommitOrRollback(tx)

	// It's checking if the username is already used or not.
	userData, repoErr := service.Repository.FindByUsername(ctx, tx, request.Username)
	if repoErr == nil && !errors.Is(repoErr, gorm.ErrRecordNotFound) {
		service.Logger.PanicIfErr(repoErr, ErrUserService)
	}

	userResponse := helper.UserDomainToResponse(&userData)
	return userResponse
}

func (service *UserServiceImpl) GetCurrentProfile(ctx context.Context, jwtToken string) (web.UserResponse, error) {
	// It's getting the claims from the token.
	claims := service.Jwt.GetClaims(jwtToken)

	// It's a transaction.
	tx := service.DB.Begin()
	defer helper.CommitOrRollback(tx)

	// It's checking if the username is already used or not.
	userData, repoErr := service.Repository.FindByEmail(ctx, tx, claims.Email)
	if repoErr == nil && !errors.Is(repoErr, gorm.ErrRecordNotFound) {
		service.Logger.PanicIfErr(repoErr, ErrUserService)
	}

	userResponse := helper.UserDomainToResponse(&userData)
	return userResponse, nil

}
