package config

import "os"

type Config struct {
	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string
	TimeZone   string
}

func Load() *Config {
	return &Config{
		DBHost:     getEnv("DB_HOST", "postgres"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBName:     getEnv("DB_NAME", "subsdb"),
		DBUser:     getEnv("DB_USER", "subsuser"),
		DBPassword: getEnv("DB_PASSWORD", "subspass"),
		TimeZone:   getEnv("DB_TZ", "UTC"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
