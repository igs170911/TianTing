package TianTingSDK

import (
	"TianTing/Logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math/rand"
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

// 最原始的會員系統，創造一個會員資訊，使得他可以綁定第三方會員

type User struct{
	// ColNameMember Number
	AutoId     string    `json:"auto_id" bson:"auto_id"`
	// 會員邀請碼 驗證這個會員是不是正確的會員
	InviteCode string    `json:"invite_code" bson:"invite_code"`
	// 新增時間
	Create     time.Time `json:"create" bson:"create"`
	// 修改時間
	Update     time.Time `json:"update" bson:"update"`
	// 該會員所榜定的第3方平台
	FromType   string    `json:"from_type" bson:"from_type"`
	// 該平台所對應的會員帳號
	FromId     string    `json:"from_id" bson:"from_id"`
	// 該平台所對應的TOKEN
	FromToken  string    `json:"from_token" bson:"from_token"`
}

 // 新增使用者
func CreateUser(fromType ,fromID string) *User{
	autoID := GetAutoID()
	if autoID =="" {
		Logger.SysLog.Errorf("[CMD][CreateUser] Can not get User ID" )
		return nil
	}
	CreateTime := time.Now()
	inviteCode := getAutoIdToInviteCode(autoID)
	if fromID == ""{
		fromID = autoID
	}

	TianTingUser := &User{
		AutoId: autoID,
		InviteCode: inviteCode,
		Create:  CreateTime,
		Update: CreateTime,
		FromType: fromType,
		FromId: fromID,
	}
	_,insertErr := GetServer().GetGameUserCol().InsertOne(nil,TianTingUser)
	if insertErr != nil{
		Logger.SysLog.Errorf("[CMD][CreateUser] Insert To DB error" )
		return nil
	}
	// TODO 看看要不要順便創立遊戲使用者
	// 用撈的可以驗證一下資料是否跟新增的一樣
	TianTingUser = FindGuestUser(autoID)
	return TianTingUser
}

// 確認是否有這個會員帳號（沒驗證
func FindGuestUser(autoID string) *User{
	filter := bson.M{"auto_id": autoID}
	var user *User
	findUserErr := GetServer().GetGameUserCol().FindOne(nil, filter).Decode(&user)
	if findUserErr != nil {
		Logger.SysLog.Errorf("[CMD][CreateUser] FindGuestUser error!" )
		return nil
	}
	return user
}

// 來平台上找該平台使用者
func FindUserByFrom(fromType string, fromId string) *User {
	filter := bson.M{"from_id": fromId, "from_type": fromType}
	var user *User
	findUserErr := GetServer().GetUserCol().FindOne(nil, filter).Decode(&user)
	if findUserErr != nil {
		Logger.SysLog.Errorf("[CMD][CreateUser] FindLiquidUserFromType error!" )
		return nil
	}
	return user
}

// 使用 AutoID + Invite Code 來找尋使用者
func FindUserByInviteCode(autoID ,inviteCode string) *User{
	filter := bson.M{"auto_id":autoID, "invite_code": inviteCode}
	var user *User
	findUserErr := GetServer().GetGameUserCol().FindOne(nil,filter)
	if findUserErr!= nil{
		Logger.SysLog.Errorf("[CMD][CreateUser] FindUserByInviteCode error!" )
		return nil
	}
	return user
}

// 將帳號綁定到第3方平台上面
func BindUser(AutoId string, FromId string, FromType string, FromToken string) (*mongo.UpdateResult, error) {
	filter := bson.M{"auto_id": AutoId}
	setBindData := bson.M{
		"$set": bson.M{
			"from_id":    FromId,
			"from_type":  FromType,
			"from_token": FromToken,
		},
	}
	setBindResult, setBindResultErr := GetServer().GetUserCol().UpdateOne(
		nil,
		filter,
		setBindData,
	)
	if setBindResultErr != nil {
		return nil, setBindResultErr
	}
	return setBindResult, nil
}



// 新增邀請碼
var convertTable = [...]string{
	"U", "V", "W", "X", "Y",
	"A", "B", "C", "D", "E",
	"F", "G", "H", "I", "J",
	"L", "M", "N", "Z", "K",
	"P", "Q", "R", "S", "T",
}
func getAutoIdToInviteCode(autoId string) string {

	inviteCode := ""
	InviteCodeList := make([]int, 0)
	autoIdInt, _ := strconv.ParseInt(autoId, 10, 64)
	rand.Seed(autoIdInt + time.Now().UnixNano())
	newId := strconv.FormatInt(time.Now().Unix()+rand.Int63n(time.Now().Unix()), 10)
	newId = newId + autoId
	newId2Int, _ := strconv.Atoi(newId[len(newId)-13:])

	resultToBase := decimalToBase(InviteCodeList, newId2Int)
	for _, base := range resultToBase {
		inviteCode += convertTable[base]
	}
	return inviteCode

}

func decimalToBase(baseList []int, decimal int) []int {
	base := len(convertTable)
	baseList = append(baseList, decimal%base)
	div := decimal / base
	if div == 0 {
		return baseList
	}
	return decimalToBase(baseList, div)
}