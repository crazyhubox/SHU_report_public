package report

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Reporter interface {
	//默认就是应用类型
	Report(string, string)
}

type ApiReport struct {
	client    *http.Client
	cookier   CookieGenerator
	viewstate string
	cookies   string
}

//Cookie generator
type CookieGenerator interface {
	GetCookies(*http.Client, string, string) string
}


func NewApiReporter() *ApiReport {
	var reporter *ApiReport = new(ApiReport) //接口类型默认就是应用类型, 但是结构体不是,  new创建的是一个指针, 这里需要注意
	reporter.cookier = new(remoteCookie)
	reporter.init()

	return reporter
}


func (r *ApiReport) Report(uid, password string) {
	r.GetCookies(uid, password)
	r.GetViewState()
	r.Submit()
}

func (r *ApiReport) init() {
	testCookieJar, _ := cookiejar.New(nil)
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			for i := 0; i < len(via); i++ {
				url := via[i].URL.String()
				index_ := strings.Index(url, "LoginSSO.aspx")
				if index_ == 30 {
					return http.ErrUseLastResponse
				}
			}
			return nil
		},
		Jar: testCookieJar,
	}
	r.client = client
}

func (r *ApiReport) GetCookies(uid, password string) {
	//这里放账号和密码，代码可以自己重构，现在这个可能有点难看，哈哈
	cookies := r.cookier.GetCookies(r.client, uid, password)
	fmt.Println(cookies)
	r.cookies = cookies
}

func (r *ApiReport) GetViewState() {
	req, err := http.NewRequest("GET", "https://selfreport.shu.edu.cn/DayReport.aspx", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Cookie", r.cookies)
	resp, err := r.client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	doc.Find("input[name=__VIEWSTATE]").Each(func(i int, s *goquery.Selection) {
		value, e := s.Attr("value")
		if !e {
			panic("There is no viewstate.")
		}
		r.viewstate = value
	})
}

func (r *ApiReport) Submit() {
	date := getDate()
	data := r.getPostDataReader(date)
	req, err := http.NewRequest("POST", "https://selfreport.shu.edu.cn/DayReport.aspx", data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Cookie", r.cookies)
	resp, err := r.client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	regResult := regexp.MustCompile(`F.alert\((.+?)\)`)
	if regResult == nil {
		fmt.Println("MustCompile err. There are some errors during reporting.")
		return
	}
	result := regResult.FindSubmatch(bodyText)

	reportResult := r.showResult(date, result)
	fmt.Printf("%v\n", reportResult)
}

func (r *ApiReport) showResult(date string, result [][]byte) *ReportResult {
	reportResult := NewReportResult(date)
	resultStr := strings.TrimRight(strings.TrimLeft(fmt.Sprintf("%s", result[1]), "{"), "}")
	res := strings.Split(resultStr, ",")
	for _, eachLine := range res {
		var s, e string
		eachLineList := strings.Split(eachLine, ":")
		s = eachLineList[0]
		e = strings.TrimRight(strings.TrimLeft(eachLineList[1], "'"), "'")
		reportResult.result[s] = e
	}
	return reportResult
}

func (r *ApiReport) getPostDataReader(date string) *strings.Reader {
	// 这个建议在自己的每日一报页面开发者模式自己抓一下这个post的表单
	// 然后照这个格式把date日期和r.viewstate传进来就好了
	// 推荐使用这个网站https://curl.trillworks.com/,里面有教程
	postData := getPostData(r.viewstate,date)
	d := url.Values{}
	for k, v := range postData {
		d.Set(k, v)
	}
	var data = strings.NewReader(d.Encode())
	return data
}


func getDate() string {
	// 得到当前的时间,字符串形式
	now := time.Now()                  //获取当前时间
	timestamp := now.Unix()            //时间戳
	timeObj := time.Unix(timestamp, 0) //将时间戳转为时间格式
	year := timeObj.Year()             //年
	month := timeObj.Month()           //月
	day := timeObj.Day()               //日
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}