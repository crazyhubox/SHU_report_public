package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

type Setting struct {
	Mail        Mail       `json:"mail"`
	ReportInfo  ReportInfo `json:"report_info"`
	RuntimeCron string     `json:"runtime_cron"`
}

type Mail struct {
	Sender            string   `json:"sender"`
	Receivers         []string `json:"receivers"`
	AuthorizationCode string   `json:"authorization_code"`
}

type ReportInfo struct {
	StuNum   string `json:"stu_num"`
	Password string `json:"password"`
}

//这里需要注意,json中的字段名字必须和结构体中的tag值相同
//如果不小心修改了json当中的字段名字,则不会报错,但是程序卡死
func GetSettings() (setting Setting, err error) {
	byte, _ := ioutil.ReadFile("setting.json")
	// fmt.Printf("%s\n", byte)
	err = json.Unmarshal(byte, &setting)
	if err != nil {
		log.Fatalln(err)
	}
	return
}

func main() {
	s, _ := GetSettings()
	fmt.Println(s.RuntimeCron)
}
