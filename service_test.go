package ft_config

import (
	"errors"
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	mapping := ConfigMapping{
		"SupabaseKey": "SUPABASE_KEY",
		"DatabaseURL": "DATABASE_URL",
	}

	service := New(mapping)

	if service == nil {
		t.Fatal("New() returned nil")
	}

	if service.config == nil {
		t.Error("service.config is nil")
	}

	if service.mapping == nil {
		t.Error("service.mapping is nil")
	}

	if len(service.mapping) != 2 {
		t.Errorf("Expected 2 mappings, got %d", len(service.mapping))
	}

	t.Log("✓ New() creates valid service instance")
}

func TestLoad(t *testing.T) {
	mapping := ConfigMapping{
		"SupabaseKey": "SUPABASE_KEY",
		"SupabaseURL": "SUPABASE_URL",
		"DatabaseURL": "DATABASE_URL",
		"ApiToken":    "API_TOKEN",
		"Port":        "PORT",
		"AppName":     "APP_NAME",
	}

	service := New(mapping)

	err := service.Load(".env.test")
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	// Verify values were loaded
	tests := []struct {
		key           string
		expectedValue string
	}{
		{"SupabaseKey", "test_supabase_key_123"},
		{"SupabaseURL", "https://test.supabase.co"},
		{"DatabaseURL", "postgres://localhost:5432/testdb"},
		{"ApiToken", "test_api_token_456"},
		{"Port", "8080"},
		{"AppName", "ft_config_test"},
	}

	for _, tt := range tests {
		value, err := service.Get(tt.key)
		if err != nil {
			t.Errorf("Failed to get '%s': %v", tt.key, err)
			continue
		}

		if value != tt.expectedValue {
			t.Errorf("Key '%s': expected '%s', got '%s'", tt.key, tt.expectedValue, value)
		}
	}

	t.Log("✓ Load() successfully loads environment variables")
}

func TestLoadNonExistentFile(t *testing.T) {
	mapping := ConfigMapping{
		"Key": "VALUE",
	}

	service := New(mapping)

	err := service.Load("nonexistent.env")
	if err == nil {
		t.Error("Expected error when loading non-existent file")
	}

	if !errors.Is(err, ErrLoadEnv) {
		t.Errorf("Expected ErrLoadEnv, got: %v", err)
	}

	t.Log("✓ Load() returns error for non-existent file")
}

func TestGet(t *testing.T) {
	mapping := ConfigMapping{
		"TestKey": "TEST_KEY",
	}

	service := New(mapping)
	service.Set("TestKey", "TestValue")

	// Test successful get
	value, err := service.Get("TestKey")
	if err != nil {
		t.Fatalf("Get() failed: %v", err)
	}

	if value != "TestValue" {
		t.Errorf("Expected 'TestValue', got '%s'", value)
	}

	// Test get non-existent key
	_, err = service.Get("NonExistent")
	if !errors.Is(err, ErrKeyNotFound) {
		t.Errorf("Expected ErrKeyNotFound, got: %v", err)
	}

	// Test get with empty key
	_, err = service.Get("")
	if !errors.Is(err, ErrEmptyKey) {
		t.Errorf("Expected ErrEmptyKey, got: %v", err)
	}

	t.Log("✓ Get() retrieves values correctly")
}

func TestSet(t *testing.T) {
	mapping := ConfigMapping{}
	service := New(mapping)

	// Test setting a value
	err := service.Set("Key1", "Value1")
	if err != nil {
		t.Fatalf("Set() failed: %v", err)
	}

	value, _ := service.Get("Key1")
	if value != "Value1" {
		t.Errorf("Expected 'Value1', got '%s'", value)
	}

	// Test updating a value
	err = service.Set("Key1", "UpdatedValue")
	if err != nil {
		t.Fatalf("Set() update failed: %v", err)
	}

	value, _ = service.Get("Key1")
	if value != "UpdatedValue" {
		t.Errorf("Expected 'UpdatedValue', got '%s'", value)
	}

	// Test setting empty key
	err = service.Set("", "value")
	if !errors.Is(err, ErrEmptyKey) {
		t.Errorf("Expected ErrEmptyKey, got: %v", err)
	}

	t.Log("✓ Set() stores and updates values correctly")
}

func TestGetOrDefault(t *testing.T) {
	mapping := ConfigMapping{}
	service := New(mapping)
	service.Set("ExistingKey", "ExistingValue")

	// Test get existing key
	value := service.GetOrDefault("ExistingKey", "DefaultValue")
	if value != "ExistingValue" {
		t.Errorf("Expected 'ExistingValue', got '%s'", value)
	}

	// Test get non-existent key (should return default)
	value = service.GetOrDefault("NonExistent", "DefaultValue")
	if value != "DefaultValue" {
		t.Errorf("Expected 'DefaultValue', got '%s'", value)
	}

	t.Log("✓ GetOrDefault() returns correct values")
}

func TestHas(t *testing.T) {
	mapping := ConfigMapping{}
	service := New(mapping)
	service.Set("ExistingKey", "value")

	// Test existing key
	if !service.Has("ExistingKey") {
		t.Error("Has() returned false for existing key")
	}

	// Test non-existent key
	if service.Has("NonExistent") {
		t.Error("Has() returned true for non-existent key")
	}

	t.Log("✓ Has() correctly checks key existence")
}

func TestDelete(t *testing.T) {
	mapping := ConfigMapping{}
	service := New(mapping)
	service.Set("Key1", "Value1")
	service.Set("Key2", "Value2")

	// Verify key exists
	if !service.Has("Key1") {
		t.Fatal("Key1 should exist before deletion")
	}

	// Delete key
	service.Delete("Key1")

	// Verify key is deleted
	if service.Has("Key1") {
		t.Error("Key1 should not exist after deletion")
	}

	// Verify other keys are not affected
	if !service.Has("Key2") {
		t.Error("Key2 should still exist after deleting Key1")
	}

	t.Log("✓ Delete() removes keys correctly")
}

func TestGetAll(t *testing.T) {
	mapping := ConfigMapping{}
	service := New(mapping)
	service.Set("Key1", "Value1")
	service.Set("Key2", "Value2")
	service.Set("Key3", "Value3")

	all := service.GetAll()

	if len(all) != 3 {
		t.Errorf("Expected 3 entries, got %d", len(all))
	}

	// Verify all values
	if all["Key1"] != "Value1" {
		t.Error("Key1 value mismatch")
	}
	if all["Key2"] != "Value2" {
		t.Error("Key2 value mismatch")
	}
	if all["Key3"] != "Value3" {
		t.Error("Key3 value mismatch")
	}

	t.Log("✓ GetAll() returns all entries")
}

func TestClear(t *testing.T) {
	mapping := ConfigMapping{}
	service := New(mapping)
	service.Set("Key1", "Value1")
	service.Set("Key2", "Value2")
	service.Set("Key3", "Value3")

	// Verify service has entries
	if len(service.GetAll()) != 3 {
		t.Fatal("Service should have 3 entries before clearing")
	}

	// Clear config
	service.Clear()

	// Verify config is empty
	if len(service.GetAll()) != 0 {
		t.Error("Service should be empty after Clear()")
	}

	// Verify keys are not accessible
	_, err := service.Get("Key1")
	if !errors.Is(err, ErrKeyNotFound) {
		t.Error("Keys should not be accessible after Clear()")
	}

	t.Log("✓ Clear() removes all entries")
}

func TestConfig(t *testing.T) {
	mapping := ConfigMapping{
		"SupabaseKey": "SUPABASE_KEY",
		"DatabaseURL": "DATABASE_URL",
		"Port":        "PORT",
	}

	service := New(mapping)
	service.Load(".env.test")

	// Get config
	cfg := service.Config()

	// Test config methods
	supabaseKey, err := cfg.Get("SupabaseKey")
	if err != nil {
		t.Fatalf("cfg.Get() failed: %v", err)
	}

	if supabaseKey != "test_supabase_key_123" {
		t.Errorf("Expected 'test_supabase_key_123', got '%s'", supabaseKey)
	}

	// Test GetOrDefault
	port := cfg.GetOrDefault("Port", "3000")
	if port != "8080" {
		t.Errorf("Expected '8080', got '%s'", port)
	}

	// Test Has
	if !cfg.Has("DatabaseURL") {
		t.Error("cfg.Has() returned false for existing key")
	}

	t.Log("✓ Config() returns working config instance")
}

func TestCompleteWorkflow(t *testing.T) {
	// Create mapping
	mapping := ConfigMapping{
		"SupabaseKey": "SUPABASE_KEY",
		"DatabaseURL": "DATABASE_URL",
		"Port":        "PORT",
	}

	// Create service
	service := New(mapping)

	// Load .env file
	err := service.Load(".env.test")
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	// Verify loaded values
	supabaseKey, _ := service.Get("SupabaseKey")
	if supabaseKey != "test_supabase_key_123" {
		t.Errorf("SupabaseKey mismatch: got '%s'", supabaseKey)
	}

	// Add custom value
	service.Set("CustomKey", "CustomValue")

	// Update loaded value
	service.Set("Port", "9090")

	// Verify all values
	port, _ := service.Get("Port")
	if port != "9090" {
		t.Errorf("Port should be updated to '9090', got '%s'", port)
	}

	customValue, _ := service.Get("CustomKey")
	if customValue != "CustomValue" {
		t.Errorf("CustomKey mismatch: got '%s'", customValue)
	}

	// Get config instance
	cfg := service.Config()
	cfgPort, _ := cfg.Get("Port")
	if cfgPort != "9090" {
		t.Errorf("cfg.Get('Port') should return '9090', got '%s'", cfgPort)
	}

	t.Log("✓ Complete workflow test passed")
}

func TestConcurrency(t *testing.T) {
	mapping := ConfigMapping{}
	service := New(mapping)

	// Test concurrent reads and writes
	done := make(chan bool)

	// Writer goroutine
	go func() {
		for i := 0; i < 100; i++ {
			service.Set("Key1", "Value")
		}
		done <- true
	}()

	// Reader goroutine
	go func() {
		for i := 0; i < 100; i++ {
			service.Get("Key1")
		}
		done <- true
	}()

	// Wait for both goroutines
	<-done
	<-done

	t.Log("✓ Concurrent access works without race conditions")
}

func TestLogging(t *testing.T) {
	// Test logging enable/disable
	if IsLoggingEnabled() {
		t.Error("Logging should be disabled by default")
	}

	EnableLogging()
	if !IsLoggingEnabled() {
		t.Error("Logging should be enabled after EnableLogging()")
	}

	DisableLogging()
	if IsLoggingEnabled() {
		t.Error("Logging should be disabled after DisableLogging()")
	}

	t.Log("✓ Logging enable/disable works correctly")
}

// Cleanup after tests
func TestMain(m *testing.M) {
	// Run tests
	code := m.Run()

	os.Exit(code)
}