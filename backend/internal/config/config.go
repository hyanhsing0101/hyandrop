package config

import (
	"os"
	"strconv"
)

type Config struct {
	AppEnv        string
	AppName       string
	DatabaseURL   string
	RedisAddr     string
	RedisPassword string
	RedisDB       int
	HTTPAddr      string
	PublicBaseURL string
	RoomTTLHours  int
	MaxFileSizeMB int
	UploadRoot    string
}

func Load() Config {
	return Config{
		AppEnv:        env("APP_ENV", "development"),
		AppName:       env("APP_NAME", "HyanDrop"),
		DatabaseURL:   env("DATABASE_URL", "postgres://hyandrop:hyandrop_dev_password@localhost:5432/hyandrop?sslmode=disable"),
		RedisAddr:     env("REDIS_ADDR", "localhost:6379"),
		RedisPassword: env("REDIS_PASSWORD", ""),
		RedisDB:       envInt("REDIS_DB", 0),
		HTTPAddr:      env("HTTP_ADDR", ":8080"),
		PublicBaseURL: env("PUBLIC_BASE_URL", "http://localhost:5173"),
		RoomTTLHours:  envInt("ROOM_TTL_HOURS", 24),
		MaxFileSizeMB: envInt("MAX_FILE_SIZE_MB", 100),
		UploadRoot:    env("UPLOAD_ROOT", "./storage/uploads"),
	}
}

func env(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func envInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}
