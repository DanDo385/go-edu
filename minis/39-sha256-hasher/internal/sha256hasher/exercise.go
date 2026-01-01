//go:build !solution && !reference

package sha256hasher

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"io"
	"os"
	"strings"
)

// HashString implements the exercise.
//
// TODO: Implement this function
func HashString(s string) string {
	// TODO: Implement
	return ""
}

// HashFile implements the exercise.
//
// TODO: Implement this function
func HashFile(filename string) ([]byte, error) {
	// TODO: Implement
	return nil, nil
}

// VerifyFile implements the exercise.
//
// TODO: Implement this function
func VerifyFile(filename string, expectedHashHex string) (bool, error) {
	// TODO: Implement
	return false, nil
}

// HashIncremental implements the exercise.
//
// TODO: Implement this function
func HashIncremental(parts ...string) string {
	// TODO: Implement
	return ""
}

// CompareHashes implements the exercise.
//
// TODO: Implement this function
func CompareHashes(hash1 string, hash2 string) bool {
	// TODO: Implement
	return false
}
