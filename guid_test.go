package guid

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	// used for benchmarking - commented out to avoid taking dependencies
	//"github.com/sixafter/nanoid"
	//"github.com/google/uuid"
)

var testcases = []struct {
	guidAsHex string
	base64Url string
}{
	// Same test cases: guid hex string, and expected Base64Url string
	{"00000000000000000000000000000000", "AAAAAAAAAAAAAAAAAAAAAA"},
	{"ffffffffffffffffffffffffffffffff", "_____________________w"},
	{"1234567890abcdef1234567890abcdef", "EjRWeJCrze8SNFZ4kKvN7w"},
	{"deadbeefdeadbeefdeadbeefdeadbeef", "3q2-796tvu_erb7v3q2-7w"},
	{"0102030405060708090a0b0c0d0e0f10", "AQIDBAUGBwgJCgsMDQ4PEA"},
	{"abcdefabcdefabcdefabcdefabcdefab", "q83vq83vq83vq83vq83vqw"},
	{"00112233445566778899aabbccddeeff", "ABEiM0RVZneImaq7zN3u_w"},
	{"112233445566778899aabbccddeeff00", "ESIzRFVmd4iZqrvM3e7_AA"},
	{"cafebabecafebabecafebabecafebabe", "yv66vsr-ur7K_rq-yv66vg"},
	{"0f1e2d3c4b5a69788796a5b4c3d2e1f0", "Dx4tPEtaaXiHlqW0w9Lh8A"},
	{"4B2201F58645249A4B73596A9FACBC30", "SyIB9YZFJJpLc1lqn6y8MA"},
	{"9F8FFDFF0E51C078BDA3F774A0674714", "n4_9_w5RwHi9o_d0oGdHFA"},
	{"2C1FBA58F91CDB55372BB51A55B6B3F2", "LB-6WPkc21U3K7UaVbaz8g"},
	{"77E555B922679136D944FADF705DEED5", "d-VVuSJnkTbZRPrfcF3u1Q"},
	{"95AFF5DE37A2CA5AD243DE05A4CDC948", "la_13jeiylrSQ94FpM3JSA"},
	{"F9BFE865B2AE755239C02648CB7A344B", "-b_oZbKudVI5wCZIy3o0Sw"},
	{"E9DF3B1BA91D4133C406C2EF9EB252C6", "6d87G6kdQTPEBsLvnrJSxg"},
	{"5FF246A5B6F629003BA2B4059CE0F77D", "X_JGpbb2KQA7orQFnOD3fQ"},
	{"2AB693B0F0AE0703C92F7E7C6FF3A2BA", "KraTsPCuBwPJL358b_Oiug"},
	{"8F29EAE3EF26F279E43AEAE0CD5EE306", "jynq4-8m8nnkOurgzV7jBg"},
}

func TestGuidLength(t *testing.T) {
	guidLength := len(New())
	if guidLength != GuidByteSize {
		t.Errorf("Generated Guid should have length [%d], got [%d]", GuidByteSize, guidLength)
	}
} //TestGuidLength()

func TestGenerateGuids(t *testing.T) {
	guids := make([]Guid, 1_000_000)
	for i := range guids {
		guids[i] = New()
	}
	guids[len(guids)-1] = Nil
	//guids[len(guids)-2] = Nil // simulate a failed test
	if duplicatesFound(guids) {
		t.Errorf("Duplicate Guids found")
	}
} //TestGenerateGuids()

func TestGenerateGuidsInParallel(t *testing.T) {
	t.Parallel()
	guids := make([]Guid, 4_000_000)

	guidIndex := int64(-1) // Use int64 for atomic operations
	var wg sync.WaitGroup
	numCPUs := runtime.NumCPU()
	t.Logf("Using %d CPUs for parallel GUID generation", numCPUs)
	wg.Add(numCPUs)

	for range numCPUs { // Concurrently generate GUIDs
		go func() {
			defer wg.Done()
			for {
				idx := int(atomic.AddInt64(&guidIndex, 1))
				if idx >= len(guids) {
					break
				}
				guids[idx] = New()
			}
		}()
	}

	wg.Wait()

	guids[len(guids)-1] = Nil
	//guids[len(guids)-2] = Nil // simulate a failed test
	if duplicatesFound(guids) {
		t.Errorf("Duplicate Guids found")
	}
}

func duplicatesFound(guids []Guid) bool {
	lenGuids := len(guids)
	guidMap := make(map[Guid]struct{}, lenGuids)

	for i := range guids {
		guidMap[guids[i]] = struct{}{}
	}
	return len(guidMap) != lenGuids
} //duplicatesFound

func TestGuid_ToBase64Url_RoundTrip(t *testing.T) {
	for _, tc := range testcases {
		// Decode hex string to bytes
		bytes, err := hex.DecodeString(tc.guidAsHex)
		if err != nil {
			t.Fatalf("Failed to decode hex string %q: %v", tc.guidAsHex, err)
		}
		if len(bytes) != GuidByteSize {
			t.Fatalf("Decoded hex string %q to %d bytes, want %d", tc.guidAsHex, len(bytes), GuidByteSize)
		}

		// Convert bytes to Guid
		var g1 Guid
		copy(g1[:], bytes)

		// Convert Guid to Base64Url string
		b64url := g1.String()
		if b64url != tc.base64Url {
			t.Errorf("Guid %q: got Base64Url %q, want %q", tc.guidAsHex, b64url, tc.base64Url)
		}

		// Round-trip: parse Base64Url string back to Guid
		g2, err := Parse(b64url)
		if err != nil {
			t.Errorf("Failed to parse Base64Url string %q: %v", b64url, err)
			continue
		}
		if g1 != g2 {
			t.Errorf("Round-trip mismatch: original Guid %x, after parse %x", g1, g2)
		}
	}
}

func TestGuid_Marshalling_ZeroAndNilInputs(t *testing.T) {
	// UnmarshalBinary: zero-length and nil
	var g Guid
	if err := g.UnmarshalBinary(nil); err == nil {
		t.Error("UnmarshalBinary(nil) should fail")
	}
	if err := g.UnmarshalBinary([]byte{}); err == nil {
		t.Error("UnmarshalBinary(empty) should fail")
	}
	// MarshalBinary: always returns 16 bytes
	bin, err := g.MarshalBinary()
	if err != nil {
		t.Errorf("MarshalBinary failed: %v", err)
	}
	if len(bin) != GuidByteSize {
		t.Errorf("MarshalBinary returned %d bytes, want %d", len(bin), GuidByteSize)
	}

	// UnmarshalText: zero-length and nil
	if err := g.UnmarshalText(nil); err == nil {
		t.Error("UnmarshalText(nil) should fail")
	}
	if err := g.UnmarshalText([]byte{}); err == nil {
		t.Error("UnmarshalText(empty) should fail")
	}
	// MarshalText: always returns 22 bytes
	txt, err := g.MarshalText()
	if err != nil {
		t.Errorf("MarshalText failed: %v", err)
	}
	if len(txt) != 22 {
		t.Errorf("MarshalText returned %d bytes, want 22", len(txt))
	}

	// Parse: zero-length and nil
	_, err = Parse("")
	if err == nil {
		t.Error("Parse(\"\") should fail")
	}
	// nil string not possible in Go, but test with wrong length
	_, err = Parse("short")
	if err == nil {
		t.Error("Parse(\"short\") should fail")
	}

	// ParseBytes: zero-length and nil
	_, err = ParseBytes(nil)
	if err == nil {
		t.Error("ParseBytes(nil) should fail")
	}
	_, err = ParseBytes([]byte{})
	if err == nil {
		t.Error("ParseBytes(empty) should fail")
	}

	// ToBase64URL_Buffer: buffer too small
	small := make([]byte, 10)
	defer func() {
		if r := recover(); r == nil {
			t.Error("ToBase64URL_Buffer with small buffer should panic or fail")
		}
	}()
	g.EncodeBase64URL(small)
}

func TestGuidStringEqualMarshalText(t *testing.T) {
	for range 10000 {
		g := New()
		s := g.String()
		txt, _ := g.MarshalText()

		// String should match MarshalText
		if s != string(txt) {
			t.Errorf("String() != MarshalText(): %q vs %q", s, string(txt))
		}
	}
}

func TestParseAndDecodeBase64URL(t *testing.T) {
	for _, tc := range testcases {
		g1, err := Parse(tc.base64Url)
		if err != nil {
			t.Errorf("Parse(%q) failed: %v", tc.base64Url, err)
			continue
		}

		// Test with DecodeBase64URL
		var g2 Guid
		ok := DecodeBase64URL(g2[:], []byte(tc.base64Url))
		if !ok {
			t.Errorf("DecodeBase64URL(%q) failed", tc.base64Url)
		}
		if string(g2[:]) != string(g1[:]) {
			t.Errorf("DecodeBase64URL result mismatch")
		}
	}

	// Test with invalid input
	var g Guid
	ok := DecodeBase64URL(g[:], []byte(""))
	if ok {
		t.Error("DecodeBase64URL(\"\") should fail")
	}
	ok = DecodeBase64URL(g[:], []byte("short"))
	if ok {
		t.Error("DecodeBase64URL(\"short\") should fail")
	}
	ok = DecodeBase64URL(g[:], []byte("!@#$%^&*()_+{}|"))
	if ok {
		t.Error("DecodeBase64URL(invalid chars) should fail")
	}

	// Test with non-ASCII/Unicode input
	unicodeStr := "こんaにちは世ち"
	unicodeStrBytes := []byte(unicodeStr)
	if len(unicodeStrBytes) != GuidBase64UrlByteSize {
		panic("Unicode string length is not 22 bytes - we want this length to be 22 bytes")
	}
	ok = DecodeBase64URL(g[:], unicodeStrBytes)
	if ok {
		t.Errorf("DecodeBase64URL(%q) should fail", unicodeStr)
	}
	_, err := Parse(unicodeStr)
	if err == nil {
		t.Errorf("Parse(%q) should fail", unicodeStr)
	}
}

//*******************
// Benchmarks // to run: go test -bench=".*" -benchmem -benchtime=4s
//*******************

// BenchmarkNew benchmarks the New function of the guid package.
func Benchmark_guid_New_x10(b *testing.B) {
	for b.Loop() {
		_ = New()
		_ = New()
		_ = New()
		_ = New()
		_ = New()
		_ = New()
		_ = New()
		_ = New()
		_ = New()
		_ = New()
	}
}

func Benchmark_guid_NewString_x10(b *testing.B) {
	for b.Loop() {
		_ = NewString()
		_ = NewString()
		_ = NewString()
		_ = NewString()
		_ = NewString()
		_ = NewString()
		_ = NewString()
		_ = NewString()
		_ = NewString()
		_ = NewString()
	}
}

func Benchmark_guid_String_x10(b *testing.B) {
	guid01 := New()
	guid02 := New()
	guid03 := New()
	guid04 := New()
	guid05 := New()
	guid06 := New()
	guid07 := New()
	guid08 := New()
	guid09 := New()
	guid10 := New()

	for b.Loop() {
		_ = guid01.String()
		_ = guid02.String()
		_ = guid03.String()
		_ = guid04.String()
		_ = guid05.String()
		_ = guid06.String()
		_ = guid07.String()
		_ = guid08.String()
		_ = guid09.String()
		_ = guid10.String()
	}
}

func Benchmark_guid_New_Parallel_x10(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = New()
			_ = New()
			_ = New()
			_ = New()
			_ = New()
			_ = New()
			_ = New()
			_ = New()
			_ = New()
			_ = New()
		}
	})
}

func Benchmark_guid_NewString_Parallel_x10(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = NewString()
			_ = NewString()
			_ = NewString()
			_ = NewString()
			_ = NewString()
			_ = NewString()
			_ = NewString()
			_ = NewString()
			_ = NewString()
			_ = NewString()
		}
	})
}

/* commented out to avoid taking dependencies
func Benchmark_nanoid_New_x10(b *testing.B) {
	for b.Loop() {
		_, _ = nanoid.New()
		_, _ = nanoid.New()
		_, _ = nanoid.New()
		_, _ = nanoid.New()
		_, _ = nanoid.New()
		_, _ = nanoid.New()
		_, _ = nanoid.New()
		_, _ = nanoid.New()
		_, _ = nanoid.New()
		_, _ = nanoid.New()
	}
}

func Benchmark_nanoid_New_Parallel_x10(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = nanoid.New()
			_, _ = nanoid.New()
			_, _ = nanoid.New()
			_, _ = nanoid.New()
			_, _ = nanoid.New()
			_, _ = nanoid.New()
			_, _ = nanoid.New()
			_, _ = nanoid.New()
			_, _ = nanoid.New()
			_, _ = nanoid.New()
		}
	})
}

func Benchmark_uuid_New_x10(b *testing.B) {
	uuid.DisableRandPool()
	for b.Loop() {
		_ = uuid.New()
		_ = uuid.New()
		_ = uuid.New()
		_ = uuid.New()
		_ = uuid.New()
		_ = uuid.New()
		_ = uuid.New()
		_ = uuid.New()
		_ = uuid.New()
		_ = uuid.New()
	}
}

func Benchmark_uuid_New_RandPool_x10(b *testing.B) {
	uuid.EnableRandPool()
	for b.Loop() {
		_ = uuid.New()
		_ = uuid.New()
		_ = uuid.New()
		_ = uuid.New()
		_ = uuid.New()
		_ = uuid.New()
		_ = uuid.New()
		_ = uuid.New()
		_ = uuid.New()
		_ = uuid.New()
	}
}

func Benchmark_uuid_New_Parallel_x10(b *testing.B) {
	uuid.DisableRandPool()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = uuid.New()
			_ = uuid.New()
			_ = uuid.New()
			_ = uuid.New()
			_ = uuid.New()
			_ = uuid.New()
			_ = uuid.New()
			_ = uuid.New()
			_ = uuid.New()
			_ = uuid.New()
		}
	})
}

func Benchmark_uuid_New_RandPool_Parallel_x10(b *testing.B) {
	uuid.EnableRandPool()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = uuid.New()
			_ = uuid.New()
			_ = uuid.New()
			_ = uuid.New()
			_ = uuid.New()
			_ = uuid.New()
			_ = uuid.New()
			_ = uuid.New()
			_ = uuid.New()
			_ = uuid.New()
		}
	})
}
*/

var benchGuids []Guid

func setupBenchGuids() {
	if len(benchGuids) == 0 {
		benchGuids = make([]Guid, len(testcases))
		for i, tc := range testcases {
			bytes, err := hex.DecodeString(tc.guidAsHex)
			if err != nil {
				panic(fmt.Sprintf("Failed to decode hex string %q: %v", tc.guidAsHex, err))
			}
			var g Guid
			copy(g[:], bytes)
			benchGuids[i] = g
		}
	}
}

func Benchmark_guid_ToBase64UrlString(b *testing.B) {
	setupBenchGuids()
	for b.Loop() {
		for _, g := range benchGuids {
			_ = g.String()
		}
	}
}

func Benchmark_base64_RawURLEncoding_EncodeToString(b *testing.B) {
	setupBenchGuids()
	for b.Loop() {
		for _, g := range benchGuids {
			_ = base64.RawURLEncoding.EncodeToString(g[:])
		}
	}
}

func Benchmark_guid_EncodeBase64URL(b *testing.B) {
	setupBenchGuids()
	buffer := make([]byte, GuidBase64UrlByteSize)
	for b.Loop() {
		for _, g := range benchGuids {
			g.EncodeBase64URL(buffer)
		}
	}
}

func Benchmark_base64_RawURLEncoding_Encode(b *testing.B) {
	setupBenchGuids()
	buffer := make([]byte, GuidBase64UrlByteSize)
	for b.Loop() {
		for _, g := range benchGuids {
			base64.RawURLEncoding.Encode(buffer, g[:])
		}
	}
}

func ExampleNew() {
	g := New()      // new random Guid
	fmt.Println(&g) // calls g.String(), which returns the Base64Url encoded string
}

func ExampleGuid_String() {
	// g is a 16-byte Guid represented as a hex string "0123456789abcdef0123456789abcdef"
	var g Guid = [16]byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef, 0x10, 0x32, 0x54, 0x76, 0x98, 0xba, 0xdc, 0xfe}
	fmt.Println(&g) // calls g.String(), which returns the Base64Url encoded string
	// Output: ASNFZ4mrze8QMlR2mLrc_g
}
