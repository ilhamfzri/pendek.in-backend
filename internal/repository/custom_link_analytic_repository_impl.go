package repository

import (
	"context"
	"time"

	"github.com/ilhamfzri/pendek.in/app/logger"
	"github.com/ilhamfzri/pendek.in/internal/model/domain"
	"gorm.io/gorm"
)

type CustomLinkAnalyticRepositoryImpl struct {
	Logger *logger.Logger
}

func NewCustomLinkAnalyticRepository(logger *logger.Logger) CustomLinkAnalyticRepository {
	return &CustomLinkAnalyticRepositoryImpl{
		Logger: logger,
	}
}

func (repository *CustomLinkAnalyticRepositoryImpl) Create(ctx context.Context, tx *gorm.DB, customLinkAnalytic domain.CustomLinkAnalytic) (domain.CustomLinkAnalytic, error) {
	result := tx.WithContext(ctx).Create(&customLinkAnalytic)
	return customLinkAnalytic, result.Error
}

func (repository *CustomLinkAnalyticRepositoryImpl) Update(ctx context.Context, tx *gorm.DB, customLinkAnalytic domain.CustomLinkAnalytic) (domain.CustomLinkAnalytic, error) {
	result := tx.WithContext(ctx).Save(&customLinkAnalytic)
	return customLinkAnalytic, result.Error
}

func (repository *CustomLinkAnalyticRepositoryImpl) FindByLinkIDAndDate(ctx context.Context, tx *gorm.DB, customLinkID uint, date time.Time) (domain.CustomLinkAnalytic, error) {
	var customLinkAnalytic domain.CustomLinkAnalytic
	result := tx.WithContext(ctx).Preload("DeviceAnalytic").Where("custom_link_id = ? AND date = ?", customLinkID, date).First(&customLinkAnalytic)
	return customLinkAnalytic, result.Error
}
