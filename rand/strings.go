package rand

import (
	"crypto/rand"
	"encoding/base64"
)

const RememberTokenBytes = 32

// Bytes will help to generate n random bytes, or will return an error if
// there was one. This uses the crypto/rand package so it is safe to use
// with things like remember tokens.
func Bytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// String will create a byte slice of size nBytes and then return a string
// that is a base64 url encoded version of that byte slice.
func String(nBytes int) (string, error) {
	b, err := Bytes(nBytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// RememberToken is a helper function designed to generate remember tokens
// of a predefined byte size.
func RememberToken() (string, error) {
	return String(RememberTokenBytes)
}
