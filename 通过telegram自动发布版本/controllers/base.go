package controllers

import (
	"strings"
	"tsEngine/tsTime"

	"github.com/astaxie/beego"
	"github.com/beego/i18n"
)

type BaseController struct {
	beego.Controller
	i18n.Locale
	Code   int64
	Msg    string
	Result interface{}
}

//定义全局变量
//var langTypes []string // Languages that are supported.

func init() {
	//logs.Trace("初始化控制器")
	////获取语言包列表
	//langTypes = strings.Split(beego.AppConfig.String("LangTypes"), "|")
	//
	//for _, lang := range langTypes {
	//	if lang != "" {
	//		logs.Trace("载入语言包: " + lang)
	//		if err := i18n.SetMessage(lang, "conf/"+"locale_"+lang+".ini"); err != nil {
	//			logs.Error("Fail to set message file:", err)
	//			return
	//		}
	//	}
	//
	//}
}

func (this *BaseController) Display(tpl string, layout bool) {
	//检测语言
	this.CheckLanguage()

	this.Data["Lang"] = this.Lang
	this.Data["Version"] = beego.AppConfig.String("Version")

	if beego.AppConfig.String("runmode") == "dev" {
		this.Data["Version"] = tsTime.CurrSe()
	}

	this.Data["Appname"] = beego.AppConfig.String("AppName")
	this.Data["Website"] = beego.AppConfig.String("WebSite")
	this.Data["Weburl"] = beego.AppConfig.String("WebUrl")
	this.Data["Email"] = beego.AppConfig.String("Email")
	if layout {
		this.Layout = "layout/main.html"
	}
	this.TplName = tpl + ".html"
}

//json 输出
func (this *BaseController) TraceJson() {
	//检测语言
	//this.CheckLanguage()
	//if this.Code > 0 {
	//	this.Msg = this.Tr("error." + fmt.Sprintf("%d", this.Code))
	//}
	this.Data["json"] = &map[string]interface{}{"Code": this.Code, "Msg": this.Msg, "Data": this.Result}
	this.ServeJSON()
	this.StopRun()
}

//检测语言包
func (this *BaseController) CheckLanguage() {

	//设置语言
	this.Lang = ""

	// 1. 获取 'Accept-Language' 值
	al := this.Ctx.Request.Header.Get("Accept-Language")
	if len(al) > 4 {
		al = al[:5] // Only compare first 5 letters.
		if i18n.IsExist(al) {
			this.Lang = al
		}
	}

	// 2. 默认为中文
	if len(this.Lang) == 0 {
		this.Lang = "zh-Cn"
	}

}

//判断是否是手机
func (this *BaseController) IsMobile() bool {

	agent := this.Ctx.Request.UserAgent()

	rule := []string{"Android", "iPhone", "SymbianOS", "Windows Phone", "iPad", "iPod"}

	for i := 0; i < len(rule); i++ {
		if strings.Contains(agent, rule[i]) {
			return true
		}

	}
	return false
}

type ErrorController struct {
	BaseController
}

func (this *ErrorController) Error404() {
	this.Data["Content"] = "page not found"
	this.TplName = "error.html"
}

func (this *ErrorController) Error501() {
	this.Data["Content"] = "server error"
	this.TplName = "error.html"
}

func (this *ErrorController) ErrorDb() {
	this.Data["Content"] = "database is now down"
	//this.Display("error")
	this.TplName = "error.html"
}
