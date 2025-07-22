package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
)

// ComputeHash computes a hash of the given data using the specified algorithm
func ComputeHash(data interface{}, algorithm string) (string, error) {
	// Convert data to canonical JSON
	canonical, err := CanonicalJSON(data)
	if err != nil {
		return "", fmt.Errorf("failed to canonicalize data: %w", err)
	}

	// Select hash algorithm
	var hasher hash.Hash
	switch algorithm {
	case "sha256":
		hasher = sha256.New()
	case "sha512":
		hasher = sha512.New()
	default:
		return "", fmt.Errorf("unsupported hash algorithm: %s", algorithm)
	}

	// Compute hash
	hasher.Write([]byte(canonical))
	hashBytes := hasher.Sum(nil)
	
	return hex.EncodeToString(hashBytes), nil
}

// ComputeSHA256 computes a SHA256 hash of the given data (convenience function)
func ComputeSHA256(data interface{}) (string, error) {
	return ComputeHash(data, "sha256")
}

// ComputeSHA512 computes a SHA512 hash of the given data (convenience function)
func ComputeSHA512(data interface{}) (string, error) {
	return ComputeHash(data, "sha512")
}

// VerifyHash verifies that the given data matches the provided hash
func VerifyHash(data interface{}, expectedHash string, algorithm string) (bool, error) {
	computedHash, err := ComputeHash(data, algorithm)
	if err != nil {
		return false, err
	}
	
	return computedHash == expectedHash, nil
}

// VerifySHA256 verifies that the given data matches the provided SHA256 hash (convenience function)
func VerifySHA256(data interface{}, expectedHash string) (bool, error) {
	return VerifyHash(data, expectedHash, "sha256")
}

// GenerateRandomHash generates a random 32-byte hash as a hex string
func GenerateRandomHash() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	
	return hex.EncodeToString(bytes), nil
}

// GenerateRandomHashWithSize generates a random hash of specified byte size
func GenerateRandomHashWithSize(size int) (string, error) {
	if size <= 0 {
		return "", fmt.Errorf("hash size must be positive, got: %d", size)
	}
	
	bytes := make([]byte, size)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	
	return hex.EncodeToString(bytes), nil
} 