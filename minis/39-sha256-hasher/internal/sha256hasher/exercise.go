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

func HashString(s string) string {
	// TODO: Implement this function
	// TODO: Handle errors appropriately
	// TODO: Add necessary validations
	panic("not implemented")
}

func HashFile(filename string) ([]byte, error) {
	// TODO: Implement this function
	panic("not implemented")
}

func VerifyFile(filename, expectedHashHex string) (bool, error) {
	// TODO: Implement this function
	panic("not implemented")
}

func HashIncremental(parts ...string) string {
	// TODO: Implement this function
	panic("not implemented")
}

func CompareHashes(hash1, hash2 string) bool {
	// TODO: Implement this function
	panic("not implemented")
}
