package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	Database  sqlx.SqlConf
	Cache     cache.CacheConf `json:",optional"`
	JWTSecret string          `json:",default=your-secret-key-change-in-production"`
}

// LoadFromEnv loads configuration from environment variables
func (c *Config) LoadFromEnv() {
	// Server configuration
	if host := os.Getenv("SERVER_HOST"); host != "" {
		c.Host = host
	}
	if port := os.Getenv("SERVER_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			c.Port = p
		}
	}

	// JWT Secret
	if jwtSecret := os.Getenv("JWT_SECRET"); jwtSecret != "" {
		c.JWTSecret = jwtSecret
	}

	// Database configuration - only override if env vars are set
	if os.Getenv("DB_HOST") != "" || os.Getenv("DB_USER") != "" || os.Getenv("DB_NAME") != "" {
		dbHost := getEnvOrDefault("DB_HOST", "localhost")
		dbPort := getEnvOrDefault("DB_PORT", "3306")
		dbUser := getEnvOrDefault("DB_USER", "root")
		dbPassword := getEnvOrDefault("DB_PASSWORD", "")
		dbName := getEnvOrDefault("DB_NAME", "wata_bot")
		dbCharset := getEnvOrDefault("DB_CHARSET", "utf8mb4")
		dbTimezone := getEnvOrDefault("DB_TIMEZONE", "Asia/Ho_Chi_Minh")

		// Build DataSource string (URL encode timezone to handle special characters like /)
		encodedTimezone := url.QueryEscape(dbTimezone)
		c.Database.DataSource = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true&loc=%s",
			dbUser, dbPassword, dbHost, dbPort, dbName, dbCharset, encodedTimezone)
	}

	// Log configuration
	if logServiceName := os.Getenv("LOG_SERVICE_NAME"); logServiceName != "" {
		c.Log.ServiceName = logServiceName
	}
	if logMode := os.Getenv("LOG_MODE"); logMode != "" {
		c.Log.Mode = logMode
	}
	if logPath := os.Getenv("LOG_PATH"); logPath != "" {
		c.Log.Path = logPath
	}
	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		c.Log.Level = logLevel
	}
	if logCompress := os.Getenv("LOG_COMPRESS"); logCompress != "" {
		if compress, err := strconv.ParseBool(logCompress); err == nil {
			c.Log.Compress = compress
		}
	}
	if logKeepDays := os.Getenv("LOG_KEEP_DAYS"); logKeepDays != "" {
		if days, err := strconv.Atoi(logKeepDays); err == nil {
			c.Log.KeepDays = days
		}
	}
}

// getEnvOrDefault returns environment variable value or default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

