package repository

import (
	"context"
	"time"

	"github.com/ilhamfzri/pendek.in/app/logger"
	"github.com/ilhamfzri/pendek.in/internal/model/domain"
	"gorm.io/gorm"
)

type SocialMediaInteractionImpl struct {
	Log *logger.Logger
}

func NewSocialMediaInteractionRepository(log *logger.Logger) SocialMediaInteractionRepository {
	return &SocialMediaInteractionImpl{
		Log: log,
	}
}

func (repository *SocialMediaInteractionImpl) Create(ctx context.Context, tx *gorm.DB, socialMediaInteraction domain.SocialMediaInteraction) error {
	result := tx.WithContext(ctx).Create(&socialMediaInteraction)
	return result.Error
}

func (repository *SocialMediaInteractionImpl) FindBySocialMediaLinkIDAndDate(ctx context.Context, tx *gorm.DB, socialMediaLinkID uint, date time.Time) ([]domain.SocialMediaInteraction, error) {
	var socialMediaInteractions []domain.SocialMediaInteraction
	dateFirstRange := date
	dateEndRange := date.Add(1 * time.Second * 60 * 60 * 24)
	result := tx.WithContext(ctx).Where("social_media_link_id = ?", socialMediaLinkID).
		Where("created_at BETWEEN ? AND ?", dateFirstRange, dateEndRange).Find(&socialMediaInteractions)
	return socialMediaInteractions, result.Error
}
