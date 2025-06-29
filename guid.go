package guid

import (
	cryptoRand "crypto/rand"
	"sync"
)

const (
	GuidByteSize      = 16                               // Size of a Guid in bytes
	GuidCacheByteSize = 4096                             // 4096 bytes per cache
	GuidsPerCache     = GuidCacheByteSize / GuidByteSize // 256 GUIDs per cache (4096/16)
)

// Ensure that the constants are correct.
var _ = map[bool]int{false: 0, GuidsPerCache == 256: 1}
var _ = map[bool]int{false: 0, GuidsPerCache*GuidByteSize == GuidCacheByteSize: 1}

// 16-byte (128-bit) cryptographically random value.
type Guid [GuidByteSize]byte

// Empty Guid (zero value for Guid)
var Nil Guid

// guidCache holds a 4096-byte buffer and a byte index for GUID allocation
type guidCache struct {
	buffer [GuidCacheByteSize]byte
	index  uint8
}

// guidCachePool is a sync.Pool that holds guidCache instances.
var guidCachePool = sync.Pool{
	New: func() any {
		return &guidCache{}
	},
}

// New generates a new cryptographically secure Guid.
func New() (guid Guid) {
	guidCacheRef := guidCachePool.Get().(*guidCache)

	if guidCacheRef.index == 0 {
		// Refill buffer if index wraps (Go 1.24+: cryptoRand.Read is guaranteed to succeed)
		cryptoRand.Read(guidCacheRef.buffer[:])

	}

	// Extract GUID at current index
	startPos := int(guidCacheRef.index) * GuidByteSize
	copy(guid[:], guidCacheRef.buffer[startPos:startPos+GuidByteSize])

	guidCacheRef.index++ // Increment index for next call, uint8 wraps from 255 to 0 automatically

	guidCachePool.Put(guidCacheRef)
	return
} //New()
