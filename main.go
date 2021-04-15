package main

import (
	"SHU/mail"
	"SHU/parseSetting"
	"SHU/report"
	"fmt"
	"log"
	"runtime/debug"

	"github.com/robfig/cron/v3"
)

func main() {
	settings, err := parseSetting.GetSettings()
	if err != nil {
		panic("Error on loading Settings.")
	}
	defer sendError(settings)
	// 以上是异常处理
	c := cron.New()
	c.AddFunc(settings.RuntimeCron, func() {
		// 这里需要重复写
		// 因为cron是新开了一个go程,如果这个go程中出错信息将无法进入主线程
		// 而且每次执行都重新加载配置,就可以不用关闭程序来修改配置了
		settings, err := parseSetting.GetSettings()
		if err != nil {
			panic("Error on loading Settings.")
		}
		defer sendError(settings)
		everdayReport(settings)
	})
	// 每天自动报送
	c.Start()

	select {}
}

func everdayReport(settings parseSetting.Setting) {
	report.VocationReportRun(settings.ReportInfo.StuNum, settings.ReportInfo.Password)
}

func sendError(settings parseSetting.Setting) {
	// 这个函数负责异常处理
	// 若发生异常则会发送一份邮件到qq邮箱
	// 若不使用注释掉即可
	errorInfo := recover()
	if errorInfo == nil {
		return
	}
	errorTraceback := debug.Stack()
	info := fmt.Sprintf("%s", errorInfo)
	traceback := fmt.Sprintf("%s", errorTraceback)
	fmt.Printf("%s\n", info)
	fmt.Printf("%s\n", traceback)
	err := mail.SendErro(info, traceback, settings)
	if err != nil {
		fmt.Println("Send fail! - ", err)
		return
	}
	fmt.Println("Send successfully!")

}

func test(u, p string) {
	report.VocationReportRun(u, p)
}

func testError() {
	settings, err := parseSetting.GetSettings()
	settings.ReportInfo.Password = ""
	if err != nil {
		log.Fatalln(err)
	}
	reporter := new(report.VocationReport)
	reporter.Init()
	reporter.GetCookies(settings.ReportInfo.StuNum, settings.ReportInfo.Password)
	reporter.GetViewState()
	reporter.Report()
}
