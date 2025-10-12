package utils

import (
	"crypto/rand"
	"encoding/base64"
	"strings"

	"golang.org/x/crypto/argon2"
)

// HashPassword generates an argon2id hash for a password and returns encoded salt.hash
func HashPassword(password string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	return base64.RawStdEncoding.EncodeToString(salt) + "." + base64.RawStdEncoding.EncodeToString(hash), nil
}

// ComparePassword compares encoded hash (salt.hash) with a plaintext password
func ComparePassword(encoded, password string) bool {
	parts := strings.Split(encoded, ".")
	if len(parts) != 2 {
		return false
	}
	salt, err := base64.RawStdEncoding.DecodeString(parts[0])
	if err != nil {
		return false
	}
	expected, err := base64.RawStdEncoding.DecodeString(parts[1])
	if err != nil {
		return false
	}
	h := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	if len(h) != len(expected) {
		return false
	}
	var res byte
	for i := range h {
		res |= h[i] ^ expected[i]
	}
	return res == 0
}
