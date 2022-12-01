package service

import (
	"context"
	"errors"
	"time"

	"github.com/ilhamfzri/pendek.in/app/logger"
	"github.com/ilhamfzri/pendek.in/helper"
	"github.com/ilhamfzri/pendek.in/internal/model/domain"
	"github.com/ilhamfzri/pendek.in/internal/model/web"
	"github.com/ilhamfzri/pendek.in/internal/repository"
	"gorm.io/gorm"
)

type CustomLinkAnalyticServiceImpl struct {
	CustomLinkRepository            repository.CustomLinkRepository
	CustomLinkAnalyticRepository    repository.CustomLinkAnalyticRepository
	CustomLinkInteractionRepository repository.CustomLinkInteractionRepository
	DeviceAnalyticRepository        repository.DeviceAnalyticRepository
	DB                              *gorm.DB
	Logger                          *logger.Logger
	Jwt                             helper.IJwt
}

var (
	IntervalDurationUpdateCustomLinkAnalytic = 1 * time.Hour // today analytic threshold duration, please set this to 1s in dev mode
	ErrCustomLinkAnalyticService             = "[Custom Link Analytic Service] Failed to execute"
	ErrCustomLinkAnalyticInvalidEndDate      = errors.New("end date value atleast today, not in the future")
	ErrCustomLinkAnalyticInvalidStartDate    = errors.New("start date format invaled, start date up to last 30 days")
)

func NewCustomLinkAnalyticService(
	customLinkRepo repository.CustomLinkRepository,
	customLinkAnalyticRepo repository.CustomLinkAnalyticRepository,
	customLinkInteractionRepo repository.CustomLinkInteractionRepository,
	deviceAnalyticRepo repository.DeviceAnalyticRepository,
	db *gorm.DB,
	logger *logger.Logger,
	jwt helper.IJwt) CustomLinkAnalyticService {

	return &CustomLinkAnalyticServiceImpl{
		CustomLinkRepository:            customLinkRepo,
		CustomLinkAnalyticRepository:    customLinkAnalyticRepo,
		CustomLinkInteractionRepository: customLinkInteractionRepo,
		DeviceAnalyticRepository:        deviceAnalyticRepo,
		DB:                              db,
		Logger:                          logger,
		Jwt:                             jwt,
	}

}

func (service *CustomLinkAnalyticServiceImpl) SaveInteraction(ctx context.Context, request web.CustomLinkAnalyticInteractionRequest) error {
	// It's a transaction.
	tx := service.DB.Begin()
	defer helper.CommitOrRollback(tx)

	customLinkInteractionDomain := domain.CustomLinkInteraction{
		ClientIP:     request.ClientIP,
		UserAgent:    request.UserAgent,
		CustomLinkID: request.CustomLinkID,
	}

	repoErr := service.CustomLinkInteractionRepository.Create(ctx, tx, customLinkInteractionDomain)
	return repoErr
}

func (service *CustomLinkAnalyticServiceImpl) GetLinkAnalytic(ctx context.Context, request web.CustomLinkAnalyticGetRequest, jwtToken string) ([]web.CustomLinkAnalyticResponse, error) {
	// It's getting the claims from the token.
	claims := service.Jwt.GetClaims(jwtToken)

	// It's a transaction.
	tx := service.DB.Begin()
	defer helper.CommitOrRollback(tx)

	startDate := helper.ToDate(request.StartDate)
	endDate := helper.ToDate(request.EndDate)

	if helper.IsFutureDate(endDate) {
		return []web.CustomLinkAnalyticResponse{}, ErrCustomLinkAnalyticInvalidEndDate
	}

	if !helper.IsLast30Days(startDate) {
		return []web.CustomLinkAnalyticResponse{}, ErrCustomLinkAnalyticInvalidStartDate
	}

	customLink, errRepo := service.CustomLinkRepository.FindByIdAndUserID(ctx, tx, request.LinkID, claims.Id)
	if errRepo != nil && errors.Is(errRepo, gorm.ErrRecordNotFound) {
		return []web.CustomLinkAnalyticResponse{}, ErrCustomLinkNotRegistered
	}
	service.Logger.PanicIfErr(errRepo, ErrCustomLinkAnalyticService)

	var costumLinkAnalyticResponses []web.CustomLinkAnalyticResponse

	for requestDate := startDate; !requestDate.After(endDate); requestDate = requestDate.AddDate(0, 0, 1) {
		requestDate = helper.ToDate(requestDate)

		customLinkAnalytic, errAnalyticRepo := service.CustomLinkAnalyticRepository.FindByLinkIDAndDate(ctx, tx, customLink.ID, requestDate)
		if errAnalyticRepo != nil && !errors.Is(errAnalyticRepo, gorm.ErrRecordNotFound) {
			service.Logger.PanicIfErr(errAnalyticRepo, ErrCustomLinkAnalyticService)
		}

		if errors.Is(errAnalyticRepo, gorm.ErrRecordNotFound) {

			customLinkInteractions, errInteractionRepo := service.CustomLinkInteractionRepository.
				FindByLinkIdAndDate(ctx, tx, int(customLink.ID), requestDate)

			if errInteractionRepo != nil && !errors.Is(errInteractionRepo, gorm.ErrRecordNotFound) {
				service.Logger.PanicIfErr(errInteractionRepo, ErrCustomLinkAnalyticService)
			}

			deviceAnalytic := helper.CustomLinkInteractionsToDeviceAnalytic(&customLinkInteractions)
			clickCount := len(customLinkInteractions)

			deviceAnalytic, errDeviceAnalyticRepo := service.DeviceAnalyticRepository.Create(ctx, tx, deviceAnalytic)
			service.Logger.PanicIfErr(errDeviceAnalyticRepo, ErrCustomLinkAnalyticService)

			customLinkAnalytic = domain.CustomLinkAnalytic{
				ClickCount:       clickCount,
				CustomLinkID:     customLink.ID,
				DeviceAnalyticID: deviceAnalytic.ID,
				Date:             requestDate,
			}

			customLinkAnalytic, errAnalyticRepo = service.CustomLinkAnalyticRepository.Create(ctx, tx, customLinkAnalytic)
			service.Logger.PanicIfErr(errAnalyticRepo, ErrCustomLinkAnalyticService)

			customLinkAnalytic.DeviceAnalytic = deviceAnalytic
			customLinkAnalyticResponse := helper.CustomLinkAnalyticDomainToResponse(&customLinkAnalytic)
			costumLinkAnalyticResponses = append(costumLinkAnalyticResponses, customLinkAnalyticResponse)

			continue

		} else {
			endDate := requestDate.Add(time.Hour * 24)
			lastUpdate := customLinkAnalytic.UpdatedAt
			isNeedUpdate := helper.IsNeedUpdate(lastUpdate, IntervalDurationUpdateCustomLinkAnalytic)
			isToday := helper.IsToday(requestDate)

			if (lastUpdate.Before(endDate) && !isToday) || (lastUpdate.Before(endDate) && isToday && isNeedUpdate) {
				customLinkInteractions, errInteractionRepo := service.CustomLinkInteractionRepository.
					FindByLinkIdAndDate(ctx, tx, int(customLink.ID), requestDate)

				if errInteractionRepo != nil && !errors.Is(errInteractionRepo, gorm.ErrRecordNotFound) {
					service.Logger.PanicIfErr(errInteractionRepo, ErrCustomLinkAnalyticService)
				}

				deviceAnalytic := helper.CustomLinkInteractionsToDeviceAnalytic(&customLinkInteractions)
				clickCount := len(customLinkInteractions)

				deviceAnalytic.ID = customLinkAnalytic.DeviceAnalyticID
				deviceAnalytic, errDeviceAnalyticRepo := service.DeviceAnalyticRepository.Update(ctx, tx, deviceAnalytic)
				service.Logger.PanicIfErr(errDeviceAnalyticRepo, ErrCustomLinkAnalyticService)

				customLinkAnalytic.ClickCount = clickCount
				customLinkAnalytic, errAnalyticRepo = service.CustomLinkAnalyticRepository.Update(ctx, tx, customLinkAnalytic)
				service.Logger.PanicIfErr(errAnalyticRepo, ErrCustomLinkAnalyticService)

				customLinkAnalytic.DeviceAnalytic = deviceAnalytic
			}

			customLinkAnalyticResponse := helper.CustomLinkAnalyticDomainToResponse(&customLinkAnalytic)
			costumLinkAnalyticResponses = append(costumLinkAnalyticResponses, customLinkAnalyticResponse)

			continue
		}
	}
	return costumLinkAnalyticResponses, nil
}

func (service *CustomLinkAnalyticServiceImpl) GetSummaryLinkAnalytic(ctx context.Context, jwtToken string) (web.CustomLinkAnalyticSummaryResponse, error) {
	// It's getting the claims from the token.
	claims := service.Jwt.GetClaims(jwtToken)

	// It's a transaction.
	tx := service.DB.Begin()
	defer helper.CommitOrRollback(tx)

	customLinks, errRepo := service.CustomLinkRepository.FetchAllByUserID(ctx, tx, claims.Id)
	if errRepo != nil && errors.Is(errRepo, gorm.ErrRecordNotFound) {
		return web.CustomLinkAnalyticSummaryResponse{}, nil
	}
	service.Logger.PanicIfErr(errRepo, ErrCustomLinkAnalyticService)

	startDate := helper.GetLast30Days()
	endDate := helper.ToDate(time.Now())

	deviceAnalyticTotal := domain.DeviceAnalytic{}
	totalCustomLinkResponses := []web.TotalCustomLinkAnalyticResponse{}

	for _, customLink := range customLinks {
		totalCustomLinkResponse := web.TotalCustomLinkAnalyticResponse{}
		totalCustomLinkResponse.LinkID = int(customLink.ID)

		for requestDate := startDate; !requestDate.After(endDate); requestDate = requestDate.AddDate(0, 0, 1) {
			requestDate = helper.ToDate(requestDate)

			customLinkAnalytic, errAnalyticRepo := service.CustomLinkAnalyticRepository.FindByLinkIDAndDate(ctx, tx, customLink.ID, requestDate)
			if errAnalyticRepo != nil && !errors.Is(errAnalyticRepo, gorm.ErrRecordNotFound) {
				service.Logger.PanicIfErr(errAnalyticRepo, ErrCustomLinkAnalyticService)
			}

			if errors.Is(errAnalyticRepo, gorm.ErrRecordNotFound) {

				customLinkInteractions, errInteractionRepo := service.CustomLinkInteractionRepository.
					FindByLinkIdAndDate(ctx, tx, int(customLink.ID), requestDate)

				if errInteractionRepo != nil && !errors.Is(errInteractionRepo, gorm.ErrRecordNotFound) {
					service.Logger.PanicIfErr(errInteractionRepo, ErrCustomLinkAnalyticService)
				}

				deviceAnalytic := helper.CustomLinkInteractionsToDeviceAnalytic(&customLinkInteractions)
				clickCount := len(customLinkInteractions)

				deviceAnalytic, errDeviceAnalyticRepo := service.DeviceAnalyticRepository.Create(ctx, tx, deviceAnalytic)
				service.Logger.PanicIfErr(errDeviceAnalyticRepo, ErrCustomLinkAnalyticService)

				customLinkAnalytic = domain.CustomLinkAnalytic{
					ClickCount:       clickCount,
					CustomLinkID:     customLink.ID,
					DeviceAnalyticID: deviceAnalytic.ID,
					Date:             requestDate,
				}

				customLinkAnalytic.DeviceAnalytic = deviceAnalytic
				customLinkAnalytic, errAnalyticRepo = service.CustomLinkAnalyticRepository.Create(ctx, tx, customLinkAnalytic)
				service.Logger.PanicIfErr(errAnalyticRepo, ErrCustomLinkAnalyticService)

				totalCustomLinkResponse.TotalClickCount += customLinkAnalytic.ClickCount
				totalCustomLinkResponse.TotalViewCount += customLinkAnalytic.ViewCount

				deviceAnalyticTotal.Desktop += customLinkAnalytic.DeviceAnalytic.Desktop
				deviceAnalyticTotal.Tablet += customLinkAnalytic.DeviceAnalytic.Tablet
				deviceAnalyticTotal.Mobile += customLinkAnalytic.DeviceAnalytic.Mobile
				deviceAnalyticTotal.Other += customLinkAnalytic.DeviceAnalytic.Other

				continue

			} else {
				endDate := requestDate.Add(time.Hour * 24)
				lastUpdate := customLinkAnalytic.UpdatedAt
				isNeedUpdate := helper.IsNeedUpdate(lastUpdate, IntervalDurationUpdateCustomLinkAnalytic)
				isToday := helper.IsToday(requestDate)

				if (lastUpdate.Before(endDate) && !isToday) || (lastUpdate.Before(endDate) && isToday && isNeedUpdate) {
					customLinkInteractions, errInteractionRepo := service.CustomLinkInteractionRepository.
						FindByLinkIdAndDate(ctx, tx, int(customLink.ID), requestDate)

					if errInteractionRepo != nil && !errors.Is(errInteractionRepo, gorm.ErrRecordNotFound) {
						service.Logger.PanicIfErr(errInteractionRepo, ErrCustomLinkAnalyticService)
					}

					deviceAnalytic := helper.CustomLinkInteractionsToDeviceAnalytic(&customLinkInteractions)
					clickCount := len(customLinkInteractions)

					deviceAnalytic.ID = customLinkAnalytic.DeviceAnalyticID
					deviceAnalytic, errDeviceAnalyticRepo := service.DeviceAnalyticRepository.Update(ctx, tx, deviceAnalytic)
					service.Logger.PanicIfErr(errDeviceAnalyticRepo, ErrCustomLinkAnalyticService)

					customLinkAnalytic.ClickCount = clickCount
					customLinkAnalytic, errAnalyticRepo = service.CustomLinkAnalyticRepository.Update(ctx, tx, customLinkAnalytic)
					service.Logger.PanicIfErr(errAnalyticRepo, ErrCustomLinkAnalyticService)

					customLinkAnalytic.DeviceAnalytic = deviceAnalytic
				}

				totalCustomLinkResponse.TotalClickCount += customLinkAnalytic.ClickCount
				totalCustomLinkResponse.TotalViewCount += customLinkAnalytic.ViewCount

				deviceAnalyticTotal.Desktop += customLinkAnalytic.DeviceAnalytic.Desktop
				deviceAnalyticTotal.Tablet += customLinkAnalytic.DeviceAnalytic.Tablet
				deviceAnalyticTotal.Mobile += customLinkAnalytic.DeviceAnalytic.Mobile
				deviceAnalyticTotal.Other += customLinkAnalytic.DeviceAnalytic.Other

				continue
			}
		}

		totalCustomLinkResponses = append(totalCustomLinkResponses, totalCustomLinkResponse)
	}

	deviceAnalyticResponse := helper.DeviceAnalyticDomainToResponse(&deviceAnalyticTotal)
	customLinkSummaryResponse := web.CustomLinkAnalyticSummaryResponse{
		CustomLink:     totalCustomLinkResponses,
		DeviceAnalytic: deviceAnalyticResponse,
		LastUpdated:    time.Now(),
	}
	return customLinkSummaryResponse, nil
}
