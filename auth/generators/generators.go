package generators

import (
	"crypto/sha256"
)

func sha256Hash(token string) (string, error) {
	h := sha256.New()

	if _, err := h.Write([]byte(token)); err != nil {
		return "", err
	}

	return string(h.Sum(nil)), nil
}
