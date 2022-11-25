package service

import (
	"context"
	"testing"
	"time"

	"github.com/ilhamfzri/pendek.in/app/database"
	"github.com/ilhamfzri/pendek.in/app/logger"
	"github.com/ilhamfzri/pendek.in/helper"
	"github.com/ilhamfzri/pendek.in/internal/model/domain"
	"github.com/ilhamfzri/pendek.in/internal/model/web"
	"github.com/ilhamfzri/pendek.in/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

var ctx = context.Background()
var db = database.NewDatabaseConnectionMock()
var log = new(logger.Logger)

func TestMain(m *testing.M) {
	m.Run()
}

var userNotFound = domain.User{
	Username: "testuser01",
	Email:    "testuser01@mail.com",
	Password: "testpassword01",
	Verified: false,
}

var userFound = domain.User{
	Username: "testuser02",
	Email:    "testuser02@mail.com",
	Password: "testpassword02", // Hashed Value : $2a$14$SIxTHeN2csRDv.WqW2H5M.0pDPli7p1OAsikanREUi2B5tt.KQy.i
	Verified: false,
}

func TestUserServiceRegister(t *testing.T) {
	var jwt = new(helper.JwtMock)
	var userRepository = new(repository.UserRepositoryMock)
	var userService = NewUserService(userRepository, db, log, jwt)

	userRepository.Mock.On("Create", mock.Anything, mock.Anything, mock.Anything).Return(userNotFound, nil)
	userRepository.Mock.On("FindByUsername", mock.Anything, mock.Anything, userNotFound.Username).Return(domain.User{}, gorm.ErrRecordNotFound)
	userRepository.Mock.On("FindByUsername", mock.Anything, mock.Anything, userFound.Username).Return(userFound, nil)
	userRepository.Mock.On("FindByUsername", mock.Anything, mock.Anything, userFound.Username+"1").Return(domain.User{}, gorm.ErrRecordNotFound)
	userRepository.Mock.On("FindByEmail", mock.Anything, mock.Anything, userFound.Email).Return(userFound, nil)
	userRepository.Mock.On("FindByEmail", mock.Anything, mock.Anything, userNotFound.Email).Return(domain.User{}, gorm.ErrRecordNotFound)

	t.Run(
		"[Register][Success]", func(t *testing.T) {
			request := web.UserRegisterRequest{
				Username: userNotFound.Username,
				Email:    userNotFound.Email,
				Password: userNotFound.Password,
			}

			user, err := userService.Register(ctx, request)

			assert.Equal(t, request.Username, user.Username)
			assert.Equal(t, request.Email, user.Email)
			assert.Nil(t, err)
		},
	)

	t.Run(
		"[Register][Failed:Username Used]", func(t *testing.T) {
			request := web.UserRegisterRequest{
				Username: userFound.Username,
				Email:    userFound.Email,
				Password: userFound.Password,
			}

			user, err := userService.Register(ctx, request)

			assert.Equal(t, web.UserResponse{}, user)
			assert.NotNil(t, err)
			assert.Equal(t, ErrUsernameFound, err)
		},
	)

	t.Run(
		"[Register][Failed:Email Used]", func(t *testing.T) {
			request := web.UserRegisterRequest{
				Username: userFound.Username + "1",
				Email:    userFound.Email,
				Password: userFound.Password,
			}

			user, err := userService.Register(ctx, request)

			assert.Equal(t, web.UserResponse{}, user)
			assert.NotNil(t, err)
			assert.Equal(t, ErrEmailFound, err)
		},
	)

	userRepository.AssertExpectations(t)
}

func TestUserServiceLogin(t *testing.T) {
	var jwt = new(helper.JwtMock)
	var userRepository = new(repository.UserRepositoryMock)
	var userService = NewUserService(userRepository, db, log, jwt)

	newUserFound := userFound
	newUserFound.Password = "$2a$14$SIxTHeN2csRDv.WqW2H5M.0pDPli7p1OAsikanREUi2B5tt.KQy.i"

	newUserFoundVerified := userFound
	newUserFoundVerified.Email = newUserFoundVerified.Email + "verified"
	newUserFoundVerified.Password = "$2a$14$SIxTHeN2csRDv.WqW2H5M.0pDPli7p1OAsikanREUi2B5tt.KQy.i"
	newUserFoundVerified.Verified = true

	userRepository.Mock.On("FindByEmail", mock.Anything, mock.Anything, userNotFound.Email).Return(domain.User{}, gorm.ErrRecordNotFound)
	userRepository.Mock.On("FindByEmail", mock.Anything, mock.Anything, newUserFound.Email).Return(newUserFound, nil)
	userRepository.Mock.On("FindByEmail", mock.Anything, mock.Anything, newUserFoundVerified.Email).Return(newUserFoundVerified, nil)

	jwt.Mock.On("NewToken", mock.Anything, mock.Anything, mock.Anything).Return("SIxTHeN2csRDv.WqW2H5M.0pDPli7p1OAsikanREUi2B5tt", time.Now(), nil)

	t.Run("[Login][Success]", func(t *testing.T) {
		request := web.UserLoginRequest{
			Email:    newUserFoundVerified.Email,
			Password: "testpassword02",
		}
		token, err := userService.Login(ctx, request)
		assert.NotNil(t, token.AccessToken)
		assert.Nil(t, err)
		assert.IsType(t, token, web.TokenResponse{})
		assert.IsType(t, token.ValidUntil, time.Time{})
	})

	t.Run("[Failed: Email Not Found", func(t *testing.T) {
		request := web.UserLoginRequest{
			Email:    userNotFound.Email,
			Password: userNotFound.Password,
		}
		token, err := userService.Login(ctx, request)
		assert.NotNil(t, err)
		assert.Equal(t, ErrEmailNotFound, err)
		assert.IsType(t, token, web.TokenResponse{})
	})

	t.Run("[Failed: Password Incorrect", func(t *testing.T) {
		request := web.UserLoginRequest{
			Email:    newUserFound.Email,
			Password: "testpassword03",
		}
		token, err := userService.Login(ctx, request)
		assert.NotNil(t, err)
		assert.Equal(t, ErrPasswordIncorrect, err)
		assert.IsType(t, token, web.TokenResponse{})
	})

	t.Run("[Failed: Email Not Verified", func(t *testing.T) {
		request := web.UserLoginRequest{
			Email:    newUserFound.Email,
			Password: "testpassword02",
		}
		token, err := userService.Login(ctx, request)
		assert.NotNil(t, err)
		assert.Equal(t, ErrEmailNotVerified, err)
		assert.IsType(t, token, web.TokenResponse{})
	})

	userRepository.AssertExpectations(t)
}

func TestUserChangePassword(t *testing.T) {
	var jwt = new(helper.JwtMock)
	var userRepository = new(repository.UserRepositoryMock)
	var userService = NewUserService(userRepository, db, log, jwt)

	newUserFound := userFound
	newUserFound.Password = "$2a$14$SIxTHeN2csRDv.WqW2H5M.0pDPli7p1OAsikanREUi2B5tt.KQy.i"

	dummyJwt := "ASDEFGHJKDSANEQWENEWNQENWN"

	userRepository.Mock.On("FindByEmail", mock.Anything, mock.Anything, newUserFound.Email).Return(newUserFound, nil)
	userRepository.Mock.On("UpdatePassword", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	jwt.Mock.On("GetClaims", dummyJwt).Return(helper.JwtUserClaims{
		Id:       "123456",
		Username: newUserFound.Username,
		Email:    newUserFound.Email,
	})

	t.Run("[Change-Password][Success]", func(t *testing.T) {
		request := web.UserChangePasswordRequest{
			CurrentPassword: "testpassword02",
			NewPassword:     "testpassword03",
		}
		err := userService.ChangePassword(ctx, request, dummyJwt)
		assert.Nil(t, err)
	})

	t.Run("[Change-Password][Failed: Current Password Incorrect]", func(t *testing.T) {

		request := web.UserChangePasswordRequest{
			CurrentPassword: "testpassword04",
			NewPassword:     "testpassword03",
		}
		err := userService.ChangePassword(ctx, request, dummyJwt)
		assert.NotNil(t, err)
		assert.Equal(t, err, ErrCurrentPasswordIncorrect)
	})

	userRepository.AssertExpectations(t)
}

func TestUserEmailVerification(t *testing.T) {
	var jwt = new(helper.JwtMock)
	var userRepository = new(repository.UserRepositoryMock)
	var userService = NewUserService(userRepository, db, log, jwt)

	var newUserFound = userFound
	newUserFound.VerificationCode = "ABCDEF"

	dummyJwt := "ASDEFGHJKDSANEQWENEWNQENWN"

	jwt.Mock.On("GetClaims", dummyJwt).Return(helper.JwtUserClaims{
		Id:       "123456",
		Username: newUserFound.Username,
		Email:    newUserFound.Email,
	})

	userRepository.Mock.On("FindByEmail", mock.Anything, mock.Anything, newUserFound.Email).Return(newUserFound, nil)
	userRepository.Mock.On("FindByEmail", mock.Anything, mock.Anything, userNotFound.Email).Return(domain.User{}, gorm.ErrRecordNotFound)
	userRepository.Mock.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(newUserFound, nil)

	t.Run("[EmailVerification][Success]", func(t *testing.T) {
		request := web.UserEmailVerificationRequest{
			Email:            newUserFound.Email,
			VerificationCode: newUserFound.VerificationCode,
		}
		userResponse, err := userService.EmailVerification(ctx, request)
		assert.Nil(t, err)
		assert.IsType(t, userResponse, web.UserResponse{})
	})

	t.Run("[EmailVerification][Failed: Email Not Registered", func(t *testing.T) {
		request := web.UserEmailVerificationRequest{
			Email:            userNotFound.Email,
			VerificationCode: userNotFound.VerificationCode,
		}
		userResponse, err := userService.EmailVerification(ctx, request)
		assert.NotNil(t, err)
		assert.Equal(t, ErrEmailNotFound, err)
		assert.IsType(t, userResponse, web.UserResponse{})
	})

	t.Run("[EmailVerification][Failed: Verification Code Expired Or Invalid", func(t *testing.T) {
		request := web.UserEmailVerificationRequest{
			Email:            newUserFound.Email,
			VerificationCode: "ABCDWQ",
		}
		userResponse, err := userService.EmailVerification(ctx, request)
		assert.NotNil(t, err)
		assert.Equal(t, ErrVerificationCodeInvalid, err)
		assert.IsType(t, userResponse, web.UserResponse{})
	})
}
