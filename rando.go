package main

import (
	"crypto/rand"
	"encoding/base64"
)

func generateRandomBytes(length int) []byte {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		panic("Failed to generate secure random bytes! This should never happen.")
	}
	return bytes
}

func generateRandomString(bytesLength int) string {
	bytes := generateRandomBytes(bytesLength)
	return base64.URLEncoding.EncodeToString(bytes)
}
