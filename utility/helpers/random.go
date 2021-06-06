package helpers

import (
	"math/rand"
	"time"
)

var seededRand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func stringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func RandomNumber(length int) string {
	return stringWithCharset(length, "1234567890")
}

func RandomAlphabet(length int) string {
	return stringWithCharset(length, "ABCDEFGHIJKLMNOPQRSTUVWXYZ")
}

func RandomAlphaNum(length int) string {
	return stringWithCharset(length, "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
}

func RandomLowerAlphaNum(length int) string {
	return stringWithCharset(length, "abcdefghijklmnopqrstuvwxyz0123456789")
}

// without 0 to reduce O and 0 confusion
func RandomAlphaNumID(length int) string {
	return stringWithCharset(length, "ABCDEFGHIJKLMNOPQRSTUVWXYZ123456789")
}
