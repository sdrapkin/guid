package guid

import (
	cryptoRand "crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"testing"
	// used for benchmarking - commented out to avoid taking dependencies
	//"github.com/sixafter/nanoid"
	//"github.com/google/uuid"
)

//*******************
// Benchmarks // to run: go test -bench=".*" -benchmem -benchtime=4s
//*******************

/****************************************************************************
C:\Code\guid>go test -bench="(guid.*|base64.*)" -benchmem -benchtime=4s
goos: windows
goarch: amd64
pkg: github.com/sdrapkin/guid
cpu: Intel(R) Core(TM) i7-10510U CPU @ 1.80GHz
Benchmark_guid_New_x10-8                                14520387               313.1 ns/op             0 B/op          0 allocs/op
Benchmark_guid_NewPG_x10-8                              12206565               397.0 ns/op             0 B/op          0 allocs/op
Benchmark_guid_NewSS_x10-8                              12072081               401.1 ns/op             0 B/op          0 allocs/op
Benchmark_guid_New_Parallel_x10-8                       79158543                99.52 ns/op            0 B/op          0 allocs/op
Benchmark_guid_NewString_x10-8                           5486835               848.9 ns/op           240 B/op         10 allocs/op
Benchmark_guid_String_x10-8                              9691195               469.4 ns/op           240 B/op         10 allocs/op
Benchmark_guid_NewString_Parallel_x10-8                 10971610               610.0 ns/op           241 B/op         10 allocs/op
Benchmark_guid_String_x20-8                              4035253              1038 ns/op             480 B/op         20 allocs/op
Benchmark_base64_RawURLEncoding_EncodeToString_x20-8     2678386              1792 ns/op             960 B/op         40 allocs/op
Benchmark_guid_EncodeBase64URL_x20-8                    10934779               418.1 ns/op             0 B/op          0 allocs/op
Benchmark_base64_RawURLEncoding_Encode_x20-8            10476289               479.6 ns/op             0 B/op          0 allocs/op
****************************************************************************/

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

func Benchmark_guid_NewPG_x10(b *testing.B) {
	for b.Loop() {
		_ = NewPG()
		_ = NewPG()
		_ = NewPG()
		_ = NewPG()
		_ = NewPG()
		_ = NewPG()
		_ = NewPG()
		_ = NewPG()
		_ = NewPG()
		_ = NewPG()
	}
}

func Benchmark_guid_NewSS_x10(b *testing.B) {
	for b.Loop() {
		_ = NewSS()
		_ = NewSS()
		_ = NewSS()
		_ = NewSS()
		_ = NewSS()
		_ = NewSS()
		_ = NewSS()
		_ = NewSS()
		_ = NewSS()
		_ = NewSS()
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
*/

/* commented out to avoid taking dependencies
func Benchmark_uuid_New_x10(b *testing.B) {
	uuid.SetRand(nil)
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

func Benchmark_uuid_New_guidRand_x10(b *testing.B) {
	uuid.SetRand(Reader)
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
	uuid.SetRand(nil)
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

func Benchmark_uuid_New_RandPool_guidRand_x10(b *testing.B) {
	uuid.SetRand(Reader)
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
	uuid.SetRand(nil)
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

func Benchmark_uuid_New_Parallel_guidRand_x10(b *testing.B) {
	uuid.SetRand(Reader)
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

func Benchmark_uuid_New_Parallel_RandPool_x10(b *testing.B) {
	uuid.SetRand(nil)
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

func Benchmark_uuid_New_Parallel_RandPool_guidRand_x10(b *testing.B) {
	uuid.SetRand(Reader)
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

func Benchmark_uuid_NewV7_RandPool_x10(b *testing.B) {
	uuid.SetRand(nil)
	uuid.EnableRandPool()
	for b.Loop() {
		_, _ = uuid.NewV7()
		_, _ = uuid.NewV7()
		_, _ = uuid.NewV7()
		_, _ = uuid.NewV7()
		_, _ = uuid.NewV7()
		_, _ = uuid.NewV7()
		_, _ = uuid.NewV7()
		_, _ = uuid.NewV7()
		_, _ = uuid.NewV7()
		_, _ = uuid.NewV7()
	}
}

func Benchmark_uuid_NewV7_Parallel_RandPool_x10(b *testing.B) {
	uuid.SetRand(nil)
	uuid.EnableRandPool()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = uuid.NewV7()
			_, _ = uuid.NewV7()
			_, _ = uuid.NewV7()
			_, _ = uuid.NewV7()
			_, _ = uuid.NewV7()
			_, _ = uuid.NewV7()
			_, _ = uuid.NewV7()
			_, _ = uuid.NewV7()
			_, _ = uuid.NewV7()
			_, _ = uuid.NewV7()
		}
	})
}
*/

func Benchmark_guid_NewPG_Parallel_x10(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = NewPG()
			_ = NewPG()
			_ = NewPG()
			_ = NewPG()
			_ = NewPG()
			_ = NewPG()
			_ = NewPG()
			_ = NewPG()
			_ = NewPG()
			_ = NewPG()
		}
	})
}

//*/

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

func Benchmark_guid_String_x20(b *testing.B) {
	setupBenchGuids()
	for b.Loop() {
		for _, g := range benchGuids {
			_ = g.String()
		}
	}
}

func Benchmark_base64_RawURLEncoding_EncodeToString_x20(b *testing.B) {
	setupBenchGuids()
	for b.Loop() {
		for _, g := range benchGuids {
			_ = base64.RawURLEncoding.EncodeToString(g[:])
		}
	}
}

func Benchmark_guid_EncodeBase64URL_x20(b *testing.B) {
	setupBenchGuids()
	buffer := make([]byte, GuidBase64UrlByteSize)
	for b.Loop() {
		for _, g := range benchGuids {
			g.EncodeBase64URL(buffer)
		}
	}
}

func Benchmark_base64_RawURLEncoding_Encode_x20(b *testing.B) {
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
