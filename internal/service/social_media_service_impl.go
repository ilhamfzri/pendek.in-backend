package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/ilhamfzri/pendek.in/app/logger"
	"github.com/ilhamfzri/pendek.in/helper"
	"github.com/ilhamfzri/pendek.in/internal/model/domain"
	"github.com/ilhamfzri/pendek.in/internal/model/web"
	"github.com/ilhamfzri/pendek.in/internal/repository"
	"gorm.io/gorm"
)

type SocialMediaLinkServiceImpl struct {
	UserRepository            repository.UserRepository
	SocialMediaLinkRepository repository.SocialMediaLinkRepository
	SocialMediaTypeRepository repository.SocialMediaTypeRepository
	DB                        *gorm.DB
	Logger                    *logger.Logger
	Jwt                       helper.IJwt
}

var (
	ErrSocialMediaLinkService        = "[Social Media Link Service] Failed Execute Social Media Link Service"
	ErrSocialMediaLinkFound          = errors.New("social media link is registered, please use update instead")
	ErrSocialMediaLinkNotFound       = errors.New("social media link is not registered, please use post instead")
	ErrSocialMediaTypeInvalid        = errors.New("social media type id invalid")
	ErrSocialMediaInvalidLink        = errors.New("invalid link")
	ErrSocialMediaLinkUsernameOrLink = errors.New("username or link format is not valid")
)

func NewSocialMediaLinkService(userRepository repository.UserRepository, socialMediaLinkRepository repository.SocialMediaLinkRepository, socialMediaTypeRepository repository.SocialMediaTypeRepository, DB *gorm.DB, logger *logger.Logger, jwt helper.IJwt) SocialMediaLinkService {
	return &SocialMediaLinkServiceImpl{
		UserRepository:            userRepository,
		SocialMediaLinkRepository: socialMediaLinkRepository,
		SocialMediaTypeRepository: socialMediaTypeRepository,
		DB:                        DB,
		Logger:                    logger,
		Jwt:                       jwt,
	}
}

func (service *SocialMediaLinkServiceImpl) GetAllTypes(ctx context.Context) ([]web.SocialMediaTypeResponse, error) {
	// It's a transaction.
	tx := service.DB.Begin()
	defer helper.CommitOrRollback(tx)

	// It's getting all social media type data from the database.
	socialMediaTypes, repoErr := service.SocialMediaTypeRepository.FetchAll(ctx, tx)
	service.Logger.PanicIfErr(repoErr, ErrSocialMediaLinkService)

	var webResponse []web.SocialMediaTypeResponse
	for _, socialMediaType := range socialMediaTypes {
		entry := web.SocialMediaTypeResponse{
			ID:      socialMediaType.ID,
			Name:    socialMediaType.Name,
			Example: socialMediaType.Example,
			IconUrl: socialMediaType.IconUrl,
		}
		webResponse = append(webResponse, entry)
	}
	return webResponse, nil
}

func (service *SocialMediaLinkServiceImpl) CreateLink(ctx context.Context, request web.SocialMediaLinkCreateRequest, host string, jwtToken string) (web.SocialMediaLinkResponse, error) {
	// It's getting the claims from the token.
	claims := service.Jwt.GetClaims(jwtToken)

	// It's a transaction.
	tx := service.DB.Begin()
	defer helper.CommitOrRollback(tx)

	// It's checking if the social media type id is valid or not.
	socialMediaType, repoErr := service.SocialMediaTypeRepository.FindByID(ctx, tx, request.TypeID)
	if repoErr != nil && errors.Is(repoErr, gorm.ErrRecordNotFound) {
		return web.SocialMediaLinkResponse{}, ErrSocialMediaTypeInvalid
	}
	service.Logger.PanicIfErr(repoErr, ErrSocialMediaLinkService)

	// It's checking if the social media link or username is valid or not.
	if !helper.SocialMediaValidator(socialMediaType.Name, request.LinkOrUsername) {
		return web.SocialMediaLinkResponse{}, ErrSocialMediaLinkUsernameOrLink
	}

	// It's checking if the social media link is registered or not.
	socialMediaLink, repoErr := service.SocialMediaLinkRepository.FindByTypeAndUserID(ctx, tx, uint(request.TypeID), claims.Id)
	if repoErr == nil && socialMediaLink.ID != 0 {
		return web.SocialMediaLinkResponse{}, ErrSocialMediaLinkFound
	}

	if !errors.Is(repoErr, gorm.ErrRecordNotFound) && repoErr != nil {
		service.Logger.PanicIfErr(repoErr, ErrSocialMediaLinkService)
	}

	socialMediaLinkData := domain.SocialMediaLink{
		TypeID:         uint(request.TypeID),
		UserID:         claims.Id,
		LinkOrUsername: request.LinkOrUsername,
		Activate:       true,
	}

	// It's creating a new social media link data.
	socialMediaLink, repoErr = service.SocialMediaLinkRepository.Create(ctx, tx, socialMediaLinkData)
	socialMediaLink.SocialMediaType = socialMediaType
	service.Logger.PanicIfErr(repoErr, ErrSocialMediaLinkService)

	socialMediaLinkResponse := helper.SocialMediaLinkDomainToResponse(&socialMediaLink, host, claims.Username)
	return socialMediaLinkResponse, nil
}

func (service *SocialMediaLinkServiceImpl) UpdateLink(ctx context.Context, request web.SocialMediaLinkUpdateRequest, host string, jwtToken string) (web.SocialMediaLinkResponse, error) {
	// It's getting the claims from the token.
	claims := service.Jwt.GetClaims(jwtToken)

	// It's a transaction.
	tx := service.DB.Begin()
	defer helper.CommitOrRollback(tx)

	// It's checking if the social media type id is valid or not.
	_, repoErr := service.SocialMediaTypeRepository.FindByID(ctx, tx, request.TypeID)
	if repoErr != nil && errors.Is(repoErr, gorm.ErrRecordNotFound) {
		return web.SocialMediaLinkResponse{}, ErrSocialMediaTypeInvalid
	}
	service.Logger.PanicIfErr(repoErr, ErrSocialMediaLinkService)

	// It's checking if the social media link is registered or not.
	socialMediaLink, repoErr := service.SocialMediaLinkRepository.FindByTypeAndUserID(ctx, tx, uint(request.TypeID), claims.Id)
	if repoErr != nil && errors.Is(repoErr, gorm.ErrRecordNotFound) {
		return web.SocialMediaLinkResponse{}, ErrSocialMediaLinkNotFound
	}
	service.Logger.PanicIfErr(repoErr, ErrSocialMediaLinkService)

	if !helper.SocialMediaValidator(socialMediaLink.SocialMediaType.Name, request.NewLinkOrUsername) {
		return web.SocialMediaLinkResponse{}, ErrSocialMediaLinkUsernameOrLink
	}

	// It's checking if the request has a new link or username and activate value. If it has, it will
	// update the value.
	if request.NewLinkOrUsername != "" {
		socialMediaLink.LinkOrUsername = request.NewLinkOrUsername
	}

	if request.Activate != nil {
		socialMediaLink.Activate = *request.Activate
	}

	// It's updating the social media link data.
	fmt.Println(socialMediaLink)
	socialMediaLink, repoErr = service.SocialMediaLinkRepository.Update(ctx, tx, socialMediaLink)
	service.Logger.PanicIfErr(repoErr, ErrSocialMediaLinkService)
	fmt.Println("Here")
	socialMediaLinkResponse := helper.SocialMediaLinkDomainToResponse(&socialMediaLink, host, claims.Username)
	return socialMediaLinkResponse, nil
}

func (service *SocialMediaLinkServiceImpl) GetAllLink(ctx context.Context, host string, jwtToken string) ([]web.SocialMediaLinkResponse, error) {
	// It's getting the claims from the token.
	claims := service.Jwt.GetClaims(jwtToken)

	// It's a transaction.
	tx := service.DB.Begin()
	defer helper.CommitOrRollback(tx)

	// It's getting all social media link data from the database.
	socialMediaLinks, repoErr := service.SocialMediaLinkRepository.FindByUserID(ctx, tx, claims.Id)
	service.Logger.PanicIfErr(repoErr, ErrSocialMediaLinkService)

	var socialMediaLinksReponse []web.SocialMediaLinkResponse
	for _, socialMediaLink := range socialMediaLinks {
		socialMediaLinkReponse := helper.SocialMediaLinkDomainToResponse(&socialMediaLink, host, claims.Username)
		socialMediaLinksReponse = append(socialMediaLinksReponse, socialMediaLinkReponse)
	}

	return socialMediaLinksReponse, nil
}

func (service *SocialMediaLinkServiceImpl) RedirectLink(ctx context.Context, request web.SocialMediaLinkRedirectRequest) (string, error) {
	// It's a transaction.
	tx := service.DB.Begin()
	defer helper.CommitOrRollback(tx)

	// It's checking if the social media name is valid or not.
	socialMediaName := helper.SocialMediaUrlToNameFormat(request.SocialMediaName)
	socialMediaType, repoErr := service.SocialMediaTypeRepository.FindByName(ctx, tx, socialMediaName)
	if repoErr != nil && errors.Is(repoErr, gorm.ErrRecordNotFound) {
		return "", ErrSocialMediaInvalidLink
	}
	service.Logger.PanicIfErr(repoErr, ErrSocialMediaLinkService)

	// It's checking if the username is valid or not.
	userData, repoErr := service.UserRepository.FindByUsername(ctx, tx, request.Username)
	if repoErr != nil && errors.Is(repoErr, gorm.ErrRecordNotFound) {
		return "", ErrSocialMediaInvalidLink
	}
	service.Logger.PanicIfErr(repoErr, ErrSocialMediaLinkService)

	// It's checking if the social media link is registered or not.
	socialMediaLink, repoErr := service.SocialMediaLinkRepository.FindByTypeAndUserID(ctx, tx, socialMediaType.ID, userData.ID)
	if (repoErr != nil && errors.Is(repoErr, gorm.ErrRecordNotFound)) || !socialMediaLink.Activate {
		return "", ErrSocialMediaInvalidLink
	}
	service.Logger.PanicIfErr(repoErr, ErrSocialMediaLinkService)

	linkResponse := helper.GenerateLinkResponse(socialMediaLink.SocialMediaType.Name, socialMediaLink.LinkOrUsername)
	return linkResponse, nil
}
