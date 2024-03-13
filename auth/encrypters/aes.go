package encrypters

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"io"

	"github.com/lukeshay/g/auth"
)

// AesEncrypter is an implementation of the Encrypter interface that uses AES
// encryption.
type AesEncrypter struct {
	gcmInstance cipher.AEAD
}

// NewAesEncrypter returns a new instance of AesEncrypter.
func NewAesEncrypter(secret string) (auth.Encrypter, error) {
	aesBlock, err := aes.NewCipher([]byte(md5Hash(secret)))
	if err != nil {
		return nil, err
	}

	gcmInstance, err := cipher.NewGCM(aesBlock)
	if err != nil {
		return nil, err
	}

	return &AesEncrypter{
		gcmInstance: gcmInstance,
	}, nil
}

func (e *AesEncrypter) Encrypt(plaintext string) (string, error) {
	nonce := make([]byte, e.gcmInstance.NonceSize())
	_, err := io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return "", err
	}

	value := e.gcmInstance.Seal(nonce, nonce, []byte(plaintext), nil)
	result := base64.StdEncoding.EncodeToString(value)

	return result, nil
}

func (e *AesEncrypter) Decrypt(encrypted string) (string, error) {
	ciphered, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", err
	}

	nonceSize := e.gcmInstance.NonceSize()
	nonce, ciphertext := ciphered[:nonceSize], ciphered[nonceSize:]

	originalText, err := e.gcmInstance.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(originalText), nil
}

func md5Hash(input string) string {
	byteInput := []byte(input)
	md5Hash := md5.Sum(byteInput)
	return hex.EncodeToString(md5Hash[:])
}
