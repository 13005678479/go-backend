package middleware

import (
	"blogV2/config"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger 结构体封装日志功能
type Logger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
	debugLogger *log.Logger
	warnLogger  *log.Logger
	config      *config.LogConfig
}

// NewLogger 创建新的日志实例
func NewLogger(cfg *config.LogConfig) *Logger {
	// 创建日志文件
	logFile, err := os.OpenFile(cfg.LogFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
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
		config:      cfg,
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

// LoggerMiddleware Gin中间件，记录HTTP请求日志
func LoggerMiddleware(logger *Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录请求开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 计算请求处理时间
		latency := time.Since(startTime)

		// 获取响应状态码
		statusCode := c.Writer.Status()

		// 记录请求日志
		logger.Info("HTTP请求: %s %s | 状态码: %d | 耗时: %v | 客户端IP: %s",
			c.Request.Method,
			c.Request.URL.Path,
			statusCode,
			latency,
			c.ClientIP(),
		)
	}
}

// 全局日志实例
var globalLogger *Logger

// InitLogger 初始化全局日志实例
func InitLogger(cfg *config.LogConfig) {
	globalLogger = NewLogger(cfg)
}

// GetLogger 获取全局日志实例
func GetLogger() *Logger {
	if globalLogger == nil {
		panic("日志器未初始化，请先调用InitLogger")
	}
	return globalLogger
}

// 全局日志函数（兼容现有代码）
func Info(format string, v ...interface{}) {
	GetLogger().Info(format, v...)
}

func Error(format string, v ...interface{}) {
	GetLogger().Error(format, v...)
}

func Debug(format string, v ...interface{}) {
	GetLogger().Debug(format, v...)
}

func Warn(format string, v ...interface{}) {
	GetLogger().Warn(format, v...)
}