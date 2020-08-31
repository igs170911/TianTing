package TianTingSDK

import (
	"TianTing/Logger"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"
)

var KeyTemplate = "MyNameIs_%s"

func (server *TianTingServer) SetCodeName(codename string) {
	Logger.SysLog.Infof("[Engine] Codename -> %s", codename)
	server.CodeName = codename
}

func (server *TianTingServer) GetKeyStatic() string {
	if time.Now().Sub(server.KeyUpdate).Minutes() > 60 {
		return server.GetKey()
	}
	return server.Key
}

func (server *TianTingServer) GetKey() string {
	RedisLiquidKeyName := fmt.Sprintf(KeyTemplate, server.CodeName)
	TianTingKey, GetKeyErr := server.GetCacheDb().Get(RedisLiquidKeyName)
	if GetKeyErr != nil {
		server.GenerateKey()
		server.KeyUpdate = time.Now()
		return server.Key
	}
	ReceivedLiquidKey := string(TianTingKey)
	if ReceivedLiquidKey != server.Key {
		server.Key = ReceivedLiquidKey
	}
	server.KeyUpdate = time.Now()
	return server.Key
}

func (server *TianTingServer) InitCodenameKey() {
	RedisLiquidKeyName := fmt.Sprintf(KeyTemplate, server.CodeName)
	TianTingKey, GetKeyErr := server.GetCacheDb().Get(RedisLiquidKeyName)
	if GetKeyErr != nil {
		server.GenerateKey()
	} else {
		server.Key = string(TianTingKey)
	}
	Logger.SysLog.Infof("[Engine] System Key -> %s", server.Key)
}

func (server *TianTingServer) GenerateKey() {
	conJunctions := "*****TianTing*****"
	md5Generate := md5.New()
	var keyOriginConcat bytes.Buffer
	keyOriginConcat.Write([]byte(server.CodeName))
	keyOriginConcat.Write([]byte(conJunctions))
	keyOriginConcat.Write([]byte(strconv.FormatInt(time.Now().Unix(), 10)))
	md5Generate.Write(keyOriginConcat.Bytes())
	Key := hex.EncodeToString(md5Generate.Sum(nil))
	RedisLiquidKeyName := fmt.Sprintf(KeyTemplate, server.CodeName)
	SaveKey2RedisErr := server.GetCacheDb().SetString(RedisLiquidKeyName, Key, -1)
	if SaveKey2RedisErr != nil {
		Logger.SysLog.Errorf("[System] Save System Key To Redis Error, %s", SaveKey2RedisErr)
	}
	server.Key = Key
}
