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