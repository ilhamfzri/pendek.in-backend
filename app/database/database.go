package database

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ilhamfzri/pendek.in/app/logger"
	"github.com/ilhamfzri/pendek.in/config"
	"github.com/ilhamfzri/pendek.in/helper"
	"github.com/ilhamfzri/pendek.in/internal/model/domain"
	"github.com/ilhamfzri/pendek.in/internal/repository"
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

func NewDatabaseConnectionMock() *gorm.DB {
	sqlDB, _, _ := sqlmock.New()
	gormDB, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	return gormDB
}

func Migration(DB *gorm.DB, log *logger.Logger) {
	var errMigration = "[Database] Failed Migration"

	err := DB.AutoMigrate(&domain.User{})
	log.FatalIfErr(err, errMigration)
	log.Info().Msg("[Database] Successful Migration Users Table")

	err = DB.AutoMigrate(&domain.SocialMediaType{})
	log.FatalIfErr(err, errMigration)
	log.Info().Msg("[Database] Successful Migration SocialMediaType Table")

	err = DB.AutoMigrate(&domain.SocialMediaLink{})
	log.FatalIfErr(err, errMigration)
	log.Info().Msg("[Database] Successful Migration SocialMediaLink Table")

	err = DB.AutoMigrate(&domain.SocialMediaInteraction{})
	log.FatalIfErr(err, errMigration)
	log.Info().Msg("[Database] Successful Migration SocialMediaInteraction Table")

	err = DB.AutoMigrate(&domain.SocialMediaAnalytic{})
	log.FatalIfErr(err, errMigration)
	log.Info().Msg("[Database] Successful Migration SocialMediaAnalytic Table")

	err = DB.AutoMigrate(&domain.DeviceAnalytic{})
	log.FatalIfErr(err, errMigration)
	log.Info().Msg("[Database] Successful Migration DeviceAnalytic Table")

	CreateSocialMediaTypeEntries(DB, log)

}

func CreateSocialMediaTypeEntries(DB *gorm.DB, log *logger.Logger) {
	tx := DB.Begin()
	ctx := context.Background()
	defer helper.CommitOrRollback(tx)

	socialMediaTypeRepository := repository.NewSocialMediaTypeRepository(log)
	var socialMediaTypeEntries []domain.SocialMediaType = []domain.SocialMediaType{
		{Name: "Amazon", Example: "https://amazon.com/shop/yourshopname"},
		{Name: "Android Play Store", Example: "https://play.google.com/store/apps/details?url=.com.yourapp.app"},
		{Name: "Apple App Store", Example: "https://apps.apple.com/us/yourapp/url12346"},
		{Name: "Apple Music", Example: "https://music.apple.com/us/album/youralbum"},
		{Name: "Apple Podcast", Example: "https://podcasts.apple.com/us/podcast/yourpodcast/123456"},
		{Name: "Bandcamp", Example: "https://you.bandcamp.com/"},
		{Name: "BeReal", Example: "https://bere.al/yourusername"},
		{Name: "Cameo", Example: "https://cameo.com/"},
		{Name: "Clubhouse", Example: "https://clubhouse.com/@profile"},
		{Name: "Discord", Example: "https://discord.com/invite/yourchannel"},
		{Name: "Etsy", Example: "https://www.etsy.com/shop/yourshop"},
		{Name: "Facebook", Example: "https://facebook.com/facebookpageurl"},
		{Name: "Instagram", Example: "@yourinstagramusername"},
		{Name: "LinkedIn", Example: "https://linkedin.com/in/username"},
		{Name: "Patreon", Example: "https://patreon.com/"},
		{Name: "Payment", Example: "https://venmo.com/yourusername"},
		{Name: "Pinterest", Example: "https://pinterest.com/"},
		{Name: "Signal", Example: "https://t.me/"},
		{Name: "Snapchat", Example: "https://www.snapchat.com/add/yourusername"},
		{Name: "Soundcloud", Example: "https://soundcloud.com/username"},
		{Name: "Spotify", Example: "https://open.spotify.com/artist/artistname"},
		{Name: "Substack", Example: "https://you.substack.com/"},
		{Name: "Telegram", Example: "https://t.me/"},
		{Name: "Tiktok", Example: "@tiktokusername"},
		{Name: "Twitch", Example: "https://twitch.tv/"},
		{Name: "Twitter", Example: "@yourtwitterusername"},
		{Name: "Whatsapp", Example: "+0000000000"},
		{Name: "Youtube", Example: "https://youtube.com/channel/youtubechannelurl"},
		{Name: "Tokopedia", Example: "https://www.tokopedia.com/yourstore"},
		{Name: "Shopee", Example: "https://shopee.co.id/yourstore"},
	}

	for _, socialMediaTypeEntry := range socialMediaTypeEntries {
		_, repoErr := socialMediaTypeRepository.FindByName(ctx, tx, socialMediaTypeEntry.Name)
		if repoErr != nil && errors.Is(repoErr, gorm.ErrRecordNotFound) {
			_, err := socialMediaTypeRepository.Create(ctx, tx, socialMediaTypeEntry)
			log.PanicIfErr(err, "[Database] Failed Create SocialMediaType Entries")

		} else {
			log.PanicIfErr(repoErr, "[Database] Failed Create SocialMediaType Entries")
		}
	}
	log.Info().Msg("[Database] Successful Create SocialMediaType Entries")
}
