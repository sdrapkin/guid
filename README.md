# guid
### Fast cryptographically safe Guid generator for Go.

[Go playground](https://go.dev/play/p/H98wqbCQH_m)
```go
package main

import (
	"fmt"

	"github.com/sdrapkin/guid"
)

func main() {
	for range 100 {
		fmt.Printf("%x\n", guid.New())
	}
}
```

```
79c9779af20dcd21fbe60f3b336ed08c
da2026d38edca4371a476efd41333d23
88c3033b002b0e73321509ef26de607f
a84e961ff7f09f5210ea04585f152e73
...
```