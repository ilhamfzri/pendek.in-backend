package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ilhamfzri/pendek.in/app/logger"
	"github.com/ilhamfzri/pendek.in/helper"
	"github.com/ilhamfzri/pendek.in/internal/model/domain"
	"github.com/ilhamfzri/pendek.in/internal/model/web"
	"github.com/ilhamfzri/pendek.in/internal/repository"
	"gorm.io/gorm"
)

type SocialMediaAnalyticServiceImpl struct {
	UserRepository                   repository.UserRepository
	SocialMediaLinkRepository        repository.SocialMediaLinkRepository
	SocialMediaInteractionRepository repository.SocialMediaInteractionRepository
	SocialMediaAnalyticRepository    repository.SocialMediaAnalyticRepository
	DeviceAnalyticRepository         repository.DeviceAnalyticRepository
	DB                               *gorm.DB
	Logger                           *logger.Logger
	Jwt                              helper.IJwt
}

var (
	IntervalDurationUpdateAnalytic         = 1 * time.Hour // today analytic threshold duration
	ErrSocialMediaAnalyticService          = "[Social Media Analytic Service] Failed Execute Social Media Analytic Service"
	ErrSocialMediaAnalyticInvalidEndDate   = errors.New("end date value atleast today, not in the future")
	ErrSocialMediaAnalyticInvalidStartDate = errors.New("start date format invaled, start date up to last 30 days")
)

func NewSocialMediaAnalyticService(userRepository repository.UserRepository,
	socialMediaLinkRepository repository.SocialMediaLinkRepository,
	socialMediaInteractionRepository repository.SocialMediaInteractionRepository,
	socialMediaAnalyticRepository repository.SocialMediaAnalyticRepository,
	deviceAnalyticRepository repository.DeviceAnalyticRepository,
	DB *gorm.DB, logger *logger.Logger, jwt helper.IJwt) SocialMediaAnalytic {
	return &SocialMediaAnalyticServiceImpl{
		UserRepository:                   userRepository,
		SocialMediaLinkRepository:        socialMediaLinkRepository,
		SocialMediaInteractionRepository: socialMediaInteractionRepository,
		SocialMediaAnalyticRepository:    socialMediaAnalyticRepository,
		DeviceAnalyticRepository:         deviceAnalyticRepository,
		DB:                               DB,
		Logger:                           logger,
		Jwt:                              jwt,
	}
}

func (service *SocialMediaAnalyticServiceImpl) SaveInteraction(ctx context.Context, request web.SocialMediaAnalyticInteractionRequest) error {
	// It's a transaction.
	tx := service.DB.Begin()
	defer helper.CommitOrRollback(tx)

	socialMediaInteractionDomain := domain.SocialMediaInteraction{
		ClientIP:          request.ClientIP,
		UserAgent:         request.UserAgent,
		SocialMediaLinkID: request.SocialMediaLinkID,
	}
	repoErr := service.SocialMediaInteractionRepository.Create(ctx, tx, socialMediaInteractionDomain)
	return repoErr
}

func (service *SocialMediaAnalyticServiceImpl) GetLinkAnalytic(ctx context.Context, request web.SocialMediaAnalyticGetRequest, jwtToken string) ([]web.SocialMediaAnalyticResponse, error) {
	// It's getting the claims from the token.
	claims := service.Jwt.GetClaims(jwtToken)

	// It's a transaction.
	tx := service.DB.Begin()
	defer helper.CommitOrRollback(tx)

	startDate := helper.ToDate(request.StartDate)
	endDate := helper.ToDate(request.EndDate)

	if helper.IsFutureDate(endDate) {
		return []web.SocialMediaAnalyticResponse{}, ErrSocialMediaAnalyticInvalidEndDate
	}

	if !helper.IsLast30Days(startDate) {
		return []web.SocialMediaAnalyticResponse{}, ErrSocialMediaAnalyticInvalidStartDate
	}

	socialMediaLink, repoErr := service.SocialMediaLinkRepository.FindByTypeAndUserID(ctx, tx, uint(request.TypeID), claims.Id)
	if repoErr != nil && errors.Is(repoErr, gorm.ErrRecordNotFound) {
		return []web.SocialMediaAnalyticResponse{}, ErrSocialMediaLinkNotFound
	}
	service.Logger.PanicIfErr(repoErr, ErrSocialMediaAnalyticService)

	var socialMediaAnalyticResponses []web.SocialMediaAnalyticResponse

	for requestDate := startDate; !requestDate.After(endDate); requestDate = requestDate.AddDate(0, 0, 1) {
		requestDate = helper.ToDate(requestDate)

		socialMediaAnalytic, repoAnalyticErr := service.SocialMediaAnalyticRepository.FindBySocialMediaLinkIDAndDate(ctx, tx, socialMediaLink.ID, requestDate)
		if repoAnalyticErr != nil && !errors.Is(repoAnalyticErr, gorm.ErrRecordNotFound) {
			service.Logger.PanicIfErr(repoAnalyticErr, ErrSocialMediaAnalyticService)
		}

		if errors.Is(repoAnalyticErr, gorm.ErrRecordNotFound) {
			socialMediaInteractions, repoInteractionErr := service.SocialMediaInteractionRepository.
				FindBySocialMediaLinkIDAndDate(ctx, tx, socialMediaLink.ID, requestDate)

			if repoInteractionErr != nil && !errors.Is(repoErr, gorm.ErrRecordNotFound) {
				service.Logger.PanicIfErr(repoInteractionErr, ErrSocialMediaAnalyticService)
			}

			deviceAnalytic := helper.SocialMediaInteractionsToDeviceAnalytic(&socialMediaInteractions)
			clickCount := len(socialMediaInteractions)

			deviceAnalytic, repoDeviceAnalyticErr := service.DeviceAnalyticRepository.Create(ctx, tx, deviceAnalytic)
			service.Logger.PanicIfErr(repoDeviceAnalyticErr, ErrSocialMediaAnalyticService)

			socialMediaAnalytic = domain.SocialMediaAnalytic{
				ClickCount:        clickCount,
				SocialMediaLinkID: socialMediaLink.ID,
				DeviceAnalyticID:  deviceAnalytic.ID,
				Date:              requestDate,
			}

			socialMediaAnalytic, repoAnalyticErr = service.SocialMediaAnalyticRepository.Create(ctx, tx, socialMediaAnalytic)
			service.Logger.PanicIfErr(repoAnalyticErr, ErrSocialMediaAnalyticService)

			socialMediaAnalytic.DeviceAnalytic = deviceAnalytic
			socialMediaName := socialMediaLink.SocialMediaType.Name
			socialMediaAnalyticResponse := helper.SocialMediaAnalyticDomainToResponse(&socialMediaAnalytic, socialMediaName)

			socialMediaAnalyticResponses = append(socialMediaAnalyticResponses, socialMediaAnalyticResponse)
			continue

		} else {
			endDate := requestDate.Add(time.Hour * 24)
			lastUpdate := socialMediaAnalytic.UpdatedAt
			isNeedUpdate := helper.IsNeedUpdate(lastUpdate, IntervalDurationUpdateAnalytic)
			isToday := helper.IsToday(requestDate)

			if (lastUpdate.Before(endDate) && !isToday) || (lastUpdate.Before(endDate) && isToday && isNeedUpdate) {
				socialMediaInteractions, repoInteractionErr := service.SocialMediaInteractionRepository.
					FindBySocialMediaLinkIDAndDate(ctx, tx, socialMediaLink.ID, requestDate)
				service.Logger.PanicIfErr(repoInteractionErr, ErrSocialMediaAnalyticService)

				deviceAnalytic := helper.SocialMediaInteractionsToDeviceAnalytic(&socialMediaInteractions)
				clickCount := len(socialMediaInteractions)

				deviceAnalytic.ID = socialMediaAnalytic.DeviceAnalyticID
				deviceAnalytic, repoDeviceAnalyticErr := service.DeviceAnalyticRepository.Update(ctx, tx, deviceAnalytic)
				service.Logger.PanicIfErr(repoDeviceAnalyticErr, ErrSocialMediaAnalyticService)

				socialMediaAnalytic.ClickCount = clickCount
				socialMediaAnalytic, repoAnalyticErr = service.SocialMediaAnalyticRepository.Update(ctx, tx, socialMediaAnalytic)
				service.Logger.PanicIfErr(repoAnalyticErr, ErrSocialMediaAnalyticService)

				socialMediaAnalytic.DeviceAnalytic = deviceAnalytic
			}

			socialMediaName := socialMediaLink.SocialMediaType.Name
			socialMediaAnalyticResponse := helper.SocialMediaAnalyticDomainToResponse(&socialMediaAnalytic, socialMediaName)
			socialMediaAnalyticResponses = append(socialMediaAnalyticResponses, socialMediaAnalyticResponse)
			continue
		}
	}

	return socialMediaAnalyticResponses, nil
}

func (service *SocialMediaAnalyticServiceImpl) GetSummaryLinkAnalytic(ctx context.Context, jwtToken string) (web.SocialMediaAnalyticSummaryResponse, error) {
	// It's getting the claims from the token.
	claims := service.Jwt.GetClaims(jwtToken)

	// It's a transaction.
	tx := service.DB.Begin()
	defer helper.CommitOrRollback(tx)

	socialMediaLinks, repoErr := service.SocialMediaLinkRepository.FindByUserID(ctx, tx, claims.Id)
	if repoErr != nil && errors.Is(repoErr, gorm.ErrRecordNotFound) {
		return web.SocialMediaAnalyticSummaryResponse{}, ErrSocialMediaLinkNotFound
	}
	service.Logger.PanicIfErr(repoErr, ErrSocialMediaAnalyticService)

	startDate := helper.GetLast30Days()
	endDate := helper.ToDate(time.Now())

	deviceAnalyticTotal := domain.DeviceAnalytic{}
	totalSocialMediaResponses := []web.TotalSocialMediaAnalyticResponse{}

	for _, socialMediaLink := range socialMediaLinks {
		totalSocialMediaResponse := web.TotalSocialMediaAnalyticResponse{}
		totalSocialMediaResponse.SocialMediaName = socialMediaLink.SocialMediaType.Name

		for requestDate := startDate; !requestDate.After(endDate); requestDate = requestDate.AddDate(0, 0, 1) {
			requestDate = helper.ToDate(requestDate)

			fmt.Println(requestDate, socialMediaLink.SocialMediaType.Name)

			socialMediaAnalytic, repoAnalyticErr := service.SocialMediaAnalyticRepository.FindBySocialMediaLinkIDAndDate(ctx, tx, socialMediaLink.ID, requestDate)
			if repoAnalyticErr != nil && !errors.Is(repoAnalyticErr, gorm.ErrRecordNotFound) {
				service.Logger.PanicIfErr(repoAnalyticErr, ErrSocialMediaAnalyticService)
			}

			if errors.Is(repoAnalyticErr, gorm.ErrRecordNotFound) {
				socialMediaInteractions, repoInteractionErr := service.SocialMediaInteractionRepository.
					FindBySocialMediaLinkIDAndDate(ctx, tx, socialMediaLink.ID, requestDate)

				if repoInteractionErr != nil && !errors.Is(repoErr, gorm.ErrRecordNotFound) {
					service.Logger.PanicIfErr(repoInteractionErr, ErrSocialMediaAnalyticService)
				}

				deviceAnalytic := helper.SocialMediaInteractionsToDeviceAnalytic(&socialMediaInteractions)
				clickCount := len(socialMediaInteractions)

				deviceAnalytic, repoDeviceAnalyticErr := service.DeviceAnalyticRepository.Create(ctx, tx, deviceAnalytic)
				service.Logger.PanicIfErr(repoDeviceAnalyticErr, ErrSocialMediaAnalyticService)

				socialMediaAnalytic = domain.SocialMediaAnalytic{
					ClickCount:        clickCount,
					SocialMediaLinkID: socialMediaLink.ID,
					DeviceAnalyticID:  deviceAnalytic.ID,
					Date:              requestDate,
				}

				socialMediaAnalytic, repoAnalyticErr = service.SocialMediaAnalyticRepository.Create(ctx, tx, socialMediaAnalytic)
				service.Logger.PanicIfErr(repoAnalyticErr, ErrSocialMediaAnalyticService)

				socialMediaAnalytic.DeviceAnalytic = deviceAnalytic

				totalSocialMediaResponse.TotalClickCount += socialMediaAnalytic.ClickCount
				totalSocialMediaResponse.TotalViewCount += socialMediaAnalytic.ViewCount

				deviceAnalyticTotal.Desktop += socialMediaAnalytic.DeviceAnalytic.Desktop
				deviceAnalyticTotal.Mobile += socialMediaAnalytic.DeviceAnalytic.Mobile
				deviceAnalyticTotal.Tablet += socialMediaAnalytic.DeviceAnalytic.Tablet
				deviceAnalyticTotal.Other += socialMediaAnalytic.DeviceAnalytic.Other
				continue

			} else {
				endDate := requestDate.Add(time.Hour * 24)
				lastUpdate := socialMediaAnalytic.UpdatedAt
				isNeedUpdate := helper.IsNeedUpdate(lastUpdate, IntervalDurationUpdateAnalytic)
				isToday := helper.IsToday(requestDate)

				if (lastUpdate.Before(endDate) && !isToday) || (lastUpdate.Before(endDate) && isToday && isNeedUpdate) {
					socialMediaInteractions, repoInteractionErr := service.SocialMediaInteractionRepository.
						FindBySocialMediaLinkIDAndDate(ctx, tx, socialMediaLink.ID, requestDate)
					service.Logger.PanicIfErr(repoInteractionErr, ErrSocialMediaAnalyticService)

					deviceAnalytic := helper.SocialMediaInteractionsToDeviceAnalytic(&socialMediaInteractions)
					clickCount := len(socialMediaInteractions)

					deviceAnalytic.ID = socialMediaAnalytic.DeviceAnalyticID
					deviceAnalytic, repoDeviceAnalyticErr := service.DeviceAnalyticRepository.Update(ctx, tx, deviceAnalytic)
					service.Logger.PanicIfErr(repoDeviceAnalyticErr, ErrSocialMediaAnalyticService)

					socialMediaAnalytic.ClickCount = clickCount
					socialMediaAnalytic, repoAnalyticErr = service.SocialMediaAnalyticRepository.Update(ctx, tx, socialMediaAnalytic)
					service.Logger.PanicIfErr(repoAnalyticErr, ErrSocialMediaAnalyticService)

					socialMediaAnalytic.DeviceAnalytic = deviceAnalytic
				}

				totalSocialMediaResponse.TotalClickCount += socialMediaAnalytic.ClickCount
				totalSocialMediaResponse.TotalViewCount += socialMediaAnalytic.ViewCount

				deviceAnalyticTotal.Desktop += socialMediaAnalytic.DeviceAnalytic.Desktop
				deviceAnalyticTotal.Mobile += socialMediaAnalytic.DeviceAnalytic.Mobile
				deviceAnalyticTotal.Tablet += socialMediaAnalytic.DeviceAnalytic.Tablet
				deviceAnalyticTotal.Other += socialMediaAnalytic.DeviceAnalytic.Other
				continue
			}
		}

		totalSocialMediaResponses = append(totalSocialMediaResponses, totalSocialMediaResponse)
	}

	deviceAnalyticResponse := helper.DeviceAnalyticDomainToResponse(&deviceAnalyticTotal)
	socialMediaAnalyticSummaryResponse := web.SocialMediaAnalyticSummaryResponse{
		SocialMedia:    totalSocialMediaResponses,
		DeviceAnalytic: deviceAnalyticResponse,
		LastUpdated:    time.Now(),
	}
	return socialMediaAnalyticSummaryResponse, nil

}
