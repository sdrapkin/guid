package guid

import (
	cryptoRand "crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
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
				g := New()
				_ = g.String() // also test that .String() never panics
				guids[idx] = g
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
	err := g.UnmarshalBinary([]byte{1, 2, 3})
	if err == nil {
		t.Error("UnmarshalBinary should fail on short slice")
	}
	// Too long (should succeed, only first 16 bytes used)
	data := make([]byte, 32)
	for i := range data {
		data[i] = byte(i)
	}
	err = g.UnmarshalBinary(data)
	if err != nil {
		t.Errorf("UnmarshalBinary failed on long slice: %v", err)
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
	// Wrong length
	if err := g.UnmarshalText([]byte("short")); err == nil {
		t.Error("UnmarshalText should fail on short input")
	}
	// Invalid chars
	if err = g.UnmarshalText([]byte("!@#$%^&*()_+{}|")); err == nil {
		t.Error("UnmarshalText should fail on invalid chars")
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

func TestNilGuidBehavior(t *testing.T) {
	var nilGuid Guid
	if nilGuid != Nil {
		t.Errorf("Zero Guid should equal Nil constant")
	}
	if Nil.String() != "AAAAAAAAAAAAAAAAAAAAAA" {
		t.Errorf("Nil.String() = %q, want %q", Nil.String(), "AAAAAAAAAAAAAAAAAAAAAA")
	}
	txt, err := Nil.MarshalText()
	if err != nil {
		t.Errorf("Nil.MarshalText() error: %v", err)
	}
	if string(txt) != "AAAAAAAAAAAAAAAAAAAAAA" {
		t.Errorf("Nil.MarshalText() = %q, want %q", string(txt), "AAAAAAAAAAAAAAAAAAAAAA")
	}
	bin, err := Nil.MarshalBinary()
	if err != nil {
		t.Errorf("Nil.MarshalBinary() error: %v", err)
	}
	if len(bin) != GuidByteSize {
		t.Errorf("Nil.MarshalBinary() returned %d bytes, want %d", len(bin), GuidByteSize)
	}
}

func TestGuidEquality(t *testing.T) {
	g1 := New()
	g2 := g1
	if g1 != g2 {
		t.Error("Copied Guid should be equal to original")
	}
	g3 := New()
	if g1 == g3 {
		t.Error("Different Guids should not be equal")
	}
}

func TestGuidMarshalUnmarshalRoundTrip(t *testing.T) {
	g := New()
	bin, err := g.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary failed: %v", err)
	}
	var g2 Guid
	if err := g2.UnmarshalBinary(bin); err != nil {
		t.Fatalf("UnmarshalBinary failed: %v", err)
	}
	if g != g2 {
		t.Errorf("MarshalBinary/UnmarshalBinary round-trip mismatch")
	}

	txt, err := g.MarshalText()
	if err != nil {
		t.Fatalf("MarshalText failed: %v", err)
	}
	var g3 Guid
	if err := g3.UnmarshalText(txt); err != nil {
		t.Fatalf("UnmarshalText failed: %v", err)
	}
	if g != g3 {
		t.Errorf("MarshalText/UnmarshalText round-trip mismatch")
	}
}

func TestGuidStringLengthAndUniqueness(t *testing.T) {
	seen := make(map[string]struct{})
	for range 1000 {
		g := New()
		s := g.String()
		if len(s) != 22 {
			t.Errorf("Guid.String() length = %d, want 22", len(s))
		}
		if _, exists := seen[s]; exists {
			t.Errorf("Duplicate Guid string: %q", s)
		}
		seen[s] = struct{}{}
	}
}

func TestGuidEncodeBase64URLBufferSizes(t *testing.T) {
	g := New()
	// Correct size
	buf := make([]byte, GuidBase64UrlByteSize)
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("EncodeBase64URL panicked with correct buffer size: %v", r)
		}
	}()
	g.EncodeBase64URL(buf)
	// Too large buffer should not panic
	largeBuf := make([]byte, GuidBase64UrlByteSize+5)
	g.EncodeBase64URL(largeBuf[:GuidBase64UrlByteSize])
}

func TestGuidParsePadding(t *testing.T) {
	// Padding should fail
	_, err := Parse("AAAAAAAAAAAAAAAAAAAAAA==")
	if err == nil {
		t.Error("Parse should fail with padding")
	}
}

func TestGuidJSONMarshalling(t *testing.T) {
	type wrapper1 struct {
		ID Guid `json:"id"`
	}
	type wrapper2 struct {
		ID *Guid `json:"id"`
	}

	{
		w1 := wrapper1{ID: New()}
		data1, err := json.Marshal(w1)
		if err != nil {
			t.Fatalf("json.Marshal failed: %v", err)
		}
		w1Clone := wrapper1{ID: New()}
		if err := json.Unmarshal(data1, &w1Clone); err != nil {
			t.Fatalf("json.Unmarshal failed: %v", err)
		}
		if w1.ID != w1Clone.ID {
			t.Errorf("JSON round-trip mismatch: got %v, want %v", w1Clone.ID, w1.ID)
		}
	}
	{
		w1 := wrapper1{ID: Nil}
		data1, err := json.Marshal(w1)
		if err != nil {
			t.Fatalf("json.Marshal failed: %v", err)
		}
		w1Clone := wrapper1{ID: New()}
		if err := json.Unmarshal(data1, &w1Clone); err != nil {
			t.Fatalf("json.Unmarshal failed: %v", err)
		}
		if w1.ID != w1Clone.ID {
			t.Errorf("JSON round-trip mismatch: got %v, want %v", w1Clone.ID, w1.ID)
		}
	}
	{
		w2 := wrapper2{ID: nil}
		data2, err := json.Marshal(w2)
		if err != nil || string(data2) != "{\"id\":null}" {
			t.Fatalf("json.Marshal failed: %v", err)
		}
		g := New()
		w2Clone := wrapper2{ID: &g}
		if err := json.Unmarshal(data2, &w2Clone); err != nil {
			t.Fatalf("json.Unmarshal failed: %v", err)
		}
		if w2.ID != w2Clone.ID {
			t.Errorf("JSON round-trip mismatch: got %v, want %v", w2Clone.ID, w2.ID)
		}
	}
	{
		w2 := wrapper2{ID: &Nil}
		data2, err := json.Marshal(w2)
		if err != nil || string(data2) != "{\"id\":\"AAAAAAAAAAAAAAAAAAAAAA\"}" {
			t.Fatalf("json.Marshal failed: %v", err)
		}
		g := New()
		w2Clone := wrapper2{ID: &g}
		if err := json.Unmarshal(data2, &w2Clone); err != nil {
			t.Fatalf("json.Unmarshal failed: %v", err)
		}
		if *w2.ID != *w2Clone.ID {
			t.Errorf("JSON round-trip mismatch: got %v, want %v", *w2Clone.ID, *w2.ID)
		}
	}
	{
		g := New()
		w2 := wrapper2{ID: &g}
		data2, err := json.Marshal(w2)
		if err != nil {
			t.Fatalf("json.Marshal failed: %v", err)
		}
		w2Clone := wrapper2{ID: nil}
		if err := json.Unmarshal(data2, &w2Clone); err != nil {
			t.Fatalf("json.Unmarshal failed: %v", err)
		}
		if *w2.ID != *w2Clone.ID {
			t.Errorf("JSON round-trip mismatch: got %v, want %v", *w2Clone.ID, *w2.ID)
		}
	}
	{
		var g Guid
		// Not a string
		err := g.UnmarshalJSON([]byte("123"))
		if err == nil {
			t.Error("UnmarshalJSON should fail on non-string JSON")
		}
		// Invalid string
		err = g.UnmarshalJSON([]byte(`"not-a-guid"`))
		if err == nil {
			t.Error("UnmarshalJSON should fail on invalid Guid string")
		}
	}
}

func TestFromBytes(t *testing.T) {
	g1 := New()
	g2, err := FromBytes(g1[:])
	if err != nil {
		t.Fatalf("FromBytes failed: %v", err)
	}
	if g2 != g1 {
		t.Errorf("FromBytes mismatch: got %v, want %v", g2, g1)
	}
	// Too short
	_, err = FromBytes(g1[:15])
	if err == nil {
		t.Error("FromBytes should fail on short slice")
	}
	// Too long
	long := append(g1[:], g1[:4]...)
	g2, err = FromBytes(long)
	if err != nil {
		t.Errorf("FromBytes failed on long slice: %v", err)
	}
	if g2 != g1 {
		t.Errorf("FromBytes mismatch: got %v, want %v", g2, g1)
	}
}

func TestReader_Read(t *testing.T) {
	isZeroChunk := func(chunk []byte) bool {
		for _, v := range chunk {
			if v != 0 {
				return false
			}
		}
		return true
	}

	hasZeroChunk := func(b []byte, clen int) bool {
		for len(b) > 0 {
			limit := min(clen, len(b))
			if isZeroChunk(b[0:limit]) {
				return true
			}
			b = b[limit:]
		}
		return false
	}

	bufLens := []int{}
	for i := range 256 {
		bufLens = append(bufLens, i)
	}
	bufLens = append(bufLens, 0, 1, 8, 256, 511, 511, 511, 511, 512, 513)

	for _, bufLen := range bufLens {
		buf := make([]byte, bufLen)
		n, err := Reader.Read(buf)
		if err != nil {
			t.Fatalf("Reader.Read returned error: %v", err)
		}
		if n != bufLen {
			t.Errorf("Reader.Read returned n=%d, want %d", n, bufLen)
		}

		const chunkLen = 8
		if hasZeroChunk(buf, chunkLen) {
			t.Errorf("Reader.Read buffer contains %d consecutive zero bytes", chunkLen)

		}
	}
}

func TestReadFunction(t *testing.T) {
	const bufLen = 32
	buf := make([]byte, bufLen)
	n, err := Read(buf)
	if err != nil {
		t.Fatalf("Read returned error: %v", err)
	}
	if n != bufLen {
		t.Errorf("Read returned n=%d, want %d", n, bufLen)
	}
}

func TestDecodeBase64URL_LastByteInvalid(t *testing.T) {
	src := make([]byte, GuidBase64UrlByteSize)
	copy(src, "AAAAAAAAAAAAAAAAAAAA") // 20 valid chars
	src[20] = '!'                     // invalid Base64Url
	src[21] = '!'                     // invalid Base64Url

	var g Guid
	ok := DecodeBase64URL(g[:], src)
	if ok {
		t.Error("DecodeBase64URL should fail when final 2 chars are invalid")
	}
}

func TestNewString(t *testing.T) {
	s := NewString()
	if len(s) != GuidBase64UrlByteSize {
		t.Errorf("NewString() returned string of length %d, want %d", len(s), GuidBase64UrlByteSize)
	}
	_, err := Parse(s)
	if err != nil {
		t.Errorf("NewString() returned invalid Guid string: %v", err)
	}
}

func FuzzParse(f *testing.F) {
	// Add some valid and invalid seed cases
	f.Add("AAAAAAAAAAAAAAAAAAAAAA")   // valid (Nil)
	f.Add("_____________________w")   // valid
	f.Add("not-a-guid")               // invalid
	f.Add("")                         // invalid
	f.Add("AAAAAAAAAAAAAAAAAAAAAA==") // invalid (with padding)
	f.Add("1234567890123456789012")   // invalid (wrong chars)
	f.Add("!@#$%^&*()_+{}|")          // invalid (wrong chars, wrong length)
	f.Add("こんaにちは世ち")                 // unicode with len=22

	f.Fuzz(func(t *testing.T, s string) {
		g, err := Parse(s)
		if err != nil {
			// If Parse fails, that's fine for invalid input
			return
		}
		// If Parse succeeds, the string must be 22 chars and must round-trip
		if len(s) != GuidBase64UrlByteSize {
			t.Errorf("Parse succeeded for string of wrong length: %q", s)
		}
		// Round-trip: g.String() should equal s (case for valid input)
		s2 := g.String()
		g2, err2 := Parse(s2)
		if err2 != nil {
			t.Errorf("Parse failed on round-trip string: %q", s2)
		}
		if g != g2 {
			t.Errorf("Round-trip mismatch: got %v, want %v", g2, g)
		}
	})
}

func FuzzParseBytes(f *testing.F) {
	f.Add([]byte("AAAAAAAAAAAAAAAAAAAAAA"))   // valid
	f.Add([]byte("not-a-guid"))               // invalid
	f.Add([]byte(""))                         // invalid
	f.Add([]byte("AAAAAAAAAAAAAAAAAAAAAA==")) // invalid
	f.Add([]byte("1234567890123456789012"))   // invalid
	f.Add([]byte("!@#$%^&*()_+{}|"))          // invalid
	f.Add([]byte("こんaにちは世ち"))                 // unicode

	f.Fuzz(func(t *testing.T, b []byte) {
		g, err := ParseBytes(b)
		if err != nil {
			return
		}
		if len(b) != GuidBase64UrlByteSize {
			t.Errorf("ParseBytes succeeded for slice of wrong length: %q", b)
		}
		s := g.String()
		g2, err2 := Parse(s)
		if err2 != nil {
			t.Errorf("Parse failed on round-trip string: %q", s)
		}
		if g != g2 {
			t.Errorf("Round-trip mismatch: got %v, want %v", g2, g)
		}
	})
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

func BenchmarkReadPerf(b *testing.B) {

	sizes := []int{0, 1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 513, 1024, 2048, 4096}

	// Create a slice of slices
	var data [][]byte
	for _, size := range sizes {
		// Allocate a zero-filled slice of the desired size
		data = append(data, make([]byte, size))
	}

	separator := func() { fmt.Println("=================================") }
	separator()
	for _, buf := range data {
		benchName_guid := fmt.Sprintf("      Guid_Read([%v]byte)", len(buf))
		benchName_rand := fmt.Sprintf("cryptoRand_Read([%v]byte)", len(buf))
		b.Run(
			benchName_guid,
			func(b *testing.B) {
				for b.Loop() {
					Read(buf)
				}
			},
		)

		b.Run(
			benchName_rand,
			func(b *testing.B) {
				for b.Loop() {
					cryptoRand.Read(buf)
				}
			},
		)
		separator()
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
