package ft_config

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"sync"
	"testing"
)

// ANSI color codes for terminal output
const (
	colorReset = "\033[0m"
	colorRed   = "\033[31m"
	colorGreen = "\033[32m"
)

// TestResult tracks the outcome of a single test
type TestResult struct {
	Name   string
	Passed bool
	Output string
	Error  string
}

// Global test results tracker
var (
	testResults   []TestResult
	testResultsMu sync.Mutex
)

// recordTestResult records the result of a test
func recordTestResult(name string, passed bool, output string, err string) {
	testResultsMu.Lock()
	defer testResultsMu.Unlock()

	testResults = append(testResults, TestResult{
		Name:   name,
		Passed: passed,
		Output: output,
		Error:  err,
	})
}

// printTestSummary prints a colored summary of all test results
func printTestSummary() {
	var (
		totalTests  int
		passedTests int
		failedTests int
	)

	testResultsMu.Lock()
	defer testResultsMu.Unlock()

	totalTests = len(testResults)
	for _, result := range testResults {
		if result.Passed {
			passedTests++
		} else {
			failedTests++
		}
	}

	fmt.Println("\n\n" + string(bytes.Repeat([]byte("="), 60)))
	fmt.Println("                     TEST SUMMARY")
	fmt.Println(string(bytes.Repeat([]byte("="), 60)))

	// print individual test results
	for _, result := range testResults {
		if result.Passed {
			fmt.Printf("%s[SUCCESS]%s %s\n", colorGreen, colorReset, result.Name)
		} else {
			fmt.Printf("%s[FAIL]%s %s\n", colorRed, colorReset, result.Name)
			if result.Error != "" {
				fmt.Printf("  Error: %s\n", result.Error)
			}
			if result.Output != "" {
				fmt.Printf("  Output:\n%s\n", result.Output)
			}
		}
	}

	// print summary statistics
	fmt.Println(string(bytes.Repeat([]byte("="), 60)))
	fmt.Printf("Total:  %d tests\n", totalTests)
	fmt.Printf("%sPassed: %d tests%s\n", colorGreen, passedTests, colorReset)
	if failedTests > 0 {
		fmt.Printf("%sFailed: %d tests%s\n", colorRed, failedTests, colorReset)
	} else {
		fmt.Printf("Failed: %d tests\n", failedTests)
	}
	fmt.Println(string(bytes.Repeat([]byte("="), 60)))
}

// TestMain runs before all tests and prints summary after all tests complete
func TestMain(m *testing.M) {
	// run all tests
	exitCode := m.Run()

	// print summary
	printTestSummary()

	// exit with appropriate code
	os.Exit(exitCode)
}

func TestLoad(t *testing.T) {
	var (
		testName     = "TestLoad"
		config       TestConfig
		output       bytes.Buffer
		errorMessage string
		err          error
	)

	output.WriteString("\n========================================\n")
	output.WriteString("Testing Load\n")
	output.WriteString("========================================\n")

	err = Load(".env.test", &config)
	if err != nil {
		errorMessage = fmt.Sprintf("Load() failed: %v", err)
		recordTestResult(testName, false, output.String(), errorMessage)
		t.Fatalf("%s", errorMessage)
		return
	}

	// Verify all values were loaded correctly
	if config.SupabaseKey != "test_supabase_key_123" {
		errorMessage = fmt.Sprintf("SupabaseKey: expected 'test_supabase_key_123', got '%s'", config.SupabaseKey)
		recordTestResult(testName, false, output.String(), errorMessage)
		t.Errorf("%s", errorMessage)
		return
	}

	if config.SupabaseURL != "https://test.supabase.co" {
		errorMessage = fmt.Sprintf("SupabaseURL: expected 'https://test.supabase.co', got '%s'", config.SupabaseURL)
		recordTestResult(testName, false, output.String(), errorMessage)
		t.Errorf("%s", errorMessage)
		return
	}

	if config.DatabaseURL != "postgres://localhost:5432/testdb" {
		errorMessage = fmt.Sprintf("DatabaseURL: expected 'postgres://localhost:5432/testdb', got '%s'", config.DatabaseURL)
		recordTestResult(testName, false, output.String(), errorMessage)
		t.Errorf("%s", errorMessage)
		return
	}

	if config.ApiToken != "test_api_token_456" {
		errorMessage = fmt.Sprintf("ApiToken: expected 'test_api_token_456', got '%s'", config.ApiToken)
		recordTestResult(testName, false, output.String(), errorMessage)
		t.Errorf("%s", errorMessage)
		return
	}

	if config.Port != "8080" {
		errorMessage = fmt.Sprintf("Port: expected '8080', got '%s'", config.Port)
		recordTestResult(testName, false, output.String(), errorMessage)
		t.Errorf("%s", errorMessage)
		return
	}

	if config.AppName != "ft_config_test" {
		errorMessage = fmt.Sprintf("AppName: expected 'ft_config_test', got '%s'", config.AppName)
		recordTestResult(testName, false, output.String(), errorMessage)
		t.Errorf("%s", errorMessage)
		return
	}

	output.WriteString("✓ Load() successfully fills struct from .env file\n")
	output.WriteString("========================================\n")

	recordTestResult(testName, true, output.String(), "")
}

func TestLoadNonExistentFile(t *testing.T) {
	var (
		testName     = "TestLoadNonExistentFile"
		config       TestConfig
		output       bytes.Buffer
		errorMessage string
		err          error
	)

	output.WriteString("\n========================================\n")
	output.WriteString("Testing LoadNonExistentFile\n")
	output.WriteString("========================================\n")

	err = Load("nonexistent.env", &config)
	if err == nil {
		errorMessage = "Expected error when loading non-existent file"
		recordTestResult(testName, false, output.String(), errorMessage)
		t.Error(errorMessage)
		return
	}

	if !errors.Is(err, ErrLoadEnv) {
		errorMessage = fmt.Sprintf("Expected ErrLoadEnv, got: %v", err)
		recordTestResult(testName, false, output.String(), errorMessage)
		t.Errorf("%s", errorMessage)
		return
	}

	output.WriteString("✓ Load() returns error for non-existent file\n")
	output.WriteString("========================================\n")

	recordTestResult(testName, true, output.String(), "")
}

func TestLoadWithPartialFields(t *testing.T) {
	var (
		testName     = "TestLoadWithPartialFields"
		output       bytes.Buffer
		errorMessage string
		err          error
	)

	// Struct with only some fields
	type PartialConfig struct {
		SupabaseKey string `env:"SUPABASE_KEY"`
		Port        string `env:"PORT"`
	}

	var config PartialConfig

	output.WriteString("\n========================================\n")
	output.WriteString("Testing LoadWithPartialFields\n")
	output.WriteString("========================================\n")

	err = Load(".env.test", &config)
	if err != nil {
		errorMessage = fmt.Sprintf("Load() failed: %v", err)
		recordTestResult(testName, false, output.String(), errorMessage)
		t.Fatalf("%s", errorMessage)
		return
	}

	if config.SupabaseKey != "test_supabase_key_123" {
		errorMessage = fmt.Sprintf("SupabaseKey: expected 'test_supabase_key_123', got '%s'", config.SupabaseKey)
		recordTestResult(testName, false, output.String(), errorMessage)
		t.Errorf("%s", errorMessage)
		return
	}

	if config.Port != "8080" {
		errorMessage = fmt.Sprintf("Port: expected '8080', got '%s'", config.Port)
		recordTestResult(testName, false, output.String(), errorMessage)
		t.Errorf("%s", errorMessage)
		return
	}

	output.WriteString("✓ Load() works with partial struct fields\n")
	output.WriteString("========================================\n")

	recordTestResult(testName, true, output.String(), "")
}

func TestLoadWithoutEnvTag(t *testing.T) {
	var (
		testName     = "TestLoadWithoutEnvTag"
		output       bytes.Buffer
		errorMessage string
		err          error
	)

	// Struct with fields without env tags
	type ConfigWithoutTags struct {
		SupabaseKey  string `env:"SUPABASE_KEY"`
		IgnoredField string // No env tag - should be ignored
	}

	var config ConfigWithoutTags

	output.WriteString("\n========================================\n")
	output.WriteString("Testing LoadWithoutEnvTag\n")
	output.WriteString("========================================\n")

	err = Load(".env.test", &config)
	if err != nil {
		errorMessage = fmt.Sprintf("Load() failed: %v", err)
		recordTestResult(testName, false, output.String(), errorMessage)
		t.Fatalf("%s", errorMessage)
		return
	}

	if config.SupabaseKey != "test_supabase_key_123" {
		errorMessage = fmt.Sprintf("SupabaseKey: expected 'test_supabase_key_123', got '%s'", config.SupabaseKey)
		recordTestResult(testName, false, output.String(), errorMessage)
		t.Errorf("%s", errorMessage)
		return
	}

	if config.IgnoredField != "" {
		errorMessage = fmt.Sprintf("IgnoredField should be empty, got '%s'", config.IgnoredField)
		recordTestResult(testName, false, output.String(), errorMessage)
		t.Errorf("%s", errorMessage)
		return
	}

	output.WriteString("✓ Load() ignores fields without env tag\n")
	output.WriteString("========================================\n")

	recordTestResult(testName, true, output.String(), "")
}

func TestLoadNotAPointer(t *testing.T) {
	var (
		testName     = "TestLoadNotAPointer"
		config       TestConfig
		output       bytes.Buffer
		errorMessage string
		err          error
	)

	output.WriteString("\n========================================\n")
	output.WriteString("Testing LoadNotAPointer\n")
	output.WriteString("========================================\n")

	// Pass struct by value instead of pointer
	err = Load(".env.test", config)
	if err == nil {
		errorMessage = "Expected error when passing struct by value"
		recordTestResult(testName, false, output.String(), errorMessage)
		t.Error(errorMessage)
		return
	}

	output.WriteString("✓ Load() returns error for non-pointer argument\n")
	output.WriteString("========================================\n")

	recordTestResult(testName, true, output.String(), "")
}

func TestLoadWithMissingEnvVar(t *testing.T) {
	var (
		testName     = "TestLoadWithMissingEnvVar"
		output       bytes.Buffer
		errorMessage string
		err          error
	)

	type ConfigWithMissing struct {
		SupabaseKey    string `env:"SUPABASE_KEY"`
		NonExistentVar string `env:"NON_EXISTENT_VAR"`
	}

	var config ConfigWithMissing

	output.WriteString("\n========================================\n")
	output.WriteString("Testing LoadWithMissingEnvVar\n")
	output.WriteString("========================================\n")

	err = Load(".env.test", &config)
	if err == nil {
		errorMessage = "Expected error when env var is missing"
		recordTestResult(testName, false, output.String(), errorMessage)
		t.Error(errorMessage)
		return
	}

	expectedErrMsg := "missing required environment variables: NON_EXISTENT_VAR"
	if err.Error() != expectedErrMsg {
		errorMessage = fmt.Sprintf("Expected error message '%s', got '%s'", expectedErrMsg, err.Error())
		recordTestResult(testName, false, output.String(), errorMessage)
		t.Errorf("%s", errorMessage)
		return
	}

	output.WriteString("✓ Load() returns error for missing environment variables\n")
	output.WriteString("========================================\n")

	recordTestResult(testName, true, output.String(), "")
}

func TestRealWorldExample(t *testing.T) {
	var (
		testName     = "TestRealWorldExample"
		output       bytes.Buffer
		errorMessage string
		err          error
	)

	// Real-world usage example
	type AppConfig struct {
		SupabaseKey string `env:"SUPABASE_KEY"`
		SupabaseURL string `env:"SUPABASE_URL"`
		DatabaseURL string `env:"DATABASE_URL"`
		Port        string `env:"PORT"`
	}

	var config AppConfig

	output.WriteString("\n========================================\n")
	output.WriteString("Testing RealWorldExample\n")
	output.WriteString("========================================\n")

	// One simple call!
	err = Load(".env.test", &config)
	if err != nil {
		errorMessage = fmt.Sprintf("Load() failed: %v", err)
		recordTestResult(testName, false, output.String(), errorMessage)
		t.Fatalf("%s", errorMessage)
		return
	}

	// Now use the config directly
	if config.SupabaseKey == "" {
		errorMessage = "SupabaseKey should not be empty"
		recordTestResult(testName, false, output.String(), errorMessage)
		t.Error(errorMessage)
		return
	}

	if config.Port == "" {
		errorMessage = "Port should not be empty"
		recordTestResult(testName, false, output.String(), errorMessage)
		t.Error(errorMessage)
		return
	}

	// Config is ready to use!
	output.WriteString(fmt.Sprintf("✓ Config loaded: Port=%s, Database=%s\n", config.Port, config.DatabaseURL))
	output.WriteString("========================================\n")

	recordTestResult(testName, true, output.String(), "")
}

// Test struct matching .env.test values
type TestConfig struct {
	SupabaseKey string `env:"SUPABASE_KEY"`
	SupabaseURL string `env:"SUPABASE_URL"`
	DatabaseURL string `env:"DATABASE_URL"`
	ApiToken    string `env:"API_TOKEN"`
	Port        string `env:"PORT"`
	AppName     string `env:"APP_NAME"`
}
