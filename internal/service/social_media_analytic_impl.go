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
	ThresholdDurationUpdateAnalytic         = 1 * time.Hour // today analytic threshold duration
	ErrSocialMediaAnalyticService           = "[Social Media Analytic Service] Failed Execute Social Media Analytic Service"
	ErrSocialMediaAnalyticInvalidDateFormat = errors.New("date value atleast today, not in the future")
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

func (service *SocialMediaAnalyticServiceImpl) GetLinkAnalytic(ctx context.Context, request web.SocialMediaAnalyticGetRequest, jwtToken string) (web.SocialMediaAnalyticResponse, error) {
	// It's getting the claims from the token.
	claims := service.Jwt.GetClaims(jwtToken)

	// It's a transaction.
	tx := service.DB.Begin()
	defer helper.CommitOrRollback(tx)

	timeNow := time.Now()
	requestTime := request.Datetime

	todayDate := time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), 0, 1, 0, 0, timeNow.Location())
	requestDate := time.Date(requestTime.Year(), requestTime.Month(), requestTime.Day(), 0, 0, 0, 0, requestTime.Location())

	if requestDate.After(todayDate) {
		return web.SocialMediaAnalyticResponse{}, ErrSocialMediaAnalyticInvalidDateFormat
	}

	socialMediaLink, repoErr := service.SocialMediaLinkRepository.FindByTypeAndUserID(ctx, tx, uint(request.TypeID), claims.Id)
	if repoErr != nil && errors.Is(repoErr, gorm.ErrRecordNotFound) {
		return web.SocialMediaAnalyticResponse{}, ErrSocialMediaLinkNotFound
	}
	service.Logger.PanicIfErr(repoErr, ErrSocialMediaAnalyticService)

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

		return socialMediaAnalyticResponse, nil

	} else {
		endDate := requestDate.Add(time.Hour * 24)
		lastUpdate := socialMediaAnalytic.UpdatedAt
		isNeedUpdate := helper.IsNeedUpdate(lastUpdate, ThresholdDurationUpdateAnalytic)
		isToday := helper.IsToday(requestDate)
		fmt.Println(isToday, isNeedUpdate)

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
		return socialMediaAnalyticResponse, nil
	}
}
