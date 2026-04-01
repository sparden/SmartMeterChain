package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

// HashData computes SHA-256 hash of any data structure
func HashData(data interface{}) string {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	hash := sha256.Sum256(jsonBytes)
	return hex.EncodeToString(hash[:])
}

// HashString computes SHA-256 of a string
func HashString(s string) string {
	hash := sha256.Sum256([]byte(s))
	return hex.EncodeToString(hash[:])
}

// VerifyHash checks if data matches an expected hash
func VerifyHash(data interface{}, expectedHash string) bool {
	return HashData(data) == expectedHash
}

// GenerateReadingHash creates a deterministic hash for a meter reading
func GenerateReadingHash(meterID string, reading float64, timestamp int64) string {
	input := fmt.Sprintf("%s:%.4f:%d", meterID, reading, timestamp)
	return HashString(input)
}
