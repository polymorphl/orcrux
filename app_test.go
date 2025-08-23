package main

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

// TestNewApp tests the App constructor
func TestNewApp(t *testing.T) {
	app := NewApp()
	if app == nil {
		t.Fatal("NewApp() returned nil")
	}
	if app.ctx != nil {
		t.Error("NewApp() should not have context set initially")
	}
}

// TestAppStartup tests the startup method
func TestAppStartup(t *testing.T) {
	app := NewApp()
	ctx := &mockContext{}

	app.startup(ctx)

	if app.ctx != ctx {
		t.Error("startup() did not set the context correctly")
	}
}

// TestAppSplit tests the Split method with various inputs
func TestAppSplit(t *testing.T) {
	app := NewApp()

	tests := []struct {
		name          string
		secret        string
		shards        int
		shardsNeeded  int
		output        string
		expectError   bool
		expectedLines int
	}{
		{
			name:          "basic hex split",
			secret:        "test secret",
			shards:        3,
			shardsNeeded:  2,
			output:        "hex",
			expectError:   false,
			expectedLines: 3,
		},
		{
			name:          "basic base64 split",
			secret:        "hello world",
			shards:        5,
			shardsNeeded:  3,
			output:        "base64",
			expectError:   false,
			expectedLines: 5,
		},
		{
			name:          "minimum parameters",
			secret:        "min",
			shards:        2,
			shardsNeeded:  2,
			output:        "hex",
			expectError:   false,
			expectedLines: 2,
		},
		{
			name:          "large secret",
			secret:        strings.Repeat("a", 100),
			shards:        10,
			shardsNeeded:  5,
			output:        "hex",
			expectError:   false,
			expectedLines: 10,
		},
		{
			name:          "empty secret",
			secret:        "",
			shards:        3,
			shardsNeeded:  2,
			output:        "hex",
			expectError:   true,
			expectedLines: 0,
		},
		{
			name:          "invalid shards count",
			secret:        "test",
			shards:        1, // Too small
			shardsNeeded:  2,
			output:        "hex",
			expectError:   true,
			expectedLines: 0,
		},
		{
			name:          "invalid threshold",
			secret:        "test",
			shards:        5,
			shardsNeeded:  6, // Larger than shards
			output:        "hex",
			expectError:   true,
			expectedLines: 0,
		},
		{
			name:          "invalid output format",
			secret:        "test",
			shards:        3,
			shardsNeeded:  2,
			output:        "invalid",
			expectError:   true,
			expectedLines: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := app.Split(tt.secret, tt.shards, tt.shardsNeeded, tt.output)

			// Parse JSON response
			var response Response
			err := json.Unmarshal([]byte(result), &response)
			if err != nil {
				t.Fatalf("Failed to parse JSON response: %v, response: %s", err, result)
			}

			if tt.expectError {
				if response.Error == nil {
					t.Errorf("Split() should have returned error, got: %+v", response)
				}
				return
			}

			if response.Error != nil {
				t.Errorf("Split() returned unexpected error: %s", *response.Error)
				return
			}

			// Verify the result format
			dataStr, ok := response.Data.(string)
			if !ok {
				t.Errorf("Split() data field is not a string: %T", response.Data)
				return
			}

			lines := strings.Split(strings.TrimSpace(dataStr), "\n")
			if len(lines) != tt.expectedLines {
				t.Errorf("Split() returned %d lines, expected %d", len(lines), tt.expectedLines)
			}

			// Verify each line format: "xx:encoded_data"
			for i, line := range lines {
				if line == "" {
					continue
				}
				parts := strings.Split(line, ":")
				if len(parts) != 2 {
					t.Errorf("Line %d has invalid format: %s", i, line)
					continue
				}

				// Check x-coordinate format (2 hex chars)
				if len(parts[0]) != 2 {
					t.Errorf("Line %d x-coordinate has wrong length: %s", i, parts[0])
				}

				// Check that encoded data is not empty
				if len(parts[1]) == 0 {
					t.Errorf("Line %d has empty encoded data", i)
				}
			}
		})
	}
}

// TestAppRecompose tests the Recompose method with various inputs
func TestAppRecompose(t *testing.T) {
	app := NewApp()

	// First create some valid shares to test with
	secret := "test secret for recompose"
	sharesResult := app.Split(secret, 5, 3, "hex")

	var sharesResponse Response
	err := json.Unmarshal([]byte(sharesResult), &sharesResponse)
	if err != nil {
		t.Fatalf("Failed to parse shares JSON response: %v", err)
	}
	if sharesResponse.Error != nil {
		t.Fatalf("Failed to create test shares: %s", *sharesResponse.Error)
	}

	sharesStr, ok := sharesResponse.Data.(string)
	if !ok {
		t.Fatalf("Shares data is not a string: %T", sharesResponse.Data)
	}

	shareLines := strings.Split(strings.TrimSpace(sharesStr), "\n")
	if len(shareLines) < 3 {
		t.Fatalf("Expected at least 3 shares, got %d", len(shareLines))
	}

	tests := []struct {
		name        string
		shards      []string
		expectError bool
		description string
	}{
		{
			name:        "valid reconstruction with minimum shares",
			shards:      shareLines[:3], // Use first 3 shares
			expectError: false,
			description: "Reconstruct using exactly the required number of shares",
		},
		{
			name:        "valid reconstruction with extra shares",
			shards:      shareLines[:4], // Use first 4 shares
			expectError: false,
			description: "Reconstruct using more than the required shares",
		},
		{
			name:        "insufficient shares",
			shards:      shareLines[:2], // Use only 2 shares (need 3)
			expectError: false,          // Current implementation allows this but produces wrong result
			description: "Reconstruct with fewer shares than threshold",
		},
		{
			name:        "empty shards array",
			shards:      []string{},
			expectError: true,
			description: "Fail with no shares provided",
		},
		{
			name:        "single shard",
			shards:      shareLines[:1],
			expectError: true,
			description: "Fail with only one share",
		},
		{
			name:        "invalid share format",
			shards:      []string{"invalid:format", "01:616263"},
			expectError: true,
			description: "Fail with malformed share",
		},
		{
			name:        "mixed valid and invalid shares",
			shards:      []string{shareLines[0], "invalid:format", shareLines[2]},
			expectError: true,
			description: "Fail when some shares are invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := app.Recompose(tt.shards)

			// Parse JSON response
			var response Response
			err := json.Unmarshal([]byte(result), &response)
			if err != nil {
				t.Fatalf("Failed to parse JSON response: %v, response: %s", err, result)
			}

			if tt.expectError {
				if response.Error == nil {
					t.Errorf("Recompose() should have returned error, got: %+v", response)
				}
				return
			}

			if response.Error != nil {
				t.Errorf("Recompose() returned unexpected error: %s", *response.Error)
				return
			}

			// For insufficient shares, we expect a wrong result (not an error)
			if tt.name == "insufficient shares" {
				dataStr, ok := response.Data.(string)
				if !ok {
					t.Errorf("Recompose() data field is not a string: %T", response.Data)
					return
				}
				if dataStr == secret {
					t.Logf("Note: Reconstruction with insufficient shares produced correct result (this may indicate a security issue)")
				} else {
					t.Logf("Expected wrong result with insufficient shares: got %q, want %q", dataStr, secret)
				}
				return
			}

			// For valid reconstructions, verify the result
			dataStr, ok := response.Data.(string)
			if !ok {
				t.Errorf("Recompose() data field is not a string: %T", response.Data)
				return
			}
			if dataStr != secret {
				t.Errorf("Recompose() failed: got %q, want %q", dataStr, secret)
			}
		})
	}
}

// TestAppRecomposeConsistency tests that reconstruction is consistent
func TestAppRecomposeConsistency(t *testing.T) {
	app := NewApp()
	secret := "consistency test secret"

	// Create shares
	sharesResult := app.Split(secret, 6, 3, "hex")
	var sharesResponse Response
	err := json.Unmarshal([]byte(sharesResult), &sharesResponse)
	if err != nil {
		t.Fatalf("Failed to parse shares JSON response: %v", err)
	}
	if sharesResponse.Error != nil {
		t.Fatalf("Failed to create test shares: %s", *sharesResponse.Error)
	}

	sharesStr, ok := sharesResponse.Data.(string)
	if !ok {
		t.Fatalf("Shares data is not a string: %T", sharesResponse.Data)
	}

	shareLines := strings.Split(strings.TrimSpace(sharesStr), "\n")

	// Test reconstruction with different subsets of shares
	subsets := [][]string{
		shareLines[:3],  // First 3 shares
		shareLines[1:4], // Middle 3 shares
		shareLines[3:],  // Last 3 shares
		{shareLines[0], shareLines[2], shareLines[4]}, // Non-consecutive shares
	}

	results := make([]string, len(subsets))
	for i, subset := range subsets {
		result := app.Recompose(subset)

		var response Response
		err := json.Unmarshal([]byte(result), &response)
		if err != nil {
			t.Errorf("Failed to parse JSON response for subset %d: %v", i, err)
			continue
		}

		if response.Error != nil {
			t.Errorf("Recompose() failed for subset %d: %s", i, *response.Error)
			continue
		}

		dataStr, ok := response.Data.(string)
		if !ok {
			t.Errorf("Recompose() data field is not a string for subset %d: %T", i, response.Data)
			continue
		}

		results[i] = dataStr
	}

	// All reconstructions should produce the same result
	for i := 1; i < len(results); i++ {
		if results[i] != results[0] {
			t.Errorf("Inconsistent reconstruction: subset %d = %q, subset 0 = %q", i, results[i], results[0])
		}
	}

	// Verify the result matches the original secret
	if results[0] != secret {
		t.Errorf("Reconstruction failed: got %q, want %q", results[0], secret)
	}
}

// TestAppRecomposeMixedFormats tests reconstruction with different encoding formats
func TestAppRecomposeMixedFormats(t *testing.T) {
	app := NewApp()

	// Test hex format
	secretHex := "hex test secret"
	sharesHexResult := app.Split(secretHex, 4, 2, "hex")
	var sharesHexResponse Response
	err := json.Unmarshal([]byte(sharesHexResult), &sharesHexResponse)
	if err != nil {
		t.Fatalf("Failed to parse hex shares JSON response: %v", err)
	}
	if sharesHexResponse.Error != nil {
		t.Fatalf("Failed to create hex shares: %s", *sharesHexResponse.Error)
	}

	// Test base64 format
	secretBase64 := "base64 test secret"
	sharesBase64Result := app.Split(secretBase64, 4, 2, "base64")
	var sharesBase64Response Response
	err = json.Unmarshal([]byte(sharesBase64Result), &sharesBase64Response)
	if err != nil {
		t.Fatalf("Failed to parse base64 shares JSON response: %v", err)
	}
	if sharesBase64Response.Error != nil {
		t.Fatalf("Failed to create base64 shares: %s", *sharesBase64Response.Error)
	}

	// Test hex reconstruction
	sharesHexStr, ok := sharesHexResponse.Data.(string)
	if !ok {
		t.Fatalf("Hex shares data is not a string: %T", sharesHexResponse.Data)
	}
	shareLinesHex := strings.Split(strings.TrimSpace(sharesHexStr), "\n")
	resultHex := app.Recompose(shareLinesHex[:2])

	var resultHexResponse Response
	err = json.Unmarshal([]byte(resultHex), &resultHexResponse)
	if err != nil {
		t.Errorf("Failed to parse hex result JSON response: %v", err)
	} else if resultHexResponse.Error != nil {
		t.Errorf("Hex reconstruction failed: %s", *resultHexResponse.Error)
	} else {
		resultHexStr, ok := resultHexResponse.Data.(string)
		if !ok {
			t.Errorf("Hex result data field is not a string: %T", resultHexResponse.Data)
		} else if resultHexStr != secretHex {
			t.Errorf("Hex reconstruction failed: got %q, want %q", resultHexStr, secretHex)
		}
	}

	// Test base64 reconstruction
	sharesBase64Str, ok := sharesBase64Response.Data.(string)
	if !ok {
		t.Fatalf("Base64 shares data is not a string: %T", sharesBase64Response.Data)
	}
	shareLinesBase64 := strings.Split(strings.TrimSpace(sharesBase64Str), "\n")
	resultBase64 := app.Recompose(shareLinesBase64[:2])

	var resultBase64Response Response
	err = json.Unmarshal([]byte(resultBase64), &resultBase64Response)
	if err != nil {
		t.Errorf("Failed to parse base64 result JSON response: %v", err)
	} else if resultBase64Response.Error != nil {
		t.Errorf("Base64 reconstruction failed: %s", *resultBase64Response.Error)
	} else {
		resultBase64Str, ok := resultBase64Response.Data.(string)
		if !ok {
			t.Errorf("Base64 result data field is not a string: %T", resultBase64Response.Data)
		} else if resultBase64Str != secretBase64 {
			t.Errorf("Base64 reconstruction failed: got %q, want %q", resultBase64Str, secretBase64)
		}
	}
}

// TestAppErrorHandling tests various error conditions
func TestAppErrorHandling(t *testing.T) {
	app := NewApp()

	tests := []struct {
		name          string
		secret        string
		shards        int
		shardsNeeded  int
		output        string
		expectedError string
	}{
		{
			name:          "empty secret",
			secret:        "",
			shards:        3,
			shardsNeeded:  2,
			output:        "hex",
			expectedError: "empty secret",
		},
		{
			name:          "shards too small",
			secret:        "test",
			shards:        1,
			shardsNeeded:  1,
			output:        "hex",
			expectedError: "shards must be in [2, 255]",
		},
		{
			name:          "shards too large",
			secret:        "test",
			shards:        256,
			shardsNeeded:  3,
			output:        "hex",
			expectedError: "shards must be in [2, 255]",
		},
		{
			name:          "threshold too small",
			secret:        "test",
			shards:        5,
			shardsNeeded:  1,
			output:        "hex",
			expectedError: "shardsNeeded must be in [2, shards]",
		},
		{
			name:          "threshold larger than shards",
			secret:        "test",
			shards:        3,
			shardsNeeded:  4,
			output:        "hex",
			expectedError: "shardsNeeded must be in [2, shards]",
		},
		{
			name:          "invalid output format",
			secret:        "test",
			shards:        3,
			shardsNeeded:  2,
			output:        "invalid",
			expectedError: "output must be 'hex' or 'base64'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := app.Split(tt.secret, tt.shards, tt.shardsNeeded, tt.output)

			// Parse JSON response
			var response Response
			err := json.Unmarshal([]byte(result), &response)
			if err != nil {
				t.Fatalf("Failed to parse JSON response: %v, response: %s", err, result)
			}

			if response.Error == nil {
				t.Errorf("Split() should have returned error, got: %+v", response)
				return
			}

			errorMsg := *response.Error
			if !strings.Contains(errorMsg, tt.expectedError) {
				t.Errorf("Split() error message '%s' should contain '%s'", errorMsg, tt.expectedError)
			}
		})
	}
}

// Mock context for testing
type mockContext struct{}

func (m *mockContext) Deadline() (deadline time.Time, ok bool) { return time.Time{}, false }
func (m *mockContext) Done() <-chan struct{}                   { return nil }
func (m *mockContext) Err() error                              { return nil }
func (m *mockContext) Value(key interface{}) interface{}       { return nil }
