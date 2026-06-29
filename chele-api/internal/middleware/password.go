package middleware

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/pbkdf2"
)

// DjangoHash represents a parsed Django password hash
type DjangoHash struct {
	Algorithm  string
	Iterations int
	Salt       string
	Hash       string
}

func ParseDjangoHash(encoded string) (*DjangoHash, error) {
	parts := strings.Split(encoded, "$")
	if len(parts) != 4 {
		return nil, fmt.Errorf("invalid django hash format")
	}
	iter, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid iteration count: %w", err)
	}
	return &DjangoHash{
		Algorithm:  parts[0],
		Iterations: iter,
		Salt:       parts[2],
		Hash:       parts[3],
	}, nil
}

func CheckDjangoPassword(password, encoded string) bool {
	dh, err := ParseDjangoHash(encoded)
	if err != nil {
		return false
	}
	switch dh.Algorithm {
	case "pbkdf2_sha256":
		dk := pbkdf2.Key([]byte(password), []byte(dh.Salt), dh.Iterations, 32, sha256.New)
		// Django stores hash in base64, Go-generated hashes use hex
		computedHex := hex.EncodeToString(dk)
		if computedHex == dh.Hash {
			return true
		}
		// Django uses base64 with padding (StdEncoding)
		if base64.StdEncoding.EncodeToString(dk) == dh.Hash {
			return true
		}
		// Fallback: base64 without padding
		return base64.RawStdEncoding.EncodeToString(dk) == dh.Hash
	default:
		return false
	}
}

func MakeDjangoPassword(password string) string {
	salt := make([]byte, 12)
	rand.Read(salt)
	saltStr := base64.RawStdEncoding.EncodeToString(salt)
	dk := pbkdf2.Key([]byte(password), []byte(saltStr), 720000, 32, sha256.New)
	return fmt.Sprintf("pbkdf2_sha256$720000$%s$%s", saltStr, hex.EncodeToString(dk))
}
