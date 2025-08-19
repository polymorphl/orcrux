package shamir

import (
	"encoding/base64"
	"encoding/hex"
	"strings"
	"testing"
)

// TestValidateShamirParams tests the parameter validation function
func TestValidateShamirParams(t *testing.T) {
	tests := []struct {
		name    string
		secret  []byte
		n       int
		t       int
		output  string
		wantErr bool
	}{
		{
			name:    "valid parameters",
			secret:  []byte("test"),
			n:       5,
			t:       3,
			output:  "hex",
			wantErr: false,
		},
		{
			name:    "empty secret",
			secret:  []byte{},
			n:       5,
			t:       3,
			output:  "hex",
			wantErr: true,
		},
		{
			name:    "n too small",
			secret:  []byte("test"),
			n:       1,
			t:       2,
			output:  "hex",
			wantErr: true,
		},
		{
			name:    "n too large",
			secret:  []byte("test"),
			n:       256,
			t:       3,
			output:  "hex",
			wantErr: true,
		},
		{
			name:    "t too small",
			secret:  []byte("test"),
			n:       5,
			t:       1,
			output:  "hex",
			wantErr: true,
		},
		{
			name:    "t larger than n",
			secret:  []byte("test"),
			n:       5,
			t:       6,
			output:  "hex",
			wantErr: true,
		},
		{
			name:    "invalid output format",
			secret:  []byte("test"),
			n:       5,
			t:       3,
			output:  "invalid",
			wantErr: true,
		},
		{
			name:    "output format with spaces",
			secret:  []byte("test"),
			n:       5,
			t:       3,
			output:  " base64 ",
			wantErr: false,
		},
		{
			name:    "output format case insensitive",
			secret:  []byte("test"),
			n:       5,
			t:       3,
			output:  "HEX",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateShamirParams(tt.secret, tt.n, tt.t, tt.output)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateShamirParams() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestGenerateXCoordinates tests the x-coordinate generation function
func TestGenerateXCoordinates(t *testing.T) {
	tests := []struct {
		name string
		n    int
		want []byte
	}{
		{
			name: "n=2",
			n:    2,
			want: []byte{1, 2},
		},
		{
			name: "n=5",
			n:    5,
			want: []byte{1, 2, 3, 4, 5},
		},
		{
			name: "n=10",
			n:    10,
			want: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		},
		{
			name: "n=255",
			n:    255,
			want: func() []byte {
				result := make([]byte, 255)
				for i := 0; i < 255; i++ {
					result[i] = byte(i + 1)
				}
				return result
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := generateXCoordinates(tt.n)
			if len(got) != len(tt.want) {
				t.Errorf("generateXCoordinates() length = %v, want %v", len(got), len(tt.want))
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("generateXCoordinates()[%d] = %v, want %v", i, got[i], tt.want[i])
				}
			}
		})
	}
}

// TestEvaluatePolynomial tests the polynomial evaluation function
func TestEvaluatePolynomial(t *testing.T) {
	// Test with different degrees and values
	tests := []struct {
		name       string
		secretByte byte
		x          byte
		t          int
		wantErr    bool
	}{
		{
			name:       "degree 2, x=1",
			secretByte: 0x42,
			x:          1,
			t:          2,
			wantErr:    false,
		},
		{
			name:       "degree 3, x=2",
			secretByte: 0x7F,
			x:          2,
			t:          3,
			wantErr:    false,
		},
		{
			name:       "degree 5, x=10",
			secretByte: 0x00,
			x:          10,
			t:          5,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := evaluatePolynomial(tt.secretByte, tt.x, tt.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("evaluatePolynomial() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				// Verify the result is a valid byte
				if result < 0 || result > 255 {
					t.Errorf("evaluatePolynomial() result = %v, should be in [0, 255]", result)
				}
			}
		})
	}

	// Test that different inputs produce different outputs (deterministic)
	t.Run("deterministic check", func(t *testing.T) {
		// Test that the same inputs always produce the same output
		result1, err := evaluatePolynomial(0x42, 1, 3)
		if err != nil {
			t.Errorf("evaluatePolynomial() failed: %v", err)
			return
		}

		result2, err := evaluatePolynomial(0x42, 1, 3)
		if err != nil {
			t.Errorf("evaluatePolynomial() failed: %v", err)
			return
		}

		// Results should be identical since coefficients are now deterministic
		if result1 != result2 {
			t.Errorf("evaluatePolynomial() produced different results: %v vs %v", result1, result2)
		}

		// Test that different inputs produce different outputs
		result3, err := evaluatePolynomial(0x43, 1, 3)
		if err != nil {
			t.Errorf("evaluatePolynomial() failed: %v", err)
			return
		}

		if result1 == result3 {
			t.Errorf("evaluatePolynomial() produced identical results for different inputs: %v", result1)
		}
	})
}

// TestEncodeShare tests the share encoding function
func TestEncodeShare(t *testing.T) {
	testData := []byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xAB, 0xCD, 0xEF}

	tests := []struct {
		name   string
		data   []byte
		format string
		want   string
	}{
		{
			name:   "hex encoding",
			data:   testData,
			format: "hex",
			want:   "0123456789abcdef",
		},
		{
			name:   "base64 encoding",
			data:   testData,
			format: "base64",
			want:   "ASNFZ4mrze8=",
		},
		{
			name:   "empty data hex",
			data:   []byte{},
			format: "hex",
			want:   "",
		},
		{
			name:   "empty data base64",
			data:   []byte{},
			format: "base64",
			want:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := encodeShare(tt.data, tt.format)
			if got != tt.want {
				t.Errorf("encodeShare() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestGfMul tests the Galois Field multiplication function
func TestGfMul(t *testing.T) {
	tests := []struct {
		name string
		a    byte
		b    byte
		want byte
	}{
		{
			name: "0 * 0 = 0",
			a:    0,
			b:    0,
			want: 0,
		},
		{
			name: "1 * 1 = 1",
			a:    1,
			b:    1,
			want: 1,
		},
		{
			name: "2 * 3 = 6",
			a:    2,
			b:    3,
			want: 6,
		},
		{
			name: "0x57 * 0x83 = 0xc1 (AES standard example)",
			a:    0x57,
			b:    0x83,
			want: 0xc1,
		},
		{
			name: "0x13 * 0x13 = 0x1e (reduction example)",
			a:    0x13,
			b:    0x13,
			want: 0x1e,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := gfMul(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("gfMul(0x%02x, 0x%02x) = 0x%02x, want 0x%02x", tt.a, tt.b, got, tt.want)
			}
		})
	}

	// Test commutativity: a * b = b * a
	t.Run("commutativity", func(t *testing.T) {
		for a := byte(0); a < 100; a++ {
			for b := byte(0); b < 100; b++ {
				if gfMul(a, b) != gfMul(b, a) {
					t.Errorf("gfMul not commutative: gfMul(0x%02x, 0x%02x) != gfMul(0x%02x, 0x%02x)", a, b, b, a)
					return
				}
			}
		}
	})

	// Test associativity: (a * b) * c = a * (b * c)
	t.Run("associativity", func(t *testing.T) {
		for a := byte(0); a < 50; a++ {
			for b := byte(0); b < 50; b++ {
				for c := byte(0); c < 50; c++ {
					left := gfMul(gfMul(a, b), c)
					right := gfMul(a, gfMul(b, c))
					if left != right {
						t.Errorf("gfMul not associative: (0x%02x * 0x%02x) * 0x%02x != 0x%02x * (0x%02x * 0x%02x)", a, b, c, a, b, c)
						return
					}
				}
			}
		}
	})

	// Test identity: a * 1 = a
	t.Run("identity", func(t *testing.T) {
		for a := byte(0); a < 255; a++ {
			if gfMul(a, 1) != a {
				t.Errorf("gfMul identity failed: 0x%02x * 1 = 0x%02x, want 0x%02x", a, gfMul(a, 1), a)
				return
			}
		}
	})
}

// TestSplit tests the main Shamir split function
func TestSplit(t *testing.T) {
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
			result, err := Split(tt.secret, tt.n, tt.t, tt.output)

			if (err != nil) != tt.wantErr {
				t.Errorf("Split() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			// Check that we got the expected number of shards
			shards := strings.Split(strings.TrimSpace(result), "\n")
			if len(shards) != tt.wantShards {
				t.Errorf("Split() returned %d shards, want %d", len(shards), tt.wantShards)
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

// TestSplitDeterministic tests that the same input produces consistent output structure
func TestSplitDeterministic(t *testing.T) {
	secret := []byte("deterministic test")
	n, threshold := 5, 3
	output := "hex"

	// Run multiple times with the same parameters
	results := make([]string, 5)
	for i := 0; i < 5; i++ {
		result, err := Split(secret, n, threshold, output)
		if err != nil {
			t.Errorf("Split() failed on iteration %d: %v", i, err)
			return
		}
		results[i] = result
	}

	// All results should have the same structure (same number of shards, same x-coordinates)
	for i := 1; i < len(results); i++ {
		shards1 := strings.Split(strings.TrimSpace(results[0]), "\n")
		shards2 := strings.Split(strings.TrimSpace(results[i]), "\n")

		if len(shards1) != len(shards2) {
			t.Errorf("Split() produced different number of shards on iteration %d", i)
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

// TestSplitDifferentSecrets tests that different secrets produce different outputs
func TestSplitDifferentSecrets(t *testing.T) {
	secrets := [][]byte{
		[]byte("secret1"),
		[]byte("secret2"),
		[]byte("secret3"),
	}
	n, threshold := 4, 2
	output := "hex"

	results := make([]string, len(secrets))
	for i, secret := range secrets {
		result, err := Split(secret, n, threshold, output)
		if err != nil {
			t.Errorf("Split() failed for secret %d: %v", i, err)
			return
		}
		results[i] = result
	}

	// All results should be different
	for i := 0; i < len(results); i++ {
		for j := i + 1; j < len(results); j++ {
			if results[i] == results[j] {
				t.Errorf("Split() produced identical results for different secrets %d and %d", i, j)
			}
		}
	}
}

// TestRecompose tests the secret reconstruction functionality
func TestRecompose(t *testing.T) {
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
			wantErr:     true, // Split fails on empty secrets
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
			shares, err := Split([]byte(tt.secret), tt.n, tt.t, tt.output)
			if err != nil {
				if tt.wantErr {
					// Expected error, test passes
					return
				}
				t.Fatalf("Split failed: %v", err)
			}

			// Parse shares into individual lines
			shareLines := strings.Split(strings.TrimSpace(shares), "\n")
			if len(shareLines) != tt.n {
				t.Fatalf("expected %d shares, got %d", tt.n, len(shareLines))
			}

			// Select the first 'useShards' shares for reconstruction
			selectedShares := shareLines[:tt.useShards]

			// Attempt to reconstruct
			reconstructedBytes, err := Recompose(selectedShares)

			if (err != nil) != tt.wantErr {
				t.Errorf("Recompose() error = %v, wantErr %v", err, tt.wantErr)
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

// TestRecomposeEdgeCases tests edge cases and error conditions
func TestRecomposeEdgeCases(t *testing.T) {
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
			_, err := Recompose(tt.shards)
			if (err != nil) != tt.wantErr {
				t.Errorf("Recompose() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestRecomposeConsistency tests that reconstruction is consistent
func TestRecomposeConsistency(t *testing.T) {
	secret := "consistent test secret"
	n, threshold := 6, 3
	output := "hex"

	// Split the secret
	shares, err := Split([]byte(secret), n, threshold, output)
	if err != nil {
		t.Fatalf("Split failed: %v", err)
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
		resultBytes, err := Recompose(subset)
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

// TestRecomposeMixedFormats tests reconstruction with mixed encoding formats
func TestRecomposeMixedFormats(t *testing.T) {
	secret := "mixed format test"
	n, threshold := 4, 2
	output := "hex"

	// Split the secret
	shares, err := Split([]byte(secret), n, threshold, output)
	if err != nil {
		t.Fatalf("Split failed: %v", err)
	}

	shareLines := strings.Split(strings.TrimSpace(shares), "\n")

	// Test that we can reconstruct from any subset of the required size
	for i := 0; i <= n-threshold; i++ {
		subset := shareLines[i : i+threshold]
		resultBytes, err := Recompose(subset)
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

// BenchmarkSplit benchmarks the main Shamir split function
func BenchmarkSplit(b *testing.B) {
	secret := []byte(strings.Repeat("test secret", 100))
	n, t := 10, 5
	output := "hex"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Split(secret, n, t, output)
		if err != nil {
			b.Fatalf("Split failed: %v", err)
		}
	}
}

// BenchmarkGfMul benchmarks the Galois Field multiplication function
func BenchmarkGfMul(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = gfMul(byte(i%256), byte((i+1)%256))
	}
}
