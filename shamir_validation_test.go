package main

import (
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
