package TianTingSDK

// 從天庭傳送的資料格式定義，都由這隻檔案負責

// 登入，驗證登入流程 (Login, Verify)
type CmdAccount struct {
	FromType  string      `json:"from_type" default:""`
	FromId    string      `json:"from_id" default:""`
	FromToken string      `json:"from_token" default:""`
	Platform  string      `json:"platform" default:"main"`
	ExtraData interface{} `json:"extra_data" default:""`
}

// Middleware
type CmdSignedBody struct {
	Sign string `json:"Sign"`
	Data string `json:"Data"`
}

type CmdAccountResponse struct {
	AutoId     *string `json:"auto_id" bson:"auto_id"`
	InviteCode *string `json:"invite_code" bson:"invite_code"`
}