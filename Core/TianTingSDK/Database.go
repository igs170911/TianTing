package TianTingSDK

import (
	"TianTing/Core/Database"
	"TianTing/Logger"
	"TianTing/Settings"
)

// 初始化 Redis
func (server *TianTingServer) ConnectCacheDbService(config *Settings.CacheDbConf) {
	client, err := Database.ConnectWithCacheDB(config)
	if err != nil {
		Logger.SysLog.Errorf("[CacheDb] Try To Connect Cache Database Failed -> (%s)", err)
	}
	server.TianTingCacheDb = client
}

func (server *TianTingServer) GetCacheDb() *Database.CacheDB {
	if server.TianTingCacheDb == nil {
		return nil
	}
	return server.TianTingCacheDb
}
