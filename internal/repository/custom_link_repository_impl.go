package repository

import (
	"context"

	"github.com/ilhamfzri/pendek.in/app/logger"
	"github.com/ilhamfzri/pendek.in/internal/model/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CustomLinkRepositoryImpl struct {
	Logger *logger.Logger
}

func NewCustomLinkRepository(logger *logger.Logger) CustomLinkRepository {
	return &CustomLinkRepositoryImpl{
		Logger: logger,
	}
}

func (repository *CustomLinkRepositoryImpl) Create(ctx context.Context, tx *gorm.DB, link domain.CustomLink) (domain.CustomLink, error) {
	result := tx.WithContext(ctx).Create(&link)
	return link, result.Error
}

func (repository *CustomLinkRepositoryImpl) Update(ctx context.Context, tx *gorm.DB, link domain.CustomLink) (domain.CustomLink, error) {
	result := tx.WithContext(ctx).Model(&link).Clauses(clause.Returning{}).
		Select("*").
		Updates(
			map[string]interface{}{
				"title":           link.Title,
				"short_link_code": link.ShortLinkCode,
				"long_link":       link.LongLink,
				"show_on_profile": link.ShowOnProfile,
				"activate":        link.Activate,
			},
		)
	return link, result.Error
}

func (repository *CustomLinkRepositoryImpl) UpdateThumbnailIDFK(ctx context.Context, tx *gorm.DB, linkID uint, thumbnailID *uint) (domain.CustomLink, error) {
	customLink := domain.CustomLink{}
	result := tx.WithContext(ctx).Model(&customLink).Clauses(clause.Returning{}).
		Select("thumbnail_id").Where("id = ?", linkID).
		Updates(
			map[string]interface{}{
				"thumbnail_id": thumbnailID,
			},
		)
	return customLink, result.Error
}

func (repository *CustomLinkRepositoryImpl) UpdateCustomThumbnailIDFK(ctx context.Context, tx *gorm.DB, linkID uint, customThumbnailID *uint) (domain.CustomLink, error) {
	customLink := domain.CustomLink{}
	result := tx.WithContext(ctx).Model(&customLink).Clauses(clause.Returning{}).
		Select("custom_thumbnail_id").Where("id = ?", linkID).
		Updates(
			map[string]interface{}{
				"custom_thumbnail_id": customThumbnailID,
			},
		)
	return customLink, result.Error
}

func (repository *CustomLinkRepositoryImpl) FindByShortLinkCode(ctx context.Context, tx *gorm.DB, shortLinkCode string) (domain.CustomLink, error) {
	var link domain.CustomLink
	result := tx.WithContext(ctx).Where("short_link_code = ?", shortLinkCode).First(&link)
	return link, result.Error
}

func (repository *CustomLinkRepositoryImpl) FindByIdAndUserID(ctx context.Context, tx *gorm.DB, id int, userID string) (domain.CustomLink, error) {
	var link domain.CustomLink
	result := tx.WithContext(ctx).Preload("CustomThumbnail").Preload("Thumbnail").Where("id = ? AND user_id = ?", id, userID).First(&link)
	return link, result.Error
}

func (repository *CustomLinkRepositoryImpl) FetchAllByUserID(ctx context.Context, tx *gorm.DB, userID string) ([]domain.CustomLink, error) {
	var links []domain.CustomLink
	result := tx.WithContext(ctx).Preload("CustomThumbnail").Preload("Thumbnail").Where("user_id = ?", userID).Order("id ASC").Find(&links)
	return links, result.Error
}
