package guid

import (
	"testing"
)

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
	b.ResetTimer() // Reset the timer to exclude setup time (if any).

	for i := 0; i < b.N; i++ {
		_ = New()
	}
}
