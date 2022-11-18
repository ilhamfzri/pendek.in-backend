package repository

import (
	"context"

	"github.com/ilhamfzri/pendek.in/app/logger"
	"github.com/ilhamfzri/pendek.in/internal/model/domain"
	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	Log *logger.Logger
}

func (repository *UserRepositoryImpl) Create(ctx context.Context, tx *gorm.DB, user domain.User) (domain.User, error) {
	result := tx.WithContext(ctx).Create(&user)
	return user, result.Error
}

func (repository *UserRepositoryImpl) FindByUsername(ctx context.Context, tx *gorm.DB, username string) (domain.User, error) {
	var user domain.User
	result := tx.WithContext(ctx).Where("username = ?", username).First(&user)
	return user, result.Error
}

func (repository *UserRepositoryImpl) FindByEmail(ctx context.Context, tx *gorm.DB, email string) (domain.User, error) {
	var user domain.User
	result := tx.WithContext(ctx).Where("email = ?", email).First(&user)
	return user, result.Error
}

func (repository *UserRepositoryImpl) Update(ctx context.Context, tx *gorm.DB, user domain.User) (domain.User, error) {
	result := tx.WithContext(ctx).Model(&domain.User{}).Where("id = ?", user.ID).
		Updates(
			domain.User{
				FullName:  user.FullName,
				Bio:       user.Bio,
				LastLogin: user.LastLogin,
				Verified:  user.Verified,
			})

	return user, result.Error
}

func (repository *UserRepositoryImpl) UpdatePassword(ctx context.Context, tx *gorm.DB, userId string, newPassword string) error {
	result := tx.WithContext(ctx).Model(&domain.User{}).Where("id = ?", userId).Update("password", newPassword)
	return result.Error
}
