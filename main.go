package main

import (
	"TianTing/Core"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	app := Core.New()
	app.Serve()
}
