package guid

import (
	"runtime"
	"sync"
	"sync/atomic"
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

// BenchmarkNew benchmarks the New function of the guid package.
// to run: go test -bench=".*" -benchmem -benchtime=5s
func BenchmarkGuid_New(b *testing.B) {
	// Reset the timer to exclude setup time (if any).

	for b.Loop() {
		_ = New()
	}
}

// BenchmarkGuid_New_Parallel benchmarks the New function of the guid package in parallel.
// to run: go test -bench=".*" -cpu=8 -benchmem -benchtime=5s
func BenchmarkGuid_New_Parallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = New()
		}
	})
}
