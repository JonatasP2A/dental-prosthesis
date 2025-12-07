package uuid

import (
	"crypto/rand"
	"fmt"
)

// Generator generates UUIDs
type Generator struct{}

// NewGenerator creates a new UUID generator
func NewGenerator() *Generator {
	return &Generator{}
}

// Generate generates a new UUID v4
func (g *Generator) Generate() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		panic(err) // Should never happen
	}

	// Set version (4) and variant (RFC 4122)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80

	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

