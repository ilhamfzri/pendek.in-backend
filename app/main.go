package main

import (
	"fmt"
	"os"

	"github.com/ilhamfzri/pendek.in/app/database"
	"github.com/ilhamfzri/pendek.in/app/logger"
	"github.com/ilhamfzri/pendek.in/app/router"
	"github.com/ilhamfzri/pendek.in/config"
	"github.com/ilhamfzri/pendek.in/helper"
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

	//.- Jwt Initalize
	jwtConfig := config.GetJwtConfig()
	jwt := helper.NewJwt(jwtConfig)

	//.- User Router Initalize
	router.AddUsersRoute(server, db, logger, jwt)

	//.- Run Server
	server.Run()

}
