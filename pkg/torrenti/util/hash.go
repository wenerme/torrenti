package util

import (
	"crypto/sha256"
	"encoding/hex"
)

func ContentHashBytes(v []byte) string {
	sum := sha256.Sum256(v)
	return hex.EncodeToString(sum[:])
}
