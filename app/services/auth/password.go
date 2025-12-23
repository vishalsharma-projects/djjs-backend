package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	// Argon2id parameters - tuned for security and performance
	// Memory: 64MB, Iterations: 3, Parallelism: 4, Salt length: 16, Key length: 32
	memory      uint32 = 64 * 1024 // 64 MB
	iterations  uint32 = 3
	parallelism uint8  = 4
	saltLength  uint32 = 16
	keyLength   uint32 = 32
)

var (
	ErrInvalidHash         = errors.New("invalid hash format")
	ErrIncompatibleVersion = errors.New("incompatible argon2 version")
)

// HashPassword hashes a password using Argon2id with secure parameters
func HashPassword(password string) (string, error) {
	// Generate random salt
	salt := make([]byte, saltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	// Hash password with Argon2id
	hash := argon2.IDKey([]byte(password), salt, iterations, memory, parallelism, keyLength)

	// Encode hash with format: $argon2id$v=<version>$m=<memory>,t=<iterations>,p=<parallelism>$<salt>$<hash>
	// Base64 encode salt and hash
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, memory, iterations, parallelism, b64Salt, b64Hash)

	return encodedHash, nil
}

// VerifyPassword verifies a password against a hash using constant-time comparison
func VerifyPassword(password, encodedHash string) (bool, error) {
	// Parse the encoded hash
	parts, err := parseHash(encodedHash)
	if err != nil {
		return false, err
	}

	// Decode salt
	salt, err := base64.RawStdEncoding.DecodeString(parts.salt)
	if err != nil {
		return false, fmt.Errorf("failed to decode salt: %w", err)
	}

	// Decode expected hash
	expectedHash, err := base64.RawStdEncoding.DecodeString(parts.hash)
	if err != nil {
		return false, fmt.Errorf("failed to decode hash: %w", err)
	}

	// Hash the provided password with the same parameters
	actualHash := argon2.IDKey([]byte(password), salt, parts.iterations, parts.memory, parts.parallelism, uint32(len(expectedHash)))

	// Constant-time comparison
	if subtle.ConstantTimeCompare(expectedHash, actualHash) == 1 {
		return true, nil
	}
	return false, nil
}

type hashParts struct {
	version     int
	memory      uint32
	iterations  uint32
	parallelism uint8
	salt        string
	hash        string
}

func parseHash(encodedHash string) (*hashParts, error) {
	// Format: $argon2id$v=<version>$m=<memory>,t=<iterations>,p=<parallelism>$<salt>$<hash>
	// Split by $ to parse the hash format
	partsArr := strings.Split(encodedHash, "$")
	if len(partsArr) < 6 || partsArr[1] != "argon2id" {
		return nil, ErrInvalidHash
	}

	var parts hashParts
	var version int

	// Parse version: v=19
	_, err := fmt.Sscanf(partsArr[2], "v=%d", &version)
	if err != nil {
		return nil, ErrInvalidHash
	}

	if version != argon2.Version {
		return nil, ErrIncompatibleVersion
	}
	parts.version = version

	// Parse parameters: m=65536,t=3,p=4
	_, err = fmt.Sscanf(partsArr[3], "m=%d,t=%d,p=%d", &parts.memory, &parts.iterations, &parts.parallelism)
	if err != nil {
		return nil, ErrInvalidHash
	}

	// Salt and hash are the remaining parts (base64 encoded, may contain + and /)
	parts.salt = partsArr[4]
	parts.hash = partsArr[5]

	return &parts, nil
}

