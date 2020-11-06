package main

import (
	"fmt"
	"os"

	"github.com/TarsCloud/TarsGo/tars"
	"github.com/TarsCloud/TarsGo/tars/util/rogger"

	"github.com/TarsDemo/Tars-MiniProgramm-Service-UserInfoServer/tars-protocol/LifeService"
)

// comm 定义communicator
var comm *tars.Communicator

//SLOG 日志打印
var SLOG = rogger.GetLogger("ServerLog")

func main() {
	comm = tars.NewCommunicator() //初始化communicator
	// Get server config
	cfg := tars.GetServerConfig()

	// New servant imp
	imp := new(UserInfoServiceImp)
	err := imp.Init()
	if err != nil {
		fmt.Printf("UserInfoServiceImp init fail, err:(%s)\n", err)
		os.Exit(-1)
	}
	// New servant
	app := new(LifeService.UserInfoService)
	// Register Servant
	app.AddServantWithContext(imp, cfg.App+"."+cfg.Server+".UserInfoServiceObj")

	// Run application
	tars.Run()
}
