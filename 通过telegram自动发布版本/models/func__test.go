package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/bndr/gojenkins"
	"testing"
)

func TestDoJenkins(t *testing.T) {
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

	jobNum, err := jenkins.BuildJob("tsEngine")
	if err != nil {
		logs.Trace(err)
		return
	}

	logs.Trace("当前job编号：", jobNum)
}
