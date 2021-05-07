package report

// package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)


type cookie struct {
	Name     string `json:"name"`
	Value    string `json:"value"`
	Domain   string `json:"domain"`
	Path     string `json:"path"`
	Expires  int64  `json:"expires"`
	Size     int64  `json:"size"`
	HTTPOnly bool   `json:"httpOnly"`
	Secure   bool   `json:"secure"`
	Session  bool   `json:"session"`
}


type remoteCookie struct{}

func (r *remoteCookie) GetCookies(client *http.Client, user, password string) (cookies string) {
	url := fmt.Sprintf("http://127.0.0.1:8989/cookies?id=%s&password=%s", user, password)
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Connection","keep-alive")//持久化连接, 避免TCP的慢启动
	if err != nil {
		log.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	cookies = fmt.Sprintf("%s", bodyText)
	cookies = ClearCookie(cookies)
	return
}

func ClearCookie(cookies string) string {
	// cookies = `".ncov2019selfreport=D21B4E22E88DF45DF1BB9CC9F198FEF077BC4792CBF10AC785B68930418D5F2A0F19F54CAD34E64765933E2A5418FE62DD031355BD9771358A7D44BC4EDE17161ADFEF4ACAD749F69255AEADF6BC295C53384AADFA5E04A1BE713776DAEB48DCB5381EE543896BC774E2577D0F964C2A"`
	cookies = strings.TrimSpace(cookies)
	cookies = strings.Replace(cookies, `"`, "",2)
	return cookies
}
