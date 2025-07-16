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

<details><summary>[2025-07-16]</summary>

```
C:\Code\Go\rand-reader-bench>go test -bench="(.*)"
goos: windows
goarch: amd64
pkg: rand-reader-bench
cpu: Intel(R) Core(TM) i7-10510U CPU @ 1.80GHz
Benchmark_Readers_Nano_____Serial/Size_0-8      331600810                3.371 ns/op
Benchmark_Readers_Nano_____Serial/Size_1-8      35385078                32.33 ns/op       30.93 MB/s
Benchmark_Readers_Nano_____Serial/Size_2-8      31481024                32.11 ns/op       62.30 MB/s
Benchmark_Readers_Nano_____Serial/Size_4-8      38069278                31.12 ns/op      128.53 MB/s
Benchmark_Readers_Nano_____Serial/Size_8-8      31223816                38.83 ns/op      206.03 MB/s
Benchmark_Readers_Nano_____Serial/Size_16-8     21874976                53.30 ns/op      300.21 MB/s
Benchmark_Readers_Nano_____Serial/Size_32-8     14617371                82.15 ns/op      389.54 MB/s
Benchmark_Readers_Nano_____Serial/Size_64-8      9386101               145.0 ns/op       441.37 MB/s
Benchmark_Readers_Nano_____Serial/Size_128-8     4131087               289.9 ns/op       441.60 MB/s
Benchmark_Readers_Nano_____Serial/Size_256-8     2019602               586.5 ns/op       436.49 MB/s
Benchmark_Readers_Nano_____Serial/Size_512-8     1000000              1043 ns/op         490.67 MB/s
Benchmark_Readers_Nano_____Serial/Size_513-8     1000000              1240 ns/op         413.76 MB/s
Benchmark_Readers_Nano_____Serial/Size_1024-8     460317              2339 ns/op         437.73 MB/s
Benchmark_Readers_Nano_____Serial/Size_2048-8     291912              4166 ns/op         491.62 MB/s
Benchmark_Readers_Nano_____Serial/Size_4096-8     142443              7891 ns/op         519.05 MB/s
Benchmark_Readers_Rand_____Serial/Size_0-8       9962391               126.3 ns/op
Benchmark_Readers_Rand_____Serial/Size_1-8       7508000               163.7 ns/op         6.11 MB/s
Benchmark_Readers_Rand_____Serial/Size_2-8       7273264               166.1 ns/op        12.04 MB/s
Benchmark_Readers_Rand_____Serial/Size_4-8       6966618               173.8 ns/op        23.02 MB/s
Benchmark_Readers_Rand_____Serial/Size_8-8       6277897               181.1 ns/op        44.17 MB/s
Benchmark_Readers_Rand_____Serial/Size_16-8      6426225               195.6 ns/op        81.82 MB/s
Benchmark_Readers_Rand_____Serial/Size_32-8      4705269               240.3 ns/op       133.15 MB/s
Benchmark_Readers_Rand_____Serial/Size_64-8      3733690               301.7 ns/op       212.12 MB/s
Benchmark_Readers_Rand_____Serial/Size_128-8     3220525               367.4 ns/op       348.37 MB/s
Benchmark_Readers_Rand_____Serial/Size_256-8     2860806               409.9 ns/op       624.52 MB/s
Benchmark_Readers_Rand_____Serial/Size_512-8     2431822               492.1 ns/op      1040.49 MB/s
Benchmark_Readers_Rand_____Serial/Size_513-8     2233519               538.0 ns/op       953.46 MB/s
Benchmark_Readers_Rand_____Serial/Size_1024-8    1801224               680.2 ns/op      1505.48 MB/s
Benchmark_Readers_Rand_____Serial/Size_2048-8    1000000              1092 ns/op        1875.74 MB/s
Benchmark_Readers_Rand_____Serial/Size_4096-8     634436              1911 ns/op        2143.88 MB/s
Benchmark_Readers_Guid_____Serial/Size_0-8      203178937                5.642 ns/op
Benchmark_Readers_Guid_____Serial/Size_1-8      40372231                30.17 ns/op       33.15 MB/s
Benchmark_Readers_Guid_____Serial/Size_2-8      39767624                30.50 ns/op       65.58 MB/s
Benchmark_Readers_Guid_____Serial/Size_4-8      40736654                31.67 ns/op      126.31 MB/s
Benchmark_Readers_Guid_____Serial/Size_8-8      40896450                29.77 ns/op      268.75 MB/s
Benchmark_Readers_Guid_____Serial/Size_16-8     40024147                30.53 ns/op      524.15 MB/s
Benchmark_Readers_Guid_____Serial/Size_32-8     30952884                38.47 ns/op      831.81 MB/s
Benchmark_Readers_Guid_____Serial/Size_64-8     21578154                54.47 ns/op     1174.90 MB/s
Benchmark_Readers_Guid_____Serial/Size_128-8    12806733                94.02 ns/op     1361.47 MB/s
Benchmark_Readers_Guid_____Serial/Size_256-8     7064658               161.3 ns/op      1586.67 MB/s
Benchmark_Readers_Guid_____Serial/Size_512-8     4589568               260.1 ns/op      1968.63 MB/s
Benchmark_Readers_Guid_____Serial/Size_513-8     2268535               531.8 ns/op       964.69 MB/s
Benchmark_Readers_Guid_____Serial/Size_1024-8    1800406               666.1 ns/op      1537.30 MB/s
Benchmark_Readers_Guid_____Serial/Size_2048-8    1000000              1044 ns/op        1961.69 MB/s
Benchmark_Readers_Guid_____Serial/Size_4096-8     680055              1765 ns/op        2320.21 MB/s
Benchmark_Readers_GuidLite_Serial/Size_0-8      506311594                2.361 ns/op
Benchmark_Readers_GuidLite_Serial/Size_1-8      47393551                27.45 ns/op       36.43 MB/s
Benchmark_Readers_GuidLite_Serial/Size_2-8      40033093                29.10 ns/op       68.74 MB/s
Benchmark_Readers_GuidLite_Serial/Size_4-8      39776455                30.20 ns/op      132.44 MB/s
Benchmark_Readers_GuidLite_Serial/Size_8-8      35498126                32.05 ns/op      249.58 MB/s
Benchmark_Readers_GuidLite_Serial/Size_16-8     36818624                32.55 ns/op      491.51 MB/s
Benchmark_Readers_GuidLite_Serial/Size_32-8     27219710                45.03 ns/op      710.61 MB/s
Benchmark_Readers_GuidLite_Serial/Size_64-8     16134301                72.18 ns/op      886.61 MB/s
Benchmark_Readers_GuidLite_Serial/Size_128-8     9405862               129.0 ns/op       992.51 MB/s
Benchmark_Readers_GuidLite_Serial/Size_256-8     4865019               242.3 ns/op      1056.45 MB/s
Benchmark_Readers_GuidLite_Serial/Size_512-8     2650377               453.7 ns/op      1128.44 MB/s
Benchmark_Readers_GuidLite_Serial/Size_513-8     2222863               523.4 ns/op       980.18 MB/s
Benchmark_Readers_GuidLite_Serial/Size_1024-8    1340719               887.0 ns/op      1154.42 MB/s
Benchmark_Readers_GuidLite_Serial/Size_2048-8     703270              1695 ns/op        1208.36 MB/s
Benchmark_Readers_GuidLite_Serial/Size_4096-8     360790              3593 ns/op        1139.95 MB/s
Benchmark_Readers_Nano_____Concurrent/Size_0_G64-8              1000000000               1.499 ns/op
Benchmark_Readers_Nano_____Concurrent/Size_1_G64-8              92699883                12.70 ns/op       78.74 MB/s
Benchmark_Readers_Nano_____Concurrent/Size_2_G64-8              78789788                14.05 ns/op      142.31 MB/s
Benchmark_Readers_Nano_____Concurrent/Size_4_G64-8              70041090                17.04 ns/op      234.73 MB/s
Benchmark_Readers_Nano_____Concurrent/Size_8_G64-8              45785798                22.28 ns/op      359.14 MB/s
Benchmark_Readers_Nano_____Concurrent/Size_16_G64-8             32973645                32.86 ns/op      486.96 MB/s
Benchmark_Readers_Nano_____Concurrent/Size_32_G64-8             24374142                48.55 ns/op      659.17 MB/s
Benchmark_Readers_Nano_____Concurrent/Size_64_G64-8             17045211                66.48 ns/op      962.71 MB/s
Benchmark_Readers_Nano_____Concurrent/Size_128_G64-8            10404382               110.0 ns/op      1163.64 MB/s
Benchmark_Readers_Nano_____Concurrent/Size_256_G64-8             5703511               206.4 ns/op      1240.12 MB/s
Benchmark_Readers_Nano_____Concurrent/Size_512_G64-8             3118212               373.1 ns/op      1372.36 MB/s
Benchmark_Readers_Nano_____Concurrent/Size_513_G64-8             2986875               399.8 ns/op      1283.02 MB/s
Benchmark_Readers_Nano_____Concurrent/Size_1024_G64-8            1687882               700.3 ns/op      1462.24 MB/s
Benchmark_Readers_Nano_____Concurrent/Size_2048_G64-8             883131              1345 ns/op        1522.41 MB/s
Benchmark_Readers_Nano_____Concurrent/Size_4096_G64-8             438458              2716 ns/op        1507.99 MB/s
Benchmark_Readers_Rand_____Concurrent/Size_0_G64-8              23781117                52.60 ns/op
Benchmark_Readers_Rand_____Concurrent/Size_1_G64-8              18006662                68.29 ns/op       14.64 MB/s
Benchmark_Readers_Rand_____Concurrent/Size_2_G64-8              16359828                71.18 ns/op       28.10 MB/s
Benchmark_Readers_Rand_____Concurrent/Size_4_G64-8              16545926                76.65 ns/op       52.18 MB/s
Benchmark_Readers_Rand_____Concurrent/Size_8_G64-8              15508323                79.30 ns/op      100.89 MB/s
Benchmark_Readers_Rand_____Concurrent/Size_16_G64-8             14549739                86.20 ns/op      185.62 MB/s
Benchmark_Readers_Rand_____Concurrent/Size_32_G64-8             12469306                98.70 ns/op      324.22 MB/s
Benchmark_Readers_Rand_____Concurrent/Size_64_G64-8             10256085               125.3 ns/op       510.75 MB/s
Benchmark_Readers_Rand_____Concurrent/Size_128_G64-8             8063913               158.4 ns/op       808.04 MB/s
Benchmark_Readers_Rand_____Concurrent/Size_256_G64-8             6816264               166.8 ns/op      1534.98 MB/s
Benchmark_Readers_Rand_____Concurrent/Size_512_G64-8             6225489               188.7 ns/op      2713.29 MB/s
Benchmark_Readers_Rand_____Concurrent/Size_513_G64-8             5669592               201.0 ns/op      2552.73 MB/s
Benchmark_Readers_Rand_____Concurrent/Size_1024_G64-8            4592985               242.4 ns/op      4224.55 MB/s
Benchmark_Readers_Rand_____Concurrent/Size_2048_G64-8            3460918               342.6 ns/op      5978.60 MB/s
Benchmark_Readers_Rand_____Concurrent/Size_4096_G64-8            1762368               667.5 ns/op      6136.35 MB/s
Benchmark_Readers_Guid_____Concurrent/Size_0_G64-8              1000000000               1.051 ns/op
Benchmark_Readers_Guid_____Concurrent/Size_1_G64-8              100000000               10.80 ns/op       92.60 MB/s
Benchmark_Readers_Guid_____Concurrent/Size_2_G64-8              100000000               11.38 ns/op      175.75 MB/s
Benchmark_Readers_Guid_____Concurrent/Size_4_G64-8              100000000               11.86 ns/op      337.33 MB/s
Benchmark_Readers_Guid_____Concurrent/Size_8_G64-8              100000000               11.23 ns/op      712.17 MB/s
Benchmark_Readers_Guid_____Concurrent/Size_16_G64-8             99468670                11.32 ns/op     1412.93 MB/s
Benchmark_Readers_Guid_____Concurrent/Size_32_G64-8             82896972                13.76 ns/op     2325.84 MB/s
Benchmark_Readers_Guid_____Concurrent/Size_64_G64-8             58705542                18.46 ns/op     3466.71 MB/s
Benchmark_Readers_Guid_____Concurrent/Size_128_G64-8            39580185                28.82 ns/op     4440.60 MB/s
Benchmark_Readers_Guid_____Concurrent/Size_256_G64-8            20618024                56.18 ns/op     4557.08 MB/s
Benchmark_Readers_Guid_____Concurrent/Size_512_G64-8            11164046                98.86 ns/op     5179.30 MB/s
Benchmark_Readers_Guid_____Concurrent/Size_513_G64-8             6168399               200.2 ns/op      2562.32 MB/s
Benchmark_Readers_Guid_____Concurrent/Size_1024_G64-8            4889889               249.0 ns/op      4112.73 MB/s
Benchmark_Readers_Guid_____Concurrent/Size_2048_G64-8            3330484               349.0 ns/op      5868.65 MB/s
Benchmark_Readers_Guid_____Concurrent/Size_4096_G64-8            1715498               718.6 ns/op      5700.19 MB/s
Benchmark_Readers_GuidLite_____Concurrent/Size_0_G64-8          1000000000               0.9977 ns/op
Benchmark_Readers_GuidLite_____Concurrent/Size_1_G64-8          100000000               10.51 ns/op       95.14 MB/s
Benchmark_Readers_GuidLite_____Concurrent/Size_2_G64-8          95598486                10.96 ns/op      182.44 MB/s
Benchmark_Readers_GuidLite_____Concurrent/Size_4_G64-8          94883412                10.78 ns/op      371.02 MB/s
Benchmark_Readers_GuidLite_____Concurrent/Size_8_G64-8          121873376               10.17 ns/op      786.83 MB/s
Benchmark_Readers_GuidLite_____Concurrent/Size_16_G64-8         90446578                12.68 ns/op     1261.88 MB/s
Benchmark_Readers_GuidLite_____Concurrent/Size_32_G64-8         66514053                17.74 ns/op     1803.85 MB/s
Benchmark_Readers_GuidLite_____Concurrent/Size_64_G64-8         41432028                28.13 ns/op     2275.31 MB/s
Benchmark_Readers_GuidLite_____Concurrent/Size_128_G64-8        23327151                51.98 ns/op     2462.28 MB/s
Benchmark_Readers_GuidLite_____Concurrent/Size_256_G64-8        13203890                91.20 ns/op     2807.05 MB/s
Benchmark_Readers_GuidLite_____Concurrent/Size_512_G64-8         7134079               187.0 ns/op      2737.48 MB/s
Benchmark_Readers_GuidLite_____Concurrent/Size_513_G64-8         5937442               205.8 ns/op      2492.64 MB/s
Benchmark_Readers_GuidLite_____Concurrent/Size_1024_G64-8        3010873               381.9 ns/op      2681.05 MB/s
Benchmark_Readers_GuidLite_____Concurrent/Size_2048_G64-8        1726142               700.6 ns/op      2923.24 MB/s
Benchmark_Readers_GuidLite_____Concurrent/Size_4096_G64-8         762068              1450 ns/op        2825.14 MB/s
```
</details>