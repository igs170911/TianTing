package Froundation

import (
	"TianTing/Middleware"
	"github.com/kataras/iris/v12"
)

func Routers(Group iris.Party){
	Group.Post("/api1", Middleware.Validate, HttpLogin)
	Group.Put("/api1", Middleware.Validate, HttpLogin)
}