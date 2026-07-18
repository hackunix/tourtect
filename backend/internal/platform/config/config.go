package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port         string
	DatabaseURL  string
	RedisAddr    string
	RedisPass    string
	MinioUser    string
	MinioPass    string
	LogLevel     string
	FptApiKey    string
	FptBaseURL   string
}

func Load() (*Config, error) {
	port := getEnv("PORT", "8080")

	// Database DSN construction
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbUser := getEnv("POSTGRES_USER", "tourtect")
		dbPass := getEnv("POSTGRES_PASSWORD", "change_me_postgres")
		dbDB := getEnv("POSTGRES_DB", "tourtect")
		dbPort := getEnv("POSTGRES_PORT", "5432")
		dbHost := getEnv("POSTGRES_HOST", "localhost")
		dbURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPass, dbHost, dbPort, dbDB)
	}

	redisHost := getEnv("REDIS_HOST", "localhost")
	redisPort := getEnv("REDIS_PORT", "6379")
	redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)
	redisPass := os.Getenv("REDIS_PASSWORD")

	minioUser := getEnv("MINIO_ROOT_USER", "tourtect")
	minioPass := getEnv("MINIO_ROOT_PASSWORD", "change_me_minio")

	logLevel := getEnv("LOG_LEVEL", "info")

	fptApiKey := os.Getenv("FPT_AI_API_KEY")
	fptBaseURL := getEnv("FPT_AI_BASE_URL", "https://mkp-api.fptcloud.com")

	return &Config{
		Port:        port,
		DatabaseURL: dbURL,
		RedisAddr:   redisAddr,
		RedisPass:   redisPass,
		MinioUser:   minioUser,
		MinioPass:   minioPass,
		LogLevel:    logLevel,
		FptApiKey:   fptApiKey,
		FptBaseURL:  fptBaseURL,
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	s := os.Getenv(key)
	if s == "" {
		return defaultValue
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return defaultValue
	}
	return v
}
