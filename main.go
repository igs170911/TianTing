package main

import (
	"TianTing/Core"
	"TianTing/Core/TianTingSDK"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	app := Core.New()
	ins := app.GetService()
	ins.Get("/", func(ctx iris.Context) {
		_, _ = ctx.Text(TianTingSDK.GetServer().Key)
	})
	app.Serve()
}
