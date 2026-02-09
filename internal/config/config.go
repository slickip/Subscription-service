package config

import (
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string
	TimeZone   string
}

type HTTPServerConfig struct {
	Address     string        `yaml:"address"`
	Timeout     time.Duration `yaml:"timeout"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
}

type fileConfig struct {
	Env         string           `yaml:"env"`
	StoragePath string           `yaml:"storage_path"`
	HTTPServer  HTTPServerConfig `yaml:"http_server"`
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

func LoadHTTP() HTTPServerConfig {
	defaultCfg := HTTPServerConfig{
		Address:     "0.0.0.0:8081",
		Timeout:     4 * time.Second,
		IdleTimeout: 30 * time.Second,
	}

	path := getEnv("CONFIG_PATH", "config/subscription_local.yaml")

	data, err := os.ReadFile(path)
	if err != nil {
		log.Printf("config: unable to read YAML config %q, using defaults: %v", path, err)
		return defaultCfg
	}

	var fc fileConfig
	if err := yaml.Unmarshal(data, &fc); err != nil {
		log.Printf("config: unable to parse YAML config %q, using defaults: %v", path, err)
		return defaultCfg
	}

	if fc.HTTPServer.Address == "" {
		fc.HTTPServer.Address = defaultCfg.Address
	}
	if fc.HTTPServer.Timeout == 0 {
		fc.HTTPServer.Timeout = defaultCfg.Timeout
	}
	if fc.HTTPServer.IdleTimeout == 0 {
		fc.HTTPServer.IdleTimeout = defaultCfg.IdleTimeout
	}

	return fc.HTTPServer
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
