package database

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/ilhamfzri/pendek.in/app/logger"
	"github.com/ilhamfzri/pendek.in/config"
	"github.com/ilhamfzri/pendek.in/internal/model/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const errMsg = "[Database] Failed Connecting to Database"

func NewDatabaseConnection(cfg config.DatabaseConfig, log *logger.Logger) *gorm.DB {
	log.Info().Msg("[Database] Initializing Database Connection ..")

	val := url.Values{}
	val.Add("parseTime", "1")
	val.Add("loc", "Asia/Jakarta")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		cfg.Host, cfg.User, cfg.Password, cfg.DBName, strconv.Itoa(cfg.Port))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	log.FatalIfErr(err, errMsg)

	// Ping Database
	sqlDB, err := db.DB()
	log.FatalIfErr(err, errMsg)

	err = sqlDB.Ping()
	log.FatalIfErr(err, errMsg)

	log.Info().Msg("[Database] Sucessfull Ping Database")
	return db
}

func Migration(DB *gorm.DB, log *logger.Logger) {
	var errMigration = "[Database] Failed Migration"

	err := DB.AutoMigrate(&domain.User{})
	log.FatalIfErr(err, errMigration)
	log.Info().Msg("[Database] Successful Migration Users Table")

}
