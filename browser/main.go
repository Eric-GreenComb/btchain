package main

import (
	"github.com/astaxie/beego"
	"github.com/axengine/btchain/browser/datamanage"
	_ "github.com/axengine/btchain/browser/routers"
)

func main() {
	go datamanage.InitBlock()
	beego.Run()
}
