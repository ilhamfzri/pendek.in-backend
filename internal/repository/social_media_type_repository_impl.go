package repository

import (
	"context"

	"github.com/ilhamfzri/pendek.in/app/logger"
	"github.com/ilhamfzri/pendek.in/internal/model/domain"
	"gorm.io/gorm"
)

type SocialMediaTypeRepositoryImpl struct {
	Log *logger.Logger
}

func NewSocialMediaTypeRepository(log *logger.Logger) SocialMediaTypeRepository {
	return &SocialMediaTypeRepositoryImpl{
		Log: log,
	}
}

func (repository *SocialMediaTypeRepositoryImpl) Create(ctx context.Context, tx *gorm.DB, socialMediaType domain.SocialMediaType) (domain.SocialMediaType, error) {
	result := tx.WithContext(ctx).Create(&socialMediaType)
	return socialMediaType, result.Error
}

func (repository *SocialMediaTypeRepositoryImpl) FindByID(ctx context.Context, tx *gorm.DB, id int) (domain.SocialMediaType, error) {
	var socialMediaType domain.SocialMediaType
	result := tx.WithContext(ctx).Where("id = ?", id).First(&socialMediaType)
	return socialMediaType, result.Error
}

func (repository *SocialMediaTypeRepositoryImpl) FindByName(ctx context.Context, tx *gorm.DB, name string) (domain.SocialMediaType, error) {
	var socialMediaType domain.SocialMediaType
	result := tx.WithContext(ctx).Where("name = ?", name).First(&socialMediaType)
	return socialMediaType, result.Error
}

func (repository *SocialMediaTypeRepositoryImpl) FetchAll(ctx context.Context, tx *gorm.DB) ([]domain.SocialMediaType, error) {
	var socialMediaTypes []domain.SocialMediaType
	result := tx.WithContext(ctx).Find(&socialMediaTypes)
	return socialMediaTypes, result.Error
}
