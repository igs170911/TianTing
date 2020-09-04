package Core

import (
	"TianTing/Logger"
	"github.com/kataras/iris/v12"
	"time"
)

func (core *Core)RegIrisRouter(path string, GroupRoute func(iris.Party)){
	initTime := time.Now()
	core.CoreEngine.PartyFunc(path, GroupRoute)
	diffTime := time.Since(initTime).Microseconds()
	Logger.SysLog.Infof("[Engine] Register HTTP Feature Groups(%s) in %dµs", path, diffTime)
}