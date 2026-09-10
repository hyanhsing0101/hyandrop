package room

import (
	"crypto/rand"
	"encoding/hex"
)

func randomHex(byteLen int) (string, error) {
	bytes := make([]byte, byteLen)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}
