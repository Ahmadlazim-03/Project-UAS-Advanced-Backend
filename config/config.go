package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	// Server
	AppEnv     string
	Port       string
	APIVersion string

	// PostgreSQL
	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string
	DBSSLMode  string
	DBTimezone string

	// MongoDB
	MongoURI    string
	MongoDBName string

	// JWT
	JWTSecret           string
	JWTRefreshSecret    string
	JWTExpiresIn        time.Duration
	JWTRefreshExpiresIn time.Duration

	// File Upload
	MaxFileSize int64
	UploadPath  string

	// CORS
	CORSOrigin string

	// Rate Limit
	RateLimitMax      int
	RateLimitDuration time.Duration
}

func LoadConfig() *Config {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	return &Config{
		AppEnv:              getEnv("APP_ENV", "development"),
		Port:                getEnv("APP_PORT", "3000"),
		APIVersion:          getEnv("API_VERSION", "v1"),
		DBHost:              getEnv("DB_HOST", "localhost"),
		DBPort:              getEnv("DB_PORT", "5432"),
		DBName:              getEnv("DB_NAME", "achievement_db"),
		DBUser:              getEnv("DB_USER", "postgres"),
		DBPassword:          getEnv("DB_PASSWORD", "postgres"),
		DBSSLMode:           getEnv("DB_SSLMODE", "disable"),
		DBTimezone:          getEnv("DB_TIMEZONE", "Asia/Jakarta"),
		MongoURI:            getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDBName:         getEnv("MONGO_DB_NAME", "achievement_db"),
		JWTSecret:           getEnv("JWT_SECRET", "your-secret-key"),
		JWTRefreshSecret:    getEnv("JWT_REFRESH_SECRET", "your-refresh-secret-key"),
		JWTExpiresIn:        parseDuration(getEnv("JWT_EXPIRES_IN", "1h")),
		JWTRefreshExpiresIn: parseDuration(getEnv("JWT_REFRESH_EXPIRES_IN", "168h")),
		MaxFileSize:         5242880, // 5MB
		UploadPath:          getEnv("UPLOAD_PATH", "./uploads"),
		CORSOrigin:          getEnv("CORS_ORIGIN", "*"),
		RateLimitMax:        100,
		RateLimitDuration:   1 * time.Minute,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func parseDuration(s string) time.Duration {
	duration, err := time.ParseDuration(s)
	if err != nil {
		log.Printf("Invalid duration format: %s, using default 1h", s)
		return 1 * time.Hour
	}
	return duration
}
