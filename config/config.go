package config

import (
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

type Config struct {
	Viper *viper.Viper
}

func NewConfig(configPath string) *Config {
	cfg := viper.New()
	cfg.SetConfigFile(configPath)
	err := cfg.ReadInConfig()
	panicIfError(err)

	return &Config{
		Viper: cfg,
	}
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"db_name"`
}

func (config *Config) GetDatabaseConfig() DatabaseConfig {
	dbConfig := DatabaseConfig{}
	err := config.Viper.UnmarshalKey("database", &dbConfig)
	panicIfError(err)
	return dbConfig
}

type ServerConfig struct {
	Port         int `mapstructure:"port"`
	WriteTimeout int `mapstructure:"write_timeout"`
	ReadTimeout  int `mapstructure:"read_timeout"`
}

func (config *Config) GetServerConfig() ServerConfig {
	serverConfig := ServerConfig{}
	err := config.Viper.UnmarshalKey("server", &serverConfig)
	panicIfError(err)
	return serverConfig
}

type LoggerConfig struct {
	Level  zerolog.Level
	Output string
}

func (config *Config) GetLoggerConfig() LoggerConfig {
	loggerConfig := LoggerConfig{}

	loggerConfig.Output = config.Viper.GetString("log.output")

	level := config.Viper.Get("log.level")
	switch level {
	case "debug":
		loggerConfig.Level = zerolog.DebugLevel
	case "info":
		loggerConfig.Level = zerolog.InfoLevel
	case "warning":
		loggerConfig.Level = zerolog.WarnLevel
	case "error":
		loggerConfig.Level = zerolog.ErrorLevel
	case "fatal":
		loggerConfig.Level = zerolog.FatalLevel
	case "panic":
		loggerConfig.Level = zerolog.PanicLevel
	case "no-level":
		loggerConfig.Level = zerolog.NoLevel
	default:
		loggerConfig.Level = zerolog.InfoLevel
	}
	return loggerConfig
}

type AppConfig struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
	Stage   string `mapstructure:"stage"`
}

func (config *Config) GetAppConfig() AppConfig {
	appConfig := AppConfig{}
	err := config.Viper.UnmarshalKey("application", &appConfig)
	panicIfError(err)
	return appConfig
}

type JwtConfig struct {
	SigningKey     string `mapstructure:"signing_key"`
	ExpiredTimeDay int    `mapstructure:"expired_time_day"`
	Issuer         string `mapstructure:"issuer"`
}

func (config *Config) GetJwtConfig() JwtConfig {
	jwtConfig := JwtConfig{}
	err := config.Viper.UnmarshalKey("jwt_auth", &jwtConfig)
	panicIfError(err)
	return jwtConfig
}

func panicIfError(err error) {
	if err != nil {
		panic(err)
	}
}
