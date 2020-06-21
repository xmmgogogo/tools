package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/bndr/gojenkins"
	api "github.com/go-telegram-bot-api/telegram-bot-api"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var bot *api.BotAPI
var jenkins *gojenkins.Jenkins

func init() {
	// 初始化Jenkins
	JenkinsServerUrl := beego.AppConfig.String("JenkinsServerUrl")
	JenkinsUserId := beego.AppConfig.String("JenkinsUserId")
	JenkinsApiToken := beego.AppConfig.String("JenkinsApiToken")
	logs.Trace("初始化Jenkins：", JenkinsServerUrl, JenkinsUserId, JenkinsApiToken)
	jenkins = gojenkins.CreateJenkins(nil, JenkinsServerUrl, JenkinsUserId, JenkinsApiToken)
	_, err := jenkins.Init()
	if err != nil {
		panic("Jenkins init error:" + err.Error())
	}
}

// 启动定时
func StartTG() {
	var err error
	bot, err = api.NewBotAPI(beego.AppConfig.String("BotToken"))
	if err != nil {
		log.Panic(err)
	}
	bot.Debug, _ = beego.AppConfig.Bool("BotDebug")
	log.Printf("云守护机器人: %s  ID: %d", bot.Self.UserName, bot.Self.ID)

	u := api.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		panic("Can't get Updates")
	}

	for update := range updates {
		if update.Message == nil { // ignore any non-Message updates
			continue
		}
		go processUpdate(&update)
	}
}

// 自定义操作
func processUpdate(update *api.Update) {
	chatId := update.Message.Chat.ID
	text := strings.ToLower(update.Message.Text)

	logs.Trace("绑定消息内容", chatId, text)
	if text == "" {
		return
	}

	returnText := `
	参考格式：
		[平台名 c/s 分支名](不分大小写)：
		cloud c/s dev
		Manager c dev
		tsEngin	c dev
		GameApis s dev1.4
`
	switch text {
	case "/start":
		TgReply(update.Message.Chat.ID, update.Message.MessageID, returnText)
	default:
		textList := strings.Fields(text)

		// 区分是私聊还是群组
		if chatId < 0 {
			// 群组

			if textList[0] != beego.AppConfig.String("BotName") {
				return
			}

			// 踢掉头
			textList = textList[1:]
		}

		// 是否开启自动发布，可以关闭
		IsOpenAutoJenkins := beego.AppConfig.DefaultBool("IsOpenAutoJenkins", false)
		if IsOpenAutoJenkins == false {
			TgReply(update.Message.Chat.ID, update.Message.MessageID, "管理员已关闭自动发布，现在是手动发布模式")
			return
		}

		if len(textList) != 3 {
			returnText = `
					当前格式不正确!
				` + returnText

			TgReply(update.Message.Chat.ID, update.Message.MessageID, returnText)
			return
		}

		platformName := textList[0]
		platformTypeList := strings.Split(textList[1], "/")
		gitBranchName := textList[2]
		logs.Trace("接受tg输入：", platformName, platformTypeList, gitBranchName)

		for _, v := range platformTypeList {
			DoJenkins(platformName, v, gitBranchName, update.Message.Chat.ID, update.Message.MessageID)
		}
	}
}

/**
- 调用Jenkins
- platformName - cloud
- frontBack - c
- gitBranchName - dev1.4
*/
func DoJenkins(platformName, frontBack, gitBranchName string, chantId int64, replyToMessageID int) {
	jobName := ""

	if frontBack == "s" {
		// 后端
		switch platformName {
		case "cloud":
			jobName = "qa_hd_cloud"

			// manage
		case "cloud-manage":
			jobName = "qz_hd_manage"
		case "cloudmanage":
			jobName = "qz_hd_manage"
		case "manage":
			jobName = "qz_hd_manage"

			// CloudCronShell
		case "cloudcronshell":
			jobName = "qa_hd_CloudCronShell"

			// GameAPIs
		case "gameapis":
			jobName = "qa_hd_GameAPIs"

			// tsEngine
		case "tsengine":
			jobName = "tsEngine"
		case "tsengin":
			jobName = "tsEngine"
			// ...
		}
	} else if frontBack == "c" {
		// 前端
		switch platformName {
		case "cloud":
			jobName = "qa_web_cloud"
		case "cloud-manage":
			jobName = "qa_web_manage"
		case "cloudmanage":
			jobName = "qa_web_manage"
		case "manage":
			jobName = "qa_web_manage"
		}
	} else {
		// 其他情况处理
		TgReply(chantId, replyToMessageID, "发布前后端内容有误[c/s]："+frontBack)
		return
	}

	// 如果未支持
	if jobName == "" {
		TgReply(chantId, replyToMessageID, "项目名暂不支持(联系mm)："+platformName)
		return
	} else {
		var jobNum int64
		var err error
		if jobName == "qa_hd_cloud" || jobName == "qa_web_cloud" {
			jobNum, err = jenkins.BuildJob(strings.Trim(jobName, ""), map[string]string{"name": "Branch", "value": "*/" + gitBranchName})
		} else {
			jobNum, err = jenkins.BuildJob(strings.Trim(jobName, ""))
		}

		if err != nil {
			TgReply(chantId, replyToMessageID, fmt.Sprintf("发布失败[%s], 错误码[%s]", jobName, err.Error()))
		} else {
			TgReply(chantId, replyToMessageID, fmt.Sprintf("发布成功[%s], 序号[%d]", jobName, jobNum))
		}
	}
}

// TG回复内容
func TgReply(chantId int64, replyToMessageID int, returnText string) {
	msg := api.NewMessage(chantId, returnText)
	if replyToMessageID > 0 {
		msg.ReplyToMessageID = replyToMessageID
	}
	bot.Send(msg)
}

/**
- 调用发送接口
- chantId 渠道id，类似发送者ID
- botToken 机器人唯一编号
*/
func SendTelegram(chantId, botToken, message string) (err error) {
	// 创建错误通道
	c_err := make(chan error, 0)

	if chantId != "" && botToken != "" {
		PostUrl := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)

		go func() {
			resp, err := http.Post(PostUrl,
				"application/x-www-form-urlencoded",
				strings.NewReader("chat_id="+chantId+"&text="+message))
			if err != nil {
				logs.Error(err)
				c_err <- err
			}

			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				// handle error
				logs.Error("error!!!")
				c_err <- err
			}

			logs.Trace(string(body))
		}()
	}

	select {
	case err := <-c_err:
		logs.Error("发送请求失败：", err, chantId, botToken, message)
	default:
		logs.Trace("发送成功", chantId, botToken, message)
	}

	return
}

// 获取机器人状态，做相应动作
func GetTelegramBot(botToken string) (res string, err error) {
	PostUrl := fmt.Sprintf("https://api.telegram.org/bot%s/getUpdates", botToken)

	resp, err := http.Post(PostUrl,
		"application/x-www-form-urlencoded",
		strings.NewReader(""))
	if err != nil {
		logs.Error(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		logs.Error("[func.go] error.")
	}

	res = string(body)
	logs.Trace(res)

	return
}
