package Froundation

import (
	"TianTing/Core/TianTingSDK"
	"TianTing/Logger"
	"encoding/json"
	"github.com/kataras/iris/v12"
)
type CmdData struct{
	CmdData string
}

// 登入流程
//  1. 如果有帳號，走登錄流程
//  2. 如果沒帳號，則註冊一個新的帳號
func HttpLogin(ctx iris.Context) {
	var command *TianTingSDK.CmdAccount

	RawData := ctx.Values().Get("CommandData")
	_ = json.Unmarshal(RawData.([]byte), &command)
	Logger.SysLog.Debugf("[CMD][Login] %+v", command)

	// 登入的最後結果
	result := &TianTingSDK.CmdAccountResponse{
		AutoId:     nil,
		InviteCode: nil,
	}

	// 使用者資訊
	var user *TianTingSDK.User



	_, _ = ctx.Text("Login Function")
	//var command *TianTingSDK.CmdAccount
	// 取得資料並寫轉存到該格式當中
	//_ = json.Unmarshal(ctx.MustGet("CommandData").([]byte), &command)

}

func WithOutMid(ctx iris.Context) {
	_, _ = ctx.Text("Without MID Function")
}