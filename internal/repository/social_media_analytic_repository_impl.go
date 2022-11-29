package repository

import (
	"context"
	"time"

	"github.com/ilhamfzri/pendek.in/app/logger"
	"github.com/ilhamfzri/pendek.in/internal/model/domain"
	"gorm.io/gorm"
)

type SocialMediaAnalyticRepositoryImpl struct {
	Log *logger.Logger
}

func NewSocialMediaAnalyticRepository(log *logger.Logger) SocialMediaAnalyticRepository {
	return &SocialMediaAnalyticRepositoryImpl{
		Log: log,
	}
}
func (repository *SocialMediaAnalyticRepositoryImpl) Create(ctx context.Context, tx *gorm.DB, socialMediaAnalytic domain.SocialMediaAnalytic) (domain.SocialMediaAnalytic, error) {
	result := tx.WithContext(ctx).Create(&socialMediaAnalytic)
	return socialMediaAnalytic, result.Error
}

func (repository *SocialMediaAnalyticRepositoryImpl) Update(ctx context.Context, tx *gorm.DB, socialMediaAnalytic domain.SocialMediaAnalytic) (domain.SocialMediaAnalytic, error) {
	result := tx.WithContext(ctx).Save(&socialMediaAnalytic)
	return socialMediaAnalytic, result.Error
}

func (repository *SocialMediaAnalyticRepositoryImpl) FindBySocialMediaLinkIDAndDate(ctx context.Context, tx *gorm.DB, socialMediaLinkID uint, date time.Time) (domain.SocialMediaAnalytic, error) {
	var socialMediaAnalytic domain.SocialMediaAnalytic
	result := tx.WithContext(ctx).Preload("DeviceAnalytic").Where("social_media_link_id = ? AND date = ?", socialMediaLinkID, date).First(&socialMediaAnalytic)
	return socialMediaAnalytic, result.Error

}
