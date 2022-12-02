package main

import (
	"fmt"
	"os"

	"github.com/ilhamfzri/pendek.in/app/cache"
	"github.com/ilhamfzri/pendek.in/app/database"
	"github.com/ilhamfzri/pendek.in/app/logger"
	"github.com/ilhamfzri/pendek.in/app/mail"
	"github.com/ilhamfzri/pendek.in/app/router"
	"github.com/ilhamfzri/pendek.in/config"
	"github.com/ilhamfzri/pendek.in/helper"
	"github.com/ilhamfzri/pendek.in/internal/controller"
	"github.com/ilhamfzri/pendek.in/internal/handler"
	"github.com/ilhamfzri/pendek.in/internal/repository"
	"github.com/ilhamfzri/pendek.in/internal/service"
)

func main() {

	configPath := "config/config.json"
	config := config.NewConfig(configPath)

	appConfig := config.GetAppConfig()
	os.Setenv("APP_STAGE", appConfig.Stage)

	fmt.Printf("App Name \t: %s\n", appConfig.Name)
	fmt.Printf("Version \t: %s\n", appConfig.Version)
	fmt.Printf("Stage \t\t: %s\n", appConfig.Stage)

	//.- Logger Initalize
	logConfig := config.GetLoggerConfig()
	logger := logger.NewLogger(logConfig)

	//.- Database Initalize
	dbConfig := config.GetDatabaseConfig()
	db := database.NewDatabaseConnection(dbConfig, logger)

	//.- Database Migration
	database.Migration(db, logger)

	//.- Server Initalize
	serverConfig := config.GetServerConfig()
	server := router.NewServer(serverConfig)

	//.- Redis Initialize
	redisConfig := config.GetRedisConfig()
	redis := cache.NewRedisClient(redisConfig, logger)

	//.- Jwt Initalize
	jwtConfig := config.GetJwtConfig()
	jwt := helper.NewJwt(jwtConfig)

	//.- MailClient Initialize
	mailConfig := config.GetMailConfig()
	mailClient := mail.NewMailClient(mailConfig)

	//.- Recovery Handler
	recoveryHandler := handler.NewRecoveryHandler(logger)
	server.Router.Use(recoveryHandler)

	//.- No Method Handler
	noMethodHandler := handler.NewNoMethodHandler()
	server.Router.NoMethod(noMethodHandler)

	//.- No Route Handler
	noRouteHandler := handler.NewNoRouteHandler()
	server.Router.NoRoute(noRouteHandler)

	//.- Repository Initialize
	userRepository := repository.NewUserRepository(logger)
	socialMediaTypeRepository := repository.NewSocialMediaTypeRepository(logger)
	socialMediaLinkRepository := repository.NewSocialMediaLinkRepository(logger)
	socialMediaInteractionRepository := repository.NewSocialMediaInteractionRepository(logger)
	socialMediaAnalyticRepository := repository.NewSocialMediaAnalyticRepository(logger)
	customLinkRepository := repository.NewCustomLinkRepository(logger)
	customLinkAnalyticRepository := repository.NewCustomLinkAnalyticRepository(logger)
	customLinkInteractionRepository := repository.NewCustomLinkInteractionRepository(logger)
	customThumbnailRepository := repository.NewCustomThumbnailRepository(logger)
	thumbnailRepository := repository.NewThumbnailRepository(logger)
	deviceAnalyticRepository := repository.NewDeviceAnalyticRepository(logger)

	//.- Service Initialize
	userService := service.NewUserService(userRepository, mailClient, db, logger, jwt)
	socialMediaLinkService := service.NewSocialMediaLinkService(userRepository, socialMediaLinkRepository, socialMediaTypeRepository, db, logger, jwt)
	socialMediaAnalyticsService := service.NewSocialMediaAnalyticService(userRepository, socialMediaLinkRepository, socialMediaInteractionRepository, socialMediaAnalyticRepository, deviceAnalyticRepository, db, logger, jwt)
	customLinkService := service.NewCustomLinkService(customLinkRepository, customThumbnailRepository, thumbnailRepository, db, logger, jwt)
	customLinkAnalyticService := service.NewCustomLinkAnalyticService(customLinkRepository, customLinkAnalyticRepository, customLinkInteractionRepository, deviceAnalyticRepository, db, logger, jwt)

	//.- Controller Initialize
	userController := controller.NewUserController(userService, socialMediaLinkService, customLinkService, logger)
	socialMediaLinkController := controller.NewSocialMediaLink(socialMediaLinkService, socialMediaAnalyticsService, redis, logger)
	customLinkController := controller.NewCustomLinkController(customLinkService, customLinkAnalyticService, redis, logger)

	//.- User Router Initalize
	router.AddUsersRoute(server, userController, jwt)

	//.- Social Media Router Initialize
	router.AddSocialMediaRoute(server, socialMediaLinkController, jwt)

	//.- Custom Link Router Initialize
	router.AddCustomLinkRoute(server, customLinkController, jwt)

	//.- Run Server
	server.Run()

}
