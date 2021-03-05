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

type VocationReport struct {
	client    *http.Client
	result    string
	viewstate string
	cookies   string
}

type ReportResult struct {
	result map[string]string
	date   string
}

func NewReportResult(date string) *ReportResult {
	var tMap = make(map[string]string)
	//map映射对象必须要创建，不然就会指向nil，指向nil的map无法保存任何值
	return &ReportResult{
		result: tMap,
		date:   date,
	}
}


func VocationReportRun(uid, password string) {
	reporter := new(VocationReport)
	reporter.Init()
	reporter.GetCookies(uid,password)
	reporter.GetViewState()
	reporter.Report()
}

func (r *VocationReport) Init() {
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

func (r *VocationReport) GetCookies(uid, password string) {
	//这里放账号和密码，代码可以自己重构，现在这个可能有点难看，哈哈
	dataStr := fmt.Sprintf(`username=%s&password=%s&login_submit=`, uid, password)
	var data = strings.NewReader(dataStr)
	req, err := http.NewRequest("POST", "https://newsso.shu.edu.cn/login/eyJ0aW1lc3RhbXAiOjE2MTIzNTY1Mzk3NDA0MjQ4OTUsInJlc3BvbnNlVHlwZSI6ImNvZGUiLCJjbGllbnRJZCI6IldVSFdmcm50bldZSFpmelE1UXZYVUNWeSIsInNjb3BlIjoiMSIsInJlZGlyZWN0VXJpIjoiaHR0cHM6Ly9zZWxmcmVwb3J0LnNodS5lZHUuY24vTG9naW5TU08uYXNweD9SZXR1cm5Vcmw9JTJmRGVmYXVsdC5hc3B4Iiwic3RhdGUiOiIifQ==", data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Origin", "https://newsso.shu.edu.cn")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.96 Safari/537.36 Edg/88.0.705.56")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Referer", "https://newsso.shu.edu.cn/login/eyJ0aW1lc3RhbXAiOjE2MTIzNTY1Mzk3NDA0MjQ4OTUsInJlc3BvbnNlVHlwZSI6ImNvZGUiLCJjbGllbnRJZCI6IldVSFdmcm50bldZSFpmelE1UXZYVUNWeSIsInNjb3BlIjoiMSIsInJlZGlyZWN0VXJpIjoiaHR0cHM6Ly9zZWxmcmVwb3J0LnNodS5lZHUuY24vTG9naW5TU08uYXNweD9SZXR1cm5Vcmw9JTJmRGVmYXVsdC5hc3B4Iiwic3RhdGUiOiIifQ==")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
	resp, err := r.client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	cookie := resp.Cookies()
	// fmt.Println(cookie)
	for _, v := range cookie {
		fmt.Println(v.Name)
		if v.Name == ".ncov2019selfreport" {
			r.cookies = fmt.Sprintf(".ncov2019selfreport=%s", v.Value)
		}
	}
	fmt.Println(r.cookies)
}

func (r *VocationReport) GetViewState() {
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

func (r *VocationReport) Report() {
	now := time.Now()                  //获取当前时间
	timestamp := now.Unix()            //时间戳
	timeObj := time.Unix(timestamp, 0) //将时间戳转为时间格式
	year := timeObj.Year()             //年
	month := timeObj.Month()           //月
	day := timeObj.Day()               //日
	date := fmt.Sprintf("%d-%02d-%02d", year, month, day)
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

func (r *VocationReport) showResult(date string, result [][]byte) *ReportResult {
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

func (r *VocationReport) getPostDataReader(date string) *strings.Reader {
	// 这个建议在自己的每日一报页面开发者模式自己抓一下这个post的表单
	// 然后照这个格式把date日期和r.viewstate传进来就好了
	// 推荐使用这个网站https://curl.trillworks.com/,里面有教程
	postData := map[string]string{
		"__EVENTTARGET":              "p1$ctl00$btnSubmit",
		"__EVENTARGUMENT":            "",
		"__VIEWSTATE":                r.viewstate,
		"__VIEWSTATEGENERATOR":       "7AD7E509",
		"p1$ChengNuo":                "p1_ChengNuo",
		"p1$BaoSRQ":                  date,
		"p1$DangQSTZK":               "良好",
		"p1$TiWen":                   "",
		"p1$JiuYe_ShouJHM":           "",
		"p1$JiuYe_Email":             "",
		"p1$JiuYe_Wechat":            "",
		"p1$QiuZZT":                  "",
		"p1$JiuYKN":                  "",
		"p1$JiuYSJ":                  "",
		"p1$GuoNei":                  "国内",
		"p1$ddlGuoJia$Value":         "-1",
		"p1$ddlGuoJia":               "选择国家",
		"p1$ShiFSH":                  "是",
		"p1$ShiFZX":                  "是",
		"p1$ddlSheng$Value":          "上海",
		"p1$ddlSheng":                "上海",
		"p1$ddlShi$Value":            "上海市",
		"p1$ddlShi":                  "上海市",
		"p1$ddlXian$Value":           "宝山区",
		"p1$ddlXian":                 "宝山区",
		"p1$XiangXDZ":                "上海大学宝山校区校内s2栋207",
		"p1$ShiFZJ":                  "否",
		"p1$FengXDQDL":               "否",
		"p1$TongZWDLH":               "否",
		"p1$CengFWH":                 "否",
		"p1$CengFWH_RiQi":            "",
		"p1$CengFWH_BeiZhu":          "",
		"p1$JieChu":                  "否",
		"p1$JieChu_RiQi":             "",
		"p1$JieChu_BeiZhu":           "",
		"p1$TuJWH":                   "否",
		"p1$TuJWH_RiQi":              "",
		"p1$TuJWH_BeiZhu":            "",
		"p1$QueZHZJC$Value":          "否",
		"p1$QueZHZJC":                "否",
		"p1$DangRGL":                 "否",
		"p1$GeLDZ":                   "",
		"p1$FanXRQ":                  "",
		"p1$WeiFHYY":                 "",
		"p1$ShangHJZD":               "",
		"p1$DaoXQLYGJ":               "中国",
		"p1$DaoXQLYCS":               "贵阳",
		"p1$JiaRen_BeiZhu":           "",
		"p1$SuiSM":                   "绿色",
		"p1$LvMa14Days":              "是",
		"p1$Address2":                "",
		"p1_ContentPanel1_Collapsed": "true",
		"p1_GeLSM_Collapsed":         "false",
		"p1_Collapsed":               "false",
		"F_STATE":                    "eyJwMV9CYW9TUlEiOnsiVGV4dCI6IjIwMjEtMDMtMDEifSwicDFfRGFuZ1FTVFpLIjp7IkZfSXRlbXMiOltbIuiJr+WlvSIsIuiJr+Wlve+8iOS9k+a4qeS4jemrmOS6jjM3LjPvvIkiLDFdLFsi5LiN6YCCIiwi5LiN6YCCIiwxXV0sIlNlbGVjdGVkVmFsdWUiOiLoia/lpb0ifSwicDFfWmhlbmdaaHVhbmciOnsiSGlkZGVuIjp0cnVlLCJGX0l0ZW1zIjpbWyLmhJ/lhpIiLCLmhJ/lhpIiLDFdLFsi5ZKz5Ze9Iiwi5ZKz5Ze9IiwxXSxbIuWPkeeDrSIsIuWPkeeDrSIsMV1dLCJTZWxlY3RlZFZhbHVlQXJyYXkiOltdfSwicDFfUWl1WlpUIjp7IkZfSXRlbXMiOltdLCJTZWxlY3RlZFZhbHVlQXJyYXkiOltdfSwicDFfSml1WUtOIjp7IkZfSXRlbXMiOltdLCJTZWxlY3RlZFZhbHVlQXJyYXkiOltdfSwicDFfSml1WVlYIjp7IlJlcXVpcmVkIjpmYWxzZSwiRl9JdGVtcyI6W10sIlNlbGVjdGVkVmFsdWVBcnJheSI6W119LCJwMV9KaXVZWkQiOnsiRl9JdGVtcyI6W10sIlNlbGVjdGVkVmFsdWVBcnJheSI6W119LCJwMV9KaXVZWkwiOnsiRl9JdGVtcyI6W10sIlNlbGVjdGVkVmFsdWVBcnJheSI6W119LCJwMV9HdW9OZWkiOnsiRl9JdGVtcyI6W1si5Zu95YaFIiwi5Zu95YaFIiwxXSxbIuWbveWkliIsIuWbveWkliIsMV1dLCJTZWxlY3RlZFZhbHVlIjoi5Zu95YaFIn0sInAxX2RkbEd1b0ppYSI6eyJEYXRhVGV4dEZpZWxkIjoiWmhvbmdXZW4iLCJEYXRhVmFsdWVGaWVsZCI6Ilpob25nV2VuIiwiRl9JdGVtcyI6W1siLTEiLCLpgInmi6nlm73lrrYiLDEsIiIsIiJdLFsi6Zi/5bCU5be05bC85LqaIiwi6Zi/5bCU5be05bC85LqaIiwxLCIiLCIiXSxbIumYv+WwlOWPiuWIqeS6miIsIumYv+WwlOWPiuWIqeS6miIsMSwiIiwiIl0sWyLpmL/lr4zmsZciLCLpmL/lr4zmsZciLDEsIiIsIiJdLFsi6Zi/5qC55bu3Iiwi6Zi/5qC55bu3IiwxLCIiLCIiXSxbIumYv+aLieS8r+iBlOWQiOmFi+mVv+WbvSIsIumYv+aLieS8r+iBlOWQiOmFi+mVv+WbvSIsMSwiIiwiIl0sWyLpmL/psoHlt7QiLCLpmL/psoHlt7QiLDEsIiIsIiJdLFsi6Zi/5pu8Iiwi6Zi/5pu8IiwxLCIiLCIiXSxbIumYv+WhnuaLnOeWhiIsIumYv+WhnuaLnOeWhiIsMSwiIiwiIl0sWyLln4Plj4oiLCLln4Plj4oiLDEsIiIsIiJdLFsi5Z+D5aGe5L+E5q+U5LqaIiwi5Z+D5aGe5L+E5q+U5LqaIiwxLCIiLCIiXSxbIueIseWwlOWFsCIsIueIseWwlOWFsCIsMSwiIiwiIl0sWyLniLHmspnlsLzkupoiLCLniLHmspnlsLzkupoiLDEsIiIsIiJdLFsi5a6J6YGT5bCUIiwi5a6J6YGT5bCUIiwxLCIiLCIiXSxbIuWuieWTpeaLiSIsIuWuieWTpeaLiSIsMSwiIiwiIl0sWyLlronlnK3mi4kiLCLlronlnK3mi4kiLDEsIiIsIiJdLFsi5a6J5o+Q55Oc5ZKM5be05biD6L6+Iiwi5a6J5o+Q55Oc5ZKM5be05biD6L6+IiwxLCIiLCIiXSxbIuWlpeWcsOWIqSIsIuWlpeWcsOWIqSIsMSwiIiwiIl0sWyLlpaXlhbDnvqTlspsiLCLlpaXlhbDnvqTlspsiLDEsIiIsIiJdLFsi5r6z5aSn5Yip5LqaIiwi5r6z5aSn5Yip5LqaIiwxLCIiLCIiXSxbIuW3tOW3tOWkmuaWryIsIuW3tOW3tOWkmuaWryIsMSwiIiwiIl0sWyLlt7TluIPkuprmlrDlh6DlhoXkupoiLCLlt7TluIPkuprmlrDlh6DlhoXkupoiLDEsIiIsIiJdLFsi5be05ZOI6amsIiwi5be05ZOI6amsIiwxLCIiLCIiXSxbIuW3tOWfuuaWr+WdpiIsIuW3tOWfuuaWr+WdpiIsMSwiIiwiIl0sWyLlt7Tli5Lmlq/lnaYiLCLlt7Tli5Lmlq/lnaYiLDEsIiIsIiJdLFsi5be05p6XIiwi5be05p6XIiwxLCIiLCIiXSxbIuW3tOaLv+mprCIsIuW3tOaLv+mprCIsMSwiIiwiIl0sWyLlt7Topb8iLCLlt7Topb8iLDEsIiIsIiJdLFsi55m95L+E572X5pavIiwi55m95L+E572X5pavIiwxLCIiLCIiXSxbIueZvuaFleWkpyIsIueZvuaFleWkpyIsMSwiIiwiIl0sWyLkv53liqDliKnkupoiLCLkv53liqDliKnkupoiLDEsIiIsIiJdLFsi6LSd5a6BIiwi6LSd5a6BIiwxLCIiLCIiXSxbIuavlOWIqeaXtiIsIuavlOWIqeaXtiIsMSwiIiwiIl0sWyLlhrDlspsiLCLlhrDlspsiLDEsIiIsIiJdLFsi5rOi5aSa6buO5ZCEIiwi5rOi5aSa6buO5ZCEIiwxLCIiLCIiXSxbIuazouWFsCIsIuazouWFsCIsMSwiIiwiIl0sWyLms6Lmlq/lsLzkuprlkozpu5HloZ7lk6Xnu7TpgqMiLCLms6Lmlq/lsLzkuprlkozpu5HloZ7lk6Xnu7TpgqMiLDEsIiIsIiJdLFsi54675Yip57u05LqaIiwi54675Yip57u05LqaIiwxLCIiLCIiXSxbIuS8r+WIqeWFuSIsIuS8r+WIqeWFuSIsMSwiIiwiIl0sWyLljZrojKjnk6bnurMiLCLljZrojKjnk6bnurMiLDEsIiIsIiJdLFsi5LiN5Li5Iiwi5LiN5Li5IiwxLCIiLCIiXSxbIuW4g+Wfuue6s+azlee0oiIsIuW4g+Wfuue6s+azlee0oiIsMSwiIiwiIl0sWyLluIPpmobov6oiLCLluIPpmobov6oiLDEsIiIsIiJdLFsi5biD57u05bKbIiwi5biD57u05bKbIiwxLCIiLCIiXSxbIuacnemynCIsIuacnemynCIsMSwiIiwiIl0sWyLotaTpgZPlh6DlhoXkupoiLCLotaTpgZPlh6DlhoXkupoiLDEsIiIsIiJdLFsi5Li56bqmIiwi5Li56bqmIiwxLCIiLCIiXSxbIuW+t+WbvSIsIuW+t+WbvSIsMSwiIiwiIl0sWyLkuJzluJ3msbYiLCLkuJzluJ3msbYiLDEsIiIsIiJdLFsi5Lic5bid5rG2Iiwi5Lic5bid5rG2IiwxLCIiLCIiXSxbIuWkmuWTpSIsIuWkmuWTpSIsMSwiIiwiIl0sWyLlpJrnsbPlsLzliqAiLCLlpJrnsbPlsLzliqAiLDEsIiIsIiJdLFsi5L+E572X5pav6IGU6YKmIiwi5L+E572X5pav6IGU6YKmIiwxLCIiLCIiXSxbIuWOhOeTnOWkmuWwlCIsIuWOhOeTnOWkmuWwlCIsMSwiIiwiIl0sWyLljoTnq4vnibnph4zkupoiLCLljoTnq4vnibnph4zkupoiLDEsIiIsIiJdLFsi5rOV5Zu9Iiwi5rOV5Zu9IiwxLCIiLCIiXSxbIuazleWbveWkp+mDveS8miIsIuazleWbveWkp+mDveS8miIsMSwiIiwiIl0sWyLms5XnvZfnvqTlspsiLCLms5XnvZfnvqTlspsiLDEsIiIsIiJdLFsi5rOV5bGe5rOi5Yip5bC86KW/5LqaIiwi5rOV5bGe5rOi5Yip5bC86KW/5LqaIiwxLCIiLCIiXSxbIuazleWxnuWcreS6mumCoyIsIuazleWxnuWcreS6mumCoyIsMSwiIiwiIl0sWyLmorXokoLlhogiLCLmorXokoLlhogiLDEsIiIsIiJdLFsi6I+y5b6L5a6+Iiwi6I+y5b6L5a6+IiwxLCIiLCIiXSxbIuaWkOa1jiIsIuaWkOa1jiIsMSwiIiwiIl0sWyLoiqzlhbAiLCLoiqzlhbAiLDEsIiIsIiJdLFsi5L2b5b6X6KeSIiwi5L2b5b6X6KeSIiwxLCIiLCIiXSxbIuWGiOavlOS6miIsIuWGiOavlOS6miIsMSwiIiwiIl0sWyLliJrmnpwiLCLliJrmnpwiLDEsIiIsIiJdLFsi5Yia5p6c77yI6YeR77yJIiwi5Yia5p6c77yI6YeR77yJIiwxLCIiLCIiXSxbIuWTpeS8puavlOS6miIsIuWTpeS8puavlOS6miIsMSwiIiwiIl0sWyLlk6Xmlq/ovr7pu47liqAiLCLlk6Xmlq/ovr7pu47liqAiLDEsIiIsIiJdLFsi5qC85p6X57qz6L6+Iiwi5qC85p6X57qz6L6+IiwxLCIiLCIiXSxbIuagvOmygeWQieS6miIsIuagvOmygeWQieS6miIsMSwiIiwiIl0sWyLmoLnopb/lspsiLCLmoLnopb/lspsiLDEsIiIsIiJdLFsi5Y+k5be0Iiwi5Y+k5be0IiwxLCIiLCIiXSxbIueTnOW+t+e9l+aZruWymyIsIueTnOW+t+e9l+aZruWymyIsMSwiIiwiIl0sWyLlhbPlspsiLCLlhbPlspsiLDEsIiIsIiJdLFsi5Zyt5Lqa6YKjIiwi5Zyt5Lqa6YKjIiwxLCIiLCIiXSxbIuWTiOiQqOWFi+aWr+WdpiIsIuWTiOiQqOWFi+aWr+WdpiIsMSwiIiwiIl0sWyLmtbflnLAiLCLmtbflnLAiLDEsIiIsIiJdLFsi6Z+p5Zu9Iiwi6Z+p5Zu9IiwxLCIiLCIiXSxbIuiNt+WFsCIsIuiNt+WFsCIsMSwiIiwiIl0sWyLpu5HlsbEiLCLpu5HlsbEiLDEsIiIsIiJdLFsi5rSq6YO95ouJ5pavIiwi5rSq6YO95ouJ5pavIiwxLCIiLCIiXSxbIuWfuumHjOW3tOaWryIsIuWfuumHjOW3tOaWryIsMSwiIiwiIl0sWyLlkInluIPmj5AiLCLlkInluIPmj5AiLDEsIiIsIiJdLFsi5ZCJ5bCU5ZCJ5pav5pav5Z2mIiwi5ZCJ5bCU5ZCJ5pav5pav5Z2mIiwxLCIiLCIiXSxbIuWHoOWGheS6miIsIuWHoOWGheS6miIsMSwiIiwiIl0sWyLlh6DlhoXkuprmr5Tnu40iLCLlh6DlhoXkuprmr5Tnu40iLDEsIiIsIiJdLFsi5Yqg5ou/5aSnIiwi5Yqg5ou/5aSnIiwxLCIiLCIiXSxbIuWKoOe6syIsIuWKoOe6syIsMSwiIiwiIl0sWyLliqDok6wiLCLliqDok6wiLDEsIiIsIiJdLFsi5p+s5Z+U5a+oIiwi5p+s5Z+U5a+oIiwxLCIiLCIiXSxbIuaNt+WFiyIsIuaNt+WFiyIsMSwiIiwiIl0sWyLmtKXlt7TluIPpn6YiLCLmtKXlt7TluIPpn6YiLDEsIiIsIiJdLFsi5ZaA6bqm6ZqGIiwi5ZaA6bqm6ZqGIiwxLCIiLCIiXSxbIuWNoeWhlOWwlCIsIuWNoeWhlOWwlCIsMSwiIiwiIl0sWyLnp5Hnp5Hmlq/vvIjln7rmnpfvvInnvqTlspsiLCLnp5Hnp5Hmlq/vvIjln7rmnpfvvInnvqTlspsiLDEsIiIsIiJdLFsi56eR5pGp572XIiwi56eR5pGp572XIiwxLCIiLCIiXSxbIuenkeeJuei/queTpiIsIuenkeeJuei/queTpiIsMSwiIiwiIl0sWyLnp5HlqIHnibkiLCLnp5HlqIHnibkiLDEsIiIsIiJdLFsi5YWL572X5Zyw5LqaIiwi5YWL572X5Zyw5LqaIiwxLCIiLCIiXSxbIuiCr+WwvOS6miIsIuiCr+WwvOS6miIsMSwiIiwiIl0sWyLlupPlhYvnvqTlspsiLCLlupPlhYvnvqTlspsiLDEsIiIsIiJdLFsi5ouJ6ISx57u05LqaIiwi5ouJ6ISx57u05LqaIiwxLCIiLCIiXSxbIuiOsee0ouaJmCIsIuiOsee0ouaJmCIsMSwiIiwiIl0sWyLogIHmjJ0iLCLogIHmjJ0iLDEsIiIsIiJdLFsi6buO5be05aupIiwi6buO5be05aupIiwxLCIiLCIiXSxbIueri+mZtuWumyIsIueri+mZtuWumyIsMSwiIiwiIl0sWyLliKnmr5Tph4zkupoiLCLliKnmr5Tph4zkupoiLDEsIiIsIiJdLFsi5Yip5q+U5LqaIiwi5Yip5q+U5LqaIiwxLCIiLCIiXSxbIuWIl+aUr+aVpuWjq+eZuyIsIuWIl+aUr+aVpuWjq+eZuyIsMSwiIiwiIl0sWyLnlZnlsLzmsarlspsiLCLnlZnlsLzmsarlspsiLDEsIiIsIiJdLFsi5Y2i5qOu5aChIiwi5Y2i5qOu5aChIiwxLCIiLCIiXSxbIuWNouaXuui+viIsIuWNouaXuui+viIsMSwiIiwiIl0sWyLnvZfpqazlsLzkupoiLCLnvZfpqazlsLzkupoiLDEsIiIsIiJdLFsi6ams6L6+5Yqg5pav5YqgIiwi6ams6L6+5Yqg5pav5YqgIiwxLCIiLCIiXSxbIumprOaBqeWymyIsIumprOaBqeWymyIsMSwiIiwiIl0sWyLpqazlsJTku6PlpKsiLCLpqazlsJTku6PlpKsiLDEsIiIsIiJdLFsi6ams6ICz5LuWIiwi6ams6ICz5LuWIiwxLCIiLCIiXSxbIumprOaLiee7tCIsIumprOaLiee7tCIsMSwiIiwiIl0sWyLpqazmnaXopb/kupoiLCLpqazmnaXopb/kupoiLDEsIiIsIiJdLFsi6ams6YeMIiwi6ams6YeMIiwxLCIiLCIiXSxbIumprOWFtumhvyIsIumprOWFtumhvyIsMSwiIiwiIl0sWyLpqaznu43lsJTnvqTlspsiLCLpqaznu43lsJTnvqTlspsiLDEsIiIsIiJdLFsi6ams5o+Q5bC85YWL5bKbIiwi6ams5o+Q5bC85YWL5bKbIiwxLCIiLCIiXSxbIumprOe6pueJuSIsIumprOe6pueJuSIsMSwiIiwiIl0sWyLmr5vph4zmsYLmlq8iLCLmr5vph4zmsYLmlq8iLDEsIiIsIiJdLFsi5q+b6YeM5aGU5bC85LqaIiwi5q+b6YeM5aGU5bC85LqaIiwxLCIiLCIiXSxbIue+juWbvSIsIue+juWbvSIsMSwiIiwiIl0sWyLnvo7lsZ7okKjmkankupoiLCLnvo7lsZ7okKjmkankupoiLDEsIiIsIiJdLFsi6JKZ5Y+kIiwi6JKZ5Y+kIiwxLCIiLCIiXSxbIuiSmeeJueWhnuaLieeJuSIsIuiSmeeJueWhnuaLieeJuSIsMSwiIiwiIl0sWyLlrZ/liqDmi4kiLCLlrZ/liqDmi4kiLDEsIiIsIiJdLFsi56eY6bKBIiwi56eY6bKBIiwxLCIiLCIiXSxbIuWvhuWFi+e9l+WwvOilv+S6miIsIuWvhuWFi+e9l+WwvOilv+S6miIsMSwiIiwiIl0sWyLnvIXnlLgiLCLnvIXnlLgiLDEsIiIsIiJdLFsi5pGp5bCU5aSa55OmIiwi5pGp5bCU5aSa55OmIiwxLCIiLCIiXSxbIuaRqea0m+WTpSIsIuaRqea0m+WTpSIsMSwiIiwiIl0sWyLmkannurPlk6UiLCLmkannurPlk6UiLDEsIiIsIiJdLFsi6I6r5qGR5q+U5YWLIiwi6I6r5qGR5q+U5YWLIiwxLCIiLCIiXSxbIuWiqOilv+WTpSIsIuWiqOilv+WTpSIsMSwiIiwiIl0sWyLnurPnsbPmr5TkupoiLCLnurPnsbPmr5TkupoiLDEsIiIsIiJdLFsi5Y2X6Z2eIiwi5Y2X6Z2eIiwxLCIiLCIiXSxbIuWNl+aWr+aLieWkqyIsIuWNl+aWr+aLieWkqyIsMSwiIiwiIl0sWyLnkZnpsoEiLCLnkZnpsoEiLDEsIiIsIiJdLFsi5bC85rOK5bCUIiwi5bC85rOK5bCUIiwxLCIiLCIiXSxbIuWwvOWKoOaLieeTnCIsIuWwvOWKoOaLieeTnCIsMSwiIiwiIl0sWyLlsLzml6XlsJQiLCLlsLzml6XlsJQiLDEsIiIsIiJdLFsi5bC85pel5Yip5LqaIiwi5bC85pel5Yip5LqaIiwxLCIiLCIiXSxbIue6veWfgyIsIue6veWfgyIsMSwiIiwiIl0sWyLmjKrlqIEiLCLmjKrlqIEiLDEsIiIsIiJdLFsi6K+656aP5YWL5bKbIiwi6K+656aP5YWL5bKbIiwxLCIiLCIiXSxbIuW4leWKsyIsIuW4leWKsyIsMSwiIiwiIl0sWyLnmq7nibnlh6/mgannvqTlspsiLCLnmq7nibnlh6/mgannvqTlspsiLDEsIiIsIiJdLFsi6JGh6JCE54mZIiwi6JGh6JCE54mZIiwxLCIiLCIiXSxbIuaXpeacrCIsIuaXpeacrCIsMSwiIiwiIl0sWyLnkZ7lhbgiLCLnkZ7lhbgiLDEsIiIsIiJdLFsi55Ge5aOrIiwi55Ge5aOrIiwxLCIiLCIiXSxbIuiQqOWwlOeTpuWkmiIsIuiQqOWwlOeTpuWkmiIsMSwiIiwiIl0sWyLokKjmkankupoiLCLokKjmkankupoiLDEsIiIsIiJdLFsi5aGe5bCU57u05LqaIiwi5aGe5bCU57u05LqaIiwxLCIiLCIiXSxbIuWhnuaLieWIqeaYgiIsIuWhnuaLieWIqeaYgiIsMSwiIiwiIl0sWyLloZ7lhoXliqDlsJQiLCLloZ7lhoXliqDlsJQiLDEsIiIsIiJdLFsi5aGe5rWm6Lev5pavIiwi5aGe5rWm6Lev5pavIiwxLCIiLCIiXSxbIuWhnuiIjOWwlCIsIuWhnuiIjOWwlCIsMSwiIiwiIl0sWyLmspnnibnpmL/mi4nkvK8iLCLmspnnibnpmL/mi4nkvK8iLDEsIiIsIiJdLFsi5Zyj6K+e5bKbIiwi5Zyj6K+e5bKbIiwxLCIiLCIiXSxbIuWco+Wkmue+juWSjOaZruael+ilv+avlCIsIuWco+Wkmue+juWSjOaZruael+ilv+avlCIsMSwiIiwiIl0sWyLlnKPotavli5Lmi78iLCLlnKPotavli5Lmi78iLDEsIiIsF_STATEIiJdLFsi5Zyj5Z+66Iyo5ZKM5bC857u05pavIiwi5Zyj5Z+66Iyo5ZKM5bC857u05pavIiwxLCIiLCIiXSxbIuWco+WNouilv+S6miIsIuWco+WNouilv+S6miIsMSwiIiwiIl0sWyLlnKPpqazlipvor7oiLCLlnKPpqazlipvor7oiLDEsIiIsIiJdLFsi5Zyj5paH5qOu54m55ZKM5qC85p6X57qz5LiB5pavIiwi5Zyj5paH5qOu54m55ZKM5qC85p6X57qz5LiB5pavIiwxLCIiLCIiXSxbIuaWr+mHjOWFsOWNoSIsIuaWr+mHjOWFsOWNoSIsMSwiIiwiIl0sWyLmlq/mtJvkvJDlhYsiLCLmlq/mtJvkvJDlhYsiLDEsIiIsIiJdLFsi5pav5rSb5paH5bC85LqaIiwi5pav5rSb5paH5bC85LqaIiwxLCIiLCIiXSxbIuaWr+WogeWjq+WFsCIsIuaWr+WogeWjq+WFsCIsMSwiIiwiIl0sWyLoi4/kuLkiLCLoi4/kuLkiLDEsIiIsIiJdLFsi6IuP6YeM5Y2XIiwi6IuP6YeM5Y2XIiwxLCIiLCIiXSxbIuaJgOe9l+mXqOe+pOWymyIsIuaJgOe9l+mXqOe+pOWymyIsMSwiIiwiIl0sWyLntKLpqazph4wiLCLntKLpqazph4wiLDEsIiIsIiJdLFsi5aGU5ZCJ5YWL5pav5Z2mIiwi5aGU5ZCJ5YWL5pav5Z2mIiwxLCIiLCIiXSxbIuazsOWbvSIsIuazsOWbvSIsMSwiIiwiIl0sWyLlnabmoZHlsLzkupoiLCLlnabmoZHlsLzkupoiLDEsIiIsIiJdLFsi5rGk5YqgIiwi5rGk5YqgIiwxLCIiLCIiXSxbIueJueeri+WwvOi+vuWSjOWkmuW3tOWTpSIsIueJueeri+WwvOi+vuWSjOWkmuW3tOWTpSIsMSwiIiwiIl0sWyLnqoHlsLzmlq8iLCLnqoHlsLzmlq8iLDEsIiIsIiJdLFsi5Zu+55Om5Y2iIiwi5Zu+55Om5Y2iIiwxLCIiLCIiXSxbIuWcn+iAs+WFtiIsIuWcn+iAs+WFtiIsMSwiIiwiIl0sWyLlnJ/lupPmm7zmlq/lnaYiLCLlnJ/lupPmm7zmlq/lnaYiLDEsIiIsIiJdLFsi5omY5YWL5YqzIiwi5omY5YWL5YqzIiwxLCIiLCIiXSxbIueTpuWIqeaWr+e+pOWym+WSjOWvjOWbvue6s+e+pOWymyIsIueTpuWIqeaWr+e+pOWym+WSjOWvjOWbvue6s+e+pOWymyIsMSwiIiwiIl0sWyLnk6bliqrpmL/lm74iLCLnk6bliqrpmL/lm74iLDEsIiIsIiJdLFsi5Y2x5Zyw6ams5ouJIiwi5Y2x5Zyw6ams5ouJIiwxLCIiLCIiXSxbIuWnlOWGheeRnuaLiSIsIuWnlOWGheeRnuaLiSIsMSwiIiwiIl0sWyLmlofojrEiLCLmlofojrEiLDEsIiIsIiJdLFsi5LmM5bmy6L6+Iiwi5LmM5bmy6L6+IiwxLCIiLCIiXSxbIuS5jOWFi+WFsCIsIuS5jOWFi+WFsCIsMSwiIiwiIl0sWyLkuYzmi4nlnK0iLCLkuYzmi4nlnK0iLDEsIiIsIiJdLFsi5LmM5YW55Yir5YWL5pav5Z2mIiwi5LmM5YW55Yir5YWL5pav5Z2mIiwxLCIiLCIiXSxbIuilv+ePreeJmSIsIuilv+ePreeJmSIsMSwiIiwiIl0sWyLopb/mkpLlk4jmi4kiLCLopb/mkpLlk4jmi4kiLDEsIiIsIiJdLFsi5biM6IWKIiwi5biM6IWKIiwxLCIiLCIiXSxbIuaWsOWKoOWdoSIsIuaWsOWKoOWdoSIsMSwiIiwiIl0sWyLmlrDlloDph4zlpJrlsLzkupoiLCLmlrDlloDph4zlpJrlsLzkupoiLDEsIiIsIiJdLFsi5paw6KW/5YWwIiwi5paw6KW/5YWwIiwxLCIiLCIiXSxbIuWMiOeJmeWIqSIsIuWMiOeJmeWIqSIsMSwiIiwiIl0sWyLlj5nliKnkupoiLCLlj5nliKnkupoiLDEsIiIsIiJdLFsi54mZ5Lmw5YqgIiwi54mZ5Lmw5YqgIiwxLCIiLCIiXSxbIuS6mue+juWwvOS6miIsIuS6mue+juWwvOS6miIsMSwiIiwiIl0sWyLkuZ/pl6giLCLkuZ/pl6giLDEsIiIsIiJdLFsi5LyK5ouJ5YWLIiwi5LyK5ouJ5YWLIiwxLCIiLCIiXSxbIuS8iuaclyIsIuS8iuaclyIsMSwiIiwiIl0sWyLku6XoibLliJciLCLku6XoibLliJciLDEsIiIsIiJdLFsi5oSP5aSn5YipIiwi5oSP5aSn5YipIiwxLCIiLCIiXSxbIuWNsOW6piIsIuWNsOW6piIsMSwiIiwiIl0sWyLljbDluqblsLzopb/kupoiLCLljbDluqblsLzopb/kupoiLDEsIiIsIiJdLFsi6Iux5Zu9Iiwi6Iux5Zu9IiwxLCIiLCIiXSxbIue6puaXpiIsIue6puaXpiIsMSwiIiwiIl0sWyLotorljZciLCLotorljZciLDEsIiIsIiJdLFsi6LWe5q+U5LqaIiwi6LWe5q+U5LqaIiwxLCIiLCIiXSxbIuazveilv+WymyIsIuazveilv+WymyIsMSwiIiwiIl0sWyLkuY3lvpciLCLkuY3lvpciLDEsIiIsIiJdLFsi55u05biD572X6ZmAIiwi55u05biD572X6ZmAIiwxLCIiLCIiXSxbIuaZuuWIqSIsIuaZuuWIqSIsMSwiIiwiIl0sWyLkuK3pnZ4iLCLkuK3pnZ4iLDEsIiIsIiJdXSwiU2VsZWN0ZWRWYWx1ZUFycmF5IjpbIi0xIl19LCJwMV9TaGlGU0giOnsiSGlkZGVuIjpmYWxzZSwiRl9JdGVtcyI6W1si5pivIiwi5Zyo5LiK5rW3IiwxXSxbIuWQpiIsIuS4jeWcqOS4iua1tyIsMV1dLCJTZWxlY3RlZFZhbHVlIjoi5pivIn0sInAxX1NoaUZaWCI6eyJGX0l0ZW1zIjpbWyLmmK8iLCLkvY/moKEiLDFdLFsi5ZCmIiwi5LiN5L2P5qChIiwxXV0sIlNlbGVjdGVkVmFsdWUiOiLmmK8iLCJIaWRkZW4iOmZhbHNlfSwicDFfZGRsU2hlbmciOnsiRl9JdGVtcyI6W1siLTEiLCLpgInmi6nnnIHku70iLDEsIiIsIiJdLFsi5YyX5LqsIiwi5YyX5LqsIiwxLCIiLCIiXSxbIuWkqea0pSIsIuWkqea0pSIsMSwiIiwiIl0sWyLkuIrmtbciLCLkuIrmtbciLDEsIiIsIiJdLFsi6YeN5bqGIiwi6YeN5bqGIiwxLCIiLCIiXSxbIuays+WMlyIsIuays+WMlyIsMSwiIiwiIl0sWyLlsbHopb8iLCLlsbHopb8iLDEsIiIsIiJdLFsi6L695a6BIiwi6L695a6BIiwxLCIiLCIiXSxbIuWQieaelyIsIuWQieaelyIsMSwiIiwiIl0sWyLpu5HpvpnmsZ8iLCLpu5HpvpnmsZ8iLDEsIiIsIiJdLFsi5rGf6IuPIiwi5rGf6IuPIiwxLCIiLCIiXSxbIua1meaxnyIsIua1meaxnyIsMSwiIiwiIl0sWyLlronlvr0iLCLlronlvr0iLDEsIiIsIiJdLFsi56aP5bu6Iiwi56aP5bu6IiwxLCIiLCIiXSxbIuaxn+ilvyIsIuaxn+ilvyIsMSwiIiwiIl0sWyLlsbHkuJwiLCLlsbHkuJwiLDEsIiIsIiJdLFsi5rKz5Y2XIiwi5rKz5Y2XIiwxLCIiLCIiXSxbIua5luWMlyIsIua5luWMlyIsMSwiIiwiIl0sWyLmuZbljZciLCLmuZbljZciLDEsIiIsIiJdLFsi5bm/5LicIiwi5bm/5LicIiwxLCIiLCIiXSxbIua1t+WNlyIsIua1t+WNlyIsMSwiIiwiIl0sWyLlm5vlt50iLCLlm5vlt50iLDEsIiIsIiJdLFsi6LS15beeIiwi6LS15beeIiwxLCIiLCIiXSxbIuS6keWNlyIsIuS6keWNlyIsMSwiIiwiIl0sWyLpmZXopb8iLCLpmZXopb8iLDEsIiIsIiJdLFsi55SY6IKDIiwi55SY6IKDIiwxLCIiLCIiXSxbIumdkua1tyIsIumdkua1tyIsMSwiIiwiIl0sWyLlhoXokpnlj6QiLCLlhoXokpnlj6QiLDEsIiIsIiJdLFsi5bm/6KW/Iiwi5bm/6KW/IiwxLCIiLCIiXSxbIuilv+iXjyIsIuilv+iXjyIsMSwiIiwiIl0sWyLlroHlpI8iLCLlroHlpI8iLDEsIiIsIiJdLFsi5paw55aGIiwi5paw55aGIiwxLCIiLCIiXSxbIummmea4ryIsIummmea4ryIsMSwiIiwiIl0sWyLmvrPpl6giLCLmvrPpl6giLDEsIiIsIiJdLFsi5Y+w5rm+Iiwi5Y+w5rm+IiwxLCIiLCIiXV0sIlNlbGVjdGVkVmFsdWVBcnJheSI6WyLkuIrmtbciXSwiSGlkZGVuIjpmYWxzZSwiUmVhZG9ubHkiOnRydWV9LCJwMV9kZGxTaGkiOnsiRW5hYmxlZCI6dHJ1ZSwiRl9JdGVtcyI6W1siLTEiLCLpgInmi6nluIIiLDEsIiIsIiJdLFsi5LiK5rW35biCIiwi5LiK5rW35biCIiwxLCIiLCIiXV0sIlNlbGVjdGVkVmFsdWVBcnJheSI6WyLkuIrmtbfluIIiXSwiSGlkZGVuIjpmYWxzZSwiUmVhZG9ubHkiOnRydWV9LCJwMV9kZGxYaWFuIjp7IkVuYWJsZWQiOnRydWUsIkZfSXRlbXMiOltbIi0xIiwi6YCJ5oup5Y6/5Yy6IiwxLCIiLCIiXSxbIum7hOa1puWMuiIsIum7hOa1puWMuiIsMSwiIiwiIl0sWyLljaLmub7ljLoiLCLljaLmub7ljLoiLDEsIiIsIiJdLFsi5b6Q5rGH5Yy6Iiwi5b6Q5rGH5Yy6IiwxLCIiLCIiXSxbIumVv+WugeWMuiIsIumVv+WugeWMuiIsMSwiIiwiIl0sWyLpnZnlronljLoiLCLpnZnlronljLoiLDEsIiIsIiJdLFsi5pmu6ZmA5Yy6Iiwi5pmu6ZmA5Yy6IiwxLCIiLCIiXSxbIuiZueWPo+WMuiIsIuiZueWPo+WMuiIsMSwiIiwiIl0sWyLmnajmtabljLoiLCLmnajmtabljLoiLDEsIiIsIiJdLFsi5a6d5bGx5Yy6Iiwi5a6d5bGx5Yy6IiwxLCIiLCIiXSxbIumXteihjOWMuiIsIumXteihjOWMuiIsMSwiIiwiIl0sWyLlmInlrprljLoiLCLlmInlrprljLoiLDEsIiIsIiJdLFsi5p2+5rGf5Yy6Iiwi5p2+5rGf5Yy6IiwxLCIiLCIiXSxbIumHkeWxseWMuiIsIumHkeWxseWMuiIsMSwiIiwiIl0sWyLpnZLmtabljLoiLCLpnZLmtabljLoiLDEsIiIsIiJdLFsi5aWJ6LSk5Yy6Iiwi5aWJ6LSk5Yy6IiwxLCIiLCIiXSxbIua1puS4nOaWsOWMuiIsIua1puS4nOaWsOWMuiIsMSwiIiwiIl0sWyLltIfmmI7ljLoiLCLltIfmmI7ljLoiLDEsIiIsIiJdXSwiU2VsZWN0ZWRWYWx1ZUFycmF5IjpbIuWuneWxseWMuiJdLCJIaWRkZW4iOmZhbHNlfSwicDFfWGlhbmdYRFoiOnsiVGV4dCI6IuS4iua1t+Wkp+WtpuWuneWxseagoeWMuuagoeWGhXMy5qCLMjA3IiwiSGlkZGVuIjpmYWxzZSwiTGFiZWwiOiLmoKHlhoXlrr/oiI3lnLDlnYDvvIjmoKHljLrjgIHluaLmpbwg44CB5oi/6Ze077yJIn0sInAxX1NoaUZaSiI6eyJSZXF1aXJlZCI6dHJ1ZSwiSGlkZGVuIjpmYWxzZSwiRl9JdGVtcyI6W1si5pivIiwi5a625bqt5Zyw5Z2AIiwxXSxbIuWQpiIsIuS4jeaYr+WutuW6reWcsOWdgCIsMV1dLCJTZWxlY3RlZFZhbHVlIjoi5ZCmIn0sInAxX0NvbnRlbnRQYW5lbDFfWmhvbmdHRlhEUSI6eyJUZXh0IjoiPHNwYW4gc3R5bGU9J2NvbG9yOnJlZDsnPuaXoDwvc3Bhbj4ifSwicDFfQ29udGVudFBhbmVsMSI6eyJJRnJhbWVBdHRyaWJ1dGVzIjp7fX0sInAxX0ZlbmdYRFFETCI6eyJMYWJlbCI6IjAy5pyIMTXml6Xoh7MwM+aciDAx5pel5piv5ZCm5ZyoPHNwYW4gc3R5bGU9J2NvbG9yOnJlZDsnPuS4remrmOmjjumZqeWcsOWMujwvc3Bhbj7pgJfnlZkiLCJTZWxlY3RlZFZhbHVlIjoi5ZCmIiwiRl9JdGVtcyI6W1si5pivIiwi5pivIiwxXSxbIuWQpiIsIuWQpiIsMV1dfSwicDFfVG9uZ1pXRExIIjp7IlJlcXVpcmVkIjpmYWxzZSwiSGlkZGVuIjpmYWxzZSwiTGFiZWwiOiLkuIrmtbflkIzkvY/kurrlkZjmmK/lkKbmnIkwMuaciDE15pel6IezMDPmnIgwMeaXpeadpeiHqjxzcGFuIHN0eWxlPSdjb2xvcjpyZWQ7Jz7kuK3pq5jpo47pmanlnLDljLo8L3NwYW4+55qE5Lq6IiwiRl9JdGVtcyI6W1si5pivIiwi5pivIiwxXSxbIuWQpiIsIuWQpiIsMV1dLCJTZWxlY3RlZFZhbHVlIjoi5ZCmIn0sInAxX0NlbmdGV0giOnsiTGFiZWwiOiIwMuaciDE15pel6IezMDPmnIgwMeaXpeaYr+WQpuWcqDxzcGFuIHN0eWxlPSdjb2xvcjpyZWQ7Jz7kuK3pq5jpo47pmanlnLDljLo8L3NwYW4+6YCX55WZ6L+HIiwiRl9JdGVtcyI6W1si5pivIiwi5pivIiwxXSxbIuWQpiIsIuWQpiIsMV1dLCJTZWxlY3RlZFZhbHVlIjoi5ZCmIn0sInAxX0NlbmdGV0hfUmlRaSI6eyJIaWRkZW4iOnRydWV9LCJwMV9DZW5nRldIX0JlaVpodSI6eyJIaWRkZW4iOnRydWV9LCJwMV9KaWVDaHUiOnsiTGFiZWwiOiIwMuaciDE15pel6IezMDPmnIgwMeaXpeaYr+WQpuS4juadpeiHqjxzcGFuIHN0eWxlPSdjb2xvcjpyZWQ7Jz7kuK3pq5jpo47pmanlnLDljLo8L3NwYW4+5Y+R54Ot5Lq65ZGY5a+G5YiH5o6l6KemIiwiU2VsZWN0ZWRWYWx1ZSI6IuWQpiIsIkZfSXRlbXMiOltbIuaYryIsIuaYryIsMV0sWyLlkKYiLCLlkKYiLDFdXX0sInAxX0ppZUNodV9SaVFpIjp7IkhpZGRlbiI6dHJ1ZX0sInAxX0ppZUNodV9CZWlaaHUiOnsiSGlkZGVuIjp0cnVlfSwicDFfVHVKV0giOnsiTGFiZWwiOiIwMuaciDE15pel6IezMDPmnIgwMeaXpeaYr+WQpuS5mOWdkOWFrOWFseS6pOmAmumAlOW+hDxzcGFuIHN0eWxlPSdjb2xvcjpyZWQ7Jz7kuK3pq5jpo47pmanlnLDljLo8L3NwYW4+IiwiU2VsZWN0ZWRWYWx1ZSI6IuWQpiIsIkZfSXRlbXMiOltbIuaYryIsIuaYryIsMV0sWyLlkKYiLCLlkKYiLDFdXX0sInAxX1R1SldIX1JpUWkiOnsiSGlkZGVuIjp0cnVlfSwicDFfVHVKV0hfQmVpWmh1Ijp7IkhpZGRlbiI6dHJ1ZX0sInAxX1F1ZVpIWkpDIjp7IkZfSXRlbXMiOltbIuaYryIsIuaYryIsMSwiIiwiIl0sWyLlkKYiLCLlkKYiLDEsIiIsIiJdXSwiU2VsZWN0ZWRWYWx1ZUFycmF5IjpbIuWQpiJdfSwicDFfRGFuZ1JHTCI6eyJTZWxlY3RlZFZhbHVlIjoi5ZCmIiwiRl9JdGVtcyI6W1si5pivIiwi5pivIiwxXSxbIuWQpiIsIuWQpiIsMV1dfSwicDFfR2VMU00iOnsiSGlkZGVuIjp0cnVlLCJJRnJhbWVBdHRyaWJ1dGVzIjp7fX0sInAxX0dlTEZTIjp7IlJlcXVpcmVkIjpmYWxzZSwiSGlkZGVuIjp0cnVlLCJGX0l0ZW1zIjpbWyLlsYXlrrbpmpTnprsiLCLlsYXlrrbpmpTnprsiLDFdLFsi6ZuG5Lit6ZqU56a7Iiwi6ZuG5Lit6ZqU56a7IiwxXV0sIlNlbGVjdGVkVmFsdWUiOm51bGx9LCJwMV9HZUxEWiI6eyJIaWRkZW4iOnRydWV9LCJwMV9GYW5YUlEiOnsiSGlkZGVuIjp0cnVlfSwicDFfV2VpRkhZWSI6eyJIaWRkZW4iOnRydWV9LCJwMV9TaGFuZ0hKWkQiOnsiSGlkZGVuIjp0cnVlfSwicDFfRGFvWFFMWUdKIjp7IlRleHQiOiLkuK3lm70ifSwicDFfRGFvWFFMWUNTIjp7IlRleHQiOiLotLXpmLMifSwicDFfSmlhUmVuIjp7IkxhYmVsIjoiMDLmnIgxNeaXpeiHszAz5pyIMDHml6XlrrbkurrmmK/lkKbmnInlj5Hng63nrYnnl4fnirYifSwicDFfSmlhUmVuX0JlaVpodSI6eyJIaWRkZW4iOnRydWV9LCJwMV9TdWlTTSI6eyJSZXF1aXJlZCI6dHJ1ZSwiU2VsZWN0ZWRWYWx1ZSI6Iue7v+iJsiIsIkZfSXRlbXMiOltbIue6ouiJsiIsIue6ouiJsiIsMV0sWyLpu4ToibIiLCLpu4ToibIiLDFdLFsi57u/6ImyIiwi57u/6ImyIiwxXV19LCJwMV9Mdk1hMTREYXlzIjp7IlJlcXVpcmVkIjp0cnVlLCJTZWxlY3RlZFZhbHVlIjoi5pivIiwiRl9JdGVtcyI6W1si5pivIiwi5pivIiwxXSxbIuWQpiIsIuWQpiIsMV1dfSwicDFfY3RsMDBfYnRuUmV0dXJuIjp7Ik9uQ2xpZW50Q2xpY2siOiJkb2N1bWVudC5sb2NhdGlvbi5ocmVmPScvRGVmYXVsdC5hc3B4JztyZXR1cm47In0sInAxIjp7IlRpdGxlIjoi6IOh5b+X5a6P77yIMTYxMjMxMTPvvInnmoTmr4/ml6XkuIDmiqUiLCJJRnJhbWVBdHRyaWJ1dGVzIjp7fX19",
		"F_TARGET":                   "p1_ctl00_btnSubmit",
	}
	d := url.Values{}
	for k, v := range postData {
		d.Set(k, v)
	}
	var data = strings.NewReader(d.Encode())
	return data
}
