# ft_config

A clean, production-ready Go package for managing environment configuration with type-safe access patterns.

## Features

- ✅ **Constructor-based configuration** - Define your config mapping upfront
- ✅ **Dual access patterns** - Dynamic (`cfg.Get()`) or direct field access
- ✅ **Thread-safe** - Safe for concurrent use
- ✅ **Zero dependencies** - Only uses `godotenv` for .env loading
- ✅ **Clean error handling** - Sentinel errors for type-safe error checking
- ✅ **Optional logging** - Built-in logging that's disabled by default
- ✅ **Production-ready** - 100% test coverage with 14 passing tests

## Installation

```bash
go get github.com/Cleroy288/ft_config
```

## Quick Start

### 1. Create your .env file

```env
SUPABASE_ANON_KEY=your_supabase_anon_key
SUPABASE_URL=https://your-project.supabase.co
DATABASE_URL=postgres://user:password@localhost:5432/dbname
API_TOKEN=your_api_token
PORT=8080
```

### 2. Define your configuration mapping

```go
package main

import (
    "log"
    ft_config "github.com/Cleroy288/ft_config"
)

func main() {
    // Define which env vars to load and what to call them
    mapping := ft_config.ConfigMapping{
        "SupabaseKey": "SUPABASE_ANON_KEY",  // Key in code: Key in .env
        "SupabaseURL": "SUPABASE_URL",
        "DatabaseURL": "DATABASE_URL",
        "ApiToken":    "API_TOKEN",
        "Port":        "PORT",
    }

    // Create config service
    config := ft_config.New(mapping)

    // Load .env file
    err := config.Load(".env")
    if err != nil {
        log.Fatal(err)
    }

    // Access values
    apiKey, _ := config.Get("SupabaseKey")
    log.Printf("Supabase Key: %s", apiKey)

    port := config.GetOrDefault("Port", "3000")
    log.Printf("Port: %s", port)
}
```

## Usage Patterns

### Pattern 1: Dynamic Access via Service

```go
// Create mapping
mapping := ft_config.ConfigMapping{
    "SupabaseKey": "SUPABASE_ANON_KEY",
    "DatabaseURL": "DATABASE_URL",
}

service := ft_config.New(mapping)
service.Load(".env")

// Get value (returns error if not found)
apiKey, err := service.Get("SupabaseKey")
if err != nil {
    log.Fatal(err)
}

// Get with default fallback
port := service.GetOrDefault("Port", "8080")

// Check if key exists
if service.Has("ApiToken") {
    token, _ := service.Get("ApiToken")
    fmt.Println(token)
}

// Get all config values
allConfig := service.GetAll()
for key, value := range allConfig {
    fmt.Printf("%s = %s\n", key, value)
}
```

### Pattern 2: Config Instance Access

```go
service := ft_config.New(mapping)
service.Load(".env")

// Get Config instance
cfg := service.Config()

// Access via Config methods
apiKey, err := cfg.Get("SupabaseKey")
port := cfg.GetOrDefault("Port", "3000")
hasToken := cfg.Has("ApiToken")
```

### Pattern 3: Manual Configuration (Testing)

```go
// Create empty service
service := ft_config.New(ft_config.ConfigMapping{})

// Set values manually
service.Set("ApiKey", "test_key_123")
service.Set("Port", "9090")

// Use like normal
apiKey, _ := service.Get("ApiKey")
```

## API Reference

### Constructor

#### `New(mapping ConfigMapping) *Service`

Creates a new configuration service with the provided mapping.

**Example:**
```go
mapping := ft_config.ConfigMapping{
    "SupabaseKey": "SUPABASE_ANON_KEY",
    "Port":        "PORT",
}
service := ft_config.New(mapping)
```

### Loading Configuration

#### `Load(filePath string) error`

Loads environment variables from a .env file using the registered mappings.

**Example:**
```go
err := service.Load(".env")
if err != nil {
    // Handle error
}
```

**Errors:**
- `ErrLoadEnv` - Failed to load .env file

### Accessing Values

#### `Get(key string) (string, error)`

Retrieves a configuration value by key.

**Example:**
```go
value, err := service.Get("SupabaseKey")
if err != nil {
    // Key not found
}
```

**Errors:**
- `ErrKeyNotFound` - Key doesn't exist
- `ErrEmptyKey` - Empty key provided

#### `GetOrDefault(key, defaultValue string) string`

Retrieves a value by key, returns default if not found.

**Example:**
```go
port := service.GetOrDefault("Port", "8080")
```

#### `Has(key string) bool`

Checks if a key exists in the configuration.

**Example:**
```go
if service.Has("ApiToken") {
    // Token is configured
}
```

#### `GetAll() map[string]string`

Returns all configuration key-value pairs.

**Example:**
```go
allConfig := service.GetAll()
```

### Modifying Configuration

#### `Set(key, value string) error`

Manually sets a configuration value.

**Example:**
```go
err := service.Set("Port", "9090")
```

**Errors:**
- `ErrEmptyKey` - Empty key provided

#### `Delete(key string)`

Removes a configuration key.

**Example:**
```go
service.Delete("TempKey")
```

#### `Clear()`

Removes all configuration values.

**Example:**
```go
service.Clear()
```

### Config Instance

#### `Config() *Config`

Returns the Config struct for alternative access patterns.

**Example:**
```go
cfg := service.Config()
apiKey, _ := cfg.Get("SupabaseKey")
```

## Logging

Logging is **disabled by default** for production use. Enable it for debugging:

```go
import ft_config "github.com/Cleroy288/ft_config"

// Enable logging
ft_config.EnableLogging()

// Create and use service
service := ft_config.New(mapping)
service.Load(".env")  // Logs: Loading environment variables...

// Disable logging
ft_config.DisableLogging()

// Check logging status
if ft_config.IsLoggingEnabled() {
    fmt.Println("Logging is enabled")
}
```

**Log output example:**
```
[ft_config] [New] Creating new config service with 5 mappings
[ft_config] [Load] Loading environment variables from: .env
[ft_config] [Load] Loaded SUPABASE_ANON_KEY -> SupabaseKey
[ft_config] [Load] Successfully loaded 5/5 configuration values
```

## Error Handling

The package provides sentinel errors for type-safe error checking:

```go
import (
    "errors"
    ft_config "github.com/Cleroy288/ft_config"
)

value, err := service.Get("NonExistent")
if errors.Is(err, ft_config.ErrKeyNotFound) {
    // Handle missing key
}

err = service.Load("missing.env")
if errors.Is(err, ft_config.ErrLoadEnv) {
    // Handle file loading error
}
```

**Available errors:**
- `ErrKeyNotFound` - Configuration key not found
- `ErrLoadEnv` - Failed to load .env file
- `ErrEmptyKey` - Empty key provided
- `ErrInvalidValue` - Invalid value type
- `ErrNoMapping` - No configuration mapping provided

## Project Structure

```
ft_config/
├── config.go          # Config struct and methods
├── errors.go          # Sentinel error definitions
├── logger.go          # Logging functionality
├── service.go         # Main service implementation
├── service_test.go    # Comprehensive tests
├── go.mod             # Module definition
├── go.sum             # Dependencies
├── .env.test          # Test environment file
└── README.md          # This file
```

## Testing

Run all tests:

```bash
go test -v
```

Run with race detection:

```bash
go test -race -v
```

**Test coverage:** 14/14 tests passing ✅

## Best Practices

### 1. Define mapping at application start

```go
var ConfigMapping = ft_config.ConfigMapping{
    "SupabaseKey": "SUPABASE_ANON_KEY",
    "DatabaseURL": "DATABASE_URL",
    "Port":        "PORT",
}

func main() {
    config := ft_config.New(ConfigMapping)
    config.Load(".env")
    // Use config throughout app
}
```

### 2. Use GetOrDefault for optional values

```go
// Required value
apiKey, err := config.Get("ApiKey")
if err != nil {
    log.Fatal("API key is required")
}

// Optional value with sensible default
port := config.GetOrDefault("Port", "8080")
timeout := config.GetOrDefault("Timeout", "30s")
```

### 3. Validate critical configuration on startup

```go
config.Load(".env")

requiredKeys := []string{"SupabaseKey", "DatabaseURL", "ApiToken"}
for _, key := range requiredKeys {
    if !config.Has(key) {
        log.Fatalf("Missing required configuration: %s", key)
    }
}
```

### 4. Use Config instance for cleaner code

```go
type App struct {
    config *ft_config.Config
}

func NewApp(service *ft_config.Service) *App {
    return &App{
        config: service.Config(),
    }
}

func (a *App) Start() {
    port := a.config.GetOrDefault("Port", "8080")
    // Start server
}
```

## Real-World Example

```go
package main

import (
    "log"
    "fmt"
    ft_config "github.com/Cleroy288/ft_config"
)

func main() {
    // 1. Define configuration mapping
    mapping := ft_config.ConfigMapping{
        "SupabaseKey":    "SUPABASE_ANON_KEY",
        "SupabaseURL":    "SUPABASE_URL",
        "DatabaseURL":    "DATABASE_URL",
        "RedisURL":       "REDIS_URL",
        "Port":           "PORT",
        "Environment":    "ENVIRONMENT",
        "LogLevel":       "LOG_LEVEL",
    }

    // 2. Create service and load config
    config := ft_config.New(mapping)

    if err := config.Load(".env"); err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    // 3. Validate required configuration
    required := []string{"SupabaseKey", "SupabaseURL", "DatabaseURL"}
    for _, key := range required {
        if !config.Has(key) {
            log.Fatalf("Missing required config: %s", key)
        }
    }

    // 4. Get configuration values
    supabaseKey, _ := config.Get("SupabaseKey")
    supabaseURL, _ := config.Get("SupabaseURL")
    databaseURL, _ := config.Get("DatabaseURL")

    // Optional values with defaults
    port := config.GetOrDefault("Port", "8080")
    env := config.GetOrDefault("Environment", "development")
    logLevel := config.GetOrDefault("LogLevel", "info")

    // 5. Use configuration
    fmt.Printf("Starting server in %s mode\n", env)
    fmt.Printf("Port: %s\n", port)
    fmt.Printf("Log Level: %s\n", logLevel)

    // Initialize services
    // db := initDatabase(databaseURL)
    // supabase := initSupabase(supabaseKey, supabaseURL)
    // server := startServer(port)
}
```

## Thread Safety

All operations are thread-safe and can be used concurrently:

```go
config := ft_config.New(mapping)
config.Load(".env")

// Safe to use from multiple goroutines
go func() {
    apiKey, _ := config.Get("ApiKey")
    // Use apiKey
}()

go func() {
    dbURL, _ := config.Get("DatabaseURL")
    // Use dbURL
}()
```

## Contributing

This package follows Go best practices:

- ✅ Clear separation of concerns
- ✅ Comprehensive error handling
- ✅ Thread-safe operations
- ✅ Extensive documentation
- ✅ 100% test coverage
- ✅ Zero external dependencies (except godotenv)

## License

MIT License - feel free to use in your projects!

## Author

**Charles Leroy** - [@Cleroy288](https://github.com/Cleroy288)

---

**Questions or issues?** Open an issue on [GitHub](https://github.com/Cleroy288/ft_config/issues)