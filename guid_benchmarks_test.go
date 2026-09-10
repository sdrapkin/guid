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
	"uuid"
)

//*******************
// Benchmarks:
// set GOAMD64=v3
// go test -bench="(.*)" -benchmem -benchtime=4s
//*******************

/**************************************************************************** [2026-09-10] Benchmark results, Go 1.27.1)
C:\Code\Go\guid>go test -bench="(.*)" -benchmem -benchtime=4s
goos: windows
goarch: amd64
pkg: github.com/sdrapkin/guid
cpu: Intel(R) Core(TM) i7-10510U CPU @ 1.80GHz
Benchmark_guid_New_x10-8                                21336547               240.1 ns/op             0 B/op          0 allocs/op
Benchmark_CryptoRand_Read_16bytes_x10-8                  2873866              1657 ns/op               0 B/op          0 allocs/op
Benchmark_guid_ReaderRead_16bytes_x10-8                 16466930               312.0 ns/op             0 B/op          0 allocs/op
Benchmark_guid_NewPG_x10-8                              12806184               475.9 ns/op             0 B/op          0 allocs/op
Benchmark_guid_NewSS_x10-8                              10761759               437.9 ns/op             0 B/op          0 allocs/op
Benchmark_guid_New_Parallel_x10-8                       78815274                93.94 ns/op            0 B/op          0 allocs/op
Benchmark_guid_NewString_x10-8                           6844672               676.5 ns/op           240 B/op         10 allocs/op
Benchmark_guid_String_x10-8                             12206503               453.0 ns/op           240 B/op         10 allocs/op
Benchmark_guid_NewString_Parallel_x10-8                 12761317               488.8 ns/op           240 B/op         10 allocs/op
Benchmark_uuid_New_x10-8                                 2559798              1754 ns/op               0 B/op          0 allocs/op
Benchmark_uuid_NewV7_x10-8                               2645582              1809 ns/op               0 B/op          0 allocs/op
Benchmark_uuid_String_x10-8                              5684502              1039 ns/op             480 B/op         10 allocs/op
Benchmark_uuid_New_Parallel_x10-8                        8335995               745.2 ns/op             0 B/op          0 allocs/op
Benchmark_uuid_NewV7_Parallel_x10-8                      1625245              2813 ns/op               0 B/op          0 allocs/op
Benchmark_guid_NewPG_Parallel_x10-8                     44275540               139.3 ns/op             0 B/op          0 allocs/op
Benchmark_guid_String_x20-8                              3648266              1106 ns/op             480 B/op         20 allocs/op
Benchmark_base64_RawURLEncoding_EncodeToString_x20-8     3554299              1351 ns/op             480 B/op         20 allocs/op
Benchmark_guid_EncodeBase64URL_x20-8                    14391660               317.9 ns/op             0 B/op          0 allocs/op
Benchmark_base64_RawURLEncoding_Encode_x20-8            12024792               404.0 ns/op             0 B/op          0 allocs/op
Benchmark_Concurrent_CachePool_GetPut_x10/G1-8          100000000               56.90 ns/op            0 B/op          0 allocs/op
Benchmark_Concurrent_CachePool_GetPut_x10/G2-8          83234491                58.67 ns/op            0 B/op          0 allocs/op
Benchmark_Concurrent_CachePool_GetPut_x10/G4-8          82155907                58.77 ns/op            0 B/op          0 allocs/op
Benchmark_Concurrent_CachePool_GetPut_x10/G8-8          80082484                59.83 ns/op            0 B/op          0 allocs/op
Benchmark_Concurrent_CachePool_GetPut_x10/G16-8         84330216                62.38 ns/op            0 B/op          0 allocs/op
Benchmark_Concurrent_CachePool_GetPut_x10/G32-8         72608018                58.81 ns/op            0 B/op          0 allocs/op
Benchmark_Concurrent_CachePool_GetPut_x10/G64-8         80136230                57.19 ns/op            0 B/op          0 allocs/op
Benchmark_Concurrent_guid_New_x10/G1-8                  53833807                92.60 ns/op            0 B/op          0 allocs/op
Benchmark_Concurrent_guid_New_x10/G2-8                  49439683                93.53 ns/op            0 B/op          0 allocs/op
Benchmark_Concurrent_guid_New_x10/G4-8                  48047566                97.38 ns/op            0 B/op          0 allocs/op
Benchmark_Concurrent_guid_New_x10/G8-8                  49252844                98.78 ns/op            0 B/op          0 allocs/op
Benchmark_Concurrent_guid_New_x10/G16-8                 50590965                96.33 ns/op            0 B/op          0 allocs/op
Benchmark_Concurrent_guid_New_x10/G32-8                 50578330                93.44 ns/op            0 B/op          0 allocs/op
Benchmark_Concurrent_guid_New_x10/G64-8                 50203952                93.66 ns/op            0 B/op          0 allocs/op
=================================
BenchmarkReadPerf/______Guid_Read([0]byte)-8            1000000000               2.790 ns/op           0 B/op          0 allocs/op
BenchmarkReadPerf/cryptoRand_Read([0]byte)-8            46757557                99.08 ns/op            0 B/op          0 allocs/op
=================================
BenchmarkReadPerf/______Guid_Read([1]byte)-8            241170993               19.97 ns/op            0 B/op          0 allocs/op
BenchmarkReadPerf/cryptoRand_Read([1]byte)-8            33365262               136.4 ns/op             0 B/op          0 allocs/op
=================================
BenchmarkReadPerf/______Guid_Read([2]byte)-8            229572110               20.83 ns/op            0 B/op          0 allocs/op
BenchmarkReadPerf/cryptoRand_Read([2]byte)-8            32955669               138.6 ns/op             0 B/op          0 allocs/op
=================================
BenchmarkReadPerf/______Guid_Read([4]byte)-8            228470313               22.12 ns/op            0 B/op          0 allocs/op
BenchmarkReadPerf/cryptoRand_Read([4]byte)-8            32869620               144.4 ns/op             0 B/op          0 allocs/op
=================================
BenchmarkReadPerf/______Guid_Read([8]byte)-8            196320543               25.02 ns/op            0 B/op          0 allocs/op
BenchmarkReadPerf/cryptoRand_Read([8]byte)-8            31621969               148.8 ns/op             0 B/op          0 allocs/op
=================================
BenchmarkReadPerf/______Guid_Read([16]byte)-8           173828503               27.51 ns/op            0 B/op          0 allocs/op
BenchmarkReadPerf/cryptoRand_Read([16]byte)-8           28943839               161.3 ns/op             0 B/op          0 allocs/op
=================================
BenchmarkReadPerf/______Guid_Read([32]byte)-8           138578865               34.67 ns/op            0 B/op          0 allocs/op
BenchmarkReadPerf/cryptoRand_Read([32]byte)-8           25286101               190.7 ns/op             0 B/op          0 allocs/op
=================================
BenchmarkReadPerf/______Guid_Read([64]byte)-8           97037328                51.55 ns/op            0 B/op          0 allocs/op
BenchmarkReadPerf/cryptoRand_Read([64]byte)-8           18749991               250.9 ns/op             0 B/op          0 allocs/op
=================================
BenchmarkReadPerf/______Guid_Read([128]byte)-8          59343585                79.35 ns/op            0 B/op          0 allocs/op
BenchmarkReadPerf/cryptoRand_Read([128]byte)-8          14221682               336.5 ns/op             0 B/op          0 allocs/op
=================================
BenchmarkReadPerf/______Guid_Read([256]byte)-8          32219868               142.4 ns/op             0 B/op          0 allocs/op
BenchmarkReadPerf/cryptoRand_Read([256]byte)-8          12600915               385.9 ns/op             0 B/op          0 allocs/op
=================================
BenchmarkReadPerf/______Guid_Read([512]byte)-8          17206329               422.1 ns/op             0 B/op          0 allocs/op
BenchmarkReadPerf/cryptoRand_Read([512]byte)-8           3328891              1253 ns/op               0 B/op          0 allocs/op
=================================
BenchmarkReadPerf/______Guid_Read([513]byte)-8           4639192              1010 ns/op               0 B/op          0 allocs/op
BenchmarkReadPerf/cryptoRand_Read([513]byte)-8           4708387               941.2 ns/op             0 B/op          0 allocs/op
=================================
BenchmarkReadPerf/______Guid_Read([1024]byte)-8          4345500              1134 ns/op               0 B/op          0 allocs/op
BenchmarkReadPerf/cryptoRand_Read([1024]byte)-8          3851770              1126 ns/op               0 B/op          0 allocs/op
=================================
BenchmarkReadPerf/______Guid_Read([2048]byte)-8          3065995              1545 ns/op               0 B/op          0 allocs/op
BenchmarkReadPerf/cryptoRand_Read([2048]byte)-8          3105973              1478 ns/op               0 B/op          0 allocs/op
=================================
BenchmarkReadPerf/______Guid_Read([4096]byte)-8          1997857              2424 ns/op               0 B/op          0 allocs/op
BenchmarkReadPerf/cryptoRand_Read([4096]byte)-8          1992097              2419 ns/op               0 B/op          0 allocs/op
=================================
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

func Benchmark_CryptoRand_Read_16bytes_x10(b *testing.B) {
	var g Guid
	gSlice := g.UUID[:]
	for b.Loop() {
		cryptoRand.Read(gSlice)
		cryptoRand.Read(gSlice)
		cryptoRand.Read(gSlice)
		cryptoRand.Read(gSlice)
		cryptoRand.Read(gSlice)

		cryptoRand.Read(gSlice)
		cryptoRand.Read(gSlice)
		cryptoRand.Read(gSlice)
		cryptoRand.Read(gSlice)
		cryptoRand.Read(gSlice)
	}
}

func Benchmark_guid_ReaderRead_16bytes_x10(b *testing.B) {
	var g Guid
	gSlice := g.UUID[:]
	for b.Loop() {
		Reader.Read(gSlice)
		Reader.Read(gSlice)
		Reader.Read(gSlice)
		Reader.Read(gSlice)
		Reader.Read(gSlice)

		Reader.Read(gSlice)
		Reader.Read(gSlice)
		Reader.Read(gSlice)
		Reader.Read(gSlice)
		Reader.Read(gSlice)
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

var benchmarkStringSink string

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

	b.ResetTimer()
	for b.Loop() {
		benchmarkStringSink = guid01.String()
		benchmarkStringSink = guid02.String()
		benchmarkStringSink = guid03.String()
		benchmarkStringSink = guid04.String()
		benchmarkStringSink = guid05.String()

		benchmarkStringSink = guid06.String()
		benchmarkStringSink = guid07.String()
		benchmarkStringSink = guid08.String()
		benchmarkStringSink = guid09.String()
		benchmarkStringSink = guid10.String()
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

// commented out to avoid taking dependencies
func Benchmark_uuid_New_x10(b *testing.B) {
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

func Benchmark_uuid_NewV7_x10(b *testing.B) {
	for b.Loop() {
		_ = uuid.NewV7()
		_ = uuid.NewV7()
		_ = uuid.NewV7()
		_ = uuid.NewV7()
		_ = uuid.NewV7()

		_ = uuid.NewV7()
		_ = uuid.NewV7()
		_ = uuid.NewV7()
		_ = uuid.NewV7()
		_ = uuid.NewV7()
	}
}

func Benchmark_uuid_String_x10(b *testing.B) {
	uuid01 := uuid.New()
	uuid02 := uuid.New()
	uuid03 := uuid.New()
	uuid04 := uuid.New()
	uuid05 := uuid.New()

	uuid06 := uuid.New()
	uuid07 := uuid.New()
	uuid08 := uuid.New()
	uuid09 := uuid.New()
	uuid10 := uuid.New()

	b.ResetTimer()
	for b.Loop() {
		benchmarkStringSink = uuid01.String()
		benchmarkStringSink = uuid02.String()
		benchmarkStringSink = uuid03.String()
		benchmarkStringSink = uuid04.String()
		benchmarkStringSink = uuid05.String()

		benchmarkStringSink = uuid06.String()
		benchmarkStringSink = uuid07.String()
		benchmarkStringSink = uuid08.String()
		benchmarkStringSink = uuid09.String()
		benchmarkStringSink = uuid10.String()
	}
}

func Benchmark_uuid_New_Parallel_x10(b *testing.B) {
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

func Benchmark_uuid_NewV7_Parallel_x10(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = uuid.NewV7()
			_ = uuid.NewV7()
			_ = uuid.NewV7()
			_ = uuid.NewV7()
			_ = uuid.NewV7()

			_ = uuid.NewV7()
			_ = uuid.NewV7()
			_ = uuid.NewV7()
			_ = uuid.NewV7()
			_ = uuid.NewV7()
		}
	})
}

/*
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
			copy(g.UUID[:], bytes)
			benchGuids[i] = g
		}
	}
}

func Benchmark_guid_String_x20(b *testing.B) {
	setupBenchGuids()

	b.ResetTimer()
	for b.Loop() {
		for _, g := range benchGuids {
			benchmarkStringSink = g.String()
		}
	}
}

func Benchmark_base64_RawURLEncoding_EncodeToString_x20(b *testing.B) {
	setupBenchGuids()
	for b.Loop() {
		for _, g := range benchGuids {
			benchmarkStringSink = base64.RawURLEncoding.EncodeToString(g.UUID[:])
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
			base64.RawURLEncoding.Encode(buffer, g.UUID[:])
		}
	}
}

// Used as a benchmark baseline
func _CachePool_GetPut() {
	guidCacheRef := guidCachePool.Get().(*guidCache)
	guidCachePool.Put(guidCacheRef)
}

func Benchmark_Concurrent_CachePool_GetPut_x10(b *testing.B) {
	b.ReportAllocs()
	goroutineCounts := []int{1, 2, 4, 8, 16, 32, 64}
	for _, count := range goroutineCounts {
		benchName := fmt.Sprintf("G%d", count)
		b.Run(benchName, func(b *testing.B) {
			b.SetParallelism(count)
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					_CachePool_GetPut()
					_CachePool_GetPut()
					_CachePool_GetPut()
					_CachePool_GetPut()
					_CachePool_GetPut()

					_CachePool_GetPut()
					_CachePool_GetPut()
					_CachePool_GetPut()
					_CachePool_GetPut()
					_CachePool_GetPut()
				}
			})
		})
	}
}

func Benchmark_Concurrent_guid_New_x10(b *testing.B) {
	b.ReportAllocs()
	goroutineCounts := []int{1, 2, 4, 8, 16, 32, 64}
	for _, count := range goroutineCounts {
		benchName := fmt.Sprintf("G%d", count)
		b.Run(benchName, func(b *testing.B) {
			b.SetParallelism(count)
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					New()
					New()
					New()
					New()
					New()

					New()
					New()
					New()
					New()
					New()
				}
			})
		})
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
