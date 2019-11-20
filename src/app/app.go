package app

import (
	"environment/cfgargs"
	"environment/dump"
	"environment/logger"
	"mmapcache/cache"
	"os"
	"single/client"
	"single/database"
	"sync"
)

var once sync.Once
var _app *App

// GetApp 获取当前服务的App实例
func GetApp() *App {
	once.Do(func() {
		_app = &App{}
	})
	return _app
}

// App 当前服务的App实例，用来存储一些运行时对象
type App struct {
	SrvCfg   *cfgargs.SrvConfig
	DB       *database.DB
	SingleDB client.SingleGrpcClient
}

// InitApp 加载配置、初始化日志、构建mmap缓存池
func (a *App) InitApp(srvCfg *cfgargs.SrvConfig) {
	a.SrvCfg = srvCfg

	// 初始化日志
	logger.Info("start init log")
	logger.InitLogger(srvCfg.Log.Path, srvCfg.Log.Console, srvCfg.Log.Level)
	logger.Info("end init log")

	// 初始化dump包，用来做服务端健康度检查&汇报
	logger.Info("start init dump, addr:", srvCfg.Dump.Addr)
	dump.InitDump(true, srvCfg.Dump.Interval, srvCfg.Dump.Addr, nil)
	logger.Info("end init dump")

	// 初始化mmap缓存
	logger.Info("start init mmap cache, dir:", srvCfg.Cache.Path)
	cache.InitMMapCachePool(
		srvCfg.Cache.Path, srvCfg.Cache.MMapSize,
		srvCfg.Cache.DataSize, srvCfg.Cache.PreAlloc,
		a.errorMMapCache, a.reloadMMapCache)
	logger.Info("end init mmap cache")

	// 初始化后端服务连接
	logger.Info("start init client")
	if err := a.initClientSrv(); nil != err {
		logger.Error("init client err:", err)
		os.Exit(0)
	}
	logger.Info("end init client")

	// 初始化DB层
	logger.Info("start init db")
	if err := a.initDB(); nil != err {
		logger.Error("init db err:", err)
		os.Exit(0)
	}
	logger.Info("end init db")

	// 其他初始化
	initSingleMsgCacheMap(a.SrvCfg)
}

func (a *App) errorMMapCache(err error) {
	logger.Error("mmapcache err:", err)
}

func (a *App) reloadMMapCache(mmapCaches []*cache.MMapCache) {
}

func (a *App) initClientSrv() error {
	a.SingleDB.Connect(":10101")
	return nil
}

func (a *App) initDB() error {
	db, err := database.NewDB(a.SrvCfg)
	a.DB = db
	return err
}
