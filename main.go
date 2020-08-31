package main

import (
	"TianTing/Core"
	"TianTing/Core/TianTingSDK"
	_ "github.com/joho/godotenv/autoload"
	"github.com/kataras/iris/v12"
)

func main() {
	app := Core.New()
	ins := app.GetService()
	ins.Get("/", func(ctx iris.Context) {
		_, _ = ctx.Text(TianTingSDK.GetServer().Key)
	})
	app.Serve()
}
