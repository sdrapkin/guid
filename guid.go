package guid

import (
	cryptoRand "crypto/rand"
)

// Size of a Guid in bytes.
const GuidByteSize = 16

// 16-byte (128-bit) cryptographically random value.
type Guid [GuidByteSize]byte

// Empty Guid
var Nil Guid

// New generates a new cryptographically secure Guid.
func New() (guid Guid) {
	// cryptoRand.Read in Go 1.24 is guaranteed to never fail and always fill b fully (https://pkg.go.dev/crypto/rand#Read)
	cryptoRand.Read(guid[:])
	return
} //New()
