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
// go test -bench=".*" -benchmem -benchtime=4s
//*******************

/**************************************************************************** [2026-06-05 15:30:00] Benchmark results, Go 1.27)
C:\Code\Go\guid>go test -run="^$" -bench=(.*) -benchmem -benchtime=4s
goos: windows
goarch: amd64
pkg: github.com/sdrapkin/guid
cpu: Intel(R) Core(TM) i7-10510U CPU @ 1.80GHz
Benchmark_guid_New_x10-8                                15831030               255.6 ns/op             0 B/op          0 allocs/op
Benchmark_guid_CryptoRandRead_x10-8                      2381785              1803 ns/op               0 B/op          0 allocs/op
Benchmark_guid_NewPG_x10-8                              12887008               366.1 ns/op             0 B/op          0 allocs/op
Benchmark_guid_NewSS_x10-8                              12970501               380.4 ns/op             0 B/op          0 allocs/op
Benchmark_guid_New_Parallel_x10-8                       79680348                91.27 ns/op            0 B/op          0 allocs/op
Benchmark_guid_NewString_x10-8                           6830061               686.3 ns/op           240 B/op         10 allocs/op
Benchmark_guid_String_x10-8                             26392231               175.3 ns/op             0 B/op          0 allocs/op
Benchmark_guid_NewString_Parallel_x10-8                 14398298               498.2 ns/op           240 B/op         10 allocs/op
Benchmark_uuid_New_x10-8                                 2640710              1712 ns/op               0 B/op          0 allocs/op
Benchmark_uuid_NewV7_x10-8                               2735016              1819 ns/op               0 B/op          0 allocs/op
Benchmark_uuid_String_x10-8                              6845982               721.1 ns/op           480 B/op         10 allocs/op
Benchmark_uuid_New_Parallel_x10-8                        7924022               748.7 ns/op             0 B/op          0 allocs/op
Benchmark_uuid_NewV7_Parallel_x10-8                      1571799              2818 ns/op               0 B/op          0 allocs/op
Benchmark_guid_NewPG_Parallel_x10-8                     42497318               133.9 ns/op             0 B/op          0 allocs/op
Benchmark_guid_String_x20-8                             11149111               413.8 ns/op             0 B/op          0 allocs/op
Benchmark_base64_RawURLEncoding_EncodeToString_x20-8     4709122              1089 ns/op             480 B/op         20 allocs/op
Benchmark_guid_EncodeBase64URL_x20-8                    11887330               391.7 ns/op             0 B/op          0 allocs/op
Benchmark_base64_RawURLEncoding_Encode_x20-8            11819721               397.3 ns/op             0 B/op          0 allocs/op
Benchmark_Concurrent_CachePool_GetPut/G1-8              1000000000               6.449 ns/op           0 B/op          0 allocs/op
Benchmark_Concurrent_CachePool_GetPut/G2-8              704719772                6.283 ns/op           0 B/op          0 allocs/op
Benchmark_Concurrent_CachePool_GetPut/G4-8              791839826                6.187 ns/op           0 B/op          0 allocs/op
Benchmark_Concurrent_CachePool_GetPut/G8-8              822088278                6.085 ns/op           0 B/op          0 allocs/op
Benchmark_Concurrent_CachePool_GetPut/G16-8             816457746                5.975 ns/op           0 B/op          0 allocs/op
Benchmark_Concurrent_CachePool_GetPut/G32-8             821981848                6.336 ns/op           0 B/op          0 allocs/op
Benchmark_Concurrent_CachePool_GetPut/G64-8             784138318                6.039 ns/op           0 B/op          0 allocs/op
=================================
BenchmarkReadPerf/______Guid_Read([0]byte)-8            1000000000               3.254 ns/op           0 B/op          0 allocs/op
BenchmarkReadPerf/cryptoRand_Read([0]byte)-8            48090409                97.71 ns/op            0 B/op          0 allocs/op
=================================
BenchmarkReadPerf/______Guid_Read([1]byte)-8            179076219               27.72 ns/op            0 B/op          0 allocs/op
BenchmarkReadPerf/cryptoRand_Read([1]byte)-8            34390605               134.7 ns/op             0 B/op          0 allocs/op
=================================
BenchmarkReadPerf/______Guid_Read([2]byte)-8            171541156               27.64 ns/op            0 B/op          0 allocs/op
BenchmarkReadPerf/cryptoRand_Read([2]byte)-8            33717998               137.1 ns/op             0 B/op          0 allocs/op
=================================
BenchmarkReadPerf/______Guid_Read([4]byte)-8            179178990               26.76 ns/op            0 B/op          0 allocs/op
BenchmarkReadPerf/cryptoRand_Read([4]byte)-8            33791941               147.7 ns/op             0 B/op          0 allocs/op
=================================
BenchmarkReadPerf/______Guid_Read([8]byte)-8            178643091               26.81 ns/op            0 B/op          0 allocs/op
BenchmarkReadPerf/cryptoRand_Read([8]byte)-8            31093741               147.6 ns/op             0 B/op          0 allocs/op
=================================
BenchmarkReadPerf/______Guid_Read([16]byte)-8           174099001               29.49 ns/op            0 B/op          0 allocs/op
BenchmarkReadPerf/cryptoRand_Read([16]byte)-8           29316646               165.2 ns/op             0 B/op          0 allocs/op
=================================
BenchmarkReadPerf/______Guid_Read([32]byte)-8           137474346               34.82 ns/op            0 B/op          0 allocs/op
BenchmarkReadPerf/cryptoRand_Read([32]byte)-8           25599958               193.0 ns/op             0 B/op          0 allocs/op
=================================
BenchmarkReadPerf/______Guid_Read([64]byte)-8           90277999                49.72 ns/op            0 B/op          0 allocs/op
BenchmarkReadPerf/cryptoRand_Read([64]byte)-8           19892316               247.9 ns/op             0 B/op          0 allocs/op
=================================
BenchmarkReadPerf/______Guid_Read([128]byte)-8          58561222                81.73 ns/op            0 B/op          0 allocs/op
BenchmarkReadPerf/cryptoRand_Read([128]byte)-8          14416617               335.1 ns/op             0 B/op          0 allocs/op
=================================
BenchmarkReadPerf/______Guid_Read([256]byte)-8          34431975               142.3 ns/op             0 B/op          0 allocs/op
BenchmarkReadPerf/cryptoRand_Read([256]byte)-8          12669002               373.4 ns/op             0 B/op          0 allocs/op
=================================
BenchmarkReadPerf/______Guid_Read([512]byte)-8          19075534               260.1 ns/op             0 B/op          0 allocs/op
BenchmarkReadPerf/cryptoRand_Read([512]byte)-8          10279756               463.7 ns/op             0 B/op          0 allocs/op
=================================
BenchmarkReadPerf/______Guid_Read([513]byte)-8           7993334               526.2 ns/op             0 B/op          0 allocs/op
BenchmarkReadPerf/cryptoRand_Read([513]byte)-8           8422176               553.9 ns/op             0 B/op          0 allocs/op
=================================
BenchmarkReadPerf/______Guid_Read([1024]byte)-8          7454049               639.8 ns/op             0 B/op          0 allocs/op
BenchmarkReadPerf/cryptoRand_Read([1024]byte)-8          7387695               654.8 ns/op             0 B/op          0 allocs/op
=================================
BenchmarkReadPerf/______Guid_Read([2048]byte)-8          4842636               990.5 ns/op             0 B/op          0 allocs/op
BenchmarkReadPerf/cryptoRand_Read([2048]byte)-8          4836138              1035 ns/op               0 B/op          0 allocs/op
=================================
BenchmarkReadPerf/______Guid_Read([4096]byte)-8          2802358              1985 ns/op               0 B/op          0 allocs/op
BenchmarkReadPerf/cryptoRand_Read([4096]byte)-8          2844350              1770 ns/op               0 B/op          0 allocs/op
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

func Benchmark_guid_CryptoRandRead_x10(b *testing.B) {
	var g Guid
	for b.Loop() {
		cryptoRand.Read(g[:])
		cryptoRand.Read(g[:])
		cryptoRand.Read(g[:])
		cryptoRand.Read(g[:])
		cryptoRand.Read(g[:])
		cryptoRand.Read(g[:])
		cryptoRand.Read(g[:])
		cryptoRand.Read(g[:])
		cryptoRand.Read(g[:])
		cryptoRand.Read(g[:])
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

	b.ResetTimer()
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
		_ = uuid01.String()
		_ = uuid02.String()
		_ = uuid03.String()
		_ = uuid04.String()
		_ = uuid05.String()
		_ = uuid06.String()
		_ = uuid07.String()
		_ = uuid08.String()
		_ = uuid09.String()
		_ = uuid10.String()
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
			copy(g[:], bytes)
			benchGuids[i] = g
		}
	}
}

func Benchmark_guid_String_x20(b *testing.B) {
	setupBenchGuids()

	b.ResetTimer()
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

// Used as a benchmark baseline
func _CachePool_GetPut() {
	guidCacheRef := guidCachePool.Get().(*guidCache)
	guidCachePool.Put(guidCacheRef)
}

func Benchmark_Concurrent_CachePool_GetPut(b *testing.B) {
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
