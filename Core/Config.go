package Core

import (
	"TianTing/Settings"
	"github.com/koding/multiconfig"
	"sync"
)

type (
	// Server holds supported types by the multiconfig package
	Server struct {
		App   *Settings.AppConf
		Redis *Settings.CacheDbConf
		Mongo *Settings.DocDbConf
	}
)

// 全域唯一取得config
var Cfg Server
var once sync.Once

//TODO 修正 1. 要有default 值 2. 有需要的才載入
func GetConfig() Server {
	once.Do(func() {
		m := multiconfig.NewWithPath("config.toml") // supports TOML and JSON
		// Get an empty struct for your configuration
		serverConf := new(Server)
		// Populated the serverConf struct
		m.MustLoad(serverConf) // Check for error
		//fmt.Println("After Loading: ")
		//fmt.Printf("%+v\n", serverConf)
		Cfg = *serverConf
	})
	return Cfg
}
