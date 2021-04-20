package report

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Reporter interface {
	Report(string, string)
}


type BrowReport struct {
	ApiReport
}

func(b *BrowReport) Report(uid, password string) {
	b.Init()
	b.GetCookies(uid, password)
	b.GetViewState()
	b.Submit()
}


func (b *BrowReport) GetCookies(username, password string) {
	var url string = fmt.Sprintf("http://0.0.0.0:8989/cookies/?id=%s&password=%s", username, password)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := b.client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", bodyText)

	var cookieData cookie
	json.Unmarshal(bodyText, &cookieData)
	b.cookies = fmt.Sprintf("%s=%s", cookieData.Name, cookieData.Value)
	fmt.Println(b.cookies)
}