package guid

import (
	"testing"
)

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

func duplicatesFound(guids []Guid) bool {
	lenGuids := len(guids)
	guidMap := make(map[Guid]struct{}, lenGuids)

	for i := range guids {
		guidMap[guids[i]] = struct{}{}
	}
	return len(guidMap) != lenGuids
} //duplicatesFound

// BenchmarkNew benchmarks the New function of the guid package.
// to run: go test -bench=".*" -benchmem -benchtime=5s
func BenchmarkGuid_New(b *testing.B) {
	// Reset the timer to exclude setup time (if any).

	for b.Loop() {
		_ = New()
	}
}
