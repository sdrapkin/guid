package guid

import (
	"io"
	mathRandv2 "math/rand/v2"
	"sync"
)

type readerLite struct{}       // implements io.Reader interface
var _ io.Reader = readerLite{} //Compile-time interface assertion
var _readerLite readerLite = readerLite{}

// ReaderLite is a high-throughput source of cryptographically secure random bytes.
// It uses ChaCha8 cryptographically strong prng from "math/rand/v2".
// Each pooled instance is seeded with a 256-bit key from "crypto/rand".
// ReaderLite is automatically reseeded every 1024 bytes (key erasure).
// https://pkg.go.dev/internal/chacha8rand
// https://github.com/C2SP/C2SP/blob/main/chacha8rand.md
// https://go.dev/blog/chacha8rand#the-chacha8rand-generator
var ReaderLite readerLite = _readerLite

// liteRandPool is a sync.Pool for recycling "*mathRandv2.ChaCha8" instances.
// This reduces the overhead of repeatedly allocating and garbage collecting PRNGs.
var liteRandPool = sync.Pool{
	New: func() any {
		return newLiteRandGenerator()
	},
}

// newLiteRandGenerator creates and seeds a new mathRandv2.ChaCha8 PRNG.
// The seed is obtained from the package's primary cryptographically secure reader (_reader).
func newLiteRandGenerator() *mathRandv2.ChaCha8 {
	var seed [32]byte
	_reader.Read(seed[:])
	return mathRandv2.NewChaCha8(seed)
}

//==============================================
// readerLite Extension Methods
//==============================================

// Read fills the provided byte slice 'b' with cryptographically secure random bytes.
// It leverages a pool of pre-seeded ChaCha8 generators for high performance.
// The method always fills 'b' entirely and returns len(b) and a nil error.
func (r readerLite) Read(b []byte) (int, error) {
	n := len(b)
	if n == 0 {
		return 0, nil
	}
	chacha8Ptr := liteRandPool.Get().(*mathRandv2.ChaCha8)
	chacha8Ptr.Read(b) //chacha8.Read reads exactly len(p) bytes into p. It always returns len(p) and a nil error.
	liteRandPool.Put(chacha8Ptr)
	return n, nil
}
