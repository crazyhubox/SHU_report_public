package report

import (
	// "SHU/report"
	"fmt"
	"testing"
)

const (
	uid string = ""
	password string = ""
)

func TestGetTime(t *testing.T) {
	res := getDate()
	t.Log(res)
}

func TestCookiesGet(t *testing.T) {
	test_reporter := NewApiReporter()
	test_reporter.GetCookies(uid,password)
	if test_reporter.cookies == "nocookie"{
		t.Error("nocookie")
	}
}

func TestViewStateGet(t *testing.T)  {
	test_reporter := NewApiReporter()
	test_reporter.GetCookies(uid,password)
	test_reporter.GetViewState()
	fmt.Println(test_reporter.viewstate)
	if len(test_reporter.viewstate) < 30 {
		t.Error("View_state is tested faily!")
	}
}