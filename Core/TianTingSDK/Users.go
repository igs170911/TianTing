package TianTingSDK

import (
	"TianTing/Logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
	"time"
)

// 系統中的會員順序，確保不會拿到重複的會員編號
type MemberAutoIncrement struct{
	Index			string `json:"index" bson:"index"`
	Counter		int `json:"counter" bson:"counter"`
}
const initAutoID = 1000000

func GetAutoID() string{
	// 在資料庫中搜尋的欄位
	filter := bson.M{"index":"auto_id"}
	// 將資料拿出來自動加一
	updateSql := bson.M{"$inc":bson.M{"counter":1}}
	var newAutoIdDoc *MemberAutoIncrement
	newAutoIdError := GetServer().GetMemberAutoIncrementCol().FindOneAndUpdate(
		nil,
		filter,
		updateSql,
		options.FindOneAndUpdate().SetUpsert(true),
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&newAutoIdDoc)
	if newAutoIdError != nil {
		Logger.SysLog.Errorf("[CMD][Login] Get New Auto ID Failed, %s", newAutoIdError)
		return ""
	}
	uid := strconv.Itoa(initAutoID + newAutoIdDoc.Counter)
	Logger.SysLog.Infof("[CMD][Login] New Auto ID Created : UID(%s)", uid)
	return uid
}



type User struct{
	AutoId     string    `json:"auto_id" bson:"auto_id"`
	InviteCode string    `json:"invite_code" bson:"invite_code"`
	Create     time.Time `json:"create" bson:"create"`
	Update     time.Time `json:"update" bson:"update"`
	FromType   string    `json:"from_type" bson:"from_type"`
	FromId     string    `json:"from_id" bson:"from_id"`
	FromToken  string    `json:"from_token" bson:"from_token"`
}