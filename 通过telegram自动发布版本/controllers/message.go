package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"telegram/models"
)

type MessageController struct {
	BaseController
}

func (c *MessageController) Prepare() {

}

func (c *MessageController) Get() {
	c.Data["json"] = map[string]interface{}{"success": 0, "message": ""}
	c.ServeJSON()
}

func (c *MessageController) Post() {
	messageType := c.GetString("type")
	message := c.GetString("message")

	logs.Trace(messageType, message)
	if messageType == "" || message == "" {
		c.Code = 1
		c.Msg = "type or message empty."
		c.TraceJson()
	}

	messageType = fmt.Sprintf("telegram_%s", messageType)
	chantId := beego.AppConfig.String(messageType + "::chat_id")
	botToken := beego.AppConfig.String(messageType + "::token")

	logs.Trace(fmt.Sprintf("chantId is %s, bottoken is %s", chantId, botToken))

	if chantId == "" || botToken == "" {
		c.Code = 2
		c.Msg = "chantId or botToken empty."
		c.TraceJson()
	}

	_ = models.SendTelegram(fmt.Sprintf("-%s", chantId), botToken, message)

	c.TraceJson()
}
