package config

import (
	"os"
	"strconv"
)

// Config 应用配置
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port string
	Mode string
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret string
	Expire int
}

// LoadConfig 加载配置
func LoadConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Mode: getEnv("SERVER_MODE", "debug"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "106.52.240.187"),
			Port:     getEnv("DB_PORT", "33306"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", "password@2023"),
			DBName:   getEnv("DB_NAME", "gorm"),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "yn98cryb98y4bcr9n2u49crnu43x9cru"),
			Expire: getEnvAsInt("JWT_EXPIRE", 24),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
