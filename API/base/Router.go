package base

import (
	"github.com/kataras/iris/v12"
)

func Routers(Group iris.Party){
	Group.Get("/",RootKey)
}