package generators

import (
	"encoding/base64"

	"github.com/lukeshay/g/auth"
)

type Base64Generator struct {
	base auth.Generator
}

// NewBase64Generator creates a new base64 generator. This is used to generate
// random strings, such as session IDs.
func NewBase64Generator(bytes uint) auth.Generator {
	return &Base64Generator{
		base: newBaseGenerator(bytes, base64.StdEncoding.EncodeToString),
	}
}

func (g *Base64Generator) Generate() (string, error) {
	return g.base.Generate()
}
