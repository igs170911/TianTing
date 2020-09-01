package TianTingSDK

import (
	"TianTing/Core/Database"
	"sync"
	"time"
)

// 整個 Server 的基礎架構，構成 Server 的主要結構

type TianTingServer struct {
	CodeName        string
	Key             string
	KeyUpdate       time.Time
	TokenExpireTime int

	//systemGameDict map[string]IGameSystem
	//socketGameDict map[string]IGameSystem
	//memberDict     map[string]IMemberSystem
	//
	//enableRpcTraffic  bool
	//gameRpcConnection LiquidRpc.GameAdapterClient
	DocDb       *Database.DocDB
	CacheDb     *Database.CacheDB
	//liquidRelationDb  *Database.RDB
	//liquidMsgQueue    MsgQueue.IAMQP
}

var Instance *TianTingServer
var once sync.Once

func GetServer() *TianTingServer {
	once.Do(func() {
		Instance = &TianTingServer{
			TokenExpireTime: int(time.Hour.Seconds()),
			//enableRpcTraffic: false,
			//systemGameDict:   make(map[string]IGameSystem),
			//socketGameDict:   make(map[string]IGameSystem),
			//memberDict:       make(map[string]IMemberSystem),
		}
	})
	return Instance
}
