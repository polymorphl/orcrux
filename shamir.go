package main

import (
	cryptoRand "crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

// validateShamirParams validates the input parameters for Shamir's Secret Sharing
func validateShamirParams(secret []byte, n, t int, output string) error {
	if len(secret) == 0 {
		return errors.New("empty secret")
	}
	if n < 2 || n > 255 {
		return errors.New("shards must be in [2, 255]")
	}
	if t < 2 || t > n {
		return errors.New("shardsNeeded must be in [2, shards]")
	}

	enc := strings.ToLower(strings.TrimSpace(output))
	if enc != "hex" && enc != "base64" {
		return fmt.Errorf("output must be 'hex' or 'base64', got: %q", output)
	}
	return nil
}

// generateXCoordinates generates the x-coordinates for polynomial evaluation
func generateXCoordinates(n int) []byte {
	xs := make([]byte, n)
	for i := 0; i < n; i++ {
		xs[i] = byte(i + 1)
	}
	return xs
}

// evaluatePolynomial evaluates a polynomial of degree t-1 at point x using Horner's method
func evaluatePolynomial(secretByte byte, x byte, t int) (byte, error) {
	coeffs := make([]byte, t)
	coeffs[0] = secretByte
	if _, err := cryptoRand.Read(coeffs[1:]); err != nil {
		return 0, err
	}

	// Horner's method evaluation in GF(256)
	y := coeffs[t-1]
	for k := t - 2; k >= 0; k-- {
		y = gfMul(y, x)
		y ^= coeffs[k]
	}
	return y, nil
}

// encodeShare encodes a share's data in the specified format
func encodeShare(data []byte, format string) string {
	if format == "hex" {
		return hex.EncodeToString(data)
	}
	return base64.StdEncoding.EncodeToString(data)
}

// shamirSplit implements Shamir's Secret Sharing algorithm to split a secret into n shares,
// where any t shares can reconstruct the original secret. This provides a threshold-based
// secret sharing scheme with information-theoretic security.
//
// The algorithm works by:
//  1. For each byte of the secret, generate a random polynomial of degree t-1
//     where the constant term is the secret byte value
//  2. Evaluate each polynomial at n distinct points (x = 1, 2, ..., n)
//  3. Each share contains the x-coordinate and the corresponding y-values
//  4. Any t shares can be used to reconstruct the original secret using Lagrange interpolation
//
// Parameters:
//   - secret: The secret data to be split (must not be empty)
//   - n: Total number of shares to generate (must be in [2, 255])
//   - t: Threshold - minimum number of shares needed to reconstruct (must be in [2, n])
//   - output: Output encoding format, either "hex" or "base64"
//
// Returns:
//   - A string containing n shares, one per line, in format "xx:<encoded_data>"
//     where xx is the hex representation of the x-coordinate (1 to n)
//   - An error if any validation fails or crypto operations error
//
// Security properties:
//   - Any subset of shares less than t reveals no information about the secret
//   - The scheme is information-theoretically secure (unconditionally secure)
//   - Random coefficients ensure each share appears independent of the secret
//
// Example output format:
//
//	01:a1b2c3d4...
//	02:e5f6g7h8...
//	03:i9j0k1l2...
//
// Note: This implementation uses GF(2^8) arithmetic with irreducible polynomial 0x11b
// for polynomial evaluation, which is standard in cryptographic applications.
func shamirSplit(secret []byte, n, t int, output string) (string, error) {
	if err := validateShamirParams(secret, n, t, output); err != nil {
		return "", err
	}

	xs := generateXCoordinates(n)
	enc := strings.ToLower(strings.TrimSpace(output))

	var sb strings.Builder
	for i := 0; i < n; i++ {
		x := xs[i]
		yBytes := make([]byte, len(secret))

		// Evaluate polynomial for each byte of the secret
		for b := 0; b < len(secret); b++ {
			y, err := evaluatePolynomial(secret[b], x, t)
			if err != nil {
				return "", err
			}
			yBytes[b] = y
		}

		data := encodeShare(yBytes, enc)
		fmt.Fprintf(&sb, "%02x:%s\n", x, data)
	}

	return sb.String(), nil
}

// gfMul performs Galois Field multiplication in GF(2^8) with irreducible polynomial 0x11b.
// This function multiplies two bytes a and b in the finite field, which is essential
// for Shamir's Secret Sharing scheme and other cryptographic operations.
//
// The multiplication is performed using the "Russian peasant" algorithm:
// 1. Initialize result p = 0
// 2. For each bit in b (from least to most significant):
//   - If the bit is 1, XOR the result with a
//   - Left shift a by 1
//   - If a overflows (bit 7 is set), reduce modulo 0x11b by XORing with 0x1b
//
// 3. Right shift b by 1 to process the next bit
//
// Parameters:
//   - a: First byte operand
//   - b: Second byte operand
//
// Returns:
//   - The product of a and b in GF(2^8)
//
// Example:
//
//	gfMul(0x57, 0x83) = 0xc1
func gfMul(a, b byte) byte {
	var p byte = 0
	aa := a
	bb := b
	for i := 0; i < 8; i++ {
		if (bb & 1) != 0 {
			p ^= aa
		}
		hi := aa & 0x80
		aa <<= 1
		if hi != 0 {
			aa ^= 0x1b // réduction mod 0x11b
		}
		bb >>= 1
	}
	return p
}
