package main

import (
	"TianTing/Core"
)

type HelloWorld struct{}

func (s *HelloWorld)Hello(request string, reply *string) error{
	*reply = "Hello:" + request
	return nil
}

func main() {

	app := Core.New()
	//inst := app.GetService()
	app.Serve()
}
