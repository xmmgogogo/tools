package main

import (
	"QaVersionManage/models"
	_ "QaVersionManage/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

func main() {
	//log记录设置
	logs.Async()
	_ = logs.SetLogger("file", `{"filename":"./logs/logs.log", "perm":"0775", "maxDays":15}`)
	logs.SetLevel(logs.LevelDebug)
	logs.SetLogFuncCall(true)
	logs.SetLogFuncCallDepth(3)

	// 定时启动
	isHandOpenTGServer, _ := beego.AppConfig.Bool("HandOpenTGServer")
	if isHandOpenTGServer == true {
		go func() {
			models.StartTG()
		}()
	}

	beego.Run()
}
