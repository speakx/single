package main

import (
	"environment/cfgargs"
	"fmt"
	"single/app"
	"single/server"
)

var (
	BuildVersion = ""
)

func main() {
	srvCfg, err := cfgargs.InitSrvConfig(BuildVersion, func() {
		// user flag binding code
	})
	if nil != err {
		fmt.Println(err)
		return
	}
	app.GetApp().InitApp(srvCfg)

	srv := server.NewServer()
	srv.Run(srvCfg.Info.Addr)
}
