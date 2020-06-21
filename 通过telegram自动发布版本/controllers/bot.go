package controllers

import (
	"QaVersionManage/models"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

type BotController struct {
	BaseController
}

// 这个服务暂时不用
func (this *BotController) GetBotUpdate() {
	messageType := this.GetString("type", "cloud_customize_bot")

	logs.Trace(messageType)
	if messageType == "" {
		this.Code = 1
		this.Msg = "messageType empty."
		this.TraceJson()
	}

	messageType = fmt.Sprintf("telegram_%s", messageType)
	chantId := beego.AppConfig.String(messageType + "::chat_id")
	botToken := beego.AppConfig.String(messageType + "::token")

	logs.Trace(fmt.Sprintf("chantId is %s, bottoken is %s", chantId, botToken))

	res, err := models.GetTelegramBot(botToken)
	if err != nil {
		logs.Error(err)
	}

	this.Result = res
	this.TraceJson()
}
