package repository

import (
	"context"
	"time"

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

type SocialMediaTypeRepository interface {
	Create(ctx context.Context, tx *gorm.DB, socialMediaType domain.SocialMediaType) (domain.SocialMediaType, error)
	FindByName(ctx context.Context, tx *gorm.DB, name string) (domain.SocialMediaType, error)
	FindByID(ctx context.Context, tx *gorm.DB, id int) (domain.SocialMediaType, error)
	FetchAll(ctx context.Context, tx *gorm.DB) ([]domain.SocialMediaType, error)
}

type SocialMediaLinkRepository interface {
	Create(ctx context.Context, tx *gorm.DB, socialMediaLink domain.SocialMediaLink) (domain.SocialMediaLink, error)
	Update(ctx context.Context, tx *gorm.DB, socialMediaLink domain.SocialMediaLink) (domain.SocialMediaLink, error)
	FindByUserID(ctx context.Context, tx *gorm.DB, userId string) ([]domain.SocialMediaLink, error)
	FindByTypeAndUserID(ctx context.Context, tx *gorm.DB, typeId uint, userId string) (domain.SocialMediaLink, error)
}

type SocialMediaInteractionRepository interface {
	Create(ctx context.Context, tx *gorm.DB, socialMediaInteraction domain.SocialMediaInteraction) error
	FindBySocialMediaLinkIDAndDate(ctx context.Context, tx *gorm.DB, socialMediaLinkID uint, date time.Time) ([]domain.SocialMediaInteraction, error)
}

type SocialMediaAnalyticRepository interface {
	Create(ctx context.Context, tx *gorm.DB, socialMediaAnalytic domain.SocialMediaAnalytic) (domain.SocialMediaAnalytic, error)
	Update(ctx context.Context, tx *gorm.DB, socialMediaAnalytic domain.SocialMediaAnalytic) (domain.SocialMediaAnalytic, error)
	FindBySocialMediaLinkIDAndDate(ctx context.Context, tx *gorm.DB, socialMediaLinkID uint, date time.Time) (domain.SocialMediaAnalytic, error)
}

type DeviceAnalyticRepository interface {
	Create(ctx context.Context, tx *gorm.DB, deviceAnalytic domain.DeviceAnalytic) (domain.DeviceAnalytic, error)
	Update(ctx context.Context, tx *gorm.DB, deviceAnalytic domain.DeviceAnalytic) (domain.DeviceAnalytic, error)
}
