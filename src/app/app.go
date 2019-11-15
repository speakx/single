package app

import (
	"environment/cfgargs"
	"environment/dump"
	"environment/logger"
	"fmt"
	"mmapcache/cache"
	pb "svrdemo/proto"
	"sync"

	"github.com/golang/protobuf/proto"
)

var once sync.Once
var app *App

// GetApp 获取当前服务的App实例
func GetApp() *App {
	once.Do(func() {
		app = &App{}
	})
	return app
}

// App 当前服务的App实例，用来存储一些运行时对象
type App struct {
	SrvCfg  *cfgargs.SrvConfig
	SrvDemo SimpleGrpcClient
}

// InitApp 加载配置、初始化日志、构建mmap缓存池
func (a *App) InitApp(srvCfg *cfgargs.SrvConfig) {
	a.SrvCfg = srvCfg

	// 初始化日志
	logger.InitLogger(srvCfg.Log.Path, srvCfg.Log.Console, srvCfg.Log.Level)

	// 初始化mmap缓存池
	logger.Info("start init mmap cache pool")
	cache.InitMMapCachePool(
		srvCfg.Cache.Path, srvCfg.Cache.MMapSize,
		srvCfg.Cache.DataSize, srvCfg.Cache.PreAlloc,
		a.errorMMapCache, a.reloadMMapCache)
	logger.Info("end init mmap cache pool")

	// 初始化所有的后端服务连接
	logger.Info("start init client")
	a.initClientSrv()
	logger.Info("end init client")

	// 初始化dump包，用来做服务端健康度检查&汇报
	dump.InitDump(true, srvCfg.Dump.Interval, srvCfg.Dump.Addr,
		func(recv int64, send int64, recvHandleRate int64, sendHandleRate int64) {
			fmt.Println(sendHandleRate)
		})
}

func (a *App) errorMMapCache(err error) {
	logger.Error("mmapcache err:", err)
}

func (a *App) reloadMMapCache(mmapCaches []*cache.MMapCache) {
	logger.Info("reload mmapcache.count:", len(mmapCaches))
	for idx, mmapCache := range mmapCaches {
		logger.Info("reload mmapcache.idx:", idx, " data.count:%v", len(mmapCache.GetMMapDatas()))
		for _, mmapData := range mmapCache.GetMMapDatas() {
			var req pb.SimpleHello
			proto.Unmarshal(mmapData.GetData(), &req)
		}
		cache.DefPoolMMapCache.Collect(mmapCache)
	}
}

func (a *App) initClientSrv() {
	// a.SrvDemo.Connect(a.SrvCfg.Addr)
	a.SrvDemo.Connect("10.211.55.27:10000")
}
