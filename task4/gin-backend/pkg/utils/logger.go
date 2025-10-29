package utils

import (
	"fmt"
	"log"
	"os"
	"time"
)

// Logger 结构体封装日志功能
type Logger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
	debugLogger *log.Logger
	warnLogger  *log.Logger
}

// NewLogger 创建新的日志实例
func NewLogger() *Logger {
	// 创建日志文件
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("无法创建日志文件:", err)
	}

	// 定义日志前缀格式
	prefix := func(level string) string {
		return fmt.Sprintf("[%s] [%s] ", time.Now().Format("2006-01-02 15:04:05"), level)
	}

	return &Logger{
		infoLogger:  log.New(logFile, prefix("INFO"), log.LstdFlags),
		errorLogger: log.New(logFile, prefix("ERROR"), log.LstdFlags),
		debugLogger: log.New(logFile, prefix("DEBUG"), log.LstdFlags),
		warnLogger:  log.New(logFile, prefix("WARN"), log.LstdFlags),
	}
}

// Info 记录信息级别日志
func (l *Logger) Info(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	l.infoLogger.Println(message)
	fmt.Printf("[INFO] %s\n", message)
}

// Error 记录错误级别日志
func (l *Logger) Error(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	l.errorLogger.Println(message)
	fmt.Printf("[ERROR] %s\n", message)
}

// Debug 记录调试级别日志
func (l *Logger) Debug(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	l.debugLogger.Println(message)
	fmt.Printf("[DEBUG] %s\n", message)
}

// Warn 记录警告级别日志
func (l *Logger) Warn(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	l.warnLogger.Println(message)
	fmt.Printf("[WARN] %s\n", message)
}

// 全局日志实例
var logger = NewLogger()

// 全局日志函数
func Info(format string, v ...interface{}) {
	logger.Info(format, v...)
}

func Error(format string, v ...interface{}) {
	logger.Error(format, v...)
}

func Debug(format string, v ...interface{}) {
	logger.Debug(format, v...)
}

func Warn(format string, v ...interface{}) {
	logger.Warn(format, v...)
}

// 字符串截断（用于文章列表显示摘要）
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// uint转字符串
func UintToString(u uint) string {
	return string(rune(u))
}
