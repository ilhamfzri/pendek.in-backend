package service

import (
	"context"
	"testing"

	"github.com/ilhamfzri/pendek.in/helper"
	"github.com/ilhamfzri/pendek.in/internal/model/domain"
	"github.com/ilhamfzri/pendek.in/internal/model/web"
	"github.com/ilhamfzri/pendek.in/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

var socialMediaTypes = []domain.SocialMediaType{
	{ID: 1, Name: "Instagram", Example: "@urinstagram"},
	{ID: 2, Name: "Twitter", Example: "@urtwitter"},
	{ID: 3, Name: "Tiktok", Example: "@urtiktok"},
	{ID: 4, Name: "Youtube", Example: "https://youtube.com/channel/youtubechannelurl"},
	{ID: 5, Name: "Whatsapp", Example: "+0000000000"},
}

var socialMediaLinkFound = domain.SocialMediaLink{
	TypeID:         1,
	UserID:         "2312jhdsam21312d",
	LinkOrUsername: "testuserinstagram",
	Activate:       true,
}

var socialMediaInvalidID = 100

func TestSocialMediaServiceGetAllTypes(t *testing.T) {
	var jwt = new(helper.JwtMock)
	var userRepository = mocks.NewUserRepository(t)
	var socialMediaLinkRepository = mocks.NewSocialMediaLinkRepository(t)
	var socialMediaTypeRepository = mocks.NewSocialMediaTypeRepository(t)

	var socialMediaService = NewSocialMediaLinkService(userRepository,
		socialMediaLinkRepository, socialMediaTypeRepository, db, log, jwt)

	socialMediaTypeRepository.Mock.On("FetchAll", mock.Anything, mock.Anything).Return(socialMediaTypes, nil)

	t.Run(
		"[GetAllTypes][Success]", func(t *testing.T) {
			socialMediaTypeResponses, err := socialMediaService.GetAllTypes(ctx)
			assert.Nil(t, err)
			assert.IsType(t, []web.SocialMediaTypeResponse{}, socialMediaTypeResponses)
		},
	)
}

func TestSocialMediaServiceCreateLink(t *testing.T) {
	var jwt = new(helper.JwtMock)
	dummyJwt := "ASDEFGHJKDSANEQWENEWNQENWN"
	host := "http://pendek.in"

	jwt.Mock.On("GetClaims", dummyJwt).Return(helper.JwtUserClaims{
		Id:       "123456",
		Username: "testuser",
		Email:    "testuser@mail.com",
	})

	var userRepository = mocks.NewUserRepository(t)
	var socialMediaLinkRepository = mocks.NewSocialMediaLinkRepository(t)
	var socialMediaTypeRepository = mocks.NewSocialMediaTypeRepository(t)

	var socialMediaService = NewSocialMediaLinkService(userRepository,
		socialMediaLinkRepository, socialMediaTypeRepository, db, log, jwt)

	socialMediaTypeRepository.Mock.On("FindByID", mock.Anything, mock.Anything, socialMediaInvalidID).Return(domain.SocialMediaType{}, gorm.ErrRecordNotFound)
	socialMediaTypeRepository.Mock.On("FindByID", mock.Anything, mock.Anything, 1).Return(socialMediaTypes[0], nil)
	socialMediaTypeRepository.Mock.On("FindByID", mock.Anything, mock.Anything, 2).Return(socialMediaTypes[1], nil)
	socialMediaTypeRepository.Mock.On("FindByID", mock.Anything, mock.Anything, 3).Return(socialMediaTypes[2], nil)
	socialMediaTypeRepository.Mock.On("FindByID", mock.Anything, mock.Anything, 4).Return(socialMediaTypes[3], nil)
	socialMediaTypeRepository.Mock.On("FindByID", mock.Anything, mock.Anything, 5).Return(socialMediaTypes[4], nil)

	socialMediaLinkFound.ID = 1
	socialMediaLinkRepository.Mock.On("FindByTypeAndUserID", mock.Anything, mock.Anything, uint(1), "123456").Return(socialMediaLinkFound, nil)
	socialMediaLinkRepository.Mock.On("FindByTypeAndUserID", mock.Anything, mock.Anything, uint(2), "123456").Return(domain.SocialMediaLink{}, gorm.ErrRecordNotFound)
	socialMediaLinkRepository.Mock.On("FindByTypeAndUserID", mock.Anything, mock.Anything, uint(3), "123456").Return(domain.SocialMediaLink{}, gorm.ErrRecordNotFound)
	socialMediaLinkRepository.Mock.On("FindByTypeAndUserID", mock.Anything, mock.Anything, uint(4), "123456").Return(domain.SocialMediaLink{}, gorm.ErrRecordNotFound)
	socialMediaLinkRepository.Mock.On("FindByTypeAndUserID", mock.Anything, mock.Anything, uint(5), "123456").Return(domain.SocialMediaLink{}, gorm.ErrRecordNotFound)

	socialMediaLinkRepository.Mock.On("Create", mock.Anything, mock.Anything, mock.AnythingOfType("domain.SocialMediaLink")).Return(
		func(ctx context.Context, tx *gorm.DB, socialMediaLink domain.SocialMediaLink) domain.SocialMediaLink {
			return socialMediaLink
		},
		func(ctx context.Context, tx *gorm.DB, socialMediaLink domain.SocialMediaLink) error {
			return nil
		},
	)

	t.Run(
		"[CreateLink][Failed: Social Media Type Invalid]", func(t *testing.T) {
			request := web.SocialMediaLinkCreateRequest{
				TypeID:         socialMediaInvalidID,
				LinkOrUsername: "testuser01",
			}
			socialMediaResponse, err := socialMediaService.CreateLink(ctx, request, host, dummyJwt)
			assert.NotNil(t, err)
			assert.Equal(t, ErrSocialMediaTypeInvalid, err)
			assert.IsType(t, web.SocialMediaLinkResponse{}, socialMediaResponse)

		})

	t.Run("[CreateLink][Failed: Username Or Link Invalid]", func(t *testing.T) {
		tests := []struct {
			Request     web.SocialMediaLinkCreateRequest
			ErrExpected error
		}{
			{Request: web.SocialMediaLinkCreateRequest{
				TypeID:         1,
				LinkOrUsername: "!2testuser01",
			},
				ErrExpected: ErrSocialMediaLinkUsernameOrLink},
			{Request: web.SocialMediaLinkCreateRequest{
				TypeID:         2,
				LinkOrUsername: "thttps//estuser01",
			},
				ErrExpected: ErrSocialMediaLinkUsernameOrLink},
			{Request: web.SocialMediaLinkCreateRequest{
				TypeID:         3,
				LinkOrUsername: "testuse!!r01",
			},
				ErrExpected: ErrSocialMediaLinkUsernameOrLink},
			{Request: web.SocialMediaLinkCreateRequest{
				TypeID:         4,
				LinkOrUsername: "testuse!!r01",
			},
				ErrExpected: ErrSocialMediaLinkUsernameOrLink},
			{Request: web.SocialMediaLinkCreateRequest{
				TypeID:         5,
				LinkOrUsername: "e123wq",
			},
				ErrExpected: ErrSocialMediaLinkUsernameOrLink},
		}

		for _, test := range tests {
			socialMediaResponse, err := socialMediaService.CreateLink(ctx, test.Request, host, dummyJwt)
			assert.NotNil(t, err)
			assert.Equal(t, test.ErrExpected, err)
			assert.IsType(t, web.SocialMediaLinkResponse{}, socialMediaResponse)
		}
	})

	t.Run("[CreateLink][Failed: Social Media Link Registered", func(t *testing.T) {
		request := web.SocialMediaLinkCreateRequest{
			TypeID:         1,
			LinkOrUsername: "testuser01",
		}
		socialMediaResponse, err := socialMediaService.CreateLink(ctx, request, host, dummyJwt)
		assert.NotNil(t, err)
		assert.Equal(t, ErrSocialMediaLinkFound, err)
		assert.IsType(t, web.SocialMediaLinkResponse{}, socialMediaResponse)
	})

	t.Run("[CreateLink][Success]", func(t *testing.T) {
		tests := []struct {
			Request                 web.SocialMediaLinkCreateRequest
			SocialMediaNameExpected string
			RedirectLinkExpected    string
			ErrExpected             error
		}{
			{Request: web.SocialMediaLinkCreateRequest{
				TypeID:         2,
				LinkOrUsername: "usertwitter",
			},
				SocialMediaNameExpected: "Twitter",
				RedirectLinkExpected:    "http://pendek.in/testuser/twitter",
				ErrExpected:             ErrSocialMediaLinkUsernameOrLink},
			{Request: web.SocialMediaLinkCreateRequest{
				TypeID:         3,
				LinkOrUsername: "usertiktok",
			},
				SocialMediaNameExpected: "Tiktok",
				RedirectLinkExpected:    "http://pendek.in/testuser/tiktok",
				ErrExpected:             ErrSocialMediaLinkUsernameOrLink},
			{Request: web.SocialMediaLinkCreateRequest{
				TypeID:         4,
				LinkOrUsername: "https://youtube.com/channel/youtubechannelurl",
			},
				SocialMediaNameExpected: "Youtube",
				RedirectLinkExpected:    "http://pendek.in/testuser/youtube",
				ErrExpected:             ErrSocialMediaLinkUsernameOrLink},
			{Request: web.SocialMediaLinkCreateRequest{
				TypeID:         5,
				LinkOrUsername: "+6212345678",
			},
				SocialMediaNameExpected: "Whatsapp",
				RedirectLinkExpected:    "http://pendek.in/testuser/whatsapp",
				ErrExpected:             ErrSocialMediaLinkUsernameOrLink},
		}

		for _, test := range tests {
			socialMediaResponse, err := socialMediaService.CreateLink(ctx, test.Request, host, dummyJwt)
			assert.Nil(t, err)
			assert.Equal(t, uint(test.Request.TypeID), socialMediaResponse.TypeID)
			assert.Equal(t, test.Request.LinkOrUsername, socialMediaResponse.LinkOrUsername)
			assert.Equal(t, test.SocialMediaNameExpected, socialMediaResponse.SocialMediaName)
			assert.Equal(t, test.RedirectLinkExpected, socialMediaResponse.RedirectLink)
			assert.IsType(t, web.SocialMediaLinkResponse{}, socialMediaResponse)
		}
	})
}
