package guid

import (
	cryptoRand "crypto/rand"
	"io"
)

// Size of a Guid in bytes.
const GuidByteSize = 16

// 16-byte (128-bit) cryptographically random value.
type Guid [GuidByteSize]byte

// Empty Guid
var Nil Guid

var cryptoRandReader io.Reader // initialized in init()

func New() (guid Guid) {
	_, err := cryptoRandReader.Read(guid[:])
	if err != nil {
		panic(err) // cryptoRand.Reader.Read should never fail; if it does, there is no safe recourse
	}
	return
} //New()

func init() {
	cryptoRandReader = cryptoRand.Reader
} //init()
