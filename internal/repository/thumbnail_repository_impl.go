package repository

import (
	"context"

	"github.com/ilhamfzri/pendek.in/app/logger"
	"github.com/ilhamfzri/pendek.in/internal/model/domain"
	"gorm.io/gorm"
)

type ThumbnailRepositoryImpl struct {
	Log *logger.Logger
}

func NewThumbnailRepository(log *logger.Logger) ThumbnailRepository {
	return &ThumbnailRepositoryImpl{
		Log: log,
	}
}

func (repository *ThumbnailRepositoryImpl) Create(ctx context.Context, tx *gorm.DB, thumbnail domain.Thumbnail) (domain.Thumbnail, error) {
	result := tx.WithContext(ctx).Create(&thumbnail)
	return thumbnail, result.Error
}

func (repository *ThumbnailRepositoryImpl) FindByName(ctx context.Context, tx *gorm.DB, name string) (domain.Thumbnail, error) {
	var thumbnail domain.Thumbnail
	result := tx.WithContext(ctx).Model(domain.Thumbnail{}).Where("name = ?", name).First(&thumbnail)
	return thumbnail, result.Error
}

func (repository *ThumbnailRepositoryImpl) FindByID(ctx context.Context, tx *gorm.DB, id int) (domain.Thumbnail, error) {
	var thumbnail domain.Thumbnail
	result := tx.WithContext(ctx).Model(domain.Thumbnail{}).Where("id = ?", id).First(&thumbnail)
	return thumbnail, result.Error
}

func (repository *ThumbnailRepositoryImpl) FetchAll(ctx context.Context, tx *gorm.DB) ([]domain.Thumbnail, error) {
	var thumbnails []domain.Thumbnail
	result := tx.WithContext(ctx).Model(domain.Thumbnail{}).Find(&thumbnails)
	return thumbnails, result.Error
}
