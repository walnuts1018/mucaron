package filehash

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

func FileHash(path string) (string, error) {
	r, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}

	hash := sha256.New()
	if _, err := io.Copy(hash, r); err != nil {
		return "", fmt.Errorf("failed to copy file: %w", err)
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
