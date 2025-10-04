package ft_config

import (
	"os"
	"sync"

	"github.com/joho/godotenv"
)

// Service manages configuration loading and access.
type Service struct {
	mapping ConfigMapping
	config  *Config
	mu      sync.RWMutex
}

// New creates a new configuration service with the provided mapping.
//
// The mapping defines which environment variables to load and what to call them.
// Example:
//
//	mapping := ft_config.ConfigMapping{
//	    "SupabaseKey": "SUPABASE_ANON_KEY",
//	    "DatabaseURL": "DATABASE_URL",
//	    "Port":        "PORT",
//	}
//	service := ft_config.New(mapping)
func New(mapping ConfigMapping) *Service {
	logInfof("New", "Creating new config service with %d mappings", len(mapping))

	return &Service{
		mapping: mapping,
		config: &Config{
			values: make(map[string]string),
		},
	}
}

// Load loads environment variables from a .env file.
// Only variables defined in the ConfigMapping will be loaded.
//
// Example:
//
//	err := service.Load(".env")
//	if err != nil {
//	    log.Fatal(err)
//	}
func (s *Service) Load(filePath string) error {
	logInfof("Load", "Loading environment variables from: %s", filePath)

	s.mu.Lock()
	defer s.mu.Unlock()

	// Load .env file
	err := godotenv.Load(filePath)
	if err != nil {
		logError("Load", err)
		return ErrLoadEnv
	}

	// Load mapped environment variables
	loadedCount := 0
	for configKey, envKey := range s.mapping {
		if envValue, exists := os.LookupEnv(envKey); exists {
			s.config.values[configKey] = envValue
			loadedCount++
			logInfof("Load", "Loaded %s -> %s", envKey, configKey)
		} else {
			logInfof("Load", "Environment variable not found: %s", envKey)
		}
	}

	logInfof("Load", "Successfully loaded %d/%d configuration values", loadedCount, len(s.mapping))
	return nil
}

// Get retrieves a configuration value by key.
//
// Example:
//
//	apiKey, err := service.Get("SupabaseKey")
//	if err != nil {
//	    log.Fatal(err)
//	}
func (s *Service) Get(key string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.config.Get(key)
}

// GetOrDefault retrieves a value by key, returns default if not found.
//
// Example:
//
//	port := service.GetOrDefault("Port", "8080")
func (s *Service) GetOrDefault(key, defaultValue string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.config.GetOrDefault(key, defaultValue)
}

// Has checks if a key exists in the configuration.
//
// Example:
//
//	if service.Has("ApiToken") {
//	    token, _ := service.Get("ApiToken")
//	}
func (s *Service) Has(key string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.config.Has(key)
}

// GetAll returns all configuration key-value pairs.
//
// Example:
//
//	allConfig := service.GetAll()
//	for key, value := range allConfig {
//	    fmt.Printf("%s = %s\n", key, value)
//	}
func (s *Service) GetAll() map[string]string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.config.GetAll()
}

// Set manually sets a configuration value.
// This is useful for testing or overriding loaded values.
//
// Example:
//
//	service.Set("Port", "9090")
func (s *Service) Set(key, value string) error {
	if key == "" {
		return ErrEmptyKey
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.config.values[key] = value
	logInfof("Set", "Set %s = %s", key, value)
	return nil
}

// Delete removes a configuration key.
//
// Example:
//
//	service.Delete("TempKey")
func (s *Service) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.config.values, key)
	logInfof("Delete", "Deleted key: %s", key)
}

// Clear removes all configuration values.
//
// Example:
//
//	service.Clear()
func (s *Service) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.config.values = make(map[string]string)
	logInfo("Clear", "Cleared all configuration values")
}

// Config returns the Config struct for direct field access.
// This allows you to access values like: cfg.Get("SupabaseKey")
//
// Example:
//
//	cfg := service.Config()
//	apiKey, _ := cfg.Get("SupabaseKey")
//	port := cfg.GetOrDefault("Port", "8080")
func (s *Service) Config() *Config {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.config
}