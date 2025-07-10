// Package guid provides fast, efficient, cryptographically secure GUID generation and manipulation.
// It supports Base64Url encoding and decoding, and is optimized for performance.
package guid

import (
	cryptoRand "crypto/rand"
	"encoding"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sync"
	"unsafe"
)

//==============================================
// Constants
//==============================================

const (
	GuidByteSize          = 16                                                                 // Size of a Guid in bytes
	guidsPerCache         = 256                                                                // 256 Guids per cache - do not change this value
	guidCacheByteSize     = GuidByteSize * guidsPerCache                                       // 4096 bytes per cache (256*16)
	GuidBase64UrlByteSize = 22                                                                 // Base64Url encoding of a Guid is 22 characters
	base64UrlAlphabet     = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_" // Base64Url alphabet used for encoding
)

// Ensure that the constants are not changed without thought.
var _ = map[bool]int{false: 0, guidsPerCache == 256: 1}
var _ = map[bool]int{false: 0, guidCacheByteSize == 4096: 1}

//==============================================
// Errors and Variables
//==============================================

var (
	// Empty Guid (zero value for Guid)
	Nil Guid
	// Reader is a global, shared instance of a cryptographically secure random number generator. It is safe for concurrent use.
	Reader reader
)

var (
	// ErrInvalidBase64UrlGuidEncoding is returned when a Base64Url string does not represent a valid Guid.
	ErrInvalidBase64UrlGuidEncoding = errors.New("invalid Base64Url Guid encoding (invalid characters, or length != 22)")
	// ErrInvalidGuidSlice is returned when a byte slice cannot represent a valid Guid (length < 16 bytes).
	ErrInvalidGuidSlice = errors.New("invalid Guid slice (length < 16 bytes)")
)

//==============================================
// Compile-time interface assertions
//==============================================

var (
	_ fmt.Stringer               = (*Guid)(nil)
	_ encoding.TextMarshaler     = (*Guid)(nil)
	_ encoding.TextUnmarshaler   = (*Guid)(nil)
	_ encoding.BinaryMarshaler   = (*Guid)(nil)
	_ encoding.BinaryUnmarshaler = (*Guid)(nil)
	_ io.Reader                  = (*reader)(nil)
)

//==============================================
// Types
//==============================================

// Guid is a 16-byte (128-bit) cryptographically random value.
type Guid [GuidByteSize]byte
type reader struct{}

//==============================================
// Shared Variables
//==============================================

// guidCache holds a 4096-byte buffer and a byte index for Guid allocation.
type guidCache struct {
	buffer [guidCacheByteSize]byte
	index  uint8
	_      [63]byte // pad ensures each index is on its own cache line

}

//==============================================
// Guid Extension Methods
//==============================================

// MarshalBinary implements the encoding.BinaryMarshaler interface for Guid.
func (guid *Guid) MarshalBinary() (data []byte, err error) {
	return guid[:], nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface for Guid.
func (guid *Guid) UnmarshalBinary(data []byte) error {
	if len(data) < GuidByteSize {
		return ErrInvalidGuidSlice
	}
	copy(guid[:], data)
	return nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (guid *Guid) UnmarshalText(data []byte) error {
	if ok := DecodeBase64URL(guid[:], data); !ok {
		return ErrInvalidBase64UrlGuidEncoding
	}
	return nil
}

// MarshalText implements encoding.TextMarshaler.
func (guid *Guid) MarshalText() ([]byte, error) {
	buffer := make([]byte, GuidBase64UrlByteSize)
	guid.EncodeBase64URL(buffer)
	return buffer, nil
}

// MarshalJSON implements the json.Marshaler interface.
// It marshals the Guid to its Base64Url string representation.
func (g Guid) MarshalJSON() ([]byte, error) {
	return json.Marshal(g.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It unmarshals a JSON string into a Guid.
func (g *Guid) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("guid: cannot unmarshal JSON string %s into a Guid: %w", string(data), err)
	}

	parsedGuid, err := Parse(s)
	if err != nil {
		return err // The error from Parse is already informative.
	}

	*g = parsedGuid
	return nil
}

// String returns a Base64Url-encoded string representation of the Guid.
func (guid *Guid) String() string {
	buffer, _ := guid.MarshalText()
	return unsafe.String(&buffer[0], GuidBase64UrlByteSize)

	// Is it "safe" to use unsafe.String() here? Yes.
	// buffer will be allocated on the heap, and will not be gc'ed until string is alive.
	// This is the same approach that Golang uses in "strings.Clone()" [https://pkg.go.dev/strings#Clone],
	// which calls internal "stringslite.Clone()":
	// https://cs.opensource.google/go/go/+/refs/tags/go1.24.4:src/internal/stringslite/strings.go;l=143
	/* stringslite.Clone():
		func Clone(s string) string {
		if len(s) == 0 {
			return ""
		}
		b := make([]byte, len(s))
		copy(b, s)
		return unsafe.String(&b[0], len(b))
	}*/
}

// EncodeBase64URL encodes the Guid into the provided dst as Base64Url.
func (guid *Guid) EncodeBase64URL(dst []byte) {
	const lengthMod3 = 1                    // 16 % 3 = 1
	const limit = GuidByteSize - lengthMod3 // 15 bytes can be processed in groups of 3 bytes, leaving 1 byte at the end.

	j := 0 // Index in the output buffer

	// Process the first 15 bytes (5 groups of 3 bytes). Each 3-byte group is converted to 4 Base64Url characters.
	for i := 0; i < limit; i += 3 {
		val := uint(guid[i])<<16 | uint(guid[i+1])<<8 | uint(guid[i+2])

		// Combine 3 bytes into a 24-bit integer and extract 4 6-bit indices.
		dst[j] = base64UrlAlphabet[val>>18&0x3F]
		dst[j+1] = base64UrlAlphabet[val>>12&0x3F]
		dst[j+2] = base64UrlAlphabet[val>>6&0x3F]
		dst[j+3] = base64UrlAlphabet[val&0x3F]
		j += 4
	}

	// Handle the last byte, converted to 2 Base64Url characters.
	b0 := guid[limit]
	dst[j] = base64UrlAlphabet[b0>>2]
	dst[j+1] = base64UrlAlphabet[(b0&0x03)<<4]
}

//==============================================
// reader Extension Methods
//==============================================

// Read fills b with cryptographically secure random bytes.
// It always fills b entirely, and returns len(b) and nil error.
// guid.Read() is up to 7x faster than crypto/rand.Read() for small slices.
// if b is > 512 bytes, it simply calls crypto/rand.Read().
func (r *reader) Read(b []byte) (n int, err error) {
	const MaxBytesToFillViaGuids = 512
	n = len(b)

	if n == 0 {
		return
	}

	if n > MaxBytesToFillViaGuids {
		return cryptoRand.Read(b)
	}

	guidCacheRef := guidCachePool.Get().(*guidCache)

	if guidCacheRef.index == 0 {
		// Refill buffer if index wraps (Go 1.24+: cryptoRand.Read is guaranteed to succeed)
		cryptoRand.Read(guidCacheRef.buffer[:])
	}

	n1 := copy(b, guidCacheRef.buffer[int(guidCacheRef.index)*GuidByteSize:])

	if n1 == n {
		// Case 1: Everything fit in one copy
		guidCacheRef.index += byte((n1 + GuidByteSize - 1) / GuidByteSize)
		guidCachePool.Put(guidCacheRef)
		return
	}

	// Case 2: Need to refill for remainder
	cryptoRand.Read(guidCacheRef.buffer[:])
	n2 := copy(b[n1:], guidCacheRef.buffer[:])

	if n1+n2 != n {
		panic("guid: internal panic in reader.Read(); should never happen")
	}
	guidCacheRef.index = byte((n2 + GuidByteSize - 1) / GuidByteSize)
	guidCachePool.Put(guidCacheRef)
	return
} //func (r *reader) Read

//==============================================
// Standalone Functions
//==============================================

// New generates a new cryptographically secure Guid.
func New() (g Guid) {
	guidCacheRef := guidCachePool.Get().(*guidCache)

	if guidCacheRef.index == 0 {
		// Refill buffer if index wraps (Go 1.24+: cryptoRand.Read is guaranteed to succeed)
		cryptoRand.Read(guidCacheRef.buffer[:])
	}

	// Extract GUID at current index
	startPos := int(guidCacheRef.index) * GuidByteSize
	copy(g[:], guidCacheRef.buffer[startPos:startPos+GuidByteSize])

	guidCacheRef.index++ // Increment index for next call, uint8 wraps from 255 to 0 automatically
	guidCachePool.Put(guidCacheRef)
	return
}

// NewString generates a new cryptographically secure Guid, and returns it as a Base64Url string.
// NewString is equivalent to "g := guid.New(); return g.String();".
func NewString() string {
	g := New()
	return g.String()
}

// Parse parses a Base64Url-encoded string into the Guid.
// Returns an error if the string is not a valid Guid encoding.
func Parse(s string) (g Guid, err error) {
	if len(s) != GuidBase64UrlByteSize {
		return Guid{}, ErrInvalidBase64UrlGuidEncoding
	}

	// Zero-copy conversion of a string to a byte slice
	sBytes := unsafe.Slice(unsafe.StringData(s), GuidBase64UrlByteSize)

	if ok := DecodeBase64URL(g[:], sBytes); !ok {
		return Guid{}, ErrInvalidBase64UrlGuidEncoding
	}
	return g, nil
}

// ParseBytes parses a Base64Url-encoded string represented as a byte slice into the Guid.
// Returns an error if the string byte slice is not a valid Guid encoding.
// ParseBytes is like Parse, except it parses a string byte slice instead of a string.
func ParseBytes(src []byte) (g Guid, err error) {
	if len(src) != GuidBase64UrlByteSize {
		return Guid{}, ErrInvalidBase64UrlGuidEncoding
	}

	if ok := DecodeBase64URL(g[:], src); !ok {
		return Guid{}, ErrInvalidBase64UrlGuidEncoding
	}
	return g, nil
}

// FromBytes returns a Guid from a 16-byte slice.
func FromBytes(src []byte) (Guid, error) {
	if len(src) < GuidByteSize {
		return Guid{}, ErrInvalidGuidSlice
	}
	var g Guid
	copy(g[:], src)
	return g, nil
}

// DecodeBase64URL decodes a Base64Url-encoded src byte slice into a Guid dst byte slice.
// Does not panic on invalid input.
// dst must be at least 16 bytes long and src must be at least 22 bytes long (returns false otherwise).
// dst is modified even if the function returns false.
func DecodeBase64URL(dst []byte, src []byte) (ok bool) {
	if (len(dst) < GuidByteSize) || (len(src) < GuidBase64UrlByteSize) {
		return false
	}

	const lengthMod3 = 1 // 16 % 3 = 1
	const limit = GuidByteSize - lengthMod3

	j := 0

	// Process 5 groups of 4 characters to 3 bytes
	for i := 0; i < limit; i += 3 {
		b0 := decodeLookup[src[j]]
		b1 := decodeLookup[src[j+1]]
		b2 := decodeLookup[src[j+2]]
		b3 := decodeLookup[src[j+3]]

		if (b0 | b1 | b2 | b3) >= 64 {
			return false
		}

		dst[i] = (b0 << 2) | (b1 >> 4)
		dst[i+1] = (b1 << 4) | (b2 >> 2)
		dst[i+2] = (b2 << 6) | b3
		j += 4
	}

	// Handle the remaining 2 characters to 1 byte
	b0 := decodeLookup[src[j]]
	b1 := decodeLookup[src[j+1]]

	if (b0 | b1) >= 64 {
		return false
	}

	dst[limit] = (b0 << 2) | (b1 >> 4)
	return true
}

// Read fills b with cryptographically secure random bytes.
// It never returns an error, and always fills b entirely.
// guid.Read() is up to 7x faster than crypto/rand.Read() for small slices.
// if b is > 512 bytes, it simply calls crypto/rand.Read().
func Read(b []byte) (n int, err error) {
	return Reader.Read(b)
}

//==============================================
// Internal Variables
//==============================================

// guidCachePool is a sync.Pool that holds guidCache instances.
var guidCachePool = sync.Pool{
	New: func() any {
		return &guidCache{}
	},
}

/********************************************************
	c# code to generate the decodeLookup table:
	Span<byte> decodeLookup = stackalloc byte[byte.MaxValue+1];
	decodeLookup.Fill(0xFF); // Initialize with invalid value

	for (var i = 0; i < BASE64URL_ALPHABET_STRING.Length; i++)
	{ decodeLookup[BASE64URL_ALPHABET_STRING[i]] = (byte)i;	}

	"[".Dump();
	for (var i = 0; i < decodeLookup.Length; ++i)
		{ Console.Write($"0x{decodeLookup[i]:X2},"); if ((i + 1) % 8 == 0) "".Dump(); }
	"]".Dump();
*********************************************************/

// decodeLookup is a lookup table for decoding Base64Url characters to their byte values.
// Values outside the Base64Url alphabet are marked with 0xFF.
var decodeLookup = [256]byte{
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x3E, 0xFF, 0xFF,
	0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3A, 0x3B,
	0x3C, 0x3D, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06,
	0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E,
	0x0F, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16,
	0x17, 0x18, 0x19, 0xFF, 0xFF, 0xFF, 0xFF, 0x3F,
	0xFF, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F, 0x20,
	0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28,
	0x29, 0x2A, 0x2B, 0x2C, 0x2D, 0x2E, 0x2F, 0x30,
	0x31, 0x32, 0x33, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
}
