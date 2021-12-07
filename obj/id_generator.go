package obj

import (
	"math"
	"time"
)

var (
	IdGen = &IdGenerator{}
)

type IdGenerator struct {
	tm  uint32
	val uint32
}

// GenId 随便实现个后续调整
func (i *IdGenerator) GenId() uint64 {
	now := uint32(time.Now().Unix())
	if i.tm == uint32(time.Now().Unix()) {
		if i.val >= math.MaxInt32 {
			i.tm += 1
			i.val = 1
		} else {
			i.val++
		}
	} else {
		i.tm = now
		i.val = 1
	}
	return uint64(i.tm)<<32 | uint64(i.val)
}
