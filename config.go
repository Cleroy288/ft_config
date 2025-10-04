package ft_config

// ConfigMapping defines the mapping between environment variables and config field names.
// Key: The name you want to use in your code (e.g., "SupabaseKey")
// Value: The environment variable name in your .env file (e.g., "SUPABASE_ANON_KEY")
type ConfigMapping map[string]string

// Config represents your application configuration.
// This struct is dynamically populated based on the ConfigMapping you provide.
// You can access values in two ways:
//   1. Direct field access: cfg.SupabaseKey
//   2. Dynamic access: cfg.Get("SupabaseKey")
type Config struct {
	values map[string]string
}

// Get retrieves a configuration value by key.
// Returns the value and nil if found, empty string and ErrKeyNotFound if not found.
func (c *Config) Get(key string) (string, error) {
	if key == "" {
		return "", ErrEmptyKey
	}

	value, exists := c.values[key]
	if !exists {
		return "", ErrKeyNotFound
	}

	return value, nil
}

// GetOrDefault retrieves a value by key, returns default if not found.
func (c *Config) GetOrDefault(key, defaultValue string) string {
	value, exists := c.values[key]
	if !exists {
		return defaultValue
	}
	return value
}

// Has checks if a key exists in the configuration.
func (c *Config) Has(key string) bool {
	_, exists := c.values[key]
	return exists
}

// GetAll returns all configuration key-value pairs.
func (c *Config) GetAll() map[string]string {
	result := make(map[string]string, len(c.values))
	for k, v := range c.values {
		result[k] = v
	}
	return result
}

// GetValue is a helper method for accessing specific configured values.
// This is used internally by accessor methods.
func (c *Config) GetValue(key string) string {
	return c.values[key]
}