package id

import (
	"math/bits"
	"os"
	"time"
)

const BitLength = uint64(64)

var (
	sequence       = uint64(0)
	lastGenerateAt = time.Now()
)

// Generate はidを発番する
func Generate() uint64 {
	var (
		now       = time.Now()
		unixTime  = uint64(now.Unix())
		processID = uint64(os.Getpid())
	)

	sequence++

	if lastGenerateAt != now {
		lastGenerateAt = now
		sequence = 0
	}

	res := unixTime + processID
	sequenceBit := BitLength - uint64(bits.Len64(res))
	res = res<<sequenceBit + sequence
	return res
}
