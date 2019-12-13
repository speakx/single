package app

import (
	"environment/cfgargs"
	"environment/dump"
	"environment/logger"
	"fmt"
	"mmapcache/cache"
	"singledb/proto/pbsingledb"

	uuid "github.com/satori/go.uuid"
)

const (
	mmapStatusUseing  = 0
	mmapStatusSaving  = 1
	mmapStatusCollect = 2
)

var singleMsgCacheMap map[string]*cache.MMapCache
var mmapCacheSaveCh chan *cache.MMapCache

func initSingleMsgCacheMap(cfg *cfgargs.SrvConfig) {
	singleMsgCacheMap = make(map[string]*cache.MMapCache)
	mmapCacheSaveCh = make(chan *cache.MMapCache, 0x1000)
	// go func() {
	// 	savingMsgCacheloop()
	// }()
}

// GetSingleMsgCache 获取单人消息的mmap缓存，通过缓存批量写入后端DB
func GetSingleMsgCache(fromUID, toUID uint64) (*cache.MMapCache, error) {
	key := fmt.Sprintf("%d_%d", fromUID, toUID)
	mmapCache, _ := singleMsgCacheMap[key]
	if nil == mmapCache {
		mmapCache = cache.DefPoolMMapCache.Alloc()
		mmapCache.SetStatus(mmapStatusUseing)

		singleMsgCacheMap[key] = mmapCache
	}

	// 确保拿到的cache足够写一条记录
	if mmapCache.GetFreeContentLen() < _app.SrvCfg.Cache.DataSize {
		oldMMapCache := mmapCache
		oldMMapCache.SetStatus(mmapStatusSaving)
		if err := saveMsgCache(oldMMapCache); nil != err {
			return nil, err
		}

		mmapCache = cache.DefPoolMMapCache.Alloc()
		mmapCache.SetStatus(mmapStatusUseing)

		singleMsgCacheMap[key] = mmapCache
	}

	return mmapCache, nil
}

func saveMsgCache(mmapCache *cache.MMapCache) error {
	logger.Debug("mmap cache save len:", len(mmapCache.GetMMapDatas()))
	dump.NetEventSendIncr(0)
	key := mmapCache.GetMMapDatas()[len(mmapCache.GetMMapDatas())-1].GetKey()
	req := &pbsingledb.SingleMsgChunk{
		Transid: uuid.NewV1().String(),
		Key:     key,
		Data:    mmapCache.GetWrittenData(),
	}
	srv := _app.SingleDB
	_, err := srv.Save(srv.GetCtx(), req)
	if nil != err {
		logger.Error("save msg cache, singledb err:", err)
		dump.NetEventSendDecr(0)
		return err
	}
	dump.NetEventSendDecr(0)

	mmapCache.SetStatus(mmapStatusCollect)
	cache.DefPoolMMapCache.Collect(mmapCache)
	return nil
}
