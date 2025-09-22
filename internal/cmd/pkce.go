package cmd

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
)

// generateCodeVerifier generates a cryptographically random string for the code verifier.
func generateCodeVerifier() (string, error) {
	b := make([]byte, 32) // 32 bytes for 256 bits of entropy
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", fmt.Errorf("failed to generate random bytes for code verifier: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// generateCodeChallenge generates the code challenge from the code verifier.
// It uses SHA256 hash and base64 URL encoding.
func generateCodeChallenge(verifier string) string {
	hasher := sha256.New()
	hasher.Write([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(hasher.Sum(nil))
}
