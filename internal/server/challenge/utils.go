package challenge

import (
	"crypto/rand"
	"encoding/hex"
)

// generateRandomHex возвращает случайные N байт в hex-строке
func generateRandomHex(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
