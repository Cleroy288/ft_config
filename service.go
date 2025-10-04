// Package ft_config provides a simple, reflection-based configuration loader
// that reads .env files and populates Go structs using struct tags.
//
// This package emphasizes simplicity and ease of use - just one function call
// to load your entire configuration with type-safe struct access and IDE autocomplete.
package ft_config

import (
	"errors"
	"os"
	"reflect"

	"github.com/joho/godotenv"
)

// Load loads environment variables from a .env file and fills the provided struct.
// If envFile is an empty string, it loads from OS environment variables only.
//
// Usage:
//
//	type MyConfig struct {
//	    SupabaseKey string `env:"SUPABASE_ANON_KEY"`
//	    DatabaseURL string `env:"DATABASE_URL"`
//	    Port        string `env:"PORT"`
//	}
//
//	var config MyConfig
//	err := ft_config.Load(".env", &config)
//
// The struct fields must have an `env` tag with the environment variable name.
// Returns an error if any required environment variable is not found.
func Load(envFile string, configStruct any) error {
	// Print initialization details
	if envFile != "" {
		logInfof("Load", "[ft_config] Initialized - loading from file: %s", envFile)
	} else {
		logInfof("Load", "[ft_config] Initialized - loading from OS environment variables")
	}

	// Load .env file only if path is provided
	if envFile != "" {
		if err := godotenv.Load(envFile); err != nil {
			logError("Load", err)
			return errors.Join(ErrLoadEnv, err)
		}
	}

	// Use reflection to fill struct fields
	v := reflect.ValueOf(configStruct)
	if v.Kind() != reflect.Ptr {
		return errors.New("config must be a pointer to a struct")
	}

	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return errors.New("config must be a pointer to a struct")
	}

	t := v.Type()
	var missingVars []string

	// Iterate through struct fields
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		// Get env tag
		envKey := fieldType.Tag.Get("env")
		if envKey == "" {
			continue // Skip fields without env tag
		}

		// Check if field is settable
		if !field.CanSet() {
			logInfof("Load", "Field %s is not settable", fieldType.Name)
			continue
		}

		// Get value from environment
		envValue, exists := os.LookupEnv(envKey)
		if !exists {
			missingVars = append(missingVars, envKey)
			continue
		}

		// Set the field value
		if field.Kind() == reflect.String {
			field.SetString(envValue)
			logInfof("Load", "Loaded %s -> %s", envKey, fieldType.Name)
		} else {
			logInfof("Load", "Unsupported field type for %s: %s", fieldType.Name, field.Kind())
		}
	}

	// Return error if any variables are missing
	if len(missingVars) > 0 {
		err := errors.New("missing required environment variables: " + missingVars[0])
		for i := 1; i < len(missingVars); i++ {
			err = errors.New(err.Error() + ", " + missingVars[i])
		}
		logError("Load", err)
		return err
	}

	logInfof("Load", "Successfully loaded all configuration values")
	return nil
}