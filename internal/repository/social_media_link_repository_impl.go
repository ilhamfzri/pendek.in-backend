package repository

import (
	"context"
	"fmt"

	"github.com/ilhamfzri/pendek.in/app/logger"
	"github.com/ilhamfzri/pendek.in/internal/model/domain"
	"gorm.io/gorm"
)

type SocialMediaLinkRepositoryImpl struct {
	Log *logger.Logger
}

func NewSocialMediaLinkRepository(log *logger.Logger) SocialMediaLinkRepository {
	return &SocialMediaLinkRepositoryImpl{
		Log: log,
	}
}

func (repository *SocialMediaLinkRepositoryImpl) Create(ctx context.Context, tx *gorm.DB, socialMediaLink domain.SocialMediaLink) (domain.SocialMediaLink, error) {
	result := tx.WithContext(ctx).Create(&socialMediaLink)
	return socialMediaLink, result.Error
}

func (repository *SocialMediaLinkRepositoryImpl) Update(ctx context.Context, tx *gorm.DB, socialMediaLink domain.SocialMediaLink) (domain.SocialMediaLink, error) {
	result := tx.WithContext(ctx).Model(socialMediaLink).
		Updates(
			map[string]interface{}{
				"link_or_username": socialMediaLink.LinkOrUsername,
				"activate":         socialMediaLink.Activate,
			})
	fmt.Println(socialMediaLink.Activate)

	return socialMediaLink, result.Error
}

func (repository *SocialMediaLinkRepositoryImpl) FindByUserID(ctx context.Context, tx *gorm.DB, userId string) ([]domain.SocialMediaLink, error) {
	var socialMediaLinks []domain.SocialMediaLink
	result := tx.WithContext(ctx).Preload("SocialMediaType").Order("type_id ASC").Find(&socialMediaLinks, "social_media_links.user_id = ?", userId)
	return socialMediaLinks, result.Error
}

func (repository *SocialMediaLinkRepositoryImpl) FindByTypeAndUserID(ctx context.Context, tx *gorm.DB, typeId uint, userId string) (domain.SocialMediaLink, error) {
	var socialMediaLink domain.SocialMediaLink
	result := tx.WithContext(ctx).Where("user_id = ? AND type_id = ?", userId, typeId).First(&socialMediaLink)
	return socialMediaLink, result.Error
}
