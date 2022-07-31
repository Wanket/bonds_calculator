package util

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

const (
	TokenSize        = 32
	SharedSecretSize = 32
)

func GenerateSecret(size uint) ([]byte, error) {
	sharedSecret := make([]byte, size)
	if _, err := rand.Read(sharedSecret); err != nil {
		return nil, fmt.Errorf("failed to generate secret: %w", err)
	}

	return sharedSecret, nil
}

func GenerateSecretString(size uint) (string, error) {
	tokenBytes := make([]byte, size+uint(base64.RawURLEncoding.EncodedLen(int(size))))
	if _, err := rand.Read(tokenBytes[:size]); err != nil {
		return "", fmt.Errorf("failed to generate string secret: %w", err)
	}

	base64.RawURLEncoding.Encode(tokenBytes[size:], tokenBytes[:size])

	return ByteArrayToString(tokenBytes[size:]), nil
}
