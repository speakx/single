package app

import (
	"environment/cfgargs"
	"fmt"
	"mmapcache/cache"
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
	go func() {
		savingMsgCacheloop()
	}()
}

// GetSingleMsgCache 获取单人消息的mmap缓存，通过缓存批量写入后端DB
func GetSingleMsgCache(fromUID, toUID uint64) *cache.MMapCache {
	key := fmt.Sprint("%v_%v", fromUID, toUID)
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
		mmapCacheSaveCh <- oldMMapCache

		mmapCache = cache.DefPoolMMapCache.Alloc()
		mmapCache.SetStatus(mmapStatusUseing)
		singleMsgCacheMap[key] = mmapCache
	}

	return mmapCache
}

func savingMsgCacheloop() {
	for {
		mmapCache, _ := <-mmapCacheSaveCh
		mmapCache.SetStatus(mmapStatusCollect)

		cache.DefPoolMMapCache.Collect(mmapCache)
	}
}
