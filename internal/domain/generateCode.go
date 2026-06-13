package domain

import (
	"encoding/hex"

	"golang.org/x/crypto/blake2b"
)

func GenerateCode(url string) (string, error) {
	hash, err := blake2b.New256(nil)
	hash.Write([]byte(url))
	return hex.EncodeToString(hash.Sum(nil))[:CodeLength], err
}
