package crypt

import (
	"crypto/sha256"
	"encoding/hex"
)

func NewSHA256(data, salt string) string {
	hash := sha256.Sum256([]byte(data + salt))
	return hex.EncodeToString(hash[:])
}

func VerifyPassword(hashed, plain, salt string) bool {
	return hashed == NewSHA256(plain, salt)
}
