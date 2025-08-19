package main

import (
	"strings"
	"testing"
)

// BenchmarkShamirSplit benchmarks the main Shamir split function
func BenchmarkShamirSplit(b *testing.B) {
	secret := []byte(strings.Repeat("test secret", 100))
	n, t := 10, 5
	output := "hex"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := shamirSplit(secret, n, t, output)
		if err != nil {
			b.Fatalf("shamirSplit failed: %v", err)
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
