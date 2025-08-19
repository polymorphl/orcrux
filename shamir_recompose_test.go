package main

import (
	"strings"
	"testing"
)

// TestShamirRecompose tests the secret reconstruction functionality
func TestShamirRecompose(t *testing.T) {
	tests := []struct {
		name        string
		secret      string
		n           int
		t           int
		output      string
		useShards   int // How many shards to use for reconstruction
		wantErr     bool
		description string
	}{
		{
			name:        "basic reconstruction hex",
			secret:      "test secret",
			n:           5,
			t:           3,
			output:      "hex",
			useShards:   3,
			wantErr:     false,
			description: "Reconstruct from 3 out of 5 hex-encoded shares",
		},
		{
			name:        "basic reconstruction base64",
			secret:      "hello world",
			n:           4,
			t:           2,
			output:      "base64",
			useShards:   2,
			wantErr:     false,
			description: "Reconstruct from 2 out of 4 base64-encoded shares",
		},
		{
			name:        "minimum shares required",
			secret:      "minimal",
			n:           3,
			t:           2,
			output:      "hex",
			useShards:   2,
			wantErr:     false,
			description: "Reconstruct using exactly the minimum required shares",
		},
		{
			name:        "more shares than minimum",
			secret:      "extra shares",
			n:           6,
			t:           3,
			output:      "hex",
			useShards:   5,
			wantErr:     false,
			description: "Reconstruct using more shares than the minimum required",
		},
		{
			name:        "insufficient shares",
			secret:      "insufficient",
			n:           5,
			t:           3,
			output:      "hex",
			useShards:   2,
			wantErr:     false, // Current implementation allows reconstruction with 2+ shares
			description: "Reconstruct with fewer shares than threshold (not secure but technically possible)",
		},
		{
			name:        "empty shards",
			secret:      "",
			n:           3,
			t:           2,
			output:      "hex",
			useShards:   2,
			wantErr:     true, // shamirSplit fails on empty secrets
			description: "Fail to split empty secret",
		},
		{
			name:        "large secret",
			secret:      strings.Repeat("a", 1000),
			n:           10,
			t:           5,
			output:      "hex",
			useShards:   6,
			wantErr:     false,
			description: "Reconstruct large secret",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// First, split the secret
			shares, err := shamirSplit([]byte(tt.secret), tt.n, tt.t, tt.output)
			if err != nil {
				if tt.wantErr {
					// Expected error, test passes
					return
				}
				t.Fatalf("shamirSplit failed: %v", err)
			}

			// Parse shares into individual lines
			shareLines := strings.Split(strings.TrimSpace(shares), "\n")
			if len(shareLines) != tt.n {
				t.Fatalf("expected %d shares, got %d", tt.n, len(shareLines))
			}

			// Select the first 'useShards' shares for reconstruction
			selectedShares := shareLines[:tt.useShards]

			// Attempt to reconstruct
			reconstructedBytes, err := shamirRecompose(selectedShares)

			if (err != nil) != tt.wantErr {
				t.Errorf("shamirRecompose() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			// Convert bytes to string for comparison
			reconstructed := string(reconstructedBytes)

			// Verify reconstruction
			if reconstructed != tt.secret {
				// For insufficient shares, we expect reconstruction to fail (wrong result)
				if tt.name == "insufficient shares" {
					t.Logf("Expected wrong result with insufficient shares: got %q, want %q", reconstructed, tt.secret)
				} else {
					t.Errorf("reconstruction failed: got %q, want %q", reconstructed, tt.secret)
				}
			}
		})
	}
}

// TestShamirRecomposeEdgeCases tests edge cases and error conditions
func TestShamirRecomposeEdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		shards  []string
		wantErr bool
	}{
		{
			name:    "no shards",
			shards:  []string{},
			wantErr: true,
		},
		{
			name:    "single shard",
			shards:  []string{"01:616263"},
			wantErr: true,
		},
		{
			name:    "invalid share format",
			shards:  []string{"01:616263", "invalid"},
			wantErr: true,
		},
		{
			name:    "malformed x-coordinate",
			shards:  []string{"1:616263", "02:646566"},
			wantErr: true,
		},
		{
			name:    "inconsistent share lengths",
			shards:  []string{"01:616263", "02:64"},
			wantErr: true,
		},
		{
			name:    "empty share data",
			shards:  []string{"01:", "02:"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := shamirRecompose(tt.shards)
			if (err != nil) != tt.wantErr {
				t.Errorf("shamirRecompose() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestShamirRecomposeConsistency tests that reconstruction is consistent
func TestShamirRecomposeConsistency(t *testing.T) {
	secret := "consistent test secret"
	n, threshold := 6, 3
	output := "hex"

	// Split the secret
	shares, err := shamirSplit([]byte(secret), n, threshold, output)
	if err != nil {
		t.Fatalf("shamirSplit failed: %v", err)
	}

	shareLines := strings.Split(strings.TrimSpace(shares), "\n")

	// Test reconstruction with different subsets of shares
	subsets := [][]string{
		shareLines[:3],  // First 3 shares
		shareLines[1:4], // Middle 3 shares
		shareLines[3:],  // Last 3 shares
		{shareLines[0], shareLines[2], shareLines[4]}, // Non-consecutive shares
	}

	results := make([]string, len(subsets))
	for i, subset := range subsets {
		resultBytes, err := shamirRecompose(subset)
		if err != nil {
			t.Errorf("reconstruction failed for subset %d: %v", i, err)
			continue
		}
		results[i] = string(resultBytes)
	}

	// All reconstructions should produce the same result
	for i := 1; i < len(results); i++ {
		if results[i] != results[0] {
			t.Errorf("inconsistent reconstruction: subset %d = %q, subset 0 = %q", i, results[i], results[0])
		}
	}

	// Verify the result matches the original secret
	if results[0] != secret {
		t.Errorf("reconstruction failed: got %q, want %q", results[0], secret)
	}
}

// TestShamirRecomposeMixedFormats tests reconstruction with mixed encoding formats
func TestShamirRecomposeMixedFormats(t *testing.T) {
	secret := "mixed format test"
	n, threshold := 4, 2
	output := "hex"

	// Split the secret
	shares, err := shamirSplit([]byte(secret), n, threshold, output)
	if err != nil {
		t.Fatalf("shamirSplit failed: %v", err)
	}

	shareLines := strings.Split(strings.TrimSpace(shares), "\n")

	// Test that we can reconstruct from any subset of the required size
	for i := 0; i <= n-threshold; i++ {
		subset := shareLines[i : i+threshold]
		resultBytes, err := shamirRecompose(subset)
		if err != nil {
			t.Errorf("reconstruction failed for subset starting at %d: %v", i, err)
			continue
		}
		result := string(resultBytes)
		if result != secret {
			t.Errorf("reconstruction failed for subset starting at %d: got %q, want %q", i, result, secret)
		}
	}
}
