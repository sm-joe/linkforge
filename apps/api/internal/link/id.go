package link

import (
	"crypto/rand"
	"encoding/hex"
)

const shortCodeLength = 8

func generateID() string {
	buffer := make([]byte, 16)

	if _, err := rand.Read(buffer); err != nil {
		panic(err)
	}

	return hex.EncodeToString(buffer)
}

func generateShortCode() string {
	const alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	buffer := make([]byte, shortCodeLength)
	random := make([]byte, shortCodeLength)

	if _, err := rand.Read(random); err != nil {
		panic(err)
	}

	for i := range buffer {
		buffer[i] = alphabet[int(random[i])%len(alphabet)]
	}

	return string(buffer)
}
