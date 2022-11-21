package repository

import (
	"context"

	"github.com/ilhamfzri/pendek.in/internal/model/domain"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, tx *gorm.DB, user domain.User) (domain.User, error)
	FindByUsername(ctx context.Context, tx *gorm.DB, username string) (domain.User, error)
	FindByEmail(ctx context.Context, tx *gorm.DB, email string) (domain.User, error)
	Update(ctx context.Context, tx *gorm.DB, user domain.User) (domain.User, error)
	UpdatePassword(ctx context.Context, tx *gorm.DB, userId string, newPassword string) error
}
