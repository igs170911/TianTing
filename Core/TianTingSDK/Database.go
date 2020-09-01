package TianTingSDK

import (
	"TianTing/Core/Database"
	"TianTing/Logger"
	"TianTing/Settings"
)

// 初始化 Catch 資料庫
func (server *TianTingServer) ConnectCacheDbService(config *Settings.CacheDbConf) {
	client, err := Database.ConnectWithCacheDB(config)
	if err != nil {
		Logger.SysLog.Errorf("[CacheDb] Try To Connect Cache Database Failed -> (%s)", err)
	}
	server.CacheDb = client
}

func (server *TianTingServer) GetCacheDb() *Database.CacheDB {
	if server.CacheDb == nil {
		return nil
	}
	return server.CacheDb
}

// 初始化非關聯性資料庫 Doc System
func (server *TianTingServer) ConnectDocDbService(config *Settings.DocDbConf) {
	client, err := Database.ConnectWithDocDB(config)
	if err != nil {
		Logger.SysLog.Errorf("[DocumentDB] Try To Connect Document Database Failed -> (%s)", err)
	}
	server.DocDb = client
}

func (server *TianTingServer) GetDocDb() *Database.DocDB {
	if server.DocDb == nil {
		return nil
	}
	return server.DocDb
}
