package hash

import (
	"crypto/sha256"
	"fmt"
)

func GenerateSha256(rawData string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(rawData)))
}

func ValidateSh256(rawData, hash string) bool {
	p := fmt.Sprintf("%x", sha256.Sum256([]byte(rawData)))
	if p != hash {
		return false
	}
	return true
}
