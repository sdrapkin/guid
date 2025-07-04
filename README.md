# guid
## Fast cryptographically safe Guid generator for Go. By [Stan Drapkin](https://github.com/sdrapkin/).

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
---
## Guid is ~10x faster than `github.com/google/uuid`

* `guid.New()` is  6~10 ns 
* `guid.NewString()` is 40~60 ns
* `String()` on existing guid is ~40 ns
* multi-goroutine calls do not increase per-call latency
* if your library is faster - please let me know!

## Benchmarks
```
go test -bench=.* -benchtime=4s
goos: windows
goarch: amd64
pkg: github.com/sdrapkin/guid
cpu: Intel(R) Core(TM) i7-10510U CPU @ 1.80GHz
```
| Benchmarks guid [10 calls] | Time/op | Bytes/op | Allocs/op |
|---|---|---|---|
| guid_New_x10-8                          |  203.4 ns/op  |   0 B/op |  0 allocs/op | |
| guid_NewString_x10-8                    |  582.4 ns/op  | 240 B/op | 10 allocs/op | |
| guid_String_x10-8                       |  388.9 ns/op  | 240 B/op | 10 allocs/op | |
| guid_New_Parallel_x10-8                 |  62.45 ns/op  |   0 B/op |  0 allocs/op | |
| guid_NewString_Parallel_x10-8           |  374.2 ns/op  | 240 B/op | 10 allocs/op | |

### Alternative library benchmarks:
| Benchmarks nanoid [10 calls] | Time/op | Bytes/op | Allocs/op |
|---|---|---|---|
| nanoid_New_x10-8                        | 2493 ns/op    | 240 B/op | 10 allocs/op | |
| nanoid_New_Parallel_x10-8               | 1282 ns/op    | 240 B/op | 10 allocs/op | |

| Benchmarks uuid [10 calls] | Time/op | Bytes/op | Allocs/op |
|---|---|---|---|
| uuid_New_x10-8                          |  2216 ns/op   | 160 B/op | 10 allocs/op | |
| uuid_New_RandPool_x10-8                 |  528.2 ns/op  |   0 B/op |  0 allocs/op | |
| uuid_New_Parallel_x10-8                 |  1064 ns/op   | 160 B/op | 10 allocs/op | |
| uuid_New_RandPool_Parallel_x10-8        |  1301 ns/op   |   0 B/op |  0 allocs/op | |

| Benchmarks [20 guid encodings] | Time/op | Bytes/op | Allocs/op |
|---|---|---|---|
| guid_ToBase64UrlString-8                |  1025 ns/op   | 480 B/op | 20 allocs/op | |
| base64_RawURLEncoding_EncodeToString-8  |  1867 ns/op   | 960 B/op | 40 allocs/op | |
| guid_EncodeBase64URL-8                  |  392.0 ns/op  |   0 B/op |  0 allocs/op | |
| base64_RawURLEncoding_Encode-8          |  463.4 ns/op  |   0 B/op |  0 allocs/op | |

## Documentation
 [![Go Reference](https://pkg.go.dev/badge/github.com/sdrapkin/guid.svg)](https://pkg.go.dev/github.com/sdrapkin/guid)

Full `go doc` style documentation: https://pkg.go.dev/github.com/sdrapkin/guid

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
