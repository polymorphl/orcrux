package main

import (
	"testing"
)

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

	// Test that different inputs produce different outputs (randomness)
	t.Run("randomness check", func(t *testing.T) {
		results := make(map[byte]bool)
		for i := 0; i < 100; i++ {
			result, err := evaluatePolynomial(0x42, 1, 3)
			if err != nil {
				t.Errorf("evaluatePolynomial() failed: %v", err)
				return
			}
			results[result] = true
		}
		// With 100 random evaluations, we should get many different results
		if len(results) < 10 {
			t.Errorf("evaluatePolynomial() seems to produce too few unique results: %d", len(results))
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
