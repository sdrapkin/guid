// Package guid provides fast, efficient, cryptographically secure 128-bit GUID generation and manipulation.
// It supports Base64Url encoding/decoding, sequential sortable GUIDs (for PostgreSQL and SQL Server), and is optimized for performance.
// It includes a high-throughput, drop-in replacement for crypto/rand.Reader for generating secure random bytes.
package guid

import (
	cryptoRand "crypto/rand"
	"encoding"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/bits"
	"slices"
	"sync"
	"time"
	"unsafe"
	"uuid"
)

//==============================================
// Constants
//==============================================

const (
	GuidByteSize          = 16                           // Size of a Guid in bytes
	guidsPerCache         = 256 - 1                      // 256-1 Guids per cache to stay under 4096 bytes - do not change this value
	guidCacheByteSize     = GuidByteSize * guidsPerCache // 4096-16=4080 bytes per cache (16*255)
	GuidBase64UrlByteSize = 22                           // Base64Url encoding of a Guid is 22 characters
)

const (
	base64UrlAlphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_" // Base64Url alphabet used for encoding
)

// Ensure that the constants are not changed without thought.
var _ = map[bool]int{false: 0, guidsPerCache == 255: 1}
var _ = map[bool]int{false: 0, guidCacheByteSize == 4080: 1}

//==============================================
// Errors and Variables
//==============================================

var (
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
	_ fmt.Stringer               = Guid{}
	_ encoding.TextMarshaler     = Guid{}
	_ encoding.TextUnmarshaler   = &Guid{}
	_ encoding.BinaryMarshaler   = Guid{}
	_ encoding.BinaryUnmarshaler = &Guid{}
	_ encoding.TextAppender      = Guid{}
	_ io.Reader                  = reader{}
	_ json.Marshaler             = Guid{}
	_ json.Unmarshaler           = &Guid{}
)

//==============================================
// Types
//==============================================

// Guid is a 16-byte (128-bit) cryptographically random value.
type Guid struct {
	UUID uuid.UUID
}

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
	buffer [guidCacheByteSize]byte
	offset int
}

//==============================================
// Guid Extension Methods
//==============================================

// Compare compares the Guid with another Guid (big-endian byte order).
// Returns -1 if g < other, 0 if g == other, and 1 if g > other.
func (g Guid) Compare(other Guid) int {
	hi1 := binary.BigEndian.Uint64(g.UUID[:8])
	hi2 := binary.BigEndian.Uint64(other.UUID[:8])

	if hi1 != hi2 {
		if hi1 < hi2 {
			return -1
		}
		return 1
	}

	lo1 := binary.BigEndian.Uint64(g.UUID[8:])
	lo2 := binary.BigEndian.Uint64(other.UUID[8:])

	if lo1 < lo2 {
		return -1
	}
	if lo1 > lo2 {
		return 1
	}
	return 0
}

// MarshalBinary implements the encoding.BinaryMarshaler interface for Guid.
func (guid Guid) MarshalBinary() (data []byte, err error) {
	return guid.UUID[:], nil // value receiver creates a copy, so it's safe to return the slice directly.
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface for Guid.
func (guid *Guid) UnmarshalBinary(data []byte) error {
	if len(data) < GuidByteSize {
		return ErrInvalidGuidSlice
	}
	copy(guid.UUID[:], data)
	return nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (guid *Guid) UnmarshalText(data []byte) error {
	if ok := DecodeBase64URL(guid.UUID[:], data); !ok {
		return ErrInvalidBase64UrlGuidEncoding
	}
	return nil
}

// MarshalText implements encoding.TextMarshaler.
func (guid Guid) MarshalText() ([]byte, error) {
	buffer := make([]byte, GuidBase64UrlByteSize)
	guid.encodeBase64URL(buffer)
	return buffer, nil
}

// MarshalJSON implements the json.Marshaler interface.
// It marshals the Guid to its Base64Url string representation.
func (g Guid) MarshalJSON() ([]byte, error) {
	gStringWithQuotes := make([]byte, GuidBase64UrlByteSize+2)
	gStringWithQuotes[1+GuidBase64UrlByteSize], gStringWithQuotes[0] = '"', '"'
	g.encodeBase64URL(gStringWithQuotes[1 : 1+GuidBase64UrlByteSize])
	return gStringWithQuotes, nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It unmarshals a JSON string into a Guid.
func (g *Guid) UnmarshalJSON(data []byte) error {
	if len(data) == 4 && data[0] == 'n' && data[1] == 'u' && data[2] == 'l' && data[3] == 'l' {
		*g = Guid{}
		return nil // valid null Guid
	}

	if len(data) != (GuidBase64UrlByteSize+2) ||
		data[0] != '"' ||
		data[GuidBase64UrlByteSize+1] != '"' ||
		!DecodeBase64URL(g.UUID[:], data[1:1+GuidBase64UrlByteSize]) {
		return fmt.Errorf("guid: cannot unmarshal JSON string %q into a Guid", string(data))
	}
	return nil
}

// String returns a Base64Url-encoded string representation of the Guid.
func (guid Guid) String() string {
	buffer := make([]byte, GuidBase64UrlByteSize)
	guid.encodeBase64URL(buffer)
	return unsafe.String(&buffer[0], GuidBase64UrlByteSize)

	// Is it "safe" to use unsafe.String() here? Yes.
	// buffer will be allocated on the heap, and will not be gc'ed until string is alive.
	// This is the same approach that Golang uses in "strings.Clone()" [https://pkg.go.dev/strings#Clone],
	// which calls internal "stringslite.Clone()":
	// https://cs.opensource.google/go/go/+/refs/tags/go1.27.1:src/internal/stringslite/strings.go;l=115
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

// AppendText implements the encoding.TextAppender interface.
func (guid Guid) AppendText(b []byte) ([]byte, error) {
	n := len(b)
	b = slices.Grow(b, GuidBase64UrlByteSize)[:n+GuidBase64UrlByteSize]
	guid.encodeBase64URL(b[n:])
	return b, nil
}

// EncodeBase64URL encodes the Guid into the provided dst as Base64Url.
func (guid Guid) EncodeBase64URL(dst []byte) error {
	if len(dst) < GuidBase64UrlByteSize {
		return ErrBufferTooSmallBase64Url
	}
	guid.encodeBase64URL(dst)
	return nil
}

// Helper closure to pack 4 characters into one 32-bit integer (Little-Endian)
//
//go:inline
func encode4(val uint32) uint32 {
	c0 := uint32(base64UrlAlphabet[val>>18])
	c1 := uint32(base64UrlAlphabet[val>>12&0x3F])
	c2 := uint32(base64UrlAlphabet[val>>6&0x3F])
	c3 := uint32(base64UrlAlphabet[val&0x3F])
	return c0 | (c1 << 8) | (c2 << 16) | (c3 << 24)
}

// private - panics on undersized buffer or nil guid
func (guid *Guid) encodeBase64URL(dst []byte) {
	// Bounds Check Elimination (BCE) for the entire output slice
	_ = dst[21]
	u := guid.UUID[:]
	b15 := u[15]

	// Unroll 5 iterations: Load 4 bytes at once, shift right by 8 to get 3-byte payload

	// Offset 0
	v := binary.BigEndian.Uint32(u[0:4]) >> 8
	*(*uint32)(unsafe.Pointer(&dst[0])) = encode4(v)

	// Offset 3
	v = binary.BigEndian.Uint32(u[3:7]) >> 8
	*(*uint32)(unsafe.Pointer(&dst[4])) = encode4(v)

	// Offset 6
	v = binary.BigEndian.Uint32(u[6:10]) >> 8
	*(*uint32)(unsafe.Pointer(&dst[8])) = encode4(v)

	// Offset 9
	v = binary.BigEndian.Uint32(u[9:13]) >> 8
	*(*uint32)(unsafe.Pointer(&dst[12])) = encode4(v)

	// Offset 12 (bytes 12, 13, 14, 15 - slice upper bound 16 is safe)
	v = binary.BigEndian.Uint32(u[12:16]) >> 8
	*(*uint32)(unsafe.Pointer(&dst[16])) = encode4(v)

	// Final 16th byte (byte 15) -> 2 output characters (dst[20] and dst[21])
	dst[20] = base64UrlAlphabet[b15>>2]
	dst[21] = base64UrlAlphabet[(b15&0x03)<<4]
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

	offset := guidCacheRef.offset

	// Refills buffer if requested bytes exceed remaining capacity, or if pool object is brand new with max-offset
	if offset+n > guidCacheByteSize {
		cryptoRand.Read(guidCacheRef.buffer[:])
		offset = 0
	}

	copy(b, guidCacheRef.buffer[offset:offset+n])
	guidCacheRef.offset = offset + n

	guidCachePool.Put(guidCacheRef)
	return n, nil
} //func (r reader) Read

//==============================================
// GuidPG Extension Methods
//==============================================

// Timestamp extracts the timestamp from the PostgreSQL Guid.
// The timestamp is stored in the first 8 bytes as nanoseconds since Unix epoch.
// Returns the time.Time representation of when the Guid was created.
func (g GuidPG) Timestamp() time.Time {
	timestamp := int64(binary.BigEndian.Uint64(g.Guid.UUID[0:8])) // Extract timestamp from first 8 bytes
	return time.Unix(0, timestamp).UTC()
}

// GuidPG.Compare compares the PostgreSQL Guid with another PostgreSQL Guid using big-endian byte order.
// Returns -1 if g < other, 0 if g == other, and 1 if g > other.
func (g GuidPG) Compare(other GuidPG) int {
	return g.Guid.Compare(other.Guid)
}

//==============================================
// GuidSS Extension Methods
//==============================================

// Timestamp extracts the timestamp from the SQL Server Guid.
// The timestamp is stored in the last 8 bytes using SQL Server's Guid ordering rules.
// Returns the time.Time representation of when the Guid was created.
func (g GuidSS) Timestamp() time.Time {
	encoded := binary.BigEndian.Uint64(g.Guid.UUID[8:]) // Extract timestamp from last 8 bytes (SQL Server format)
	timestamp := int64(bits.RotateLeft64(encoded, 16))
	return time.Unix(0, timestamp).UTC()
}

// GuidSS.Compare compares the SQL Server Guid with another SQL Server Guid using SQL Server's byte ordering rules.
// Returns -1 if g < other, 0 if g == other, and 1 if g > other.
func (g GuidSS) Compare(other GuidSS) int {
	// SQL Server compares bytes in order: 10, 11, 12, 13, 14, 15, 8, 9, 6, 7, 4, 5, 0, 1, 2, 3
	// https://source.dot.net/#System.Data.Common/System/Data/SQLTypes/SQLGuid.cs,116

	// High 64 bits target order: bytes [10, 11, 12, 13, 14, 15, 8, 9]
	gHi := bits.RotateLeft64(binary.BigEndian.Uint64(g.Guid.UUID[8:]), 16)
	oHi := bits.RotateLeft64(binary.BigEndian.Uint64(other.Guid.UUID[8:]), 16)

	if gHi < oHi {
		return -1
	}
	if gHi > oHi {
		return 1
	}

	// High bits match; compute low 64 bits lazily.
	// Low 64 bits target order: bytes [6, 7, 4, 5, 0, 1, 2, 3]
	ga := binary.BigEndian.Uint64(g.Guid.UUID[:8])
	oa := binary.BigEndian.Uint64(other.Guid.UUID[:8])

	/* Rearrange ga's 8 bytes from standard layout [0,1,2,3,4,5,6,7] to SQL Server's target layout [6,7,4,5,0,1,2,3]:

	ga bit layout (MSB to LSB):
	Bits 63..32 = Bytes [0,1,2,3]
	Bits 31..16 = Bytes [4,5]
	Bits 15..0  = Bytes [6,7]

	1. (ga << 48): Shifts Bytes [6,7] left by 6 bytes (48 bits).
	2. ((ga & 0x00000000FFFF0000) << 16): Isolates Bytes [4,5] and shifts them left by 2 bytes (16 bits).
	3. (ga >> 32): Shifts Bytes [0,1,2,3] right by 4 bytes (32 bits).

	Combining with bitwise OR (|) yields a uint64 byte-ordered: [6,7,4,5,0,1,2,3] */
	gLo := (ga << 48) | ((ga & 0x00000000FFFF0000) << 16) | (ga >> 32)
	oLo := (oa << 48) | ((oa & 0x00000000FFFF0000) << 16) | (oa >> 32)

	if gLo < oLo {
		return -1
	}
	if gLo > oLo {
		return 1
	}

	return 0
	/*
		for _, i := range [16]int{
			10, 11, 12, 13, 14, 15, 8, 9, 6, 7, 4, 5, 0, 1, 2, 3,
		} {
			if g.Guid.UUID[i] < other.Guid.UUID[i] {
				return -1
			}
			if g.Guid.UUID[i] > other.Guid.UUID[i] {
				return 1
			}
		}
		return 0
	*/
}

// LoadFromSQLServerBytes loads a GuidSS from a 16-byte slice in SQL Server order.
// It swaps the bytes into the correct order for GuidSS.
func (g *GuidSS) LoadFromSQLServerBytes(src []byte) error {
	if len(src) != GuidByteSize {
		return ErrInvalidGuidSlice
	}

	*g = GuidSS{
		UUID: [16]byte{
			src[3], src[2], src[1], src[0], // Swap bytes 0-3
			src[5], src[4], // Swap bytes 4-5
			src[7], src[6], // Swap bytes 6-7
			src[8], src[9], src[10], src[11], src[12], src[13], src[14], src[15], // Bytes 8-15 remain in order
		},
	}

	return nil
}

//==============================================
// Standalone Functions
//==============================================

// Nil returns the nil Guid, with all 128 bits set to zero.
func Nil() Guid { return Guid{} }

// Max returns the maximum Guid, with all 128 bits set to one.
func Max() Guid {
	return Guid{UUID: [16]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}}
}

// New generates a new cryptographically secure Guid.
func New() (g Guid) {
	guidCacheRef := guidCachePool.Get().(*guidCache)

	offset := guidCacheRef.offset

	if offset > guidCacheByteSize-GuidByteSize {
		cryptoRand.Read(guidCacheRef.buffer[:])
		offset = 0
	}

	g.UUID = *(*[16]byte)(unsafe.Pointer(&guidCacheRef.buffer[offset]))
	guidCacheRef.offset = offset + GuidByteSize
	guidCachePool.Put(guidCacheRef)
	return g
}

// NewPG generates a new PostgreSQL sortable Guid as [8-byte time.Now() timestamp][8 random bytes]
func NewPG() GuidPG {
	return newPG(time.Now().UnixNano())
}

func newPG(ts int64) (gpg GuidPG) {
	gpg.Guid = New()
	binary.BigEndian.PutUint64(gpg.Guid.UUID[0:8], uint64(ts))
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
	encoded := bits.RotateLeft64(uint64(ts), -16)
	binary.BigEndian.PutUint64(gss.Guid.UUID[8:], encoded)
	return
}

// NewString generates a new cryptographically secure Guid, and returns it as a Base64Url string.
// NewString is equivalent to "g := guid.New(); return g.String();".
func NewString() string {
	return New().String()
}

// Parse parses a Base64Url-encoded string into the Guid.
// Returns an error if the string is not a valid Guid encoding.
func Parse(s string) (g Guid, err error) {
	if len(s) != GuidBase64UrlByteSize {
		return Guid{}, ErrInvalidBase64UrlGuidEncoding
	}

	// Zero-copy conversion of a string to a byte slice
	sBytes := unsafe.Slice(unsafe.StringData(s), GuidBase64UrlByteSize)

	if ok := DecodeBase64URL(g.UUID[:], sBytes); !ok {
		return Guid{}, ErrInvalidBase64UrlGuidEncoding
	}
	return g, nil
}

// MustParse returns the Guid represented by Base64Url-encoded string.
// It calls Parse(s) and panics if it returns an error. Use this only when you are sure the string is a valid Guid encoding.
func MustParse(s string) Guid {
	g, err := Parse(s)
	if err != nil {
		panic(err)
	}
	return g
}

// ParseBytes parses a Base64Url-encoded string represented as a byte slice into the Guid.
// Returns an error if the string byte slice is not a valid Guid encoding.
// ParseBytes is like Parse, except it parses a string byte slice instead of a string.
func ParseBytes(src []byte) (g Guid, err error) {
	if len(src) != GuidBase64UrlByteSize {
		return Guid{}, ErrInvalidBase64UrlGuidEncoding
	}

	if ok := DecodeBase64URL(g.UUID[:], src); !ok {
		return Guid{}, ErrInvalidBase64UrlGuidEncoding
	}
	return g, nil
}

// FromBytes returns a Guid from a 16-byte slice.
// If src is exactly 16 bytes, its contents are used directly.
// If src is longer than 16 bytes, only the first 16 bytes are used (excess is silently ignored).
// This differs from Parse/ParseBytes, which require an exact-length match.
// Returns ErrInvalidGuidSlice if src is shorter than 16 bytes.
func FromBytes(src []byte) (Guid, error) {
	if len(src) < GuidByteSize {
		return Guid{}, ErrInvalidGuidSlice
	}
	var g Guid
	copy(g.UUID[:], src)
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
		return &guidCache{offset: guidCacheByteSize} // Start with offset at the end to trigger a refill on first use
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
