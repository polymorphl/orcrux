package main

import (
	"encoding/base64"
	"encoding/hex"
	"strings"
	"testing"
)

// TestShamirSplit tests the main Shamir split function
func TestShamirSplit(t *testing.T) {
	tests := []struct {
		name                string
		secret              []byte
		n                   int
		t                   int
		output              string
		wantShards          int
		wantErr             bool
		checkReconstruction bool
	}{
		{
			name:                "basic hex split",
			secret:              []byte("test"),
			n:                   3,
			t:                   2,
			output:              "hex",
			wantShards:          3,
			wantErr:             false,
			checkReconstruction: false,
		},
		{
			name:                "basic base64 split",
			secret:              []byte("hello world"),
			n:                   5,
			t:                   3,
			output:              "base64",
			wantShards:          5,
			wantErr:             false,
			checkReconstruction: false,
		},
		{
			name:                "large secret",
			secret:              []byte(strings.Repeat("a", 1000)),
			n:                   10,
			t:                   5,
			output:              "hex",
			wantShards:          10,
			wantErr:             false,
			checkReconstruction: false,
		},
		{
			name:                "minimum parameters",
			secret:              []byte("min"),
			n:                   2,
			t:                   2,
			output:              "hex",
			wantShards:          2,
			wantErr:             false,
			checkReconstruction: false,
		},
		{
			name:                "maximum n",
			secret:              []byte("max"),
			n:                   255,
			t:                   128,
			output:              "hex",
			wantShards:          255,
			wantErr:             false,
			checkReconstruction: false,
		},
		{
			name:                "empty secret",
			secret:              []byte{},
			n:                   3,
			t:                   2,
			output:              "hex",
			wantShards:          0,
			wantErr:             true,
			checkReconstruction: false,
		},
		{
			name:                "n too small",
			secret:              []byte("test"),
			n:                   1,
			t:                   1,
			output:              "hex",
			wantShards:          0,
			wantErr:             true,
			checkReconstruction: false,
		},
		{
			name:                "n too large",
			secret:              []byte("test"),
			n:                   256,
			t:                   3,
			output:              "hex",
			wantShards:          0,
			wantErr:             true,
			checkReconstruction: false,
		},
		{
			name:                "t too small",
			secret:              []byte("test"),
			n:                   5,
			t:                   1,
			output:              "hex",
			wantShards:          0,
			wantErr:             true,
			checkReconstruction: false,
		},
		{
			name:                "t larger than n",
			secret:              []byte("test"),
			n:                   5,
			t:                   6,
			output:              "hex",
			wantShards:          0,
			wantErr:             true,
			checkReconstruction: false,
		},
		{
			name:                "invalid output format",
			secret:              []byte("test"),
			n:                   3,
			t:                   2,
			output:              "invalid",
			wantShards:          0,
			wantErr:             true,
			checkReconstruction: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := shamirSplit(tt.secret, tt.n, tt.t, tt.output)

			if (err != nil) != tt.wantErr {
				t.Errorf("shamirSplit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			// Check that we got the expected number of shards
			shards := strings.Split(strings.TrimSpace(result), "\n")
			if len(shards) != tt.wantShards {
				t.Errorf("shamirSplit() returned %d shards, want %d", len(shards), tt.wantShards)
			}

			// Validate each shard format
			for i, shard := range shards {
				parts := strings.Split(shard, ":")
				if len(parts) != 2 {
					t.Errorf("shard %d has invalid format: %s", i, shard)
					continue
				}

				// Check x-coordinate
				xHex := parts[0]
				if len(xHex) != 2 {
					t.Errorf("shard %d x-coordinate has wrong length: %s", i, xHex)
				}

				// Check encoded data
				encodedData := parts[1]
				if len(encodedData) == 0 {
					t.Errorf("shard %d has empty encoded data", i)
				}

				// Verify encoding format
				if tt.output == "hex" {
					if _, err := hex.DecodeString(encodedData); err != nil {
						t.Errorf("shard %d hex data is invalid: %s", i, encodedData)
					}
				} else if tt.output == "base64" {
					if _, err := base64.StdEncoding.DecodeString(encodedData); err != nil {
						t.Errorf("shard %d base64 data is invalid: %s", i, encodedData)
					}
				}
			}
		})
	}
}

// TestShamirSplitDeterministic tests that the same input produces consistent output structure
func TestShamirSplitDeterministic(t *testing.T) {
	secret := []byte("deterministic test")
	n, threshold := 5, 3
	output := "hex"

	// Run multiple times with the same parameters
	results := make([]string, 5)
	for i := 0; i < 5; i++ {
		result, err := shamirSplit(secret, n, threshold, output)
		if err != nil {
			t.Errorf("shamirSplit() failed on iteration %d: %v", i, err)
			return
		}
		results[i] = result
	}

	// All results should have the same structure (same number of shards, same x-coordinates)
	for i := 1; i < len(results); i++ {
		shards1 := strings.Split(strings.TrimSpace(results[0]), "\n")
		shards2 := strings.Split(strings.TrimSpace(results[i]), "\n")

		if len(shards1) != len(shards2) {
			t.Errorf("shamirSplit() produced different number of shards on iteration %d", i)
			continue
		}

		// Check that x-coordinates are the same (they should be deterministic)
		for j, shard1 := range shards1 {
			parts1 := strings.Split(shard1, ":")
			parts2 := strings.Split(shards2[j], ":")

			if len(parts1) != 2 || len(parts2) != 2 {
				t.Errorf("shard %d has invalid format", j)
				continue
			}

			if parts1[0] != parts2[0] {
				t.Errorf("shard %d has different x-coordinate on iteration %d: %s vs %s", j, i, parts1[0], parts2[0])
			}
		}
	}
}

// TestShamirSplitDifferentSecrets tests that different secrets produce different outputs
func TestShamirSplitDifferentSecrets(t *testing.T) {
	secrets := [][]byte{
		[]byte("secret1"),
		[]byte("secret2"),
		[]byte("secret3"),
	}
	n, threshold := 4, 2
	output := "hex"

	results := make([]string, len(secrets))
	for i, secret := range secrets {
		result, err := shamirSplit(secret, n, threshold, output)
		if err != nil {
			t.Errorf("shamirSplit() failed for secret %d: %v", i, err)
			return
		}
		results[i] = result
	}

	// All results should be different
	for i := 0; i < len(results); i++ {
		for j := i + 1; j < len(results); j++ {
			if results[i] == results[j] {
				t.Errorf("shamirSplit() produced identical results for different secrets %d and %d", i, j)
			}
		}
	}
}
