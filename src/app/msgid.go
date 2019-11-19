package app

import (
	"environment/logger"
	"sync"
	"time"
)

var incrID uint64
var timeStamp uint64
var srvidLock sync.Mutex

// MakeSrvMsgID srvid
func MakeSrvMsgID(fromUID, toUID uint64) uint64 {
	timeStamp := uint64(time.Now().Unix()) << 32 & 0xFFFFFFFF00000000
	srvID := (uint64(_app.SrvCfg.Info.ID) << 20) & 0xFFF00000
	return timeStamp | srvID | genIncrID(timeStamp)
}

func genIncrID(ts uint64) uint64 {
	srvidLock.Lock()
	if timeStamp != ts {
		incrID = 0
		timeStamp = ts
		srvidLock.Unlock()
		return incrID
	}

	incrID++
	if incrID > 0xFFFFF {
		logger.Error("db.gen.incrid overflow ts:", timeStamp, " incrid:", incrID)
	}
	srvidLock.Unlock()
	return incrID
}
