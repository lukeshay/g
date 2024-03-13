package auth

// Encrypter is responsible for encrypting and decrypting strings. This can be
// used to encrypt and decrypt session IDs, user IDs, etc. It is up to you to
// implement this interface for your specific use case.
type Encrypter interface {
	// Encrypt encrypts the given string and returns the encrypted string. The
	// returned string should be base64 encoded.
	Encrypt(string) (string, error)
	// Decrypt decrypts the given string and returns the decrypted string. The
	// given string should be base64 encoded.
	Decrypt(string) (string, error)
}
