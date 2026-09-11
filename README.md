# guid [![codecov](https://codecov.io/github/sdrapkin/guid/branch/master/graph/badge.svg?token=ARQFUQD5VP)](https://codecov.io/github/sdrapkin/guid) [![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go#uuid) 
## Fast cryptographically secure Guid generator for Go.<br>By [Stan Drapkin](https://github.com/sdrapkin/).

`Guid` is a 16-byte struct filled with 128 cryptographically strong bits. Its bytes are available through `g.UUID[:]`.

[Go playground](https://go.dev/play/p/H8xdCVR4EAw)
```go
package main

import (
	"fmt"

	"github.com/sdrapkin/guid"
)

func main() {
	fmt.Printf("%-36s %-32s %s\n", "UUID string:", "Hex:", ".String():")
	for range 4 {
		g := guid.New()
		fmt.Printf("%s %x %s\n", g.UUID, g.UUID[:], g)
	}
}
```

```
UUID string:                         Hex:                             .String():
5c61d893-95e1-b7b3-95ba-a215a1bafe84 5c61d89395e1b7b395baa215a1bafe84 XGHYk5Xht7OVuqIVobr-hA
da304a44-29d0-9779-41bc-5d011878fa6c da304a4429d0977941bc5d011878fa6c 2jBKRCnQl3lBvF0BGHj6bA
6874ff53-f7e8-8a4d-5c35-d47ec7de4509 6874ff53f7e88a4d5c35d47ec7de4509 aHT_U_foik1cNdR-x95FCQ
4bc70e6f-78f0-25a0-bdf4-76a38665cf9d 4bc70e6f78f025a0bdf476a38665cf9d S8cOb3jwJaC99HajhmXPnQ
```

## Why `guid`? 🔥

`guid` is a high-performance, cryptographically secure UUID/GUID (Globally Unique Identifier) generator for Go. It is built for speed without compromising on security.

Beyond raw speed, `guid` offers:

* **Cryptographically Strong**: Generates 128 cryptographically secure bits for robust, unique identifiers.
* **Optimized for Databases**: Includes special `GuidPG` and `GuidSS` types that generate sequential Guids, dramatically improving `INSERT` performance and preventing index fragmentation in **PostgreSQL** and **SQL Server** databases.
* **Seamless Interoperability**: Easily integrate with existing `uuid.UUID` codebases, and use `guid.Reader` as an `io.Reader` for workloads that benefit from its fast random-byte implementation.
* **FIPS 140 Compatibility**: Can be used in Go environments configured for FIPS 140 mode.
* **Zero Allocations for Core Operations**: `guid.New()` generates new Guids with no allocations; string conversion via `String()` allocates the returned string as expected.

## Performance : `Guid` is ~8x faster than `uuid.UUID` 🔥

Performance depends on the Go version, platform, CPU, workload, and benchmark settings. The benchmark source contains the command and environment for the recorded results; run `go test -run=^$ -bench=. -benchmem -benchtime=4s` to measure locally.

This package uses pooled random bytes for small reads and falls back to `crypto/rand.Read` for requests larger than 512 bytes. All generation and reader operations remain cryptographically backed by `crypto/rand`.

## API Overview
Functions and value-receiver methods are safe for concurrent use when called with independent values and buffers. Do not call methods that mutate the same pointer receiver concurrently.
| Functions | Description |
|---|---|
| `guid.New()` `Guid`           | Generate a new cryptographically secure Guid |
| `guid.NewString()` `string`   | Generate a new Guid as a Base64Url string |
| `guid.NewPG()` `GuidPG`       | Generate a new PostgreSQL sequential Guid |
| `guid.NewSS()` `GuidSS`       | Generate a new SQL Server sequential Guid |
| `guid.Parse(s string)` `(Guid, error)` | Parse a Base64Url string into a Guid |
| `guid.MustParse(s string)` `Guid` | Parse a Base64Url string or panic on invalid input |
| `guid.ParseBytes(src []byte)` `(Guid, error)` | Parse Base64Url bytes to a Guid |
| `guid.FromBytes(src []byte)` `(Guid, error)` | Create a Guid from a slice |
| `guid.DecodeBase64URL(dst []byte, src []byte)` `(ok bool)` | Decode the first 22 bytes of a Base64Url input into `dst` |
| `guid.Read(p []byte)` `(int, error)` | Fill a byte slice with secure random bytes |
| `guid.Reader` 🔥 implements `io.Reader` | Faster alternative to `crypto/rand` |
| `guid.Nil()` `Guid` | The zero-value Guid |
| `guid.Max()` `Guid` | The maximum Guid |

| `Guid` methods | Description |
|---|---|
| `.String()` `string` | Encodes the Guid into a 22-char Base64Url string (`fmt.Stringer`) |
| `.EncodeBase64URL(dst []byte)` `error` | Encodes into a destination slice |
| `.Compare(other Guid)` `int` | Lexicographic comparison using big-endian byte order |
| `.MarshalBinary()` | Implements `encoding.BinaryMarshaler` |
| `.UnmarshalBinary()` | Implements `encoding.BinaryUnmarshaler` |
| `.MarshalText()` | Implements `encoding.TextMarshaler` |
| `.UnmarshalText()` | Implements `encoding.TextUnmarshaler` |
| `.AppendText(b []byte)` `([]byte, error)` | Appends the Base64Url encoding to `b` (`encoding.TextAppender`) |
| `.MarshalJSON()` `([]byte, error)` | Encodes the Guid as a JSON string |
| `.UnmarshalJSON(data []byte)` `error` | Decodes a JSON string or `null` into the Guid |

| `GuidPG`, `GuidSS` methods | Description |
|---|---|
| `.Timestamp()` `time.Time` | Extracts the UTC timestamp |
| `GuidPG.Compare(other GuidPG)` `int` | Lexicographic comparison using big-endian byte order |
| `GuidSS.Compare(other GuidSS)` `int` | Comparison using SQL Server's Guid byte ordering rules |
| `GuidSS.LoadFromSQLServerBytes(src []byte)` `error` | Loads 16 bytes into `GuidSS` in SQL Server byte order |

## Sequential Guids 🔥
`guid` includes two special types, `GuidPG` and `GuidSS`, with time-ordered layouts intended for database keys. They can improve locality for workloads whose database ordering matches the corresponding layout, but actual index and `INSERT` behavior depends on the database, schema, and workload. Ordering is only guaranteed at the precision of the `time.Now()` timestamp used to create each value.

* **`guid.NewPG()`**: Generates a `GuidPG`, which is sortable in **PostgreSQL**.
 	- It is structured as `[8-byte timestamp][8 random bytes]`.
* **`guid.NewSS()`**: Generates a `GuidSS`, which is sortable in **SQL Server**.
	- It is structured as `[8 random bytes][8-byte SQL Server-ordered timestamp]`.
* `.Timestamp()` on `GuidPG`/`GuidSS` returns Guid creation time as UTC `time.Time`.

Both `GuidPG` and `GuidSS` contain an embedded `Guid`, so they expose its value methods and interfaces in addition to their timestamp and comparison methods.

***

### Sequential Guid Example:

```go
package main

import (
	"encoding/hex"
	"fmt"

	"github.com/sdrapkin/guid"
)

func main() {
	fmt.Printf("%s\t       %s\t\t\t\t%s\t       %s\n",
	"gpg.String()", "hex(gpg)", "gss.String()", "hex(gss)")
	for range 10 {
		gpg := guid.NewPG()
		gss := guid.NewSS()
		fmt.Println(&gpg, hex.EncodeToString(gpg.UUID[:]), &gss, hex.EncodeToString(gss.UUID[:]))
	}

	gpg := guid.NewPG()
	gss := guid.NewSS()
	fmt.Println(gpg.Timestamp()) // time.Time
	fmt.Println(gss.Timestamp()) // time.Time
}
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

## Interoperability with Go's `uuid.UUID`

Go 1.27's built-in [`uuid.UUID`](https://pkg.go.dev/uuid) is a `[16]byte` type. `Guid.UUID`, and the promoted `UUID` field on `GuidPG` and `GuidSS`, can be converted by value:
```go
package main

import (
	"fmt"

	"github.com/sdrapkin/guid"
	"uuid"
)

func main() {
	g := guid.New()
	gpg := guid.NewPG()
	gss := guid.NewSS()

	fmt.Println(g.UUID) // exposes as UUID
	fmt.Println(gpg.UUID) // exposes as UUID
	fmt.Println(gss.UUID) // exposes as UUID

	g.UUID[0], g.UUID[1] = 0xAB, 0xCD
	fmt.Println(g)
}
```
```go
05166521-a124-9d0c-cb11-7f0cbf3a030c
1852e32a-5aac-bb9c-bffc-b330606813af
7e8badae-57f8-c88d-bb9c-1852e32a5aac
abcd6521-a124-9d0c-cb11-7f0cbf3a030c
```

## FIPS Ready
* **FIPS-140 compatible** (https://go.dev/doc/security/fips140)
	* set `GODEBUG=fips140=on` environment variable
	* https://go.dev/blog/fips140

## Recorded Benchmarks

These results were recorded on Windows amd64 with Go 1.27. They are reference measurements, not performance guarantees. Each `x10` benchmark performs ten calls per benchmark iteration.

| Benchmark | Time per Op | Bytes per Op | Allocs per Op |
|---|---|---|---|
| `guid.New()` x10 | 255.6 ns/op | 0 B/op | 0 allocs/op |
| `guid.NewString()` x10 | 686.3 ns/op | 240 B/op | 10 allocs/op |
| `Guid.String()` x10 | 175.3 ns/op | 0 B/op | 0 allocs/op |
| `guid.New()` x10 parallel | 91.27 ns/op | 0 B/op | 0 allocs/op |
| `guid.NewString()` x10 parallel | 498.2 ns/op | 240 B/op | 10 allocs/op |
| `guid.NewPG()` x10 | 366.1 ns/op | 0 B/op | 0 allocs/op |
| `guid.NewSS()` x10 | 380.4 ns/op | 0 B/op | 0 allocs/op |
| `uuid.New()` x10 | 1712 ns/op | 0 B/op | 0 allocs/op |
| `uuid.NewV7()` x10 | 1819 ns/op | 0 B/op | 0 allocs/op |
| `uuid.New()` x10 parallel | 748.7 ns/op | 0 B/op | 0 allocs/op |
| `uuid.NewV7()` x10 parallel | 2818 ns/op | 0 B/op | 0 allocs/op |

For historical reader and alternative-library measurements, see [BENCHMARKS.md](BENCHMARKS.md). Those results were collected from earlier code and environments and should not be compared directly with the current table.

## Documentation
 [![Go Reference](https://pkg.go.dev/badge/github.com/sdrapkin/guid.svg)](https://pkg.go.dev/github.com/sdrapkin/guid)

Full `go doc` style documentation: https://pkg.go.dev/github.com/sdrapkin/guid

## Requirements
- Go 1.27+

## Installation
### Using `go get`

To install the `guid` package, run the following command:

```sh
go get github.com/sdrapkin/guid
```

To use the `guid` package in your Go project, import it as follows:

```go
import "github.com/sdrapkin/guid"
```
## JSON Support

`Guid` supports JSON marshalling and unmarshalling for both value and pointer types:

- Value fields serialize as 22-character Base64Url strings.
- Pointer fields serialize as strings or `null` (for nil pointers).
- Zero-value Guids (`guid.Nil()`) are handled correctly.

### Example: JSON Marshalling
```go
package main

import (
	"encoding/json"
	"fmt"

	"github.com/sdrapkin/guid"
)

type User struct {
	ID        guid.Guid  `json:"id"`
	ManagerID *guid.Guid `json:"mid"`
}

func main() {
	u, u2 := User{ID: guid.New()}, User{}
	data, err := json.Marshal(u)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data)) // {"id":"tI0EMdDXpOcvvGLktob4Ug","mid":null}

	if err := json.Unmarshal(data, &u2); err != nil {
		panic(err)
	}
	fmt.Println(u2.ID == u.ID) // true
}
```
