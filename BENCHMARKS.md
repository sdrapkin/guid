<details>
  <summary>[2025-07-11]</summary>

```
C:\Code\guid>go test -run=$^ -bench="(.*)" -benchmem -benchtime=4s .
goos: windows
goarch: amd64
pkg: github.com/sdrapkin/guid
cpu: Intel(R) Core(TM) i7-10510U CPU @ 1.80GHz
Benchmark_guid_New_x10-8                                20314392               230.1 ns/op             0 B/op          0 allocs/op
Benchmark_guid_NewPG_x10-8                              16182096               359.1 ns/op             0 B/op          0 allocs/op
Benchmark_guid_NewSS_x10-8                              12714096               396.9 ns/op             0 B/op          0 allocs/op
Benchmark_guid_New_Parallel_x10-8                       79524312                97.65 ns/op            0 B/op          0 allocs/op
Benchmark_guid_NewString_x10-8                           5694988               825.0 ns/op           240 B/op         10 allocs/op
Benchmark_guid_String_x10-8                             10582099               460.6 ns/op           240 B/op         10 allocs/op
Benchmark_guid_NewString_Parallel_x10-8                 14056472               532.1 ns/op           240 B/op         10 allocs/op
Benchmark_nanoid_New_x10-8                               2056082              2203 ns/op             240 B/op         10 allocs/op
Benchmark_nanoid_New_Parallel_x10-8                      5134498              1299 ns/op             240 B/op         10 allocs/op
Benchmark_uuid_New_x10-8                                 1838532              2524 ns/op             160 B/op         10 allocs/op
Benchmark_uuid_New_guidRand_x10-8                        6633441               781.4 ns/op           160 B/op         10 allocs/op
Benchmark_uuid_New_RandPool_x10-8                        8128683               539.8 ns/op             0 B/op          0 allocs/op
Benchmark_uuid_New_RandPool_guidRand_x10-8              13006801               404.5 ns/op             0 B/op          0 allocs/op
Benchmark_uuid_New_Parallel_x10-8                        4987186              1082 ns/op             160 B/op         10 allocs/op
Benchmark_uuid_New_Parallel_guidRand_x10-8              10500265               493.0 ns/op           160 B/op         10 allocs/op
Benchmark_uuid_New_Parallel_RandPool_x10-8               3167839              1311 ns/op               0 B/op          0 allocs/op
Benchmark_uuid_New_Parallel_RandPool_guidRand_x10-8      4164766              1236 ns/op               0 B/op          0 allocs/op
Benchmark_uuid_NewV7_RandPool_x10-8                      5048066               936.4 ns/op             0 B/op          0 allocs/op
Benchmark_uuid_NewV7_Parallel_RandPool_x10-8             2024038              2483 ns/op               0 B/op          0 allocs/op
Benchmark_guid_NewPG_Parallel_x10-8                     44894469               143.0 ns/op             0 B/op          0 allocs/op
Benchmark_guid_String_x20-8                              3530682              1216 ns/op             480 B/op         20 allocs/op
Benchmark_base64_RawURLEncoding_EncodeToString_x20-8     2832835              1688 ns/op             960 B/op         40 allocs/op
Benchmark_guid_EncodeBase64URL_x20-8                    11753940               430.7 ns/op             0 B/op          0 allocs/op
Benchmark_base64_RawURLEncoding_Encode_x20-8            10445786               453.7 ns/op             0 B/op          0 allocs/op
```
</details>
<details>
  <summary>[2025-07-12]</summary>

```
C:\Code\guid>go test -run="$^" -bench="()"
goos: windows
goarch: amd64
pkg: github.com/sdrapkin/guid
cpu: Intel(R) Core(TM) i7-10510U CPU @ 1.80GHz
Benchmark_guid_New_x10-8                                 4545321               268.7 ns/op
Benchmark_guid_NewPG_x10-8                               2994448               374.6 ns/op
Benchmark_guid_NewSS_x10-8                               3776004               298.5 ns/op
Benchmark_guid_New_Parallel_x10-8                       19469299                60.83 ns/op
Benchmark_guid_NewString_x10-8                           1901104               592.7 ns/op
Benchmark_guid_String_x10-8                              3122826               402.2 ns/op
Benchmark_guid_NewString_Parallel_x10-8                  2929981               357.4 ns/op
Benchmark_guid_NewPG_Parallel_x10-8                     12018243                90.94 ns/op
Benchmark_guid_String_x20-8                              1480188               802.3 ns/op
Benchmark_base64_RawURLEncoding_EncodeToString_x20-8      885295              1257 ns/op
Benchmark_guid_EncodeBase64URL_x20-8                     3792573               313.3 ns/op
Benchmark_base64_RawURLEncoding_Encode_x20-8             3437497               346.2 ns/op
=================================
BenchmarkReadPerf/______Guid_Read([0]byte)-8            438159452                2.821 ns/op
BenchmarkReadPerf/cryptoRand_Read([0]byte)-8            10378575               104.8 ns/op
=================================
BenchmarkReadPerf/______Guid_Read([1]byte)-8            51085567                22.81 ns/op
BenchmarkReadPerf/cryptoRand_Read([1]byte)-8             8212927               143.0 ns/op
=================================
BenchmarkReadPerf/______Guid_Read([2]byte)-8            50509088                22.72 ns/op
BenchmarkReadPerf/cryptoRand_Read([2]byte)-8             9037321               131.8 ns/op
=================================
BenchmarkReadPerf/______Guid_Read([4]byte)-8            44612982                22.96 ns/op
BenchmarkReadPerf/cryptoRand_Read([4]byte)-8             8493836               137.2 ns/op
=================================
BenchmarkReadPerf/______Guid_Read([8]byte)-8            45245286                23.30 ns/op
BenchmarkReadPerf/cryptoRand_Read([8]byte)-8             8389718               148.2 ns/op
=================================
BenchmarkReadPerf/______Guid_Read([16]byte)-8           39388945                28.61 ns/op
BenchmarkReadPerf/cryptoRand_Read([16]byte)-8            7003928               168.5 ns/op
=================================
BenchmarkReadPerf/______Guid_Read([32]byte)-8           35560718                34.28 ns/op
BenchmarkReadPerf/cryptoRand_Read([32]byte)-8            5312515               225.3 ns/op
=================================
BenchmarkReadPerf/______Guid_Read([64]byte)-8           21756630                53.73 ns/op
BenchmarkReadPerf/cryptoRand_Read([64]byte)-8            4016278               291.9 ns/op
=================================
BenchmarkReadPerf/______Guid_Read([128]byte)-8          13214228                85.15 ns/op
BenchmarkReadPerf/cryptoRand_Read([128]byte)-8           3098580               385.1 ns/op
=================================
BenchmarkReadPerf/______Guid_Read([256]byte)-8           7831952               149.4 ns/op
BenchmarkReadPerf/cryptoRand_Read([256]byte)-8           2553850               451.6 ns/op
=================================
BenchmarkReadPerf/______Guid_Read([512]byte)-8           4317787               277.8 ns/op
BenchmarkReadPerf/cryptoRand_Read([512]byte)-8           2003790               580.7 ns/op
=================================
BenchmarkReadPerf/______Guid_Read([513]byte)-8           2075565               576.8 ns/op
BenchmarkReadPerf/cryptoRand_Read([513]byte)-8           2075202               554.8 ns/op
=================================
BenchmarkReadPerf/______Guid_Read([1024]byte)-8          1914433               620.0 ns/op
BenchmarkReadPerf/cryptoRand_Read([1024]byte)-8          1950085               615.4 ns/op
=================================
BenchmarkReadPerf/______Guid_Read([2048]byte)-8          1269326               947.6 ns/op
BenchmarkReadPerf/cryptoRand_Read([2048]byte)-8          1253407               956.1 ns/op
=================================
BenchmarkReadPerf/______Guid_Read([4096]byte)-8           727717              1759 ns/op
BenchmarkReadPerf/cryptoRand_Read([4096]byte)-8           693541              1680 ns/op
=================================
```
</details>