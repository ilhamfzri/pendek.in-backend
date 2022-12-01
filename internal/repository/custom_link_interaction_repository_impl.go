package repository

import (
	"context"
	"time"

	"github.com/ilhamfzri/pendek.in/app/logger"
	"github.com/ilhamfzri/pendek.in/internal/model/domain"
	"gorm.io/gorm"
)

type CustomLinkInteractionRepositoryImpl struct {
	Logger *logger.Logger
}

func NewCustomLinkInteractionRepository(logger *logger.Logger) CustomLinkInteractionRepository {
	return &CustomLinkInteractionRepositoryImpl{
		Logger: logger,
	}
}

func (repository *CustomLinkInteractionRepositoryImpl) Create(ctx context.Context, tx *gorm.DB, linkInteraction domain.CustomLinkInteraction) error {
	result := tx.WithContext(ctx).Create(&linkInteraction)
	return result.Error
}

func (repository *CustomLinkInteractionRepositoryImpl) FindByLinkIdAndDate(ctx context.Context, tx *gorm.DB, linkId int, date time.Time) ([]domain.CustomLinkInteraction, error) {
	var customLinkInteractions []domain.CustomLinkInteraction
	dateFirstRange := date
	dateEndRange := date.Add(1 * time.Second * 60 * 60 * 24)
	result := tx.WithContext(ctx).Where("custom_link_id = ?", linkId).
		Where("created_at BETWEEN ? AND ?", dateFirstRange, dateEndRange).Find(&customLinkInteractions)
	return customLinkInteractions, result.Error
}
