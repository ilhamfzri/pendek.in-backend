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

func TestSocialMediaServiceUpdateLink(t *testing.T) {
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

	socialMediaLinkRepository.Mock.On("FindByTypeAndUserID", mock.Anything, mock.Anything, uint(1), "123456").Return(domain.SocialMediaLink{}, gorm.ErrRecordNotFound)
	socialMediaLinkRepository.Mock.On("FindByTypeAndUserID", mock.Anything, mock.Anything, uint(2), "123456").Return(domain.SocialMediaLink{
		TypeID:          socialMediaTypes[1].ID,
		UserID:          "123456",
		LinkOrUsername:  "testuser99",
		Activate:        true,
		SocialMediaType: socialMediaTypes[1],
	}, nil)

	socialMediaLinkRepository.Mock.On("Update", mock.Anything, mock.Anything, mock.AnythingOfType("domain.SocialMediaLink")).Return(
		func(ctx context.Context, tx *gorm.DB, socialMediaLink domain.SocialMediaLink) domain.SocialMediaLink {
			return socialMediaLink
		},
		func(ctx context.Context, tx *gorm.DB, socialMediaLink domain.SocialMediaLink) error {
			return nil
		},
	)

	tests := []struct {
		TestName                    string
		Request                     web.SocialMediaLinkUpdateRequest
		ErrExpected                 error
		SocialMediaResponseExpected web.SocialMediaLinkResponse
	}{
		{TestName: "[Update Link][Failed: Social Media Type Invalid]",
			Request: web.SocialMediaLinkUpdateRequest{
				TypeID:            100,
				NewLinkOrUsername: "testuser01",
				Activate:          toBoolPointer(true),
			},
			ErrExpected:                 ErrSocialMediaTypeInvalid,
			SocialMediaResponseExpected: web.SocialMediaLinkResponse{},
		},
		{TestName: "[Update Link][Failed: Social Media Not Registered]",
			Request: web.SocialMediaLinkUpdateRequest{
				TypeID:            1,
				NewLinkOrUsername: "testuser01",
				Activate:          toBoolPointer(true),
			},
			ErrExpected:                 ErrSocialMediaLinkNotFound,
			SocialMediaResponseExpected: web.SocialMediaLinkResponse{},
		},
		{TestName: "[Update Link][Failed: Social Media Not Registered]",
			Request: web.SocialMediaLinkUpdateRequest{
				TypeID:            1,
				NewLinkOrUsername: "testuser01",
				Activate:          toBoolPointer(true),
			},
			ErrExpected:                 ErrSocialMediaLinkNotFound,
			SocialMediaResponseExpected: web.SocialMediaLinkResponse{},
		},
		{TestName: "[Update Link][Failed: Username Or Link Invalid]",
			Request: web.SocialMediaLinkUpdateRequest{
				TypeID:            2,
				NewLinkOrUsername: "te!tus//er01",
				Activate:          toBoolPointer(true),
			},
			ErrExpected:                 ErrSocialMediaLinkUsernameOrLink,
			SocialMediaResponseExpected: web.SocialMediaLinkResponse{},
		},
		{TestName: "[Update Link][Success]",
			Request: web.SocialMediaLinkUpdateRequest{
				TypeID:            2,
				NewLinkOrUsername: "testuser01",
				Activate:          toBoolPointer(false),
			},
			ErrExpected: nil,
			SocialMediaResponseExpected: web.SocialMediaLinkResponse{
				TypeID:          2,
				SocialMediaName: "Twitter",
				LinkOrUsername:  "testuser01",
				Activate:        false,
				RedirectLink:    "http://pendek.in/testuser/twitter",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.TestName, func(t *testing.T) {
			socialMediaResponse, err := socialMediaService.UpdateLink(ctx, test.Request, host, dummyJwt)
			assert.Equal(t, test.ErrExpected, err)
			assert.Equal(t, test.SocialMediaResponseExpected, socialMediaResponse)
		})
	}
}

func TestSocialMediaServiceGetAllLink(t *testing.T) {
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

	var socialMediaLinks = []domain.SocialMediaLink{
		{TypeID: socialMediaTypes[0].ID, SocialMediaType: socialMediaTypes[0], UserID: "123456", LinkOrUsername: "testuserinstagram", Activate: true},
		{TypeID: socialMediaTypes[1].ID, SocialMediaType: socialMediaTypes[1], UserID: "123456", LinkOrUsername: "testtwitter", Activate: true},
		{TypeID: socialMediaTypes[2].ID, SocialMediaType: socialMediaTypes[2], UserID: "123456", LinkOrUsername: "testtiktok", Activate: true},
	}
	socialMediaLinks[0].ID = 0
	socialMediaLinks[0].ID = 1
	socialMediaLinks[0].ID = 2

	socialMediaLinkRepository.Mock.On("FindByUserID", mock.Anything, mock.Anything, "123456").Return(socialMediaLinks, nil)

	t.Run("[Get All Link][Success]", func(t *testing.T) {
		socialMediaLinkResponses, err := socialMediaService.GetAllLink(ctx, host, dummyJwt)
		assert.Nil(t, err)
		assert.IsType(t, []web.SocialMediaLinkResponse{}, socialMediaLinkResponses)
	})

}

func TestSocialMediaServiceRedirectLink(t *testing.T) {
	var jwt = new(helper.JwtMock)
	var userRepository = mocks.NewUserRepository(t)
	var socialMediaLinkRepository = mocks.NewSocialMediaLinkRepository(t)
	var socialMediaTypeRepository = mocks.NewSocialMediaTypeRepository(t)

	var socialMediaService = NewSocialMediaLinkService(userRepository,
		socialMediaLinkRepository, socialMediaTypeRepository, db, log, jwt)

	socialMediaTypeRepository.Mock.On("FindByName", mock.Anything, mock.Anything, "Invalid").Return(domain.SocialMediaType{}, gorm.ErrRecordNotFound)
	socialMediaTypeRepository.Mock.On("FindByName", mock.Anything, mock.Anything, socialMediaTypes[0].Name).Return(socialMediaTypes[0], nil)
	socialMediaTypeRepository.Mock.On("FindByName", mock.Anything, mock.Anything, socialMediaTypes[1].Name).Return(socialMediaTypes[1], nil)
	userRepository.Mock.On("FindByUsername", mock.Anything, mock.Anything, "notusername").Return(domain.User{}, gorm.ErrRecordNotFound)
	userRepository.Mock.On("FindByUsername", mock.Anything, mock.Anything, "testusername").Return(domain.User{ID: "123456"}, nil)
	socialMediaLinkRepository.Mock.On("FindByTypeAndUserID", mock.Anything, mock.Anything, uint(socialMediaTypes[0].ID), "123456").Return(
		domain.SocialMediaLink{}, gorm.ErrRecordNotFound,
	)
	socialMediaLinkRepository.Mock.On("FindByTypeAndUserID", mock.Anything, mock.Anything, uint(socialMediaTypes[1].ID), "123456").Return(
		domain.SocialMediaLink{
			TypeID:          socialMediaTypes[1].ID,
			SocialMediaType: socialMediaTypes[1],
			UserID:          "123456",
			LinkOrUsername:  "testusertwitter",
			Activate:        true,
		}, nil,
	)

	var tests = []struct {
		TestName             string
		Request              web.SocialMediaLinkRedirectRequest
		LinkResponseExpected string
		ErrResponseExpected  error
	}{
		{
			TestName: "[Failed : Social Media Type Invalid]",
			Request: web.SocialMediaLinkRedirectRequest{
				Username:        "testusername",
				SocialMediaName: "invalid",
			},
			LinkResponseExpected: "",
			ErrResponseExpected:  ErrSocialMediaInvalidLink,
		},
		{
			TestName: "[Failed : Username Not Registered]",
			Request: web.SocialMediaLinkRedirectRequest{
				Username:        "notusername",
				SocialMediaName: "Instagram",
			},
			LinkResponseExpected: "",
			ErrResponseExpected:  ErrSocialMediaInvalidLink,
		},
		{
			TestName: "[Failed : Social Media Link Not Registered]",
			Request: web.SocialMediaLinkRedirectRequest{
				Username:        "testusername",
				SocialMediaName: "Instagram",
			},
			LinkResponseExpected: "",
			ErrResponseExpected:  ErrSocialMediaInvalidLink,
		},
		{
			TestName: "[Success]",
			Request: web.SocialMediaLinkRedirectRequest{
				Username:        "testusername",
				SocialMediaName: "Twitter",
			},
			LinkResponseExpected: "https://www.twitter.com/testusertwitter",
			ErrResponseExpected:  nil,
		},
	}

	for _, test := range tests {
		t.Run(test.TestName, func(t *testing.T) {
			linkResponse, err := socialMediaService.RedirectLink(ctx, test.Request)
			assert.Equal(t, err, test.ErrResponseExpected)
			assert.Equal(t, test.LinkResponseExpected, linkResponse)
		})
	}

}

func toBoolPointer(b bool) *bool {
	return &b
}
