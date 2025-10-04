package ft_config

import (
	"fmt"
	"log"
	"sync"
)

// Logger manages logging for the ft_config package.
type Logger struct {
	enabled bool
	mu      sync.RWMutex
}

// globalLogger is the package-level logger instance.
var globalLogger = &Logger{
	enabled: false, // Disabled by default for production use
}

// EnableLogging enables logging for the ft_config package.
func EnableLogging() {
	globalLogger.mu.Lock()
	defer globalLogger.mu.Unlock()
	globalLogger.enabled = true
}

// DisableLogging disables logging for the ft_config package.
func DisableLogging() {
	globalLogger.mu.Lock()
	defer globalLogger.mu.Unlock()
	globalLogger.enabled = false
}

// IsLoggingEnabled returns whether logging is currently enabled.
func IsLoggingEnabled() bool {
	globalLogger.mu.RLock()
	defer globalLogger.mu.RUnlock()
	return globalLogger.enabled
}

// logInfo logs an informational message with context.
func logInfo(context, message string) {
	globalLogger.mu.RLock()
	enabled := globalLogger.enabled
	globalLogger.mu.RUnlock()

	if !enabled {
		return
	}

	log.Printf("[ft_config] [%s] %s", context, message)
}

// logInfof logs a formatted informational message with context.
func logInfof(context, format string, args ...any) {
	globalLogger.mu.RLock()
	enabled := globalLogger.enabled
	globalLogger.mu.RUnlock()

	if !enabled {
		return
	}

	message := fmt.Sprintf(format, args...)
	log.Printf("[ft_config] [%s] %s", context, message)
}

// logError logs an error message with context.
func logError(context string, err error) {
	globalLogger.mu.RLock()
	enabled := globalLogger.enabled
	globalLogger.mu.RUnlock()

	if !enabled {
		return
	}

	log.Printf("[ft_config] [%s] ERROR: %v", context, err)
}