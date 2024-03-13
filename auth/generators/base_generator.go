package generators

import (
	"crypto/rand"

	"github.com/lukeshay/g/auth"
)

type baseGenerator struct {
	encoder func([]byte) string
	bytes   uint
}

func newBaseGenerator(bytes uint, encoder func([]byte) string) auth.Generator {
	return &baseGenerator{
		encoder: encoder,
		bytes:   bytes,
	}
}

func (g *baseGenerator) Generate() (string, error) {
	sessionIdBytes := make([]byte, g.bytes)
	rand.Read(sessionIdBytes)
	value := g.encoder(sessionIdBytes)

	return value, nil
}
