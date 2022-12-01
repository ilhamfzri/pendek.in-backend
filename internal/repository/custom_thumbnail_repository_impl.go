package repository

import (
	"context"

	"github.com/ilhamfzri/pendek.in/app/logger"
	"github.com/ilhamfzri/pendek.in/internal/model/domain"
	"gorm.io/gorm"
)

type CustomThumbnailRepositoryImpl struct {
	Logger *logger.Logger
}

func NewCustomThumbnailRepository(logger *logger.Logger) CustomThumbnailRepository {
	return &CustomThumbnailRepositoryImpl{
		Logger: logger,
	}
}

func (repository *CustomThumbnailRepositoryImpl) Create(ctx context.Context, tx *gorm.DB, thumbnail domain.CustomThumbnail) (domain.CustomThumbnail, error) {
	result := tx.WithContext(ctx).Create(&thumbnail)
	return thumbnail, result.Error
}

func (repository *CustomThumbnailRepositoryImpl) FetchAllByUserID(ctx context.Context, tx *gorm.DB, userID string) ([]domain.CustomThumbnail, error) {
	var thumbnails []domain.CustomThumbnail
	result := tx.WithContext(ctx).Where("user_id = ?", userID).Find(&thumbnails)
	return thumbnails, result.Error
}

func (repository *CustomThumbnailRepositoryImpl) FindByThumbnailIDAndUserID(ctx context.Context, tx *gorm.DB, thumbnailID int, userID string) (domain.CustomThumbnail, error) {
	var thumbnail domain.CustomThumbnail
	result := tx.WithContext(ctx).Where("id = ? AND user_id = ?", thumbnailID, userID).First(&thumbnail)
	return thumbnail, result.Error
}
