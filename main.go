package main

import (
	_ "cloud-platform-ua/routers"

	"github.com/astaxie/beego"
	_ "github.com/astaxie/beego/session/redis"
)

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}
