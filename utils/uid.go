package utils

import "sync/atomic"

var Uid uint32

func GetUid() uint32 {
	return atomic.AddUint32(&Uid, 1)
}
