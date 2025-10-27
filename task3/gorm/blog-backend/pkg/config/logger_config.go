package config

import (
	"os"
	"strconv"
)

// LogConfig 日志配置结构体
type LogConfig struct {
	LogLevel    string
	LogFilePath string
	MaxFileSize int64
	MaxBackups  int
	MaxAge      int
}

// GetLogConfig 获取日志配置
func GetLogConfig() *LogConfig {
	logLevel := getEnv("LOG_LEVEL", "INFO")
	logFilePath := getEnv("LOG_FILE_PATH", "app.log")
	maxFileSize, _ := strconv.ParseInt(getEnv("LOG_MAX_FILE_SIZE", "10485760"), 10, 64) // 10MB
	maxBackups, _ := strconv.Atoi(getEnv("LOG_MAX_BACKUPS", "5"))
	maxAge, _ := strconv.Atoi(getEnv("LOG_MAX_AGE", "30"))

	return &LogConfig{
		LogLevel:    logLevel,
		LogFilePath: logFilePath,
		MaxFileSize: maxFileSize,
		MaxBackups:  maxBackups,
		MaxAge:      maxAge,
	}
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}