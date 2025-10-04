package ft_config

import (
	"errors"
	"os"
	"testing"
)

// Test struct matching .env.test values
type TestConfig struct {
	SupabaseKey string `env:"SUPABASE_KEY"`
	SupabaseURL string `env:"SUPABASE_URL"`
	DatabaseURL string `env:"DATABASE_URL"`
	ApiToken    string `env:"API_TOKEN"`
	Port        string `env:"PORT"`
	AppName     string `env:"APP_NAME"`
}

func TestLoad(t *testing.T) {
	var config TestConfig

	err := Load(".env.test", &config)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	// Verify all values were loaded correctly
	if config.SupabaseKey != "test_supabase_key_123" {
		t.Errorf("SupabaseKey: expected 'test_supabase_key_123', got '%s'", config.SupabaseKey)
	}

	if config.SupabaseURL != "https://test.supabase.co" {
		t.Errorf("SupabaseURL: expected 'https://test.supabase.co', got '%s'", config.SupabaseURL)
	}

	if config.DatabaseURL != "postgres://localhost:5432/testdb" {
		t.Errorf("DatabaseURL: expected 'postgres://localhost:5432/testdb', got '%s'", config.DatabaseURL)
	}

	if config.ApiToken != "test_api_token_456" {
		t.Errorf("ApiToken: expected 'test_api_token_456', got '%s'", config.ApiToken)
	}

	if config.Port != "8080" {
		t.Errorf("Port: expected '8080', got '%s'", config.Port)
	}

	if config.AppName != "ft_config_test" {
		t.Errorf("AppName: expected 'ft_config_test', got '%s'", config.AppName)
	}

	t.Log("✓ Load() successfully fills struct from .env file")
}

func TestLoadNonExistentFile(t *testing.T) {
	var config TestConfig

	err := Load("nonexistent.env", &config)
	if err == nil {
		t.Error("Expected error when loading non-existent file")
	}

	if !errors.Is(err, ErrLoadEnv) {
		t.Errorf("Expected ErrLoadEnv, got: %v", err)
	}

	t.Log("✓ Load() returns error for non-existent file")
}

func TestLoadWithPartialFields(t *testing.T) {
	// Struct with only some fields
	type PartialConfig struct {
		SupabaseKey string `env:"SUPABASE_KEY"`
		Port        string `env:"PORT"`
	}

	var config PartialConfig

	err := Load(".env.test", &config)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if config.SupabaseKey != "test_supabase_key_123" {
		t.Errorf("SupabaseKey: expected 'test_supabase_key_123', got '%s'", config.SupabaseKey)
	}

	if config.Port != "8080" {
		t.Errorf("Port: expected '8080', got '%s'", config.Port)
	}

	t.Log("✓ Load() works with partial struct fields")
}

func TestLoadWithoutEnvTag(t *testing.T) {
	// Struct with fields without env tags
	type ConfigWithoutTags struct {
		SupabaseKey string `env:"SUPABASE_KEY"`
		IgnoredField string // No env tag - should be ignored
	}

	var config ConfigWithoutTags

	err := Load(".env.test", &config)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if config.SupabaseKey != "test_supabase_key_123" {
		t.Errorf("SupabaseKey: expected 'test_supabase_key_123', got '%s'", config.SupabaseKey)
	}

	if config.IgnoredField != "" {
		t.Errorf("IgnoredField should be empty, got '%s'", config.IgnoredField)
	}

	t.Log("✓ Load() ignores fields without env tag")
}

func TestLoadNotAPointer(t *testing.T) {
	var config TestConfig

	// Pass struct by value instead of pointer
	err := Load(".env.test", config)
	if err == nil {
		t.Error("Expected error when passing struct by value")
	}

	t.Log("✓ Load() returns error for non-pointer argument")
}

func TestLoadWithMissingEnvVar(t *testing.T) {
	type ConfigWithMissing struct {
		SupabaseKey    string `env:"SUPABASE_KEY"`
		NonExistentVar string `env:"NON_EXISTENT_VAR"`
	}

	var config ConfigWithMissing

	err := Load(".env.test", &config)
	if err == nil {
		t.Error("Expected error when env var is missing")
	}

	expectedErrMsg := "missing required environment variables: NON_EXISTENT_VAR"
	if err.Error() != expectedErrMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrMsg, err.Error())
	}

	t.Log("✓ Load() returns error for missing environment variables")
}

func TestRealWorldExample(t *testing.T) {
	// Real-world usage example
	type AppConfig struct {
		SupabaseKey string `env:"SUPABASE_KEY"`
		SupabaseURL string `env:"SUPABASE_URL"`
		DatabaseURL string `env:"DATABASE_URL"`
		Port        string `env:"PORT"`
	}

	var config AppConfig

	// One simple call!
	err := Load(".env.test", &config)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	// Now use the config directly
	if config.SupabaseKey == "" {
		t.Error("SupabaseKey should not be empty")
	}

	if config.Port == "" {
		t.Error("Port should not be empty")
	}

	// Config is ready to use!
	t.Logf("✓ Config loaded: Port=%s, Database=%s", config.Port, config.DatabaseURL)
}

// Cleanup after tests
func TestMain(m *testing.M) {
	// Run tests
	code := m.Run()

	os.Exit(code)
}