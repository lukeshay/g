package auth

// Generator is an interface for generating random strings. This can be used to
// generate session IDs or random tokens. These implementations are based of of
// [The Copenhagen Book](https://thecopenhagenbook.com/). We recommend reading
// the book to better understand authentication and security, and why these are
// implemented the way they are.
type Generator interface {
	// Generate generates a random string with the given number of bytes. This can
	// be used to generate a session ID or other random tokens. It is recommended
	// that the tokens have an entropy of at least 120 bits. This is outlined in
	// The Coopenhagen Book.
	Generate() (string, error)
}
