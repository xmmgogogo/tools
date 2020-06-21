package routers

import (
	"QaVersionManage/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/send_message", &controllers.MessageController{})
	beego.AutoRouter(&controllers.BotController{})
}
