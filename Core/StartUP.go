package Core

import (
	"TianTing/Settings"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"github.com/kataras/iris/v12"
	"time"
)

type ICore interface{}

// 啟動的核心
// TODO 核心 Engine 用注入的
type Core struct{
	Config *Config
	CoreEngine *iris.Application
	StartTime time.Time
}

var _ ICore = &Core{}

func New() *Core {
	// 初始化 Config
	startTime := time.Now()
	// TODO 加入 Log Info("[Engine] Starting")
	core := &Core{
		Config: &Config{
			App: &Settings.AppConf{},
			custom:  make(map[string]interface{}),
		},
	}
	core.StartTime = startTime
	core.Config.Raw, _ = godotenv.Read()
	core.Config.SystemExternalEnv("app", core.Config.App)
	fmt.Println(core.Config.App.Codename)
	return core
}

