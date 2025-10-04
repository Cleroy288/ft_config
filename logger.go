package ft_config

import (
	"fmt"
	"log"
)

// logInfof logs a formatted informational message with context.
// All logs are prefixed with [ft_config] for easy filtering.
// Logging is always enabled and uses the standard log package.
func logInfof(context, format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	log.Printf("[ft_config] [%s] %s", context, message)
}

// logError logs an error message with context.
// Errors are prefixed with ERROR: for visibility.
func logError(context string, err error) {
	log.Printf("[ft_config] [%s] ERROR: %v", context, err)
}
