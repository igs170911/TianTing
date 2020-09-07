package Middleware

import (
	"TianTing/Core/TianTingSDK"
	"TianTing/Logger"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"github.com/kataras/iris/v12"
	"net/http"
)

// 傳入資料的驗證,沒有驗證的資料就被會擋下來
func Validate(ctx iris.Context){
	key := TianTingSDK.GetServer().GetKeyStatic()
	// 無論用神麼 Http Methods, 戴上來的資料都要解密, 都使用 x-www-form-urlencoded
	// Step 1 收資料
	RawDataBody, GetRawBodyErr := ctx.GetBody()
	if GetRawBodyErr != nil {
		Logger.SysLog.Errorf("[Engine][Middleware] Validate Failed, %s", GetRawBodyErr)
		ctx.StopWithStatus(http.StatusBadRequest)
		return
	}
	// Step 2 解壓縮資料
	DataBody, decodeBodyErr := base64.StdEncoding.DecodeString(string(RawDataBody))
	if decodeBodyErr != nil {
		Logger.SysLog.Errorf("[Engine][Middleware] Decode Data Failed, %s", decodeBodyErr)
		ctx.StopWithStatus(http.StatusBadRequest)
		return
	}

	var StructureOfTianTing *TianTingSDK.CmdSignedBody
	// byte data to string
	DataUnmarshalError := json.Unmarshal(DataBody, &StructureOfTianTing)
	if DataUnmarshalError != nil {
		Logger.SysLog.Errorf("[Engine][Middleware] Unmarshal Data Failed, %s", DataUnmarshalError)
		ctx.StopWithStatus(http.StatusBadRequest)
		return
	}

	// 驗證資料是否有被竄改過
	DataVerify := hmac.New(sha1.New, []byte(key))
	DataVerify.Write([]byte(StructureOfTianTing.Data))
	DataVerifyHexDigest := hex.EncodeToString(DataVerify.Sum(nil))
	if StructureOfTianTing.Sign != DataVerifyHexDigest {
		Logger.SysLog.Error("[Engine][Middleware] Verify Data Failed")
		ctx.StopWithStatus(http.StatusBadRequest)
		return
	}

	DecodedCommandData, DecodedCommandDataError := base64.StdEncoding.DecodeString(StructureOfTianTing.Data)
	if DecodedCommandDataError != nil {
		Logger.SysLog.Errorf("[Engine][Middleware] Decode Command Data Failed, %s", DecodedCommandDataError)
		ctx.StopWithStatus(http.StatusBadRequest)
		return
	}
	// 將資料傳遞到下一個步驟
	ctx.Values().Set("CommandData",DecodedCommandData)
	ctx.Next()
}