# ft_config

**The simplest Go configuration package.** One function. One struct. Done.

## Features

- ✅ **Ultra-simple** - One function call to load your entire config
- ✅ **Struct-based** - Define your config as a Go struct with tags
- ✅ **Type-safe** - Your IDE autocompletes field names
- ✅ **Zero boilerplate** - No mapping, no registration, just load
- ✅ **OS environment support** - Load from .env file or OS environment variables
- ✅ **Required field validation** - Automatic error if any env var is missing
- ✅ **Thread-safe** - Safe for concurrent use
- ✅ **Clean errors** - Sentinel errors for type-safe error handling
- ✅ **Built-in logging** - Simple logging for debugging
- ✅ **Production-ready** - 100% test coverage with 7 passing tests

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

### 2. Define your config struct

```go
package main

import (
    "log"
    ft_config "github.com/Cleroy288/ft_config"
)

// Define your config struct with env tags
type Config struct {
    SupabaseKey string `env:"SUPABASE_ANON_KEY"`
    SupabaseURL string `env:"SUPABASE_URL"`
    DatabaseURL string `env:"DATABASE_URL"`
    ApiToken    string `env:"API_TOKEN"`
    Port        string `env:"PORT"`
}

func main() {
    var config Config

    // One function call - that's it!
    err := ft_config.Load(".env", &config)
    if err != nil {
        log.Fatal(err)
    }

    // Use your config directly
    log.Printf("Starting on port: %s", config.Port)
    log.Printf("Database: %s", config.DatabaseURL)
}
```

That's it! No mapping, no service creation, no complexity. Just **one function call**.

## Why ft_config?

### Before (complicated):
```go
// Other packages require:
mapping := SomePackage.ConfigMapping{
    "SupabaseKey": "SUPABASE_ANON_KEY",
    "DatabaseURL": "DATABASE_URL",
    // ... repeat for every field
}
service := SomePackage.New(mapping)
service.Load(".env")
value, err := service.Get("SupabaseKey")
```

### With ft_config (simple):
```go
type Config struct {
    SupabaseKey string `env:"SUPABASE_ANON_KEY"`
    DatabaseURL string `env:"DATABASE_URL"`
}

var config Config
ft_config.Load(".env", &config)

// Direct access - your IDE autocompletes!
value := config.SupabaseKey
```

## Complete Example

```go
package main

import (
    "fmt"
    "log"
    ft_config "github.com/Cleroy288/ft_config"
)

type AppConfig struct {
    // Required fields
    SupabaseKey string `env:"SUPABASE_ANON_KEY"`
    SupabaseURL string `env:"SUPABASE_URL"`
    DatabaseURL string `env:"DATABASE_URL"`

    // Optional fields (check if empty)
    Port        string `env:"PORT"`
    Environment string `env:"ENVIRONMENT"`
    LogLevel    string `env:"LOG_LEVEL"`
}

func main() {
    var config AppConfig

    // Load config
    if err := ft_config.Load(".env", &config); err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    // Validate required fields
    if config.SupabaseKey == "" {
        log.Fatal("SUPABASE_ANON_KEY is required")
    }

    // Use defaults for optional fields
    if config.Port == "" {
        config.Port = "8080"
    }
    if config.Environment == "" {
        config.Environment = "development"
    }

    // Config is ready!
    fmt.Printf("🚀 Starting server in %s mode\n", config.Environment)
    fmt.Printf("📍 Port: %s\n", config.Port)
    fmt.Printf("💾 Database: %s\n", config.DatabaseURL)

    // Initialize services with config
    // db := initDatabase(config.DatabaseURL)
    // supabase := initSupabase(config.SupabaseKey, config.SupabaseURL)
    // server := startServer(config.Port)
}
```

## API Reference

### `Load(envFile string, configStruct any) error`

Loads environment variables from a `.env` file and fills the provided struct. If `envFile` is an empty string (`""`), it loads from OS environment variables only.

**Parameters:**
- `envFile` - Path to .env file (e.g., `".env"`, `"config/.env"`), or `""` for OS environment only
- `configStruct` - Pointer to your config struct with `env` tags

**Returns:**
- `error` - `ErrLoadEnv` if file cannot be loaded, or error listing all missing required variables

**Example:**
```go
type Config struct {
    ApiKey string `env:"API_KEY"`
    Port   string `env:"PORT"`
}

// Load from .env file
var config Config
err := ft_config.Load(".env", &config)

// Or load from OS environment variables only
err := ft_config.Load("", &config)
```

**Rules:**
- ✅ Struct fields must have an `env` tag to be loaded
- ✅ Fields without `env` tag are ignored
- ✅ Missing environment variables return an error listing all missing vars
- ✅ Only string fields are currently supported
- ✅ Must pass a pointer to struct (`&config`, not `config`)
- ✅ Pass empty string `""` as envFile to load from OS environment only

## Struct Tags

Use the `env` tag to map environment variables to struct fields:

```go
type Config struct {
    // Maps to SUPABASE_ANON_KEY in .env
    SupabaseKey string `env:"SUPABASE_ANON_KEY"`

    // Maps to DATABASE_URL in .env
    DatabaseURL string `env:"DATABASE_URL"`

    // No tag = ignored, won't be loaded
    ComputedValue string
}
```

## Error Handling

```go
import (
    "errors"
    ft_config "github.com/Cleroy288/ft_config"
)

var config Config
err := ft_config.Load(".env", &config)

if errors.Is(err, ft_config.ErrLoadEnv) {
    log.Fatal("Failed to load .env file")
}

if err != nil {
    log.Fatal(err)
}
```

**Available errors:**
- `ErrLoadEnv` - Failed to load .env file
- `ErrEmptyKey` - Empty key provided
- `ErrInvalidValue` - Invalid value type
- `ErrNoMapping` - No configuration mapping provided

## Logging

The package includes simple logging that outputs information about loaded configuration:

```
[ft_config] [Load] Loading configuration from: .env
[ft_config] [Load] Loaded SUPABASE_KEY -> SupabaseKey = xxx
[ft_config] [Load] Successfully loaded 5 configuration values
```

Logs are sent to the standard Go logger and can be controlled with `log.SetOutput()` if needed.

## Advanced Usage

### Load from OS Environment Variables

Skip the .env file and load directly from OS environment:

```go
type Config struct {
    ApiKey string `env:"API_KEY"`
    Port   string `env:"PORT"`
}

var config Config
// Pass empty string to load from OS environment only
err := ft_config.Load("", &config)
```

This is useful for:
- Docker containers with environment variables
- Kubernetes ConfigMaps and Secrets
- CI/CD pipelines
- Cloud platform environment variables (Heroku, AWS, etc.)

### Partial Configuration

Only load the fields you need:

```go
type MinimalConfig struct {
    ApiKey string `env:"API_KEY"`
    Port   string `env:"PORT"`
    // Other env vars in .env are ignored
}

var config MinimalConfig
ft_config.Load(".env", &config)
```

### Multiple Config Structs

Load different configs for different purposes:

```go
type DatabaseConfig struct {
    URL      string `env:"DATABASE_URL"`
    MaxConns string `env:"DB_MAX_CONNECTIONS"`
}

type APIConfig struct {
    Key     string `env:"API_KEY"`
    BaseURL string `env:"API_BASE_URL"`
}

var dbConfig DatabaseConfig
var apiConfig APIConfig

ft_config.Load(".env", &dbConfig)
ft_config.Load(".env", &apiConfig)
```

## Project Structure

```
ft_config/
├── service.go         # Main Load function
├── service_test.go    # Comprehensive tests
├── errors.go          # Sentinel error definitions
├── logger.go          # Simple logging
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

**Test coverage:** 7/7 tests passing ✅

## Best Practices

### 1. Define config struct at package level

```go
package main

import ft_config "github.com/Cleroy288/ft_config"

type Config struct {
    SupabaseKey string `env:"SUPABASE_ANON_KEY"`
    DatabaseURL string `env:"DATABASE_URL"`
    Port        string `env:"PORT"`
}

var AppConfig Config

func init() {
    if err := ft_config.Load(".env", &AppConfig); err != nil {
        panic(err)
    }
}
```

### 2. Validate required fields on startup

```go
func main() {
    var config Config
    ft_config.Load(".env", &config)

    // Validate
    if config.SupabaseKey == "" {
        log.Fatal("SUPABASE_ANON_KEY is required")
    }
}
```

### 3. Use defaults for optional values

```go
if config.Port == "" {
    config.Port = "8080"
}

if config.LogLevel == "" {
    config.LogLevel = "info"
}
```

### 4. Group related config

```go
type Config struct {
    // Database
    DatabaseURL      string `env:"DATABASE_URL"`
    DatabaseMaxConns string `env:"DB_MAX_CONNECTIONS"`

    // Supabase
    SupabaseKey string `env:"SUPABASE_ANON_KEY"`
    SupabaseURL string `env:"SUPABASE_URL"`

    // Server
    Port        string `env:"PORT"`
    Environment string `env:"ENVIRONMENT"`
}
```

## Comparison with Other Packages

| Feature | ft_config | viper | envconfig | godotenv |
|---------|-----------|-------|-----------|----------|
| Lines of code to use | **3** | 10+ | 5+ | 5+ |
| Mapping required | ❌ | ✅ | ❌ | ✅ |
| Type-safe access | ✅ | ❌ | ✅ | ❌ |
| IDE autocomplete | ✅ | ❌ | ✅ | ❌ |
| Dependencies | 1 | Many | 0 | 0 |
| Complexity | **Low** | High | Medium | Medium |

## Thread Safety

The `Load()` function is thread-safe and can be called from multiple goroutines. However, it's recommended to load your config once at startup.

## Contributing

This package follows Go best practices:

- ✅ Minimal API surface
- ✅ Clear separation of concerns
- ✅ Comprehensive error handling
- ✅ Extensive documentation
- ✅ 100% test coverage
- ✅ Zero complexity

## License

MIT License - feel free to use in your projects!

## Author

**Charles Leroy** - [@Cleroy288](https://github.com/Cleroy288)

---

**Simple is better than complex. Complex is better than complicated.**

Questions or issues? Open an issue on [GitHub](https://github.com/Cleroy288/ft_config/issues)