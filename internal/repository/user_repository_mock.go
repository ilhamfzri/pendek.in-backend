package repository

import (
	"context"

	"github.com/ilhamfzri/pendek.in/internal/model/domain"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type UserRepositoryMock struct {
	mock.Mock
}

func (repository UserRepositoryMock) Create(ctx context.Context, tx *gorm.DB, user domain.User) (domain.User, error) {
	args := repository.Mock.Called(ctx, tx, user)
	return args.Get(0).(domain.User), args.Error(1)
}

func (repository UserRepositoryMock) FindByUsername(ctx context.Context, tx *gorm.DB, username string) (domain.User, error) {
	args := repository.Mock.Called(ctx, tx, username)
	return args.Get(0).(domain.User), args.Error(1)
}

func (repository UserRepositoryMock) FindByEmail(ctx context.Context, tx *gorm.DB, email string) (domain.User, error) {
	args := repository.Mock.Called(ctx, tx, email)
	return args.Get(0).(domain.User), args.Error(1)
}

func (repository UserRepositoryMock) Update(ctx context.Context, tx *gorm.DB, user domain.User) (domain.User, error) {
	return domain.User{}, nil
}

func (repository UserRepositoryMock) UpdatePassword(ctx context.Context, tx *gorm.DB, userId string, newPassword string) error {
	args := repository.Mock.Called(ctx, tx, userId, newPassword)
	return args.Error(0)
}
