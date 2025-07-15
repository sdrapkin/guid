# guid [![name](https://goreportcard.com/badge/github.com/sdrapkin/guid)](https://goreportcard.com/report/github.com/sdrapkin/guid) [![codecov](https://codecov.io/github/sdrapkin/guid/branch/master/graph/badge.svg?token=ARQFUQD5VP)](https://codecov.io/github/sdrapkin/guid) [![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go#uuid) 
## Fast cryptographically secure Guid generator for Go.<br>By [Stan Drapkin](https://github.com/sdrapkin/).

`Guid` is defined as `type Guid [16]byte` and filled with 128 cryptographically strong bits.

[Go playground](https://go.dev/play/p/l_Yj74HUpgl)
```go
package main

import (
	"fmt"

	"github.com/sdrapkin/guid"
)

func main() {
	for range 4 {
		fmt.Printf("%x\n", guid.New())
	}
	fmt.Println()
	for range 4 {
		g := guid.New()
		fmt.Println(&g) // calls g.String()
	}
}
```

```
79c9779af20dcd21fbe60f3b336ed08c
da2026d38edca4371a476efd41333d23
88c3033b002b0e73321509ef26de607f
a84e961ff7f09f5210ea04585f152e73

WF8MvK5CUOrI-enEuvS0jw
AOp8Voi5knpu1mg3RjzmSg
gxOQRIVR4B_uGHD6OP76XA
Zo_hpnDxkOsAWLk1tIS6DA
```

## Why `guid`? 🔥

`guid` is a high-performance, cryptographically secure UUID/GUID (Globally Unique Identifier) generator for Go. It's built for speed without compromising on security, offering a significant performance advantage—up to **10x faster** than `github.com/google/uuid`.

Beyond raw speed, `guid` offers:

* **Cryptographically Strong**: Generates 128 cryptographically secure bits for robust, unique identifiers.
* **Optimized for Databases**: Includes special `GuidPG` and `GuidSS` types that generate sequential Guids, dramatically improving `INSERT` performance and preventing index fragmentation in PostgreSQL and SQL Server databases.
* **Seamless Interoperability**: Easily integrate with existing `google/uuid` codebases, and even boost `uuid`'s performance by up to **4x** using `guid.Reader`.
* **FIPS 140 Compliant**: Ensures adherence to stringent security standards.
* **Zero Allocations for Core Operations**: `guid.New()` generates new Guids with no memory allocations, making it incredibly efficient.

## Guid is ~10x faster than `github.com/google/uuid` 🔥

* `guid.New()` is  6~10 ns 
* `guid.NewString()` is 40~60 ns
* `String()` on existing guid is ~40 ns
* multi-goroutine calls do not increase per-call latency
* if your library is faster - please let me know!

## API Overview
**All APIs are safe for concurrent use by multiple goroutines.**
| Functions | Description |
|---|---|
| `guid.New()` `Guid`           | Generate a new Guid |
| `guid.NewString()` `string`   | Generate a new Guid as a Base64Url string |
| `guid.NewPG()` `GuidPG`       | Generate a new PostgreSQL sequential Guid |
| `guid.NewSS()` `GuidSS`       | Generate a new SQL Server sequential Guid |
| `guid.Parse(s string)` `(Guid, error)` | Parse a Base64Url string into a Guid |
| `guid.ParseBytes(src []byte)` `(Guid, error)` | Parse Base64Url bytes to a Guid |
| `guid.FromBytes(src []byte)` `(Guid, error)`  | Parse 16-byte slice to a Guid |
| `guid.DecodeBase64URL(dst []byte, src []byte)` `(ok bool)` | Decode a Base64Url slice into a Guid slice |
| `guid.Reader` 🔥 implements `io.Reader`    | Faster alternative to `crypto/rand` |
| guid.Nil                    | The zero-value Guid |

| `Guid` methods | Description |
|---|---|
| `.String()` `string` | Encodes the Guid into Base64Url 22-char string `fmt.Stringer` |
| `.EncodeBase64URL(dst []byte)` `error` | Like `.String()` but encodes into len(22) byte slice |
| .MarshalBinary() | Implements `encoding.BinaryMarshaler` |
| .UnmarshalBinary() | Implements `encoding.BinaryUnmarshaler` |
| .MarshalText() | Implements `encoding.TextMarshaler` |
| .UnmarshalText() | Implements `encoding.TextUnmarshaler` |

| `GuidPG`, `GuidSS` methods | Description |
|---|---|
| `.Timestamp()` `time.Time` | Extracts the UTC timestamp |

## Sequential Guids 🔥
`guid` includes two special types `GuidPG` and `GuidSS` optimized for use as database primary keys (PostgreSQL and SQL Server). Their time-ordered composition helps prevent index fragmentation and improves `INSERT` performance compared to fully random Guids. Note that sequential sorting is only across `time.Now()` timestamp precision.

* **`guid.NewPG()`**: Generates a `GuidPG`, which is sortable in **PostgreSQL**.
 	- It is structured as `[8-byte timestamp][8 random bytes]`.
* **`guid.NewSS()`**: Generates a `GuidSS`, which is sortable in **SQL Server**.
	- It is structured as `[8 random bytes][8-byte SQL Server-ordered timestamp]`.
* `.Timestamp()` on `GuidPG`/`GuidSS` returns Guid creation time as UTC `time.Time`.

Both `GuidPG` and `GuidSS` are nearly as fast as `guid.New()`. They can be used as a standard `Guid` and support the same interfaces.

***

### Sequential Guid Example:

```go
fmt.Printf("%s\t       %s\t\t\t\t%s\t       %s\n",
	"gpg.String()", "hex(gpg)", "gss.String()", "hex(gss)")
for range 10 {
	gpg := guid.NewPG()
	gss := guid.NewSS()
	fmt.Println(&gpg, hex.EncodeToString(gpg.Guid[:]), &gss, hex.EncodeToString(gss.Guid[:]))
}

gpg := guid.NewPG()
gss := guid.NewSS()
fmt.Println(gpg.Timestamp()) // time.Time
fmt.Println(gss.Timestamp()) // time.Time
```
```
gpg.String()           hex(gpg)                         gss.String()           hex(gss)
GFEU88wgQvDlahOowSGTKA 185114f3cc2042f0e56a13a8c1219328 9SurLKL6ti2l0BhRFPPMKA f52bab2ca2fab62da5d0185114f3cc28
GFEU88wopdChlFba89-4yg 185114f3cc28a5d0a19456daf3dfb8ca yTRE6Rr1gISl0BhRFPPMKA c93444e91af58084a5d0185114f3cc28
GFEU88ww9fA01GntVDQ_4w 185114f3cc30f5f034d469ed54343fe3 8SaILyee6q718BhRFPPMMA f126882f279eeaaef5f0185114f3cc30
GFEU88ww9fASNFzZQJpv7Q 185114f3cc30f5f012345cd9409a6fed xZ3KYLzqJ0f18BhRFPPMMA c59dca60bcea2747f5f0185114f3cc30
GFEU88ww9fAHgWvjAmkQJw 185114f3cc30f5f007816be302691027 yEif2kTQBcD18BhRFPPMMA c8489fda44d005c0f5f0185114f3cc30
GFEU88ww9fD4_Vm3PG5Vuw 185114f3cc30f5f0f8fd59b73c6e55bb SRKgSiCc-gL18BhRFPPMMA 4912a04a209cfa02f5f0185114f3cc30
GFEU88ww9fDzO_One7T6BA 185114f3cc30f5f0f33bf3a77bb4fa04 rGr2czgQcmr18BhRFPPMMA ac6af6733810726af5f0185114f3cc30
GFEU88w5PqQAifEi5tqoWQ 185114f3cc393ea40089f122e6daa859 5YYbiI3p7P4-pBhRFPPMOQ e5861b888de9ecfe3ea4185114f3cc39
GFEU88w5PqSFkX4bmxSvMQ 185114f3cc393ea485917e1b9b14af31 PqUPeiyessU-pBhRFPPMOQ 3ea50f7a2c9eb2c53ea4185114f3cc39
GFEU88w5PqTsYX0kcZzL6Q 185114f3cc393ea4ec617d24719ccbe9 yFIlRwKZJNo-pBhRFPPMOQ c8522547029924da3ea4185114f3cc39
2025-07-11 03:32:47.3597457 +0000 UTC
2025-07-11 03:32:47.3597457 +0000 UTC
```

## Interoperability with `google/uuid` 🔥
* If you must keep using `google/uuid`, use `guid` to increase performance by **2~4x**:
```go
// do this before using google/uuid
uuid.SetRand(guid.Reader)
```
* Quick conversions between `guid` and `google/uuid` if you need `uuid` behavior:
```go
g := guid.New()
gpg := guid.NewPG()
gss := guid.NewSS()

var u uuid.UUID

u = uuid.UUID(g)
fmt.Println(u)

u = uuid.UUID(gpg.Guid)
fmt.Println(u)

u = uuid.UUID(gss.Guid)
fmt.Println(u)
```
```go
2dfc2275-71e1-776b-e6a3-5818c9b16976
18527f09-d2d9-e458-2611-7c8f416e2e8b
c4abc00d-bea3-7626-e458-18527f09d2d9
```
## FIPS Compliant
* **FIPS-140 compliant** (https://go.dev/doc/security/fips140)
	* set `GODEBUG=fips140=on` environment variable

## uuid Benchmarks with and without `guid.Reader`

| Benchmark Name | Time per Op | Bytes per Op  | Allocs per Op  |
|---|---|---|---|
| Benchmark_uuid_New_x10-8                                   | 3031 ns/op  | 160 B/op      | 10 allocs/op   |
| Benchmark_uuid_New_**guidRand**_x10-8 🔥                   | 862.0 ns/op | 160 B/op      | 10 allocs/op   |
| Benchmark_uuid_New_RandPool_x10-8                          | 747.6 ns/op | 0 B/op        | 0 allocs/op    |
| Benchmark_uuid_New_RandPool_**guidRand**_x10-8 🔥          | 516.8 ns/op | 0 B/op        | 0 allocs/op    |
| Benchmark_uuid_New_Parallel_x10-8                          | 1230 ns/op  | 160 B/op      | 10 allocs/op   |
| Benchmark_uuid_New_Parallel_**guidRand**_x10-8 🔥          | 510.0 ns/op | 160 B/op      | 10 allocs/op   |
| Benchmark_uuid_New_Parallel_RandPool_x10-8                 | 1430 ns/op  | 0 B/op        | 0 allocs/op    |
| Benchmark_uuid_New_Parallel_RandPool_**guidRand**_x10-8 🔥 | 1185 ns/op  | 0 B/op        | 0 allocs/op    |


## Guid Benchmarks [[raw](BENCHMARKS.md)]
```
go test -bench=.* -benchtime=4s
goos: windows
goarch: amd64
pkg: github.com/sdrapkin/guid
cpu: Intel(R) Core(TM) i7-10510U CPU @ 1.80GHz
```
| Benchmarks guid [10 calls] | Time/op | Bytes/op | Allocs/op |
|---|---|---|---|
| guid_New_x10-8                          |  203.4 ns/op  |   0 B/op |  0 allocs/op |
| guid_NewString_x10-8                    |  582.4 ns/op  | 240 B/op | 10 allocs/op |
| guid_String_x10-8                       |  388.9 ns/op  | 240 B/op | 10 allocs/op |
| guid_New_Parallel_x10-8 🔥               |  62.45 ns/op  |   0 B/op |  0 allocs/op |
| guid_NewString_Parallel_x10-8           |  374.2 ns/op  | 240 B/op | 10 allocs/op |

## Sequential Guid Benchmarks
| `guid.NewPG()` vs `uuid.NewV7()` [10 calls] | Time/op | |
|---|---|---|
| **guid.NewPG()_x10_Sequential** | **386.4 ns/op** |
| uuid.NewV7()_x10_Sequential | 887.9 ns/op | 2.3x slower ⏳
| **guid.NewPG()_x10_Parallel** | **144.3 ns/op** |
| uuid.NewV7()_x10_Parallel | 2575 ns/op | 18x slower ⏳


### Alternative library benchmarks:
| Benchmarks nanoid v1.35 [10 calls] | Time/op | Bytes/op | Allocs/op |
|---|---|---|---|
| `guid.NewString()` x10 Sequential       | **609.9 ns/op**   | 240 B/op | 10 allocs/op |
| `guid.NewString()` x10 Parallel (8 CPU) | **384.0 ns/op**   | 240 B/op | 10 allocs/op |
| `nanoid.New()` x10 Sequential           | 2257 ns/op        | 240 B/op | 10 allocs/op |
| `nanoid.New()` x10 Parallel (8 CPU)     | 1337 ns/op        | 240 B/op | 10 allocs/op |

| Benchmarks uuid [10 calls] | Time/op | Bytes/op | Allocs/op |
|---|---|---|---|
| uuid_New_x10-8                          |  2216 ns/op   | 160 B/op | 10 allocs/op |
| uuid_New_RandPool_x10-8                 |  528.2 ns/op  |   0 B/op |  0 allocs/op |
| uuid_New_Parallel_x10-8                 |  1064 ns/op   | 160 B/op | 10 allocs/op |
| uuid_New_RandPool_Parallel_x10-8        |  1301 ns/op   |   0 B/op |  0 allocs/op |

| Benchmarks [20 guid encodings] | Time/op | Bytes/op | Allocs/op |
|---|---|---|---|
| g.String-8                |  1025 ns/op   | 480 B/op | 20 allocs/op |
| base64.RawURLEncoding.EncodeToString-8  |  1867 ns/op   | 960 B/op | 40 allocs/op |
| g.EncodeBase64URL-8                  |  392.0 ns/op  |   0 B/op |  0 allocs/op |
| base64.RawURLEncoding.Encode-8          |  463.4 ns/op  |   0 B/op |  0 allocs/op |

## Documentation
 [![Go Reference](https://pkg.go.dev/badge/github.com/sdrapkin/guid.svg)](https://pkg.go.dev/github.com/sdrapkin/guid)

Full `go doc` style documentation: https://pkg.go.dev/github.com/sdrapkin/guid

## Requirements
- Go 1.24+

## Installation
### Using `go get`

To install the `guid` package, run the following command:

```sh
go get -u github.com/sdrapkin/guid
```

To use the `guid` package in your Go project, import it as follows:

```go
import "github.com/sdrapkin/guid"
```
## JSON Support

`Guid` supports JSON marshalling and unmarshalling for both value and pointer types:

- Value fields serialize as 22-character Base64Url strings.
- Pointer fields serialize as strings or `null` (for nil pointers).
- Zero-value Guids (`guid.Nil`) are handled correctly.

### Example: JSON Marshalling
```go
type User struct {
	ID        guid.Guid  `json:"id"`
	ManagerID *guid.Guid `json:"mid"`
}

u, u2 := User{ID: guid.New()}, User{}
data, _ := json.Marshal(u)
fmt.Println(string(data)) // {"id":"tI0EMdDXpOcvvGLktob4Ug","mid":null}

_ = json.Unmarshal(data, &u2)
fmt.Println(u2.ID == u.ID) // true
```
