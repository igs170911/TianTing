package main

import (
	"TianTing/Core"
)

func main() {
	app := Core.New()
	////inst := app.GetService()
	//testData := map[string]string{}
	//testData["from_type"] = "IGSMember"
	//testData["from_id"] = "700209"
	//testData["from_token"] = ""
	//testData["platform"] = "main"
	//testData["extra_data"] = "{}"
	//
	//da, er := json.Marshal(testData)
	//if er != nil {
	//	panic(er)
	//}
	//b64da := base64.StdEncoding.EncodeToString(da)
	//DataVerify := hmac.New(sha1.New, []byte(TianTingSDK.GetServer().GetKeyStatic()))
	//DataVerify.Write([]byte(b64da))
	//DataVerifyHexDigest := hex.EncodeToString(DataVerify.Sum(nil))
	//u := map[string]string{}
	//u["Sign"] = DataVerifyHexDigest
	//u["Data"] = b64da
	//b, err := json.Marshal(u)
	//if err != nil {
	//	panic(err)
	//}
	//data := base64.StdEncoding.EncodeToString(b)
	//fmt.Println(data)


	app.Serve()
}

