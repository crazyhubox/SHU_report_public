package main

import (
	"SHU/mail"
	"SHU/parseSetting"
	"SHU/report"
	"fmt"
	"runtime/debug"
	"strings"
	"time"
)

func checkTime(settings parseSetting.Setting) bool {
	t := time.Now()
	second := t.Second()
	hour := t.Hour()
	min := t.Minute()
	currentTime := fmt.Sprintf("%02d:%02d:%02d", hour, min, second)

	setting_time := settings.RuntimeCron
	s_t := strings.Split(setting_time, " ")
	setting_time = fmt.Sprintf("%02s:%02s:00", s_t[1], s_t[0])
	fmt.Println(currentTime,setting_time)
	return currentTime == setting_time
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


func everdayReport(reporter report.Reporter, settings parseSetting.Setting) {
	defer sendError(settings)
	for {
		settings, err := parseSetting.GetSettings()
		if err != nil {
			panic("Error on loading Settings.")
		}

		if checkTime(settings) {
			reporter.Report(settings.ReportInfo.StuNum, settings.ReportInfo.Password)
		}
		time.Sleep(time.Second)
	}
}


func main() {
	settings, err := parseSetting.GetSettings()
	if err != nil {
		panic("Error on loading Settings.")
	}
	
	var reporter report.Reporter = new(report.BrowReport)
	everdayReport(reporter,settings)
}
