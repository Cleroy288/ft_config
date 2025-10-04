package ft_config

import "errors"

// Sentinel errors for configuration operations.
var (
	// ErrKeyNotFound is returned when a configuration key is not found.
	ErrKeyNotFound = errors.New("configuration key not found")

	// ErrLoadEnv is returned when the .env file cannot be loaded.
	ErrLoadEnv = errors.New("failed to load .env file")

	// ErrEmptyKey is returned when an empty key is provided.
	ErrEmptyKey = errors.New("key cannot be empty")

	// ErrInvalidValue is returned when a value has an invalid type.
	ErrInvalidValue = errors.New("invalid value type")

	// ErrNoMapping is returned when no configuration mapping is provided.
	ErrNoMapping = errors.New("no configuration mapping provided")
)