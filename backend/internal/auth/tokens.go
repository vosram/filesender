package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

// Creates a 32 byte random hexcidecimal string.
func CreateOpaqueToken() string {
	refreshToken := make([]byte, 32)
	rand.Read(refreshToken)
	return hex.EncodeToString(refreshToken)
}

// Creates a SHA256 hash of the string input. Should be
// used for 32 byte random opaque tokens. Like those saved
// to the database.
func CreateTokenHash(token string) (string, error) {
	byteToken, err := hex.DecodeString(token)
	if err != nil {
		return "", err
	}
	tokenHash := sha256.Sum256(byteToken)
	return hex.EncodeToString(tokenHash[:]), nil
}
