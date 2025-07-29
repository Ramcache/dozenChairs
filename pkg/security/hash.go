package security

import (
	"crypto/sha256"
	"encoding/hex"
)

func SHA256Sum(s string) string {
	hash := sha256.Sum256([]byte(s))
	return hex.EncodeToString(hash[:])
}
