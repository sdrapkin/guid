// Package guid provides fast, efficient, cryptographically secure 128-bit GUID generation and manipulation.
// It supports Base64Url encoding/decoding, sequential sortable GUIDs (for PostgreSQL and SQL Server), and is optimized for performance.
// It includes a high-throughput, drop-in replacement for crypto/rand.Reader for generating secure random bytes.
package guid

import (
	cryptoRand "crypto/rand"
	"encoding"
	"errors"
	"fmt"
	"io"
	"math/bits"
	"sync"
	"time"
	"unsafe"

	"golang.org/x/sys/cpu"
)

//==============================================
// Constants
//==============================================

const (
	GuidByteSize          = 16                           // Size of a Guid in bytes
	guidsPerCache         = 256                          // 256 Guids per cache - do not change this value
	guidCacheByteSize     = GuidByteSize * guidsPerCache // 4096 bytes per cache (256*16)
	GuidBase64UrlByteSize = 22                           // Base64Url encoding of a Guid is 22 characters
)

const (
	base64UrlAlphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_" // Base64Url alphabet used for encoding
)

// Ensure that the constants are not changed without thought.
var _ = map[bool]int{false: 0, guidsPerCache == 256: 1}
var _ = map[bool]int{false: 0, guidCacheByteSize == 4096: 1}

//==============================================
// Errors and Variables
//==============================================

var (
	_minGuid Guid
	_maxGuid Guid = Guid{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}
)

var (
	// Nil is the nil Guid, with all 128 bits set to zero.
	Nil Guid = _minGuid
	// Max is the maximum Guid, with all 128 bits set to one.
	Max Guid = _maxGuid
	// Reader is a global, shared instance of a cryptographically secure random number generator. It is safe for concurrent use.
	Reader reader = _reader
)

var (
	// ErrInvalidBase64UrlGuidEncoding is returned when a Base64Url string does not represent a valid Guid.
	ErrInvalidBase64UrlGuidEncoding = errors.New("invalid Base64Url Guid encoding (invalid characters, or length != 22)")
	// ErrInvalidGuidSlice is returned when a byte slice cannot represent a valid Guid (length < 16 bytes).
	ErrInvalidGuidSlice = errors.New("invalid Guid slice (length < 16 bytes)")
	// ErrBufferTooSmallBase64Url is returned when a destination slice is too small to receive the text-encoded Guid.
	ErrBufferTooSmallBase64Url = fmt.Errorf("buffer is too small (length < %d bytes)", GuidBase64UrlByteSize)
)

//==============================================
// Compile-time interface assertions
//==============================================

var (
	_ fmt.Stringer               = &Guid{}
	_ encoding.TextMarshaler     = &Guid{}
	_ encoding.TextUnmarshaler   = &Guid{}
	_ encoding.BinaryMarshaler   = Guid{}
	_ encoding.BinaryUnmarshaler = &Guid{}
	_ io.Reader                  = reader{}
)

//==============================================
// Types
//==============================================

// Guid is a 16-byte (128-bit) cryptographically random value.
type Guid [GuidByteSize]byte

// GuidPG is a 16-byte (128-bit) PostgreSQL sortable Guid formed as [8-byte time.Now() timestamp][8 random bytes]
// GuidPG is optimized for use as a PostgreSQL index key.
type GuidPG struct {
	Guid // embedded
}

// GuidSS is a 16-byte (128-bit) SQL Server sortable Guid formed as [8 random bytes][8 bytes of SQL Server ordered time.Now() timestamp]
// GuidSS is optimized for use as a SQL Server index or clustered key.
type GuidSS struct {
	Guid // embedded
}

type reader struct{} // implements io.Reader interface
var _reader reader = reader{}

//==============================================
// Shared Variables
//==============================================

// guidCache holds a 4096-byte buffer and a byte index for Guid allocation.
type guidCache struct {
	buffer []byte
	index  uint8
	_      [32]byte // pad ensures each index is on its own cache line
}

//==============================================
// Guid Extension Methods
//==============================================

// MarshalBinary implements the encoding.BinaryMarshaler interface for Guid.
func (guid Guid) MarshalBinary() (data []byte, err error) {
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
	guid.encodeBase64URL(buffer)
	return buffer, nil
}

// MarshalJSON implements the json.Marshaler interface.
// It marshals the Guid to its Base64Url string representation.
func (g Guid) MarshalJSON() ([]byte, error) {
	//return json.Marshal(g.String())
	gStringWithQuotes := make([]byte, GuidBase64UrlByteSize+2)
	gStringWithQuotes[1+GuidBase64UrlByteSize], gStringWithQuotes[0] = '"', '"'
	g.encodeBase64URL(gStringWithQuotes[1 : 1+GuidBase64UrlByteSize])
	return gStringWithQuotes, nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It unmarshals a JSON string into a Guid.
func (g *Guid) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*g = Guid{}
		return nil // valid null Guid
	}

	if len(data) != (GuidBase64UrlByteSize+2) || !DecodeBase64URL(g[:], data[1:1+GuidBase64UrlByteSize]) {
		return fmt.Errorf("guid: cannot unmarshal JSON string %q into a Guid", string(data))
	}
	return nil
}

// String returns a Base64Url-encoded string representation of the Guid.
func (guid *Guid) String() string {
	buffer := make([]byte, GuidBase64UrlByteSize)
	guid.encodeBase64URL(buffer)
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
func (guid *Guid) EncodeBase64URL(dst []byte) error {
	if len(dst) < GuidBase64UrlByteSize {
		return ErrBufferTooSmallBase64Url
	}
	guid.encodeBase64URL(dst)
	return nil
}

// private - panics on undersized buffer or nil guid
func (guid *Guid) encodeBase64URL(dst []byte) {
	const lengthMod3 = 1                    // 16 % 3 = 1
	const limit = GuidByteSize - lengthMod3 // 15 bytes can be processed in groups of 3 bytes, leaving 1 byte at the end.

	// Bounds Check Elimination
	_ = guid[GuidByteSize-1]
	_ = dst[GuidBase64UrlByteSize-1]

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
func (r reader) Read(b []byte) (int, error) {
	const MaxBytesToFillViaGuids = 512
	n := len(b)

	if n == 0 {
		return 0, nil
	}

	if n > MaxBytesToFillViaGuids {
		return cryptoRand.Read(b)
	}

	guidCacheRef := guidCachePool.Get().(*guidCache)

	if n > (guidCacheByteSize - int(guidCacheRef.index)*GuidByteSize) {
		cryptoRand.Read(guidCacheRef.buffer) // Not enough bytes remaining: refill completely. Go 1.24+ guarantees crypto/rand.Read succeeds.
		guidCacheRef.index = 0
	} else if guidCacheRef.index == 0 {
		cryptoRand.Read(guidCacheRef.buffer) // Refill buffer if index wraps (Go 1.24+: cryptoRand.Read is guaranteed to succeed)
	}

	copy(b, guidCacheRef.buffer[int(guidCacheRef.index)*GuidByteSize:])

	// Update the index based on the number of Guids consumed.
	// The ceiling division ensures the index increments correctly for partial Guid consumption.
	guidCacheRef.index += byte((n + GuidByteSize - 1) / GuidByteSize)

	guidCachePool.Put(guidCacheRef)
	return n, nil
} //func (r reader) Read

//==============================================
// GuidPG Extension Methods
//==============================================

// Timestamp extracts the timestamp from the PostgreSQL Guid.
// The timestamp is stored in the first 8 bytes as nanoseconds since Unix epoch.
// Returns the time.Time representation of when the Guid was created.
func (g *GuidPG) Timestamp() time.Time {
	timestamp := *(*int64)(unsafe.Pointer(&g.Guid[0])) // Extract timestamp from first 8 bytes

	if !cpu.IsBigEndian {
		timestamp = int64(bits.ReverseBytes64(uint64(timestamp)))
	}
	return time.Unix(0, timestamp).UTC()
}

//==============================================
// GuidSS Extension Methods
//==============================================

// Timestamp extracts the timestamp from the SQL Server Guid.
// The timestamp is stored in the last 8 bytes using SQL Server's Guid ordering rules.
// Returns the time.Time representation of when the Guid was created.
func (g *GuidSS) Timestamp() time.Time {
	encoded := *(*uint64)(unsafe.Pointer(&g.Guid[8])) // Extract timestamp from last 8 bytes (SQL Server format)
	timestamp := int64(bits.RotateLeft64(bits.ReverseBytes64(encoded), 16))
	return time.Unix(0, timestamp).UTC()
}

//==============================================
// Standalone Functions
//==============================================

// New generates a new cryptographically secure Guid.
func New() (g Guid) {
	guidCacheRef := guidCachePool.Get().(*guidCache)

	if guidCacheRef.index == 0 {
		cryptoRand.Read(guidCacheRef.buffer) // Refill buffer if index wraps (Go 1.24+: cryptoRand.Read is guaranteed to succeed)
	}

	copy(g[:], guidCacheRef.buffer[int(guidCacheRef.index)*GuidByteSize:]) // Extract GUID at current index

	guidCacheRef.index++ // Increment index for next call, uint8 wraps from 255 to 0 automatically
	guidCachePool.Put(guidCacheRef)
	return
}

// Used as a benchmark baseline
func _CachePool_GetPut() {
	guidCacheRef := guidCachePool.Get().(*guidCache)
	guidCachePool.Put(guidCacheRef)
}

var _ = _CachePool_GetPut

// NewPG generates a new PostgreSQL sortable Guid as [8-byte time.Now() timestamp][8 random bytes]
func NewPG() GuidPG {
	return newPG(time.Now().UnixNano())
}

func newPG(ts int64) (gpg GuidPG) {
	gpg.Guid = New()
	if !cpu.IsBigEndian {
		ts = int64(bits.ReverseBytes64(uint64(ts)))
	}
	*(*uint64)(unsafe.Pointer(&gpg.Guid[0])) = uint64(ts)
	return
}

// NewSS generates a new SQL Server sortable Guid as [8 random bytes][8 bytes of SQL Server ordered time.Now() timestamp]
func NewSS() GuidSS {
	return newSS(time.Now().UnixNano())
}

func newSS(ts int64) (gss GuidSS) {
	// based on Microsoft SqlGuid.cs
	// https://github.com/microsoft/referencesource/blob/5697c29004a34d80acdaf5742d7e699022c64ecd/System.Data/System/Data/SQLTypes/SQLGuid.cs
	gss.Guid = New()
	// we don't worry about big-endian, because SQL Server does not run on big-endian
	*(*uint64)(unsafe.Pointer(&gss.Guid[8])) = bits.ReverseBytes64(bits.RotateLeft64(uint64(ts), -16))
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

	// Bounds Check Elimination:
	_ = dst[GuidByteSize-1]
	_ = src[GuidBase64UrlByteSize-1]

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
		return &guidCache{buffer: make([]byte, guidCacheByteSize)}
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
