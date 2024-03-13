package generators

import (
	"encoding/hex"

	"github.com/lukeshay/g/auth"
)

type HexGenerator struct {
	base auth.Generator
}

// NewHexGenerator creates a new hex generator. This is used to generate
// random strings, such as session IDs.
func NewHexGenerator(bytes uint) auth.Generator {
	return &HexGenerator{
		base: newBaseGenerator(bytes, hex.EncodeToString),
	}
}

func (g *HexGenerator) Generate() (string, error) {
	return g.base.Generate()
}
