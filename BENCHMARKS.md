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
C:\Code\guid>go test -run=$^ -bench="(.*)" -benchmem -benchtime=4s .
goos: windows
goarch: amd64
pkg: github.com/sdrapkin/guid
cpu: Intel(R) Core(TM) i7-10510U CPU @ 1.80GHz
Benchmark_guid_New_x10-8                                17702509               227.6 ns/op             0 B/op          0 allocs/op
Benchmark_guid_NewPG_x10-8                              15895365               290.2 ns/op             0 B/op          0 allocs/op
Benchmark_guid_NewSS_x10-8                              16453174               302.7 ns/op             0 B/op          0 allocs/op
Benchmark_guid_New_Parallel_x10-8                       79849218                60.87 ns/op            0 B/op          0 allocs/op
Benchmark_guid_NewString_x10-8                           7053037               642.6 ns/op           240 B/op         10 allocs/op
Benchmark_guid_String_x10-8                             12385693               396.8 ns/op           240 B/op         10 allocs/op
Benchmark_guid_NewString_Parallel_x10-8                 14497309               393.8 ns/op           240 B/op         10 allocs/op
Benchmark_guid_NewPG_Parallel_x10-8                     50365094               104.6 ns/op             0 B/op          0 allocs/op
Benchmark_guid_String_x20-8                              5118271               870.7 ns/op           480 B/op         20 allocs/op
Benchmark_base64_RawURLEncoding_EncodeToString_x20-8     3445861              1392 ns/op             960 B/op         40 allocs/op
Benchmark_guid_EncodeBase64URL_x20-8                    13566518               352.6 ns/op             0 B/op          0 allocs/op
Benchmark_base64_RawURLEncoding_Encode_x20-8            10455889               406.0 ns/op             0 B/op          0 allocs/op
=================================
BenchmarkReadPerf/______Guid_Read([0]byte)-8            1000000000               4.000 ns/op           0 B/op          0 allocs/op
BenchmarkReadPerf/cryptoRand_Read([0]byte)-8            40660288               123.3 ns/op             0 B/op          0 allocs/op
=================================
BenchmarkReadPerf/______Guid_Read([1]byte)-8            250025704               19.94 ns/op            0 B/op          0 allocs/op
BenchmarkReadPerf/cryptoRand_Read([1]byte)-8            30590725               166.8 ns/op             0 B/op          0 allocs/op
=================================
BenchmarkReadPerf/______Guid_Read([2]byte)-8            241538685               20.61 ns/op            0 B/op          0 allocs/op
BenchmarkReadPerf/cryptoRand_Read([2]byte)-8            29171120               178.4 ns/op             0 B/op          0 allocs/op
=================================
BenchmarkReadPerf/______Guid_Read([4]byte)-8            232269996               20.47 ns/op            0 B/op          0 allocs/op
BenchmarkReadPerf/cryptoRand_Read([4]byte)-8            28933527               174.2 ns/op             0 B/op          0 allocs/op
=================================
BenchmarkReadPerf/______Guid_Read([8]byte)-8            222574704               21.23 ns/op            0 B/op          0 allocs/op
BenchmarkReadPerf/cryptoRand_Read([8]byte)-8            28294078               178.4 ns/op             0 B/op          0 allocs/op
=================================
BenchmarkReadPerf/______Guid_Read([16]byte)-8           205734214               22.06 ns/op            0 B/op          0 allocs/op
BenchmarkReadPerf/cryptoRand_Read([16]byte)-8           26099806               188.0 ns/op             0 B/op          0 allocs/op
=================================
BenchmarkReadPerf/______Guid_Read([32]byte)-8           194641297               23.37 ns/op            0 B/op          0 allocs/op
BenchmarkReadPerf/cryptoRand_Read([32]byte)-8           19012092               220.7 ns/op             0 B/op          0 allocs/op
=================================
BenchmarkReadPerf/______Guid_Read([64]byte)-8           215419802               24.22 ns/op            0 B/op          0 allocs/op
BenchmarkReadPerf/cryptoRand_Read([64]byte)-8           15618176               276.8 ns/op             0 B/op          0 allocs/op
=================================
BenchmarkReadPerf/______Guid_Read([128]byte)-8          203188810               23.91 ns/op            0 B/op          0 allocs/op
BenchmarkReadPerf/cryptoRand_Read([128]byte)-8           9112450               455.6 ns/op             0 B/op          0 allocs/op
=================================
BenchmarkReadPerf/______Guid_Read([256]byte)-8          146299468               33.30 ns/op            0 B/op          0 allocs/op
BenchmarkReadPerf/cryptoRand_Read([256]byte)-8           9549841               552.6 ns/op             0 B/op          0 allocs/op
=================================
BenchmarkReadPerf/______Guid_Read([512]byte)-8          133246482               35.12 ns/op            0 B/op          0 allocs/op
BenchmarkReadPerf/cryptoRand_Read([512]byte)-8           9626215               514.6 ns/op             0 B/op          0 allocs/op
=================================
BenchmarkReadPerf/______Guid_Read([513]byte)-8           9261601               531.1 ns/op             0 B/op          0 allocs/op
BenchmarkReadPerf/cryptoRand_Read([513]byte)-8           8733603               559.3 ns/op             0 B/op          0 allocs/op
=================================
BenchmarkReadPerf/______Guid_Read([1024]byte)-8          6781808               667.4 ns/op             0 B/op          0 allocs/op
BenchmarkReadPerf/cryptoRand_Read([1024]byte)-8          7346364               673.7 ns/op             0 B/op          0 allocs/op
=================================
BenchmarkReadPerf/______Guid_Read([2048]byte)-8          4135266              1078 ns/op               0 B/op          0 allocs/op
BenchmarkReadPerf/cryptoRand_Read([2048]byte)-8          4490186              1031 ns/op               0 B/op          0 allocs/op
=================================
BenchmarkReadPerf/______Guid_Read([4096]byte)-8          2629774              1822 ns/op               0 B/op          0 allocs/op
BenchmarkReadPerf/cryptoRand_Read([4096]byte)-8          2831491              1744 ns/op               0 B/op          0 allocs/op
=================================
```
</details>