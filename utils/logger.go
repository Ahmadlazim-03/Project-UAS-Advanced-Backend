package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
)

// LogLevel represents the severity of a log message
type LogLevel string

const (
	LogLevelInfo  LogLevel = "INFO"
	LogLevelWarn  LogLevel = "WARN"
	LogLevelError LogLevel = "ERROR"
	LogLevelDebug LogLevel = "DEBUG"
)

// Logger provides structured logging functionality
type Logger struct {
	AppName string
}

// NewLogger creates a new logger instance
func NewLogger(appName string) *Logger {
	return &Logger{
		AppName: appName,
	}
}

// log writes a structured log message
func (l *Logger) log(level LogLevel, message string, context map[string]interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	
	logMessage := fmt.Sprintf("[%s] [%s] [%s] %s", 
		timestamp, 
		l.AppName, 
		level, 
		message,
	)

	// Add context if provided
	if len(context) > 0 {
		logMessage += fmt.Sprintf(" | Context: %v", context)
	}

	// Write to stdout (in production, consider using a logging library like logrus or zap)
	fmt.Fprintln(os.Stdout, logMessage)
}

// Info logs an informational message
func (l *Logger) Info(message string, context ...map[string]interface{}) {
	ctx := make(map[string]interface{})
	if len(context) > 0 {
		ctx = context[0]
	}
	l.log(LogLevelInfo, message, ctx)
}

// Warn logs a warning message
func (l *Logger) Warn(message string, context ...map[string]interface{}) {
	ctx := make(map[string]interface{})
	if len(context) > 0 {
		ctx = context[0]
	}
	l.log(LogLevelWarn, message, ctx)
}

// Error logs an error message
func (l *Logger) Error(message string, err error, context ...map[string]interface{}) {
	ctx := make(map[string]interface{})
	if len(context) > 0 {
		ctx = context[0]
	}
	if err != nil {
		ctx["error"] = err.Error()
	}
	l.log(LogLevelError, message, ctx)
}

// Debug logs a debug message
func (l *Logger) Debug(message string, context ...map[string]interface{}) {
	ctx := make(map[string]interface{})
	if len(context) > 0 {
		ctx = context[0]
	}
	l.log(LogLevelDebug, message, ctx)
}

// LogRequest logs HTTP request details
func (l *Logger) LogRequest(c *fiber.Ctx, duration time.Duration) {
	context := map[string]interface{}{
		"method":     c.Method(),
		"path":       c.Path(),
		"status":     c.Response().StatusCode(),
		"duration":   fmt.Sprintf("%v", duration),
		"ip":         c.IP(),
		"user_agent": c.Get("User-Agent"),
	}
	
	message := fmt.Sprintf("%s %s - %d", c.Method(), c.Path(), c.Response().StatusCode())
	l.log(LogLevelInfo, message, context)
}

// LogAuth logs authentication-related events
func (l *Logger) LogAuth(event string, userID string, success bool, context ...map[string]interface{}) {
	ctx := make(map[string]interface{})
	if len(context) > 0 {
		ctx = context[0]
	}
	
	ctx["event"] = event
	ctx["user_id"] = userID
	ctx["success"] = success
	
	level := LogLevelInfo
	if !success {
		level = LogLevelWarn
	}
	
	l.log(level, fmt.Sprintf("Auth event: %s", event), ctx)
}

// LogDBOperation logs database operations
func (l *Logger) LogDBOperation(operation string, table string, duration time.Duration, err error) {
	context := map[string]interface{}{
		"operation": operation,
		"table":     table,
		"duration":  fmt.Sprintf("%v", duration),
	}
	
	if err != nil {
		context["error"] = err.Error()
		l.log(LogLevelError, fmt.Sprintf("DB %s failed on %s", operation, table), context)
	} else {
		l.log(LogLevelDebug, fmt.Sprintf("DB %s on %s", operation, table), context)
	}
}

// Global logger instance
var GlobalLogger = NewLogger("StudentAchievementSystem")
