package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	configPath := "config.json"
	config := NewConfig(configPath)
	assert.IsType(t, config, &Config{})
}

func TestGetConfig(t *testing.T) {
	configPath := "config.json"
	config := NewConfig(configPath)

	t.Run("GetDatabaseConfig", func(t *testing.T) {
		dbConfig := config.GetDatabaseConfig()
		assert.IsType(t, dbConfig, DatabaseConfig{})
	})

	t.Run("GetServerConfig", func(t *testing.T) {
		serverConfig := config.GetServerConfig()
		assert.IsType(t, ServerConfig{}, serverConfig)
	})

	t.Run("GetLoggerConfig", func(t *testing.T) {
		logConfig := config.GetLoggerConfig()
		assert.IsType(t, LoggerConfig{}, logConfig)
	})

	t.Run("GetAppConfig", func(t *testing.T) {
		appConfig := config.GetAppConfig()
		assert.IsType(t, AppConfig{}, appConfig)
	})

	t.Run("GetJwtConfig", func(t *testing.T) {
		jwtConfig := config.GetJwtConfig()
		assert.IsType(t, JwtConfig{}, jwtConfig)
	})

}
