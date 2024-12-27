package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port            string
	DBConn          string
	S3Region        string
	S3Endpoint      string
	S3BucketName    string
	S3ID            string
	S3Secret        string
	ImportDirImages string
	MaxUploadMB     int64
}

func Load() (*Config, error) {
	maxUploadMB, _ := strconv.ParseInt(getEnvOrDefault("MAX_UPLOAD_MB", "10"), 10, 64)

	cfg := &Config{
		Port:            getEnvOrDefault("PORT", "8080"),
		DBConn:          requireEnv("DB_CONN"),
		S3Region:        requireEnv("S3_REGION"),
		S3Endpoint:      requireEnv("S3_ENDPOINT"),
		S3BucketName:    requireEnv("S3_BUCKET_NAME"),
		S3ID:            requireEnv("S3_ID"),
		S3Secret:        requireEnv("S3_SECRET"),
		ImportDirImages: requireEnv("IMPORT_DIR_IMAGES"),
		MaxUploadMB:     maxUploadMB,
	}

	return cfg, nil
}

func requireEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("Required environment variable %s is not set", key))
	}
	return value
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
} 