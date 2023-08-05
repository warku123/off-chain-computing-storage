package utils

import (
	"crypto/rand"
	"encoding/base64"
	"os"
)

// GenerateRandomID generates a random ID with a given length
func GenerateRandomID(length int) (string, error) {
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(randomBytes)[:length], nil
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
