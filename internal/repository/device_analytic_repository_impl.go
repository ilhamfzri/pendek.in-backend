package repository

import (
	"context"

	"github.com/ilhamfzri/pendek.in/app/logger"
	"github.com/ilhamfzri/pendek.in/internal/model/domain"
	"gorm.io/gorm"
)

type DeviceAnalyticRepositoryImpl struct {
	Log *logger.Logger
}

func NewDeviceAnalyticRepository(log *logger.Logger) DeviceAnalyticRepository {
	return &DeviceAnalyticRepositoryImpl{
		Log: log,
	}
}

func (repository *DeviceAnalyticRepositoryImpl) Create(ctx context.Context, tx *gorm.DB, deviceAnalytic domain.DeviceAnalytic) (domain.DeviceAnalytic, error) {
	result := tx.WithContext(ctx).Create(&deviceAnalytic)
	return deviceAnalytic, result.Error
}

func (repository *DeviceAnalyticRepositoryImpl) Update(ctx context.Context, tx *gorm.DB, deviceAnalytic domain.DeviceAnalytic) (domain.DeviceAnalytic, error) {
	result := tx.WithContext(ctx).Save(&deviceAnalytic)
	return deviceAnalytic, result.Error
}
