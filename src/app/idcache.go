package app

import (
	"environment/cfgargs"
	"environment/dump"
	"environment/logger"
	"fmt"
	"idgenerator/proto/pbidgenerator"
	"os"

	uuid "github.com/satori/go.uuid"
)

var singleIDCacheMap map[string]*SingleIDCache
var singleIDCacheNum int

func initSingleIDCache(cfg *cfgargs.SrvConfig) {
	singleIDCacheMap = make(map[string]*SingleIDCache)
	singleIDCacheNum = cfg.Cache.MMapSize/cfg.Cache.DataSize - 10
	logger.Info("signle id cahce num:", singleIDCacheNum)
}

// SingleIDCache SrvID & OrderID的缓存
type SingleIDCache struct {
	srvIDs   []uint64
	orderIDs []uint64
	cursor   int
}

// PopID 获取SrvID & OrderID
func (s *SingleIDCache) PopID() (uint64, uint64) {
	s.cursor++
	return s.srvIDs[s.cursor-1], s.orderIDs[s.cursor-1]
}

// GetSingleIDCache 获取会话对应的可使用ID缓存
func GetSingleIDCache(fromUID, toUID uint64) *SingleIDCache {
	key := fmt.Sprint("%v_%v", fromUID, toUID)
	sidCache, _ := singleIDCacheMap[key]
	if nil == sidCache {
		sidCache = &SingleIDCache{
			srvIDs:   make([]uint64, singleIDCacheNum),
			orderIDs: make([]uint64, singleIDCacheNum),
			cursor:   singleIDCacheNum,
		}
	}

	if sidCache.cursor >= len(sidCache.srvIDs) {
		req := &pbidgenerator.SingleMessageId{
			Transid: uuid.NewV1().String(),
			FromUid: fromUID,
			ToUid:   toUID,
			Num:     uint32(singleIDCacheNum),
		}

		dump.NetEventSendIncr(0)
		defer dump.NetEventSendDecr(0)
		replay, err := _app.IDGenClient.GenSingleMessageId(
			_app.IDGenClient.GetCtx(), req)
		if nil != err {
			logger.Error("single id cache, id gen err:", err)
			os.Exit(0)
		}

		copy(sidCache.srvIDs, replay.SrvIds)
		copy(sidCache.orderIDs, replay.OrderIds)
		sidCache.cursor = 0
	}

	return sidCache
}
