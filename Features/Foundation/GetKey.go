package Foundation

import (
	"TianTing/Core/TianTingSDK"
	"github.com/kataras/iris/v12"
)

func RootKey(ctx iris.Context){
	_, _ = ctx.Text(TianTingSDK.GetServer().Key)
}
func Root(ctx iris.Context){
	_, _ = ctx.Text("")
}