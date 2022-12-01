package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"path"

	"github.com/google/uuid"
	"github.com/ilhamfzri/pendek.in/app/logger"
	"github.com/ilhamfzri/pendek.in/helper"
	"github.com/ilhamfzri/pendek.in/internal/model/domain"
	"github.com/ilhamfzri/pendek.in/internal/model/web"
	"github.com/ilhamfzri/pendek.in/internal/repository"
	"github.com/nfnt/resize"
	"gorm.io/gorm"
)

type CustomLinkServiceImpl struct {
	CustomLinkRepository      repository.CustomLinkRepository
	CustomThumbnailRepository repository.CustomThumbnailRepository
	ThumbnailRepository       repository.ThumbnailRepository
	DB                        *gorm.DB
	Logger                    *logger.Logger
	Jwt                       helper.IJwt
}

var (
	ErrCustomLinkService       = "[Custom Link Service] Failed To Execute "
	ErrTwoTumbnailNotNull      = errors.New("set only one value either use thumbnail_id or user_thumbnail_id")
	ErrThumbnailIDNotFound     = errors.New("thumbnail_id invalid, make sure thumbnail_id is available")
	ErrUserThumbnailIDNotFound = errors.New("user_thumbnail_id invalid, make sure user_thumbnail_id is available")
	ErrShortLinkCodeRegistered = errors.New("short_link_code is registered")
	ErrCustomLinkNotRegistered = errors.New("link is not registered")
	ErrCustomLinkInvalid       = errors.New("link is invalid")
)

func NewCustomLinkService(clr repository.CustomLinkRepository, ctr repository.CustomThumbnailRepository, tr repository.ThumbnailRepository,
	db *gorm.DB, logger *logger.Logger, jwt helper.IJwt) CustomLinkService {
	return &CustomLinkServiceImpl{
		CustomLinkRepository:      clr,
		CustomThumbnailRepository: ctr,
		ThumbnailRepository:       tr,
		DB:                        db,
		Logger:                    logger,
		Jwt:                       jwt,
	}
}

func (service *CustomLinkServiceImpl) CreateLink(ctx context.Context, request web.CustomLinkCreateRequest, domainName string, jwtToken string) (web.CustomLinkResponse, error) {
	// It's getting the claims from the token.
	claims := service.Jwt.GetClaims(jwtToken)

	// It's a transaction.
	tx := service.DB.Begin()
	defer helper.CommitOrRollback(tx)

	if request.UserThumbnailID != nil && request.ThumbnailID != nil {
		return web.CustomLinkResponse{}, ErrTwoTumbnailNotNull
	}

	var thumbnailUrl string

	if request.ThumbnailID != nil {
		thumbnailID := *request.ThumbnailID
		thumbnail, errRepo := service.ThumbnailRepository.FindByID(ctx, tx, int(thumbnailID))

		if errRepo != nil {
			if errors.Is(errRepo, gorm.ErrRecordNotFound) {
				return web.CustomLinkResponse{}, ErrThumbnailIDNotFound
			} else {
				service.Logger.PanicIfErr(errRepo, ErrCustomLinkService)
			}
		}
		thumbnailUrl = thumbnail.IconUrl
	}

	if request.UserThumbnailID != nil {
		userThumbnailID := *request.UserThumbnailID
		customThumbnail, errRepo := service.CustomThumbnailRepository.FindByThumbnailIDAndUserID(ctx, tx, int(userThumbnailID), claims.Id)
		if errRepo != nil {
			if errors.Is(errRepo, gorm.ErrRecordNotFound) {
				return web.CustomLinkResponse{}, ErrUserThumbnailIDNotFound
			} else {
				service.Logger.PanicIfErr(errRepo, ErrCustomLinkService)
			}
		}
		thumbnailUrl = helper.GetCustomThumbnailUrl(domainName, customThumbnail.ImageID)

	}

	customLink, errRepo := service.CustomLinkRepository.FindByShortLinkCode(ctx, tx, request.ShortLinkCode)
	if errRepo != nil && !errors.Is(errRepo, gorm.ErrRecordNotFound) {
		service.Logger.PanicIfErr(errRepo, ErrCustomLinkService)
	}

	if customLink.ID != 0 {
		return web.CustomLinkResponse{}, ErrShortLinkCodeRegistered
	}

	customLink = domain.CustomLink{
		UserID:        claims.Id,
		Title:         request.Title,
		ShortLinkCode: request.ShortLinkCode,
		LongLink:      request.LongLink,
		ShowOnProfile: true,
		Activate:      true,
	}

	if request.UserThumbnailID != nil {
		customLink.CustomThumbnailID = request.UserThumbnailID
	}

	if request.ThumbnailID != nil {
		customLink.ThumbnailID = request.ThumbnailID
	}

	customLink, errRepo = service.CustomLinkRepository.Create(ctx, tx, customLink)
	service.Logger.PanicIfErr(errRepo, ErrCustomLinkService)

	customLinkResponse := helper.CustomLinkDomainToResponse(&customLink)
	customLinkResponse.ThumbnailUrl = thumbnailUrl
	return customLinkResponse, nil
}

func (service *CustomLinkServiceImpl) UpdateLink(ctx context.Context, request web.CustomLinkUpdateRequest, domainName string, jwtToken string) (web.CustomLinkResponse, error) {
	// It's getting the claims from the token.
	claims := service.Jwt.GetClaims(jwtToken)

	// It's a transaction.
	tx := service.DB.Begin()
	defer helper.CommitOrRollback(tx)

	if request.ThumbnailID != nil && request.UserThumbnailID != nil {
		if *request.ThumbnailID != 0 || *request.UserThumbnailID != 0 {
			return web.CustomLinkResponse{}, ErrTwoTumbnailNotNull
		}
	}

	customLink, errRepo := service.CustomLinkRepository.FindByIdAndUserID(ctx, tx, int(request.CustomLinkID), claims.Id)
	if errRepo != nil && !errors.Is(errRepo, gorm.ErrRecordNotFound) {
		service.Logger.PanicIfErr(errRepo, ErrCustomLinkService)
	}

	if errors.Is(errRepo, gorm.ErrRecordNotFound) {
		return web.CustomLinkResponse{}, ErrCustomLinkNotRegistered
	}

	if request.ShortLinkCode != "" {
		customLinkShortLink, errRepo := service.CustomLinkRepository.FindByShortLinkCode(ctx, tx, request.ShortLinkCode)
		if errRepo != nil && !errors.Is(errRepo, gorm.ErrRecordNotFound) {
			service.Logger.PanicIfErr(errRepo, ErrCustomLinkService)
		}

		if customLinkShortLink.ID != 0 {
			return web.CustomLinkResponse{}, ErrShortLinkCodeRegistered
		}
	}

	if request.LongLink != "" {
		customLink.LongLink = request.LongLink
	}

	if request.ThumbnailID != nil && request.UserThumbnailID != nil {
		if *request.ThumbnailID == 0 || *request.UserThumbnailID == 0 {
			customLink.ThumbnailID = nil
			customLink.CustomThumbnailID = nil
		}
	}

	var thumbnailUrl string

	if request.ThumbnailID != nil && *request.ThumbnailID != 0 {
		customLink.ThumbnailID = request.ThumbnailID
		customLink.CustomThumbnailID = nil

		thumbnailID := *request.ThumbnailID
		thumbnail, errRepo := service.ThumbnailRepository.FindByID(ctx, tx, int(thumbnailID))

		if errRepo != nil {
			if errors.Is(errRepo, gorm.ErrRecordNotFound) {
				return web.CustomLinkResponse{}, ErrThumbnailIDNotFound
			} else {
				service.Logger.PanicIfErr(errRepo, ErrCustomLinkService)
			}
		}
		thumbnailUrl = thumbnail.IconUrl
	}

	if request.UserThumbnailID != nil && *request.UserThumbnailID != 0 {
		customLink.CustomThumbnailID = request.UserThumbnailID
		customLink.ThumbnailID = nil

		userThumbnailID := *request.UserThumbnailID
		customThumbnail, errRepo := service.CustomThumbnailRepository.FindByThumbnailIDAndUserID(ctx, tx, int(userThumbnailID), claims.Id)
		if errRepo != nil {
			if errors.Is(errRepo, gorm.ErrRecordNotFound) {
				return web.CustomLinkResponse{}, ErrUserThumbnailIDNotFound
			} else {
				service.Logger.PanicIfErr(errRepo, ErrCustomLinkService)
			}
		}
		thumbnailUrl = helper.GetCustomThumbnailUrl(domainName, customThumbnail.ImageID)
	}

	if request.Title != "" {
		customLink.Title = request.Title
	}

	if request.ShortLinkCode != "" {
		customLink.ShortLinkCode = request.ShortLinkCode
	}

	if request.LongLink != "" {
		customLink.LongLink = request.LongLink
	}

	if request.ShowOnProfile != nil {
		customLink.ShowOnProfile = *request.ShowOnProfile
	}

	if request.Activate != nil {
		customLink.Activate = *request.Activate
	}

	updateCustomThumbnailID := customLink.CustomThumbnailID
	updateThumbnailID := customLink.ThumbnailID

	customLink, errRepo = service.CustomLinkRepository.Update(ctx, tx, customLink)
	service.Logger.PanicIfErr(errRepo, ErrCustomLinkService)

	customLink, errRepo = service.CustomLinkRepository.UpdateThumbnailIDFK(ctx, tx, customLink.ID, updateThumbnailID)
	service.Logger.PanicIfErr(errRepo, ErrCustomLinkService)

	customLink, errRepo = service.CustomLinkRepository.UpdateCustomThumbnailIDFK(ctx, tx, customLink.ID, updateCustomThumbnailID)
	service.Logger.PanicIfErr(errRepo, ErrCustomLinkService)

	customLinkResponse := helper.CustomLinkDomainToResponse(&customLink)

	if thumbnailUrl == "" {
		if customLink.CustomThumbnailID != nil {
			thumbnailUrl = helper.GetCustomThumbnailUrl(domainName, customLink.CustomThumbnail.ImageID)
		} else if customLink.ThumbnailID != nil {
			thumbnailUrl = customLink.Thumbnail.IconUrl
		}
	}

	customLinkResponse.ThumbnailUrl = thumbnailUrl
	return customLinkResponse, nil
}

func (service *CustomLinkServiceImpl) GetLink(ctx context.Context, request web.CustomLinkGetRequest, domainName string, jwtToken string) (web.CustomLinkResponse, error) {
	// It's getting the claims from the token.
	claims := service.Jwt.GetClaims(jwtToken)

	// It's a transaction.
	tx := service.DB.Begin()
	defer helper.CommitOrRollback(tx)

	customLink, errRepo := service.CustomLinkRepository.FindByIdAndUserID(ctx, tx, int(request.LinkID), claims.Id)
	if errRepo != nil && !errors.Is(errRepo, gorm.ErrRecordNotFound) {
		service.Logger.PanicIfErr(errRepo, ErrCustomLinkService)
	}

	if errors.Is(errRepo, gorm.ErrRecordNotFound) {
		return web.CustomLinkResponse{}, ErrCustomLinkNotRegistered
	}

	var thumbnailUrl string
	if customLink.CustomThumbnailID != nil {
		thumbnailUrl = helper.GetCustomThumbnailUrl(domainName, customLink.CustomThumbnail.ImageID)
	}

	if customLink.ThumbnailID != nil {
		thumbnailUrl = customLink.Thumbnail.IconUrl
	}

	customLinkResponse := helper.CustomLinkDomainToResponse(&customLink)
	customLinkResponse.ThumbnailUrl = thumbnailUrl
	return customLinkResponse, nil
}

func (service *CustomLinkServiceImpl) GetAllLink(ctx context.Context, domainName string, jwtToken string) ([]web.CustomLinkResponse, error) {
	// It's getting the claims from the token.
	claims := service.Jwt.GetClaims(jwtToken)

	// It's a transaction.
	tx := service.DB.Begin()
	defer helper.CommitOrRollback(tx)

	customLinks, errRepo := service.CustomLinkRepository.FetchAllByUserID(ctx, tx, claims.Id)
	if errRepo != nil && !errors.Is(errRepo, gorm.ErrRecordNotFound) {
		service.Logger.PanicIfErr(errRepo, ErrCustomLinkService)
	}

	var customLinksResponse []web.CustomLinkResponse

	for _, customLink := range customLinks {
		customLinkResponse := helper.CustomLinkDomainToResponse(&customLink)
		if customLink.ThumbnailID != nil {
			customLinkResponse.ThumbnailUrl = customLink.Thumbnail.IconUrl
		}
		if customLink.CustomThumbnailID != nil {
			customLinkResponse.ThumbnailUrl = helper.GetCustomThumbnailUrl(domainName, customLink.CustomThumbnail.ImageID)
		}

		customLinksResponse = append(customLinksResponse, customLinkResponse)
	}
	return customLinksResponse, nil
}

func (service *CustomLinkServiceImpl) GetAllThumbnail(ctx context.Context) ([]web.ThumbnailResponse, error) {
	// It's a transaction.
	tx := service.DB.Begin()
	defer helper.CommitOrRollback(tx)

	thumbnails, errRepo := service.ThumbnailRepository.FetchAll(ctx, tx)
	var thumbnailsResponse []web.ThumbnailResponse
	if errRepo != nil && !errors.Is(errRepo, gorm.ErrRecordNotFound) {
		service.Logger.PanicIfErr(errRepo, ErrCustomLinkService)
	}

	for _, thumbnail := range thumbnails {
		thumbnailResponse := helper.ThumbnailDomainToResponse(&thumbnail)
		thumbnailsResponse = append(thumbnailsResponse, thumbnailResponse)
	}
	return thumbnailsResponse, nil
}

func (service *CustomLinkServiceImpl) GetUserThumbnail(ctx context.Context, domainName string, jwtToken string) ([]web.ThumbnailResponse, error) {
	// It's getting the claims from the token.
	claims := service.Jwt.GetClaims(jwtToken)

	// It's a transaction.
	tx := service.DB.Begin()
	defer helper.CommitOrRollback(tx)

	customThumbnails, errRepo := service.CustomThumbnailRepository.FetchAllByUserID(ctx, tx, claims.Id)
	if errRepo != nil && !errors.Is(errRepo, gorm.ErrRecordNotFound) {
		service.Logger.PanicIfErr(errRepo, ErrCustomLinkService)
	}

	var thumbnailsResponse []web.ThumbnailResponse
	for _, customThumbnail := range customThumbnails {
		thumbnailResponse := helper.CustomThumbnailDomainToResponse(&customThumbnail, domainName)
		thumbnailsResponse = append(thumbnailsResponse, thumbnailResponse)
	}
	return thumbnailsResponse, nil
}

func (service *CustomLinkServiceImpl) UploadCustomThumbnail(ctx context.Context, imgData []byte, domainName string, jwtToken string) (web.ThumbnailResponse, error) {
	// It's getting the claims from the token.
	claims := service.Jwt.GetClaims(jwtToken)

	// It's a transaction.
	tx := service.DB.Begin()
	defer helper.CommitOrRollback(tx)

	// It's decoding the image from the byte array.
	reader := bytes.NewReader(imgData)
	img, _, err := image.Decode(reader)
	service.Logger.PanicIfErr(err, ErrCustomLinkService)

	resizeImg := resize.Resize(150, 150, img, resize.Lanczos3)
	uuid := uuid.New().String()
	fileName := fmt.Sprintf("%s.jpg", uuid)

	thumbnailResourcePath := os.Getenv("THUMBNAIL_IMG_DIR")
	filePath := path.Join(thumbnailResourcePath, fileName)

	out, err := os.Create(filePath)
	service.Logger.PanicIfErr(err, ErrCustomLinkService)
	defer out.Close()
	jpeg.Encode(out, resizeImg, nil)

	thumbnail := domain.CustomThumbnail{
		UserID:  claims.Id,
		ImageID: uuid,
	}

	thumbnail, errRepo := service.CustomThumbnailRepository.Create(ctx, tx, thumbnail)
	service.Logger.PanicIfErr(errRepo, ErrCustomLinkService)

	thumbnailResponse := helper.CustomThumbnailDomainToResponse(&thumbnail, domainName)

	return thumbnailResponse, nil
}

func (service *CustomLinkServiceImpl) CheckShortLinkAvaibility(ctx context.Context, request web.CustomLinkCheckShortCodeAvaibilityRequest) error {
	// It's a transaction.
	tx := service.DB.Begin()
	defer helper.CommitOrRollback(tx)

	customLink, errRepo := service.CustomLinkRepository.FindByShortLinkCode(ctx, tx, request.Code)
	if errRepo != nil && !errors.Is(errRepo, gorm.ErrRecordNotFound) {
		service.Logger.PanicIfErr(errRepo, ErrCustomLinkService)
	}

	if customLink.ID != 0 {
		return ErrShortLinkCodeRegistered
	}
	return nil
}

func (service *CustomLinkServiceImpl) RedirectLink(ctx context.Context, request web.CustomLinkRedirectRequest) (string, uint, error) {
	// It's a transaction.
	tx := service.DB.Begin()
	defer helper.CommitOrRollback(tx)

	customLink, errRepo := service.CustomLinkRepository.FindByShortLinkCode(ctx, tx, request.ShortLinkCode)
	if errRepo != nil && !errors.Is(errRepo, gorm.ErrRecordNotFound) {
		service.Logger.PanicIfErr(errRepo, ErrCustomLinkService)
	}

	if errors.Is(errRepo, gorm.ErrRecordNotFound) {
		return "", 0, ErrCustomLinkInvalid
	}

	if !customLink.Activate {
		return "", 0, ErrCustomLinkInvalid
	}

	return customLink.LongLink, customLink.ID, nil
}
