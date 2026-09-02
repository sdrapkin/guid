package guid

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"
	"sync"
	"testing"
	"time"
)

//*******************
// Helpful commands:
//
// go test -v -coverprofile="coverage.txt"
// go tool cover -func=coverage.txt
// go tool cover -html="coverage.txt"
// gocyclo -over 15 .
// go test -bench=".*" -benchmem -benchtime=4s
//*******************

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
	guids[len(guids)-1] = Nil()
	//guids[len(guids)-2] = Nil() // simulate a failed test
	if duplicatesFound(guids) {
		t.Errorf("Duplicate Guids found")
	}
} //TestGenerateGuids()

func TestGenerateGuidsInParallel(t *testing.T) {
	t.Parallel()
	const total = 4_000_000
	guids := make([]Guid, total)
	numCPUs := runtime.NumCPU()
	chunkSize := total / numCPUs

	var wg sync.WaitGroup
	t.Logf("Using %d CPUs for parallel GUID generation", numCPUs)
	wg.Add(numCPUs)

	for i := range numCPUs { // Concurrently generate GUIDs
		go func(workerID int) {
			defer wg.Done()
			start := workerID * chunkSize
			end := start + chunkSize
			if workerID == numCPUs-1 {
				end = total // Ensure the last worker covers any remaining GUIDs
			}

			for j := start; j < end; j++ {
				g := New()
				_ = g.String() // also test that .String() never panics
				guids[j] = g
			}
		}(i)
	} // for
	wg.Wait()

	guids[len(guids)-1] = Nil()
	//guids[len(guids)-2] = Nil() // simulate a failed test
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

func TestGuid_AppendText(t *testing.T) {
	g := New()

	t.Run("Append to nil slice", func(t *testing.T) {
		b, err := g.AppendText(nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if string(b) != g.String() {
			t.Errorf("got %q, want %q", string(b), g.String())
		}
	})

	t.Run("Append to existing prefix", func(t *testing.T) {
		prefix := []byte("prefix-")
		b, err := g.AppendText(prefix)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := "prefix-" + g.String()
		if string(b) != expected {
			t.Errorf("got %q, want %q", string(b), expected)
		}
	})

	t.Run("Append with pre-allocated capacity", func(t *testing.T) {
		buf := make([]byte, 0, 64)
		b, err := g.AppendText(buf)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if string(b) != g.String() {
			t.Errorf("got %q, want %q", string(b), g.String())
		}
	})
}

func TestValueReceivers(t *testing.T) {
	// Verify unaddressable constructor returns allow direct method chaining
	str := New().String()
	if len(str) != GuidBase64UrlByteSize {
		t.Errorf("New().String() length = %d, want %d", len(str), GuidBase64UrlByteSize)
	}

	pgTS := NewPG().Timestamp()
	if pgTS.IsZero() {
		t.Error("NewPG().Timestamp() returned zero time")
	}

	ssTS := NewSS().Timestamp()
	if ssTS.IsZero() {
		t.Error("NewSS().Timestamp() returned zero time")
	}

	// Verify Stringer formatting directly on value type
	g := New()
	formatted := fmt.Sprintf("%s", g)
	if formatted != g.String() {
		t.Errorf("fmt.Sprintf(%%s, g) = %q, want %q", formatted, g.String())
	}
}

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

func TestGuid_Marshaling_ZeroAndNilInputs(t *testing.T) {
	var g Guid

	//--- UnmarshalBinary ---
	t.Run("UnmarshalBinary", func(t *testing.T) {
		errorCases := []struct {
			name  string
			input []byte
		}{
			{"nil input", nil},
			{"empty slice", []byte{}},
			{"short slice", []byte{1, 2, 3}},
		}

		for _, tc := range errorCases {
			t.Run(tc.name, func(t *testing.T) {
				if err := g.UnmarshalBinary(tc.input); err == nil {
					t.Error("expected an error but got nil")
				}
			})
		}

		t.Run("succeeds on long slice", func(t *testing.T) {
			longSlice := make([]byte, 32)
			if err := g.UnmarshalBinary(longSlice); err != nil {
				t.Errorf("expected no error but got: %v", err)
			}
		})
	})

	//--- MarshalBinary ---
	t.Run("MarshalBinary", func(t *testing.T) {
		g2 := New()
		marshalledSlice, err := g2.MarshalBinary()
		if err != nil {
			t.Fatalf("expected no error but got: %v", err)
		}
		if len(marshalledSlice) != GuidByteSize {
			t.Errorf("got %d bytes, want %d", len(marshalledSlice), GuidByteSize)
		}
		if !bytes.Equal(g2[:], marshalledSlice) {
			t.Errorf("marshalled slice is not equal to the original guid")
		}
		g2[0]++
		if bytes.Equal(g2[:], marshalledSlice) {
			t.Errorf("changes to original guid propagate to the marshalled slice")
		}
	})

	//--- UnmarshalText ---
	t.Run("UnmarshalText", func(t *testing.T) {
		errorCases := []struct {
			name  string
			input []byte
		}{
			{"nil input", nil},
			{"empty slice", []byte{}},
			{"short input", []byte("short")},
			{"invalid chars", []byte("!@#$%^&*()_+{}|")},
		}

		for _, tc := range errorCases {
			t.Run(tc.name, func(t *testing.T) {
				if err := g.UnmarshalText(tc.input); err == nil {
					t.Error("expected an error but got nil")
				}
			})
		}
	})

	//--- MarshalText ---
	t.Run("MarshalText", func(t *testing.T) {
		txt, err := g.MarshalText()
		if err != nil {
			t.Fatalf("expected no error but got: %v", err)
		}
		if len(txt) != 22 {
			t.Errorf("got %d bytes, want 22", len(txt))
		}
	})
}

func TestGuid_Parse_and_ParseBytes(t *testing.T) {
	//--- Parse and ParseBytes ---
	t.Run("Parse", func(t *testing.T) {
		if _, err := Parse(""); err == nil {
			t.Error("Parse(\"\") should fail")
		}
		if _, err := Parse("short"); err == nil {
			t.Error("Parse(\"short\") should fail")
		}
	})

	t.Run("ParseBytes", func(t *testing.T) {
		if _, err := ParseBytes(nil); err == nil {
			t.Error("ParseBytes(nil) should fail")
		}
		if _, err := ParseBytes([]byte{}); err == nil {
			t.Error("ParseBytes(empty) should fail")
		}
	})
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

	// Test with undersized output slice
	var undersizedDst [15]byte
	if DecodeBase64URL(undersizedDst[:], []byte("AAAAAAAAAAAAAAAAAAAAAA")) {
		t.Error("DecodeBase64URL should fail when dst slice is under 16 bytes")
	}

	// Test with invalid input
	var g Guid
	if DecodeBase64URL(g[:], []byte("")) {
		t.Error("DecodeBase64URL(\"\") should fail")
	}
	if DecodeBase64URL(g[:], []byte("short")) {
		t.Error("DecodeBase64URL(\"short\") should fail")
	}
	if DecodeBase64URL(g[:], []byte("!@#$%^&*()_+{}|")) {
		t.Error("DecodeBase64URL(invalid chars) should fail")
	}

	// Test with non-ASCII/Unicode input
	unicodeStr := "こんaにちは世ち"
	unicodeStrBytes := []byte(unicodeStr)
	if len(unicodeStrBytes) != GuidBase64UrlByteSize {
		t.Fatalf("Unicode test setup failure: expected 22 bytes, got %d", len(unicodeStrBytes))
	}
	if DecodeBase64URL(g[:], unicodeStrBytes) {
		t.Errorf("DecodeBase64URL(%q) should fail", unicodeStr)
	}
	if _, err := Parse(unicodeStr); err == nil {
		t.Errorf("Parse(%q) should fail", unicodeStr)
	}
}

func TestMax(t *testing.T) {
	gmax := Guid{}
	for i := range len(gmax) {
		gmax[i] = 0xFF
	}
	if gmax != Max() {
		t.Error("guid.Max is wrong!")
	}
}

func TestNilGuidBehavior(t *testing.T) {
	var nilGuid Guid
	if nilGuid != Nil() {
		t.Errorf("Zero Guid should equal Nil(), got %v", nilGuid)
	}
	if Nil().String() != "AAAAAAAAAAAAAAAAAAAAAA" {
		t.Errorf("Nil().String() = %q, want %q", Nil().String(), "AAAAAAAAAAAAAAAAAAAAAA")
	}
	txt, err := Nil().MarshalText()
	if err != nil {
		t.Errorf("Nil().MarshalText() error: %v", err)
	}
	if string(txt) != "AAAAAAAAAAAAAAAAAAAAAA" {
		t.Errorf("Nil().MarshalText() = %q, want %q", string(txt), "AAAAAAAAAAAAAAAAAAAAAA")
	}
	bin, err := Nil().MarshalBinary()
	if err != nil {
		t.Errorf("Nil().MarshalBinary() error: %v", err)
	}
	if len(bin) != GuidByteSize {
		t.Errorf("Nil().MarshalBinary() returned %d bytes, want %d", len(bin), GuidByteSize)
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

	if err := g.EncodeBase64URL(nil); err != ErrBufferTooSmallBase64Url {
		t.Error("EncodeBase64URL did not return an error on undersized buffer")
	}

	for bufLen := range GuidBase64UrlByteSize {
		buf := make([]byte, bufLen)
		if err := g.EncodeBase64URL(buf); err != ErrBufferTooSmallBase64Url {
			t.Error("EncodeBase64URL did not return an error on undersized buffer")
		}
	}

	for bufLen := GuidBase64UrlByteSize; bufLen < GuidBase64UrlByteSize*4; bufLen++ {
		buf := make([]byte, bufLen)
		if err := g.EncodeBase64URL(buf); err != nil {
			t.Error("EncodeBase64URL returned an error on properly sized buffer")
		}
	}
}

func TestGuidParsePadding(t *testing.T) {
	// Padding should fail
	_, err := Parse("AAAAAAAAAAAAAAAAAAAAAA==")
	if err == nil {
		t.Error("Parse should fail with padding")
	}
}

func TestGuidJSONMarshaling(t *testing.T) {
	// Define the structs with the correct JSON tags inside the test.
	type guidContainer struct {
		ID Guid `json:"id"`
	}
	type nullableGuidContainer struct {
		ID *Guid `json:"id"`
	}

	newGuidPtr := func() *Guid {
		g := New()
		return &g
	}

	nilGuid := Nil()
	nilGuidPtr := &nilGuid // Helper for a pointer to the Nil Guid.

	// testCases defines the scenarios for our table-driven test.
	testCases := []struct {
		name         string // The name of the test case.
		input        any    // The data to be marshalled.
		expectedJSON string // Optional: The expected JSON output. If empty, not checked.
		// A function to get a "clone" instance for unmarshaling.
		// This ensures we're unmarshaling into a fresh variable.
		getClone func() any
		// A function to compare the original input with the unmarshalled clone.
		isEqual func(original, clone any) bool
	}{
		{
			name:     "Wrapper with non-nil Guid",
			input:    guidContainer{ID: New()},
			getClone: func() any { return &guidContainer{} },
			isEqual: func(original, clone any) bool {
				return original.(guidContainer).ID == clone.(*guidContainer).ID
			},
		},
		{
			name:     "Wrapper with nil Guid",
			input:    guidContainer{ID: Nil()},
			getClone: func() any { return &guidContainer{} },
			isEqual: func(original, clone any) bool {
				return original.(guidContainer).ID == clone.(*guidContainer).ID
			},
		},
		{
			name:         "Wrapper with nil pointer to Guid",
			input:        nullableGuidContainer{ID: nil},
			expectedJSON: `{"id":null}`,
			getClone:     func() any { return &nullableGuidContainer{ID: newGuidPtr()} },
			isEqual: func(original, clone any) bool {
				// Both original and clone ID pointers should be nil.
				return original.(nullableGuidContainer).ID == clone.(*nullableGuidContainer).ID
			},
		},
		{
			name:         "Wrapper with pointer to nil Guid",
			input:        nullableGuidContainer{ID: nilGuidPtr},
			expectedJSON: `{"id":"AAAAAAAAAAAAAAAAAAAAAA"}`,
			getClone:     func() any { return &nullableGuidContainer{ID: newGuidPtr()} },
			isEqual: func(original, clone any) bool {
				return *original.(nullableGuidContainer).ID == *clone.(*nullableGuidContainer).ID
			},
		},
		{
			name:     "Wrapper with pointer to non-nil Guid",
			input:    nullableGuidContainer{ID: newGuidPtr()},
			getClone: func() any { return &nullableGuidContainer{} },
			isEqual: func(original, clone any) bool {
				return *original.(nullableGuidContainer).ID == *clone.(*nullableGuidContainer).ID
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Marshal the input data.
			data, err := json.Marshal(tc.input)
			if err != nil {
				t.Fatalf("json.Marshal failed: %v", err)
			}

			// Optionally, verify the marshalled JSON output.
			if tc.expectedJSON != "" && string(data) != tc.expectedJSON {
				t.Fatalf("json.Marshal produced unexpected output: got %s, want %s", string(data), tc.expectedJSON)
			}

			// Unmarshal the data into a new "clone" object.
			clone := tc.getClone()
			if err := json.Unmarshal(data, clone); err != nil {
				t.Fatalf("json.Unmarshal failed: %v", err)
			}

			// Compare the original and the cloned object.
			if !tc.isEqual(tc.input, clone) {
				t.Errorf("JSON round-trip mismatch: got %+v, want %+v", clone, tc.input)
			}
		})
	}

	// --- Direct UnmarshalJSON Error Tests ---
	t.Run("UnmarshalJSON error cases", func(t *testing.T) {
		var g Guid

		// Test with non-string JSON.
		if err := g.UnmarshalJSON([]byte("123")); err == nil {
			t.Error("UnmarshalJSON should fail on non-string JSON")
		}

		// Test with an invalid Guid string.
		if err := g.UnmarshalJSON([]byte(`"not-a-guid"`)); err == nil {
			t.Error("UnmarshalJSON should fail on invalid Guid string")
		}

		// Test with empty slice
		if err := g.UnmarshalJSON([]byte{}); err == nil {
			t.Error("UnmarshalJSON should fail on empty slice")
		}
		// Test with nil slice
		if err := g.UnmarshalJSON(nil); err == nil {
			t.Error("UnmarshalJSON should fail on nil slice")
		}
	})
}

var testUUID = MustParse("FBzYOSdp3VEBK7jzXPZleA")

func TestJSON(t *testing.T) {
	type S struct {
		ID1 Guid
		ID2 Guid
	}
	s1 := S{ID1: testUUID}
	data, err := json.Marshal(&s1)
	if err != nil {
		t.Fatal(err)
	}
	var s2 S
	if err := json.Unmarshal(data, &s2); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(&s1, &s2) {
		t.Errorf("got %#v, want %#v", s2, s1)
	}
}

func TestUnmarshalJSON(t *testing.T) {
	type S struct {
		ID1 Guid
		ID2 Guid `json:"ID2,omitempty"`
	}

	testCases := map[string]struct {
		data           []byte
		expectedError  bool
		expectedResult Guid
	}{
		"success": {
			data:           []byte(`{"ID1": "FBzYOSdp3VEBK7jzXPZleA"}`),
			expectedError:  false,
			expectedResult: testUUID,
		},
		"zero": {
			data:           []byte(`{"ID1": "AAAAAAAAAAAAAAAAAAAAAA"}`),
			expectedError:  false,
			expectedResult: Nil(),
		},
		"null": {
			data:           []byte(`{"ID1": null}`),
			expectedError:  false,
			expectedResult: Nil(),
		},
		"empty": {
			data:           []byte(`{"ID1": ""}`),
			expectedError:  true,
			expectedResult: Nil(),
		},
		"omitempty": {
			data:           []byte(`{"ID2": ""}`),
			expectedError:  true,
			expectedResult: Nil(),
		},
		"reject-unquoted": {
			data:           []byte(`{"ID1": 1234567890123456789012}`),
			expectedError:  true,
			expectedResult: Nil(),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var s S
			err := json.Unmarshal(tc.data, &s)
			if err != nil && !tc.expectedError {
				t.Errorf("unexpected error: got %v, want %v", err, tc.expectedError)
			}
			if !reflect.DeepEqual(s.ID1, tc.expectedResult) {
				t.Errorf("got %#v, want %#v", s.ID1, tc.expectedResult)
			}
		})
	}

	t.Run("UnmarshalJSON direct-call edge cases", func(t *testing.T) {
		cases := []struct {
			name string
			in   []byte
		}{
			{"unquoted 22-char token", []byte("1234567890123456789012")},
			{"missing closing quote", []byte(`"AAAAAAAAAAAAAAAAAAAAAA`)},
			{"missing opening quote", []byte(`AAAAAAAAAAAAAAAAAAAAAA"`)},
			{"single-quoted instead of double", []byte(`'AAAAAAAAAAAAAAAAAAAAAA'`)},
			{"extra leading byte before quote", []byte(`X"AAAAAAAAAAAAAAAAAAAAAA"`)},
			{"extra trailing byte after quote", []byte(`"AAAAAAAAAAAAAAAAAAAAAA"X`)},
			{"correct length, wrong delimiters", []byte(`[AAAAAAAAAAAAAAAAAAAAAA]`)},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				var g Guid
				if err := g.UnmarshalJSON(tc.in); err == nil {
					t.Errorf("UnmarshalJSON(%q) should have failed, got nil error", tc.in)
				}
			})
		}

		t.Run("UnmarshalJSON direct-call valid quoted value", func(t *testing.T) {
			var g Guid
			if err := g.UnmarshalJSON([]byte(`"AAAAAAAAAAAAAAAAAAAAAA"`)); err != nil {
				t.Errorf("UnmarshalJSON should succeed on properly quoted value: %v", err)
			}
			if g != Nil() {
				t.Errorf("got %v, want Nil", g)
			}
		})
	})
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

	b := []byte{0x7d, 0x44, 0x48, 0x40, 0x9d, 0xc0, 0x11, 0xd1, 0xb2, 0x45, 0x5f, 0xfd, 0xce, 0x74, 0xfa, 0xd2}
	g2, err = FromBytes(b)
	if err != nil {
		t.Fatalf("%s", err)
	}
	for i := range len(g2) {
		if b[i] != g2[i] {
			t.Fatalf("FromBytes() got %v expected %v\n", g2[:], b)
		}
	}
}

func TestReader_Read(t *testing.T) {
	consecutiveByteCount := func(chunk []byte) int {
		if len(chunk) < 2 {
			return len(chunk)
		}
		maxRun := 1
		currentRun := 1
		for i := 1; i < len(chunk); i++ {
			if chunk[i] == chunk[i-1] {
				// If the byte is the same as the last one, the run continues.
				currentRun++
			} else {
				// The run was broken, so reset the current count.
				currentRun = 1
			}

			// Update the max run if the current one is larger.
			if currentRun > maxRun {
				maxRun = currentRun
			}
		}
		return maxRun
	}

	bufLens := []int{}
	for range 1000 {
		for i := range 512 {
			bufLens = append(bufLens, i)
		}
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

		const chunkLen = 6
		if consecutiveByteCount(buf) >= chunkLen {
			t.Errorf("Reader.Read buffer contains %d consecutive bytes\n %v", chunkLen, buf)
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

// TestReadConcurrent tests that concurrent calls to Read(buf) will not create identical values
func TestReadConcurrent(t *testing.T) {
	groutines := 32

	doneChan := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(groutines)
	entropyMaps := make([]map[string]struct{}, groutines)
	for i := range groutines {
		entropyMaps[i] = make(map[string]struct{})
		go func(i int) {
			buf := make([]byte, 16)
			for {
				select {
				case <-doneChan:
					wg.Done()
					return
				default:
				} //select

				Read(buf)
				bufStr := string(buf)
				_, exists := entropyMaps[i][bufStr]
				if exists {
					t.Error("got the same entropy twice out of the reader")
				}
				entropyMaps[i][bufStr] = struct{}{}
			} //for
		}(i)
	}

	// Let the threads spin for a bit, then shut them down.
	time.Sleep(time.Millisecond * 750)
	close(doneChan)
	wg.Wait()

	totalGuids := 0
	for _, entropy := range entropyMaps {
		totalGuids += len(entropy)
	}

	// Compare the entropy collected and verify that no set was output twice.
	allEntropy := make(map[string]struct{}, totalGuids)
	for _, entropy := range entropyMaps {
		for str := range entropy {
			_, exists := allEntropy[str]
			if exists {
				t.Error("got the same entropy twice out of guid.Reader")
			}
			allEntropy[str] = struct{}{}
		}
	}
}

// TestReadLiteConcurrent tests that concurrent calls to ReaderLite.Read(buf) will not create identical values
func TestReadLiteConcurrent(t *testing.T) {
	_readerLite.Read(nil) // should not panic

	groutines := 32
	doneChan := make(chan struct{})

	var wg sync.WaitGroup
	wg.Add(groutines)

	entropyMaps := make([]map[string]struct{}, groutines)
	for i := range groutines {
		entropyMaps[i] = make(map[string]struct{})
		go func(i int) {
			buf := make([]byte, 16)
			for {
				select {
				case <-doneChan:
					wg.Done()
					return
				default:
				} //select

				_readerLite.Read(buf)
				bufStr := string(buf)
				_, exists := entropyMaps[i][bufStr]
				if exists {
					t.Error("got the same entropy twice out of guid.ReaderLite")
				}
				entropyMaps[i][bufStr] = struct{}{}
			} //for
		}(i)
	}

	// Let the threads spin for a bit, then shut them down.
	time.Sleep(time.Millisecond * 750)
	close(doneChan)
	wg.Wait()

	totalGuids := 0
	for _, entropy := range entropyMaps {
		totalGuids += len(entropy)
	}
	// Compare the entropy collected and verify that no set was output twice.
	allEntropy := make(map[string]struct{}, totalGuids)
	for _, entropy := range entropyMaps {
		for str := range entropy {
			_, exists := allEntropy[str]
			if exists {
				t.Error("got the same entropy twice out of the reader")
			}
			allEntropy[str] = struct{}{}
		}
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

func TestGuid_Compare(t *testing.T) {
	base := Guid{0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	tests := []struct {
		name     string
		g1, g2   Guid
		expected int
	}{
		{
			name:     "Equal GUIDs (returns 0)",
			g1:       base,
			g2:       base,
			expected: 0,
		},
		{
			name:     "High 64-bit smaller (hi1 < hi2 -> returns -1)",
			g1:       Guid{0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			g2:       base,
			expected: -1,
		},
		{
			name:     "High 64-bit larger (hi1 > hi2 -> returns 1)",
			g1:       Guid{0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			g2:       base,
			expected: 1,
		},
		{
			name:     "High 64-bit equal, Low 64-bit smaller (lo1 < lo2 -> returns -1)",
			g1:       Guid{0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			g2:       base,
			expected: -1,
		},
		{
			name:     "High 64-bit equal, Low 64-bit larger (lo1 > lo2 -> returns 1)",
			g1:       Guid{0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			g2:       base,
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.g1.Compare(tt.g2); got != tt.expected {
				t.Errorf("Compare() = %d, want %d", got, tt.expected)
			}
		})
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

func TestSortableGuids(t *testing.T) {
	t.Run("GuidPG", func(t *testing.T) {
		// Test uniqueness
		seen := make(map[Guid]bool)
		for range 1000 {
			g := NewPG()
			if seen[g.Guid] {
				t.Fatalf("Duplicate PostgreSQL-compatible Guid found: %s", g)
			}
			seen[g.Guid] = true
		}

		// Test timestamp encoding against a known value
		gFixed := newPG(0x1122334455667788)
		if hex.EncodeToString(gFixed.Guid[:8]) != "1122334455667788" {
			t.Errorf("Invalid timestamp encoding in newPG")
		}

		ts1 := time.Now().UnixNano()
		ts2 := ts1 + 1000
		g1 := newPG(ts1)
		g2 := newPG(ts2)

		if g1.Timestamp().UnixNano() != ts1 {
			t.Errorf("GuidPG timestamp mismatch: got %d, want %d", g1.Timestamp().UnixNano(), ts1)
		}

		if bytes.Compare(g1.Guid[:], g2.Guid[:]) >= 0 {
			t.Errorf("GuidPGs are not sortable. g1 should be less than g2.\ng1: %x\ng2: %x", g1.Guid, g2.Guid)
		}
	})

	t.Run("GuidSS", func(t *testing.T) {
		// Test uniqueness
		seen := make(map[Guid]bool)
		for range 1000 {
			g := NewSS()
			if seen[g.Guid] {
				t.Fatalf("Duplicate SQL Server-compatible Guid found: %s", g)
			}
			seen[g.Guid] = true
		}

		// Test timestamp encoding against a known value
		gFixed := newSS(0x1122334455667788)
		if hex := hex.EncodeToString(gFixed.Guid[8:]); hex != "7788112233445566" {
			t.Errorf("Invalid timestamp encoding in newSS: %s", hex)
		}

		ts1 := time.Now().UnixNano()
		ts2 := ts1 + 1000

		gNow1 := newSS(ts1)
		if gNow1.Timestamp().UnixNano() != ts1 {
			t.Errorf("GuidSS timestamp mismatch: got %d, want %d", gNow1.Timestamp().UnixNano(), ts1)
		}

		gNow2 := newSS(ts2)
		if bytes.Compare(gNow1.Guid[8:], gNow2.Guid[8:]) >= 0 {
			t.Errorf("GuidSS timestamps are not sortable. gNow1 should be less than gNow2.\ng1: %x\ng2: %x", gNow1.Guid[8:], gNow2.Guid[8:])
		}
	})

	// Check for immediate collision between the two types
	t.Run("CollisionCheck", func(t *testing.T) {
		if NewPG().Guid == NewSS().Guid {
			t.Error("NewPG() and NewSS() produced the same Guid")
		}
	})
}

func TestCachePoolGetPut(t *testing.T) {
	for range 10 {
		guidCacheRef := guidCachePool.Get().(*guidCache)
		guidCachePool.Put(guidCacheRef)
	}
}

func TestNew_ConcurrentPoolExhaustion(t *testing.T) {
	const numGoroutines = 100
	const opsPerRoutine = 10_000

	totalGuids := numGoroutines * opsPerRoutine

	// Pre-allocate local slices per worker to avoid mutex contention during generation
	results := make([][]Guid, numGoroutines)
	for i := range results {
		results[i] = make([]Guid, opsPerRoutine)
	}

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// Background worker repeatedly forces GC to clear sync.Pool caches
	stopGC := make(chan struct{})
	go func() {
		ticker := time.NewTicker(1 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-stopGC:
				return
			case <-ticker.C:
				runtime.GC()
			}
		}
	}()

	// Generate GUIDs concurrently
	for workerID := range numGoroutines {
		go func(id int) {
			defer wg.Done()
			for i := range opsPerRoutine {
				results[id][i] = New()
			}
		}(workerID)
	}

	wg.Wait()
	close(stopGC)

	// Validate non-nil and uniqueness post-execution
	seen := make(map[Guid]struct{}, totalGuids)
	for _, workerGuids := range results {
		for _, g := range workerGuids {
			if g == Nil() {
				t.Errorf("New() returned Nil Guid during pool exhaustion stress")
			}
			if _, exists := seen[g]; exists {
				t.Fatalf("Duplicate GUID detected under concurrent pool exhaustion: %s", g)
			}
			seen[g] = struct{}{}
		}
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

func TestMustParse_Panic(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("MustParse() did not panic on invalid input")
		}
	}()

	// Passing an invalid Base64Url string triggers Parse() to return an error,
	// which causes MustParse() to hit the panic branch.
	_ = MustParse("FBzYOSdp3VEBK7jzXP#leA")
}

func ExampleNew() {
	g := New()     // new random Guid
	fmt.Println(g) // calls g.String(), which returns the Base64Url encoded string
}

func ExampleGuid_String() {
	// g is a 16-byte Guid represented as a hex string "0123456789abcdef0123456789abcdef"
	var g Guid = [16]byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef, 0x10, 0x32, 0x54, 0x76, 0x98, 0xba, 0xdc, 0xfe}
	fmt.Println(g) // calls g.String(), which returns the Base64Url encoded string
	// Output: ASNFZ4mrze8QMlR2mLrc_g
}
