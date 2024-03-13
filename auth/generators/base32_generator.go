package generators

import (
	"encoding/base32"

	"github.com/lukeshay/g/auth"
)

type Base32Generator struct {
	base auth.Generator
}

// NewBase32Generator creates a new base32 generator. This is used to generate
// random strings, such as session IDs.
func NewBase32Generator(bytes uint) auth.Generator {
	return &Base32Generator{
		base: newBaseGenerator(bytes, base32.StdEncoding.EncodeToString),
	}
}

func (g *Base32Generator) Generate() (string, error) {
	return g.base.Generate()
}
