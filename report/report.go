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

type ApiReport struct {
	client    *http.Client
	viewstate string
	cookies   string
}

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

func (r *ApiReport) Report(uid, password string) {
	r.Init()
	r.GetCookies(uid, password)
	r.GetViewState()
	r.Submit()
}

func (r *ApiReport) Init() {
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
	fmt.Println("获取账号cookies")
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
	defer resp.Body.Close()
	cookie := resp.Cookies()
	// fmt.Println(cookie)
	for _, v := range cookie {
		fmt.Println(v.Name)
		if v.Name == ".ncov2019selfreport" {
			r.cookies = fmt.Sprintf(".ncov2019selfreport=%s", v.Value)
		}
	}
	htmlBytes, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("%+s\n", htmlBytes)
	fmt.Println(r.cookies)
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
	postData := map[string]string{
		"__EVENTTARGET":              "p1$ctl01$btnSubmit",
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
		"F_STATE":                    "eyJwMV9DaGVuZ051byI6eyJDaGVja2VkIjp0cnVlfSwicDFfQmFvU1JRIjp7IlRleHQiOiIyMDIxLTA0LTE0In0sInAxX0RhbmdRU1RaSyI6eyJGX0l0ZW1zIjpbWyLoia/lpb0iLCLoia/lpb3vvIjkvZPmuKnkuI3pq5jkuo4zNy4z77yJIiwxXSxbIuS4jemAgiIsIuS4jemAgiIsMV1dLCJTZWxlY3RlZFZhbHVlIjoi6Imv5aW9In0sInAxX1poZW5nWmh1YW5nIjp7IkhpZGRlbiI6dHJ1ZSwiRl9JdGVtcyI6W1si5oSf5YaSIiwi5oSf5YaSIiwxXSxbIuWSs+WXvSIsIuWSs+WXvSIsMV0sWyLlj5Hng60iLCLlj5Hng60iLDFdXSwiU2VsZWN0ZWRWYWx1ZUFycmF5IjpbXX0sInAxX1FpdVpaVCI6eyJGX0l0ZW1zIjpbXSwiU2VsZWN0ZWRWYWx1ZUFycmF5IjpbXX0sInAxX0ppdVlLTiI6eyJGX0l0ZW1zIjpbXSwiU2VsZWN0ZWRWYWx1ZUFycmF5IjpbXX0sInAxX0ppdVlZWCI6eyJSZXF1aXJlZCI6ZmFsc2UsIkZfSXRlbXMiOltdLCJTZWxlY3RlZFZhbHVlQXJyYXkiOltdfSwicDFfSml1WVpEIjp7IkZfSXRlbXMiOltdLCJTZWxlY3RlZFZhbHVlQXJyYXkiOltdfSwicDFfSml1WVpMIjp7IkZfSXRlbXMiOltdLCJTZWxlY3RlZFZhbHVlQXJyYXkiOltdfSwicDFfR3VvTmVpIjp7IkZfSXRlbXMiOltbIuWbveWGhSIsIuWbveWGhSIsMV0sWyLlm73lpJYiLCLlm73lpJYiLDFdXSwiU2VsZWN0ZWRWYWx1ZSI6IuWbveWGhSJ9LCJwMV9kZGxHdW9KaWEiOnsiRGF0YVRleHRGaWVsZCI6Ilpob25nV2VuIiwiRGF0YVZhbHVlRmllbGQiOiJaaG9uZ1dlbiIsIkZfSXRlbXMiOltbIi0xIiwi6YCJ5oup5Zu95a62IiwxLCIiLCIiXSxbIumYv+WwlOW3tOWwvOS6miIsIumYv+WwlOW3tOWwvOS6miIsMSwiIiwiIl0sWyLpmL/lsJTlj4rliKnkupoiLCLpmL/lsJTlj4rliKnkupoiLDEsIiIsIiJdLFsi6Zi/5a+M5rGXIiwi6Zi/5a+M5rGXIiwxLCIiLCIiXSxbIumYv+agueW7tyIsIumYv+agueW7tyIsMSwiIiwiIl0sWyLpmL/mi4nkvK/ogZTlkIjphYvplb/lm70iLCLpmL/mi4nkvK/ogZTlkIjphYvplb/lm70iLDEsIiIsIiJdLFsi6Zi/6bKB5be0Iiwi6Zi/6bKB5be0IiwxLCIiLCIiXSxbIumYv+abvCIsIumYv+abvCIsMSwiIiwiIl0sWyLpmL/loZ7mi5znloYiLCLpmL/loZ7mi5znloYiLDEsIiIsIiJdLFsi5Z+D5Y+KIiwi5Z+D5Y+KIiwxLCIiLCIiXSxbIuWfg+WhnuS/hOavlOS6miIsIuWfg+WhnuS/hOavlOS6miIsMSwiIiwiIl0sWyLniLHlsJTlhbAiLCLniLHlsJTlhbAiLDEsIiIsIiJdLFsi54ix5rKZ5bC85LqaIiwi54ix5rKZ5bC85LqaIiwxLCIiLCIiXSxbIuWuiemBk+WwlCIsIuWuiemBk+WwlCIsMSwiIiwiIl0sWyLlronlk6Xmi4kiLCLlronlk6Xmi4kiLDEsIiIsIiJdLFsi5a6J5Zyt5ouJIiwi5a6J5Zyt5ouJIiwxLCIiLCIiXSxbIuWuieaPkOeTnOWSjOW3tOW4g+i+viIsIuWuieaPkOeTnOWSjOW3tOW4g+i+viIsMSwiIiwiIl0sWyLlpaXlnLDliKkiLCLlpaXlnLDliKkiLDEsIiIsIiJdLFsi5aWl5YWw576k5bKbIiwi5aWl5YWw576k5bKbIiwxLCIiLCIiXSxbIua+s+Wkp+WIqeS6miIsIua+s+Wkp+WIqeS6miIsMSwiIiwiIl0sWyLlt7Tlt7TlpJrmlq8iLCLlt7Tlt7TlpJrmlq8iLDEsIiIsIiJdLFsi5be05biD5Lqa5paw5Yeg5YaF5LqaIiwi5be05biD5Lqa5paw5Yeg5YaF5LqaIiwxLCIiLCIiXSxbIuW3tOWTiOmprCIsIuW3tOWTiOmprCIsMSwiIiwiIl0sWyLlt7Tln7rmlq/lnaYiLCLlt7Tln7rmlq/lnaYiLDEsIiIsIiJdLFsi5be05YuS5pav5Z2mIiwi5be05YuS5pav5Z2mIiwxLCIiLCIiXSxbIuW3tOaelyIsIuW3tOaelyIsMSwiIiwiIl0sWyLlt7Tmi7/pqawiLCLlt7Tmi7/pqawiLDEsIiIsIiJdLFsi5be06KW/Iiwi5be06KW/IiwxLCIiLCIiXSxbIueZveS/hOe9l+aWryIsIueZveS/hOe9l+aWryIsMSwiIiwiIl0sWyLnmb7mhZXlpKciLCLnmb7mhZXlpKciLDEsIiIsIiJdLFsi5L+d5Yqg5Yip5LqaIiwi5L+d5Yqg5Yip5LqaIiwxLCIiLCIiXSxbIui0neWugSIsIui0neWugSIsMSwiIiwiIl0sWyLmr5TliKnml7YiLCLmr5TliKnml7YiLDEsIiIsIiJdLFsi5Yaw5bKbIiwi5Yaw5bKbIiwxLCIiLCIiXSxbIuazouWkmum7juWQhCIsIuazouWkmum7juWQhCIsMSwiIiwiIl0sWyLms6LlhbAiLCLms6LlhbAiLDEsIiIsIiJdLFsi5rOi5pav5bC85Lqa5ZKM6buR5aGe5ZOl57u06YKjIiwi5rOi5pav5bC85Lqa5ZKM6buR5aGe5ZOl57u06YKjIiwxLCIiLCIiXSxbIueOu+WIqee7tOS6miIsIueOu+WIqee7tOS6miIsMSwiIiwiIl0sWyLkvK/liKnlhbkiLCLkvK/liKnlhbkiLDEsIiIsIiJdLFsi5Y2a6Iyo55Om57qzIiwi5Y2a6Iyo55Om57qzIiwxLCIiLCIiXSxbIuS4jeS4uSIsIuS4jeS4uSIsMSwiIiwiIl0sWyLluIPln7rnurPms5XntKIiLCLluIPln7rnurPms5XntKIiLDEsIiIsIiJdLFsi5biD6ZqG6L+qIiwi5biD6ZqG6L+qIiwxLCIiLCIiXSxbIuW4g+e7tOWymyIsIuW4g+e7tOWymyIsMSwiIiwiIl0sWyLmnJ3pspwiLCLmnJ3pspwiLDEsIiIsIiJdLFsi6LWk6YGT5Yeg5YaF5LqaIiwi6LWk6YGT5Yeg5YaF5LqaIiwxLCIiLCIiXSxbIuS4uem6piIsIuS4uem6piIsMSwiIiwiIl0sWyLlvrflm70iLCLlvrflm70iLDEsIiIsIiJdLFsi5Lic5bid5rG2Iiwi5Lic5bid5rG2IiwxLCIiLCIiXSxbIuS4nOW4neaxtiIsIuS4nOW4neaxtiIsMSwiIiwiIl0sWyLlpJrlk6UiLCLlpJrlk6UiLDEsIiIsIiJdLFsi5aSa57Gz5bC85YqgIiwi5aSa57Gz5bC85YqgIiwxLCIiLCIiXSxbIuS/hOe9l+aWr+iBlOmCpiIsIuS/hOe9l+aWr+iBlOmCpiIsMSwiIiwiIl0sWyLljoTnk5zlpJrlsJQiLCLljoTnk5zlpJrlsJQiLDEsIiIsIiJdLFsi5Y6E56uL54m56YeM5LqaIiwi5Y6E56uL54m56YeM5LqaIiwxLCIiLCIiXSxbIuazleWbvSIsIuazleWbvSIsMSwiIiwiIl0sWyLms5Xlm73lpKfpg73kvJoiLCLms5Xlm73lpKfpg73kvJoiLDEsIiIsIiJdLFsi5rOV572X576k5bKbIiwi5rOV572X576k5bKbIiwxLCIiLCIiXSxbIuazleWxnuazouWIqeWwvOilv+S6miIsIuazleWxnuazouWIqeWwvOilv+S6miIsMSwiIiwiIl0sWyLms5XlsZ7lnK3kuprpgqMiLCLms5XlsZ7lnK3kuprpgqMiLDEsIiIsIiJdLFsi5qK16JKC5YaIIiwi5qK16JKC5YaIIiwxLCIiLCIiXSxbIuiPsuW+i+WuviIsIuiPsuW+i+WuviIsMSwiIiwiIl0sWyLmlpDmtY4iLCLmlpDmtY4iLDEsIiIsIiJdLFsi6Iqs5YWwIiwi6Iqs5YWwIiwxLCIiLCIiXSxbIuS9m+W+l+inkiIsIuS9m+W+l+inkiIsMSwiIiwiIl0sWyLlhojmr5TkupoiLCLlhojmr5TkupoiLDEsIiIsIiJdLFsi5Yia5p6cIiwi5Yia5p6cIiwxLCIiLCIiXSxbIuWImuaenO+8iOmHke+8iSIsIuWImuaenO+8iOmHke+8iSIsMSwiIiwiIl0sWyLlk6XkvKbmr5TkupoiLCLlk6XkvKbmr5TkupoiLDEsIiIsIiJdLFsi5ZOl5pav6L6+6buO5YqgIiwi5ZOl5pav6L6+6buO5YqgIiwxLCIiLCIiXSxbIuagvOael+e6s+i+viIsIuagvOael+e6s+i+viIsMSwiIiwiIl0sWyLmoLzpsoHlkInkupoiLCLmoLzpsoHlkInkupoiLDEsIiIsIiJdLFsi5qC56KW/5bKbIiwi5qC56KW/5bKbIiwxLCIiLCIiXSxbIuWPpOW3tCIsIuWPpOW3tCIsMSwiIiwiIl0sWyLnk5zlvrfnvZfmma7lspsiLCLnk5zlvrfnvZfmma7lspsiLDEsIiIsIiJdLFsi5YWz5bKbIiwi5YWz5bKbIiwxLCIiLCIiXSxbIuWcreS6mumCoyIsIuWcreS6mumCoyIsMSwiIiwiIl0sWyLlk4jokKjlhYvmlq/lnaYiLCLlk4jokKjlhYvmlq/lnaYiLDEsIiIsIiJdLFsi5rW35ZywIiwi5rW35ZywIiwxLCIiLCIiXSxbIumfqeWbvSIsIumfqeWbvSIsMSwiIiwiIl0sWyLojbflhbAiLCLojbflhbAiLDEsIiIsIiJdLFsi6buR5bGxIiwi6buR5bGxIiwxLCIiLCIiXSxbIua0qumDveaLieaWryIsIua0qumDveaLieaWryIsMSwiIiwiIl0sWyLln7rph4zlt7Tmlq8iLCLln7rph4zlt7Tmlq8iLDEsIiIsIiJdLFsi5ZCJ5biD5o+QIiwi5ZCJ5biD5o+QIiwxLCIiLCIiXSxbIuWQieWwlOWQieaWr+aWr+WdpiIsIuWQieWwlOWQieaWr+aWr+WdpiIsMSwiIiwiIl0sWyLlh6DlhoXkupoiLCLlh6DlhoXkupoiLDEsIiIsIiJdLFsi5Yeg5YaF5Lqa5q+U57uNIiwi5Yeg5YaF5Lqa5q+U57uNIiwxLCIiLCIiXSxbIuWKoOaLv+WkpyIsIuWKoOaLv+WkpyIsMSwiIiwiIl0sWyLliqDnurMiLCLliqDnurMiLDEsIiIsIiJdLFsi5Yqg6JOsIiwi5Yqg6JOsIiwxLCIiLCIiXSxbIuafrOWflOWvqCIsIuafrOWflOWvqCIsMSwiIiwiIl0sWyLmjbflhYsiLCLmjbflhYsiLDEsIiIsIiJdLFsi5rSl5be05biD6Z+mIiwi5rSl5be05biD6Z+mIiwxLCIiLCIiXSxbIuWWgOm6pumahiIsIuWWgOm6pumahiIsMSwiIiwiIl0sWyLljaHloZTlsJQiLCLljaHloZTlsJQiLDEsIiIsIiJdLFsi56eR56eR5pav77yI5Z+65p6X77yJ576k5bKbIiwi56eR56eR5pav77yI5Z+65p6X77yJ576k5bKbIiwxLCIiLCIiXSxbIuenkeaRqee9lyIsIuenkeaRqee9lyIsMSwiIiwiIl0sWyLnp5Hnibnov6rnk6YiLCLnp5Hnibnov6rnk6YiLDEsIiIsIiJdLFsi56eR5aiB54m5Iiwi56eR5aiB54m5IiwxLCIiLCIiXSxbIuWFi+e9l+WcsOS6miIsIuWFi+e9l+WcsOS6miIsMSwiIiwiIl0sWyLogq/lsLzkupoiLCLogq/lsLzkupoiLDEsIiIsIiJdLFsi5bqT5YWL576k5bKbIiwi5bqT5YWL576k5bKbIiwxLCIiLCIiXSxbIuaLieiEsee7tOS6miIsIuaLieiEsee7tOS6miIsMSwiIiwiIl0sWyLojrHntKLmiZgiLCLojrHntKLmiZgiLDEsIiIsIiJdLFsi6ICB5oydIiwi6ICB5oydIiwxLCIiLCIiXSxbIum7juW3tOWrqSIsIum7juW3tOWrqSIsMSwiIiwiIl0sWyLnq4vpmbblrpsiLCLnq4vpmbblrpsiLDEsIiIsIiJdLFsi5Yip5q+U6YeM5LqaIiwi5Yip5q+U6YeM5LqaIiwxLCIiLCIiXSxbIuWIqeavlOS6miIsIuWIqeavlOS6miIsMSwiIiwiIl0sWyLliJfmlK/mlablo6vnmbsiLCLliJfmlK/mlablo6vnmbsiLDEsIiIsIiJdLFsi55WZ5bC85rGq5bKbIiwi55WZ5bC85rGq5bKbIiwxLCIiLCIiXSxbIuWNouajruWgoSIsIuWNouajruWgoSIsMSwiIiwiIl0sWyLljaLml7rovr4iLCLljaLml7rovr4iLDEsIiIsIiJdLFsi572X6ams5bC85LqaIiwi572X6ams5bC85LqaIiwxLCIiLCIiXSxbIumprOi+vuWKoOaWr+WKoCIsIumprOi+vuWKoOaWr+WKoCIsMSwiIiwiIl0sWyLpqazmganlspsiLCLpqazmganlspsiLDEsIiIsIiJdLFsi6ams5bCU5Luj5aSrIiwi6ams5bCU5Luj5aSrIiwxLCIiLCIiXSxbIumprOiAs+S7liIsIumprOiAs+S7liIsMSwiIiwiIl0sWyLpqazmi4nnu7QiLCLpqazmi4nnu7QiLDEsIiIsIiJdLFsi6ams5p2l6KW/5LqaIiwi6ams5p2l6KW/5LqaIiwxLCIiLCIiXSxbIumprOmHjCIsIumprOmHjCIsMSwiIiwiIl0sWyLpqazlhbbpob8iLCLpqazlhbbpob8iLDEsIiIsIiJdLFsi6ams57uN5bCU576k5bKbIiwi6ams57uN5bCU576k5bKbIiwxLCIiLCIiXSxbIumprOaPkOWwvOWFi+WymyIsIumprOaPkOWwvOWFi+WymyIsMSwiIiwiIl0sWyLpqaznuqbnibkiLCLpqaznuqbnibkiLDEsIiIsIiJdLFsi5q+b6YeM5rGC5pavIiwi5q+b6YeM5rGC5pavIiwxLCIiLCIiXSxbIuavm+mHjOWhlOWwvOS6miIsIuavm+mHjOWhlOWwvOS6miIsMSwiIiwiIl0sWyLnvo7lm70iLCLnvo7lm70iLDEsIiIsIiJdLFsi576O5bGe6JCo5pGp5LqaIiwi576O5bGe6JCo5pGp5LqaIiwxLCIiLCIiXSxbIuiSmeWPpCIsIuiSmeWPpCIsMSwiIiwiIl0sWyLokpnnibnloZ7mi4nnibkiLCLokpnnibnloZ7mi4nnibkiLDEsIiIsIiJdLFsi5a2f5Yqg5ouJIiwi5a2f5Yqg5ouJIiwxLCIiLCIiXSxbIuenmOmygSIsIuenmOmygSIsMSwiIiwiIl0sWyLlr4blhYvnvZflsLzopb/kupoiLCLlr4blhYvnvZflsLzopb/kupoiLDEsIiIsIiJdLFsi57yF55S4Iiwi57yF55S4IiwxLCIiLCIiXSxbIuaRqeWwlOWkmueTpiIsIuaRqeWwlOWkmueTpiIsMSwiIiwiIl0sWyLmkanmtJvlk6UiLCLmkanmtJvlk6UiLDEsIiIsIiJdLFsi5pGp57qz5ZOlIiwi5pGp57qz5ZOlIiwxLCIiLCIiXSxbIuiOq+ahkeavlOWFiyIsIuiOq+ahkeavlOWFiyIsMSwiIiwiIl0sWyLloqjopb/lk6UiLCLloqjopb/lk6UiLDEsIiIsIiJdLFsi57qz57Gz5q+U5LqaIiwi57qz57Gz5q+U5LqaIiwxLCIiLCIiXSxbIuWNl+mdniIsIuWNl+mdniIsMSwiIiwiIl0sWyLljZfmlq/mi4nlpKsiLCLljZfmlq/mi4nlpKsiLDEsIiIsIiJdLFsi55GZ6bKBIiwi55GZ6bKBIiwxLCIiLCIiXSxbIuWwvOaziuWwlCIsIuWwvOaziuWwlCIsMSwiIiwiIl0sWyLlsLzliqDmi4nnk5wiLCLlsLzliqDmi4nnk5wiLDEsIiIsIiJdLFsi5bC85pel5bCUIiwi5bC85pel5bCUIiwxLCIiLCIiXSxbIuWwvOaXpeWIqeS6miIsIuWwvOaXpeWIqeS6miIsMSwiIiwiIl0sWyLnur3ln4MiLCLnur3ln4MiLDEsIiIsIiJdLFsi5oyq5aiBIiwi5oyq5aiBIiwxLCIiLCIiXSxbIuivuuemj+WFi+WymyIsIuivuuemj+WFi+WymyIsMSwiIiwiIl0sWyLluJXlirMiLCLluJXlirMiLDEsIiIsIiJdLFsi55qu54m55Yev5oGp576k5bKbIiwi55qu54m55Yev5oGp576k5bKbIiwxLCIiLCIiXSxbIuiRoeiQhOeJmSIsIuiRoeiQhOeJmSIsMSwiIiwiIl0sWyLml6XmnKwiLCLml6XmnKwiLDEsIiIsIiJdLFsi55Ge5YW4Iiwi55Ge5YW4IiwxLCIiLCIiXSxbIueRnuWjqyIsIueRnuWjqyIsMSwiIiwiIl0sWyLokKjlsJTnk6blpJoiLCLokKjlsJTnk6blpJoiLDEsIiIsIiJdLFsi6JCo5pGp5LqaIiwi6JCo5pGp5LqaIiwxLCIiLCIiXSxbIuWhnuWwlOe7tOS6miIsIuWhnuWwlOe7tOS6miIsMSwiIiwiIl0sWyLloZ7mi4nliKnmmIIiLCLloZ7mi4nliKnmmIIiLDEsIiIsIiJdLFsi5aGe5YaF5Yqg5bCUIiwi5aGe5YaF5Yqg5bCUIiwxLCIiLCIiXSxbIuWhnua1pui3r+aWryIsIuWhnua1pui3r+aWryIsMSwiIiwiIl0sWyLloZ7oiIzlsJQiLCLloZ7oiIzlsJQiLDEsIiIsIiJdLFsi5rKZ54m56Zi/5ouJ5LyvIiwi5rKZ54m56Zi/5ouJ5LyvIiwxLCIiLCIiXSxbIuWco+ivnuWymyIsIuWco+ivnuWymyIsMSwiIiwiIl0sWyLlnKPlpJrnvo7lkozmma7mnpfopb/mr5QiLCLlnKPlpJrnvo7lkozmma7mnpfopb/mr5QiLDEsIiIsIiJdLFsi5Zyj6LWr5YuS5ou/Iiwi5Zyj6LWr5YuS5ou/IiwxLCIiLCIiXSxbIuWco+WfuuiMqOWSjOWwvOe7tOaWryIsIuWco+WfuuiMqOWSjOWwvOe7tOaWryIsMSwiIiwiIl0sWyLlnKPljaLopb/kupoiLCLlnKPljaLopb/kupoiLDEsIiIsIiJdLFsi5Zyj6ams5Yqb6K+6F_STATEIiwi5Zyj6ams5Yqb6K+6IiwxLCIiLCIiXSxbIuWco+aWh+ajrueJueWSjOagvOael+e6s+S4geaWryIsIuWco+aWh+ajrueJueWSjOagvOael+e6s+S4geaWryIsMSwiIiwiIl0sWyLmlq/ph4zlhbDljaEiLCLmlq/ph4zlhbDljaEiLDEsIiIsIiJdLFsi5pav5rSb5LyQ5YWLIiwi5pav5rSb5LyQ5YWLIiwxLCIiLCIiXSxbIuaWr+a0m+aWh+WwvOS6miIsIuaWr+a0m+aWh+WwvOS6miIsMSwiIiwiIl0sWyLmlq/lqIHlo6vlhbAiLCLmlq/lqIHlo6vlhbAiLDEsIiIsIiJdLFsi6IuP5Li5Iiwi6IuP5Li5IiwxLCIiLCIiXSxbIuiLj+mHjOWNlyIsIuiLj+mHjOWNlyIsMSwiIiwiIl0sWyLmiYDnvZfpl6jnvqTlspsiLCLmiYDnvZfpl6jnvqTlspsiLDEsIiIsIiJdLFsi57Si6ams6YeMIiwi57Si6ams6YeMIiwxLCIiLCIiXSxbIuWhlOWQieWFi+aWr+WdpiIsIuWhlOWQieWFi+aWr+WdpiIsMSwiIiwiIl0sWyLms7Dlm70iLCLms7Dlm70iLDEsIiIsIiJdLFsi5Z2m5qGR5bC85LqaIiwi5Z2m5qGR5bC85LqaIiwxLCIiLCIiXSxbIuaxpOWKoCIsIuaxpOWKoCIsMSwiIiwiIl0sWyLnibnnq4vlsLzovr7lkozlpJrlt7Tlk6UiLCLnibnnq4vlsLzovr7lkozlpJrlt7Tlk6UiLDEsIiIsIiJdLFsi56qB5bC85pavIiwi56qB5bC85pavIiwxLCIiLCIiXSxbIuWbvueTpuWNoiIsIuWbvueTpuWNoiIsMSwiIiwiIl0sWyLlnJ/ogLPlhbYiLCLlnJ/ogLPlhbYiLDEsIiIsIiJdLFsi5Zyf5bqT5pu85pav5Z2mIiwi5Zyf5bqT5pu85pav5Z2mIiwxLCIiLCIiXSxbIuaJmOWFi+WKsyIsIuaJmOWFi+WKsyIsMSwiIiwiIl0sWyLnk6bliKnmlq/nvqTlspvlkozlr4zlm77nurPnvqTlspsiLCLnk6bliKnmlq/nvqTlspvlkozlr4zlm77nurPnvqTlspsiLDEsIiIsIiJdLFsi55Om5Yqq6Zi/5Zu+Iiwi55Om5Yqq6Zi/5Zu+IiwxLCIiLCIiXSxbIuWNseWcsOmprOaLiSIsIuWNseWcsOmprOaLiSIsMSwiIiwiIl0sWyLlp5TlhoXnkZ7mi4kiLCLlp5TlhoXnkZ7mi4kiLDEsIiIsIiJdLFsi5paH6I6xIiwi5paH6I6xIiwxLCIiLCIiXSxbIuS5jOW5sui+viIsIuS5jOW5sui+viIsMSwiIiwiIl0sWyLkuYzlhYvlhbAiLCLkuYzlhYvlhbAiLDEsIiIsIiJdLFsi5LmM5ouJ5ZytIiwi5LmM5ouJ5ZytIiwxLCIiLCIiXSxbIuS5jOWFueWIq+WFi+aWr+WdpiIsIuS5jOWFueWIq+WFi+aWr+WdpiIsMSwiIiwiIl0sWyLopb/nj63niZkiLCLopb/nj63niZkiLDEsIiIsIiJdLFsi6KW/5pKS5ZOI5ouJIiwi6KW/5pKS5ZOI5ouJIiwxLCIiLCIiXSxbIuW4jOiFiiIsIuW4jOiFiiIsMSwiIiwiIl0sWyLmlrDliqDlnaEiLCLmlrDliqDlnaEiLDEsIiIsIiJdLFsi5paw5ZaA6YeM5aSa5bC85LqaIiwi5paw5ZaA6YeM5aSa5bC85LqaIiwxLCIiLCIiXSxbIuaWsOilv+WFsCIsIuaWsOilv+WFsCIsMSwiIiwiIl0sWyLljIjniZnliKkiLCLljIjniZnliKkiLDEsIiIsIiJdLFsi5Y+Z5Yip5LqaIiwi5Y+Z5Yip5LqaIiwxLCIiLCIiXSxbIueJmeS5sOWKoCIsIueJmeS5sOWKoCIsMSwiIiwiIl0sWyLkuprnvo7lsLzkupoiLCLkuprnvo7lsLzkupoiLDEsIiIsIiJdLFsi5Lmf6ZeoIiwi5Lmf6ZeoIiwxLCIiLCIiXSxbIuS8iuaLieWFiyIsIuS8iuaLieWFiyIsMSwiIiwiIl0sWyLkvIrmnJciLCLkvIrmnJciLDEsIiIsIiJdLFsi5Lul6Imy5YiXIiwi5Lul6Imy5YiXIiwxLCIiLCIiXSxbIuaEj+Wkp+WIqSIsIuaEj+Wkp+WIqSIsMSwiIiwiIl0sWyLljbDluqYiLCLljbDluqYiLDEsIiIsIiJdLFsi5Y2w5bqm5bC86KW/5LqaIiwi5Y2w5bqm5bC86KW/5LqaIiwxLCIiLCIiXSxbIuiLseWbvSIsIuiLseWbvSIsMSwiIiwiIl0sWyLnuqbml6YiLCLnuqbml6YiLDEsIiIsIiJdLFsi6LaK5Y2XIiwi6LaK5Y2XIiwxLCIiLCIiXSxbIui1nuavlOS6miIsIui1nuavlOS6miIsMSwiIiwiIl0sWyLms73opb/lspsiLCLms73opb/lspsiLDEsIiIsIiJdLFsi5LmN5b6XIiwi5LmN5b6XIiwxLCIiLCIiXSxbIuebtOW4g+e9l+mZgCIsIuebtOW4g+e9l+mZgCIsMSwiIiwiIl0sWyLmmbrliKkiLCLmmbrliKkiLDEsIiIsIiJdLFsi5Lit6Z2eIiwi5Lit6Z2eIiwxLCIiLCIiXV0sIlNlbGVjdGVkVmFsdWVBcnJheSI6WyItMSJdfSwicDFfU2hpRlNIIjp7IkhpZGRlbiI6ZmFsc2UsIkZfSXRlbXMiOltbIuaYryIsIuWcqOS4iua1tyIsMV0sWyLlkKYiLCLkuI3lnKjkuIrmtbciLDFdXSwiU2VsZWN0ZWRWYWx1ZSI6IuaYryJ9LCJwMV9TaGlGWlgiOnsiRl9JdGVtcyI6W1si5pivIiwi5L2P5qChIiwxXSxbIuWQpiIsIuS4jeS9j+agoSIsMV1dLCJTZWxlY3RlZFZhbHVlIjoi5pivIiwiSGlkZGVuIjpmYWxzZX0sInAxX2RkbFNoZW5nIjp7IkZfSXRlbXMiOltbIi0xIiwi6YCJ5oup55yB5Lu9IiwxLCIiLCIiXSxbIuWMl+S6rCIsIuWMl+S6rCIsMSwiIiwiIl0sWyLlpKnmtKUiLCLlpKnmtKUiLDEsIiIsIiJdLFsi5LiK5rW3Iiwi5LiK5rW3IiwxLCIiLCIiXSxbIumHjeW6hiIsIumHjeW6hiIsMSwiIiwiIl0sWyLmsrPljJciLCLmsrPljJciLDEsIiIsIiJdLFsi5bGx6KW/Iiwi5bGx6KW/IiwxLCIiLCIiXSxbIui+veWugSIsIui+veWugSIsMSwiIiwiIl0sWyLlkInmnpciLCLlkInmnpciLDEsIiIsIiJdLFsi6buR6b6Z5rGfIiwi6buR6b6Z5rGfIiwxLCIiLCIiXSxbIuaxn+iLjyIsIuaxn+iLjyIsMSwiIiwiIl0sWyLmtZnmsZ8iLCLmtZnmsZ8iLDEsIiIsIiJdLFsi5a6J5b69Iiwi5a6J5b69IiwxLCIiLCIiXSxbIuemj+W7uiIsIuemj+W7uiIsMSwiIiwiIl0sWyLmsZ/opb8iLCLmsZ/opb8iLDEsIiIsIiJdLFsi5bGx5LicIiwi5bGx5LicIiwxLCIiLCIiXSxbIuays+WNlyIsIuays+WNlyIsMSwiIiwiIl0sWyLmuZbljJciLCLmuZbljJciLDEsIiIsIiJdLFsi5rmW5Y2XIiwi5rmW5Y2XIiwxLCIiLCIiXSxbIuW5v+S4nCIsIuW5v+S4nCIsMSwiIiwiIl0sWyLmtbfljZciLCLmtbfljZciLDEsIiIsIiJdLFsi5Zub5bedIiwi5Zub5bedIiwxLCIiLCIiXSxbIui0teW3niIsIui0teW3niIsMSwiIiwiIl0sWyLkupHljZciLCLkupHljZciLDEsIiIsIiJdLFsi6ZmV6KW/Iiwi6ZmV6KW/IiwxLCIiLCIiXSxbIueUmOiCgyIsIueUmOiCgyIsMSwiIiwiIl0sWyLpnZLmtbciLCLpnZLmtbciLDEsIiIsIiJdLFsi5YaF6JKZ5Y+kIiwi5YaF6JKZ5Y+kIiwxLCIiLCIiXSxbIuW5v+ilvyIsIuW5v+ilvyIsMSwiIiwiIl0sWyLopb/ol48iLCLopb/ol48iLDEsIiIsIiJdLFsi5a6B5aSPIiwi5a6B5aSPIiwxLCIiLCIiXSxbIuaWsOeWhiIsIuaWsOeWhiIsMSwiIiwiIl0sWyLpppnmuK8iLCLpppnmuK8iLDEsIiIsIiJdLFsi5r6z6ZeoIiwi5r6z6ZeoIiwxLCIiLCIiXSxbIuWPsOa5viIsIuWPsOa5viIsMSwiIiwiIl1dLCJTZWxlY3RlZFZhbHVlQXJyYXkiOlsi5LiK5rW3Il0sIkhpZGRlbiI6ZmFsc2UsIlJlYWRvbmx5Ijp0cnVlfSwicDFfZGRsU2hpIjp7IkVuYWJsZWQiOnRydWUsIkZfSXRlbXMiOltbIi0xIiwi6YCJ5oup5biCIiwxLCIiLCIiXSxbIuS4iua1t+W4giIsIuS4iua1t+W4giIsMSwiIiwiIl1dLCJTZWxlY3RlZFZhbHVlQXJyYXkiOlsi5LiK5rW35biCIl0sIkhpZGRlbiI6ZmFsc2UsIlJlYWRvbmx5Ijp0cnVlfSwicDFfZGRsWGlhbiI6eyJFbmFibGVkIjp0cnVlLCJGX0l0ZW1zIjpbWyItMSIsIumAieaLqeWOv+WMuiIsMSwiIiwiIl0sWyLpu4TmtabljLoiLCLpu4TmtabljLoiLDEsIiIsIiJdLFsi5Y2i5rm+5Yy6Iiwi5Y2i5rm+5Yy6IiwxLCIiLCIiXSxbIuW+kOaxh+WMuiIsIuW+kOaxh+WMuiIsMSwiIiwiIl0sWyLplb/lroHljLoiLCLplb/lroHljLoiLDEsIiIsIiJdLFsi6Z2Z5a6J5Yy6Iiwi6Z2Z5a6J5Yy6IiwxLCIiLCIiXSxbIuaZrumZgOWMuiIsIuaZrumZgOWMuiIsMSwiIiwiIl0sWyLombnlj6PljLoiLCLombnlj6PljLoiLDEsIiIsIiJdLFsi5p2o5rWm5Yy6Iiwi5p2o5rWm5Yy6IiwxLCIiLCIiXSxbIuWuneWxseWMuiIsIuWuneWxseWMuiIsMSwiIiwiIl0sWyLpl7XooYzljLoiLCLpl7XooYzljLoiLDEsIiIsIiJdLFsi5ZiJ5a6a5Yy6Iiwi5ZiJ5a6a5Yy6IiwxLCIiLCIiXSxbIuadvuaxn+WMuiIsIuadvuaxn+WMuiIsMSwiIiwiIl0sWyLph5HlsbHljLoiLCLph5HlsbHljLoiLDEsIiIsIiJdLFsi6Z2S5rWm5Yy6Iiwi6Z2S5rWm5Yy6IiwxLCIiLCIiXSxbIuWliei0pOWMuiIsIuWliei0pOWMuiIsMSwiIiwiIl0sWyLmtabkuJzmlrDljLoiLCLmtabkuJzmlrDljLoiLDEsIiIsIiJdLFsi5bSH5piO5Yy6Iiwi5bSH5piO5Yy6IiwxLCIiLCIiXV0sIlNlbGVjdGVkVmFsdWVBcnJheSI6WyLlrp3lsbHljLoiXSwiSGlkZGVuIjpmYWxzZX0sInAxX1hpYW5nWERaIjp7IlRleHQiOiLkuIrmtbflpKflrablrp3lsbHmoKHljLrmoKHlhoVzMuagizIwNyIsIkhpZGRlbiI6ZmFsc2UsIkxhYmVsIjoi5qCh5YaF5a6/6IiN5Zyw5Z2A77yI5qCh5Yy644CB5bmi5qW8IOOAgeaIv+mXtO+8iSJ9LCJwMV9TaGlGWkoiOnsiUmVxdWlyZWQiOnRydWUsIkhpZGRlbiI6ZmFsc2UsIkZfSXRlbXMiOltbIuaYryIsIuWutuW6reWcsOWdgCIsMV0sWyLlkKYiLCLkuI3mmK/lrrbluq3lnLDlnYAiLDFdXSwiU2VsZWN0ZWRWYWx1ZSI6IuaYryJ9LCJwMV9Db250ZW50UGFuZWwxX1pob25nR0ZYRFEiOnsiVGV4dCI6IjxzcGFuIHN0eWxlPSdjb2xvcjpyZWQ7Jz7kupHljZfnnIHnkZ7kuL3luILlp5DlkYrlm73pl6jnpL7ljLrvvIzlm6Lnu5PmnZHlp5TkvJrph5HlnY7jgIHlvITllorniYfljLrvvIjnkZ7kuL3lpKfpgZPku6XljZfvvInvvIzku5nlrqLlt7flkozlhYnmmI7lt7flsYXmsJHlsI/nu4TvvIzpkavnm5vml7bku6PkvbPlm63lsI/ljLrvvIznkZ7kuqzot6/nuqLnoJbljoLvvIzmmJ/msrPok53mub7lsI/ljLrvvIzlj4zlja/mnZHmsJHlsI/nu4TvvIzkuIvlvITlronmnZHmsJHlsI/nu4TvvIznj6Dlrp3ooZfogIHpo5/lk4HljoLlrrblsZ7ljLo8L3NwYW4+In0sInAxX0NvbnRlbnRQYW5lbDEiOnsiSUZyYW1lQXR0cmlidXRlcyI6e319LCJwMV9GZW5nWERRREwiOnsiTGFiZWwiOiIwM+aciDMx5pel6IezMDTmnIgxNOaXpeaYr+WQpuWcqDxzcGFuIHN0eWxlPSdjb2xvcjpyZWQ7Jz7kuK3pq5jpo47pmanlnLDljLo8L3NwYW4+6YCX55WZIiwiU2VsZWN0ZWRWYWx1ZSI6IuWQpiIsIkZfSXRlbXMiOltbIuaYryIsIuaYryIsMV0sWyLlkKYiLCLlkKYiLDFdXX0sInAxX1RvbmdaV0RMSCI6eyJSZXF1aXJlZCI6dHJ1ZSwiTGFiZWwiOiLkuIrmtbflkIzkvY/kurrlkZjmmK/lkKbmnIkwM+aciDMx5pel6IezMDTmnIgxNOaXpeadpeiHqjxzcGFuIHN0eWxlPSdjb2xvcjpyZWQ7Jz7kuK3pq5jpo47pmanlnLDljLo8L3NwYW4+55qE5Lq6IiwiU2VsZWN0ZWRWYWx1ZSI6IuWQpiIsIkZfSXRlbXMiOltbIuaYryIsIuaYryIsMV0sWyLlkKYiLCLlkKYiLDFdXX0sInAxX0NlbmdGV0giOnsiTGFiZWwiOiIwM+aciDMx5pel6IezMDTmnIgxNOaXpeaYr+WQpuWcqDxzcGFuIHN0eWxlPSdjb2xvcjpyZWQ7Jz7kuK3pq5jpo47pmanlnLDljLo8L3NwYW4+6YCX55WZ6L+HIiwiRl9JdGVtcyI6W1si5pivIiwi5pivIiwxXSxbIuWQpiIsIuWQpiIsMV1dLCJTZWxlY3RlZFZhbHVlIjoi5ZCmIn0sInAxX0NlbmdGV0hfUmlRaSI6eyJIaWRkZW4iOnRydWV9LCJwMV9DZW5nRldIX0JlaVpodSI6eyJIaWRkZW4iOnRydWV9LCJwMV9KaWVDaHUiOnsiTGFiZWwiOiIwM+aciDMx5pel6IezMDTmnIgxNOaXpeaYr+WQpuS4juadpeiHqjxzcGFuIHN0eWxlPSdjb2xvcjpyZWQ7Jz7kuK3pq5jpo47pmanlnLDljLo8L3NwYW4+5Y+R54Ot5Lq65ZGY5a+G5YiH5o6l6KemIiwiU2VsZWN0ZWRWYWx1ZSI6IuWQpiIsIkZfSXRlbXMiOltbIuaYryIsIuaYryIsMV0sWyLlkKYiLCLlkKYiLDFdXX0sInAxX0ppZUNodV9SaVFpIjp7IkhpZGRlbiI6dHJ1ZX0sInAxX0ppZUNodV9CZWlaaHUiOnsiSGlkZGVuIjp0cnVlfSwicDFfVHVKV0giOnsiTGFiZWwiOiIwM+aciDMx5pel6IezMDTmnIgxNOaXpeaYr+WQpuS5mOWdkOWFrOWFseS6pOmAmumAlOW+hDxzcGFuIHN0eWxlPSdjb2xvcjpyZWQ7Jz7kuK3pq5jpo47pmanlnLDljLo8L3NwYW4+IiwiU2VsZWN0ZWRWYWx1ZSI6IuWQpiIsIkZfSXRlbXMiOltbIuaYryIsIuaYryIsMV0sWyLlkKYiLCLlkKYiLDFdXX0sInAxX1R1SldIX1JpUWkiOnsiSGlkZGVuIjp0cnVlfSwicDFfVHVKV0hfQmVpWmh1Ijp7IkhpZGRlbiI6dHJ1ZX0sInAxX1F1ZVpIWkpDIjp7IkZfSXRlbXMiOltbIuaYryIsIuaYryIsMSwiIiwiIl0sWyLlkKYiLCLlkKYiLDEsIiIsIiJdXSwiU2VsZWN0ZWRWYWx1ZUFycmF5IjpbIuWQpiJdfSwicDFfRGFuZ1JHTCI6eyJTZWxlY3RlZFZhbHVlIjoi5ZCmIiwiRl9JdGVtcyI6W1si5pivIiwi5pivIiwxXSxbIuWQpiIsIuWQpiIsMV1dfSwicDFfR2VMU00iOnsiSGlkZGVuIjp0cnVlLCJJRnJhbWVBdHRyaWJ1dGVzIjp7fX0sInAxX0dlTEZTIjp7IlJlcXVpcmVkIjpmYWxzZSwiSGlkZGVuIjp0cnVlLCJGX0l0ZW1zIjpbWyLlsYXlrrbpmpTnprsiLCLlsYXlrrbpmpTnprsiLDFdLFsi6ZuG5Lit6ZqU56a7Iiwi6ZuG5Lit6ZqU56a7IiwxXV0sIlNlbGVjdGVkVmFsdWUiOm51bGx9LCJwMV9HZUxEWiI6eyJIaWRkZW4iOnRydWV9LCJwMV9GYW5YUlEiOnsiSGlkZGVuIjp0cnVlfSwicDFfV2VpRkhZWSI6eyJIaWRkZW4iOnRydWV9LCJwMV9TaGFuZ0hKWkQiOnsiSGlkZGVuIjp0cnVlfSwicDFfRGFvWFFMWUdKIjp7IlRleHQiOiLkuK3lm70ifSwicDFfRGFvWFFMWUNTIjp7IlRleHQiOiLotLXpmLMifSwicDFfSmlhUmVuIjp7IkxhYmVsIjoiMDPmnIgzMeaXpeiHszA05pyIMTTml6XlrrbkurrmmK/lkKbmnInlj5Hng63nrYnnl4fnirYifSwicDFfSmlhUmVuX0JlaVpodSI6eyJIaWRkZW4iOnRydWV9LCJwMV9TdWlTTSI6eyJSZXF1aXJlZCI6dHJ1ZSwiU2VsZWN0ZWRWYWx1ZSI6Iue7v+iJsiIsIkZfSXRlbXMiOltbIue6ouiJsiIsIue6ouiJsiIsMV0sWyLpu4ToibIiLCLpu4ToibIiLDFdLFsi57u/6ImyIiwi57u/6ImyIiwxXV19LCJwMV9Mdk1hMTREYXlzIjp7IlJlcXVpcmVkIjp0cnVlLCJTZWxlY3RlZFZhbHVlIjoi5pivIiwiRl9JdGVtcyI6W1si5pivIiwi5pivIiwxXSxbIuWQpiIsIuWQpiIsMV1dfSwicDFfY3RsMDFfYnRuUmV0dXJuIjp7Ik9uQ2xpZW50Q2xpY2siOiJkb2N1bWVudC5sb2NhdGlvbi5ocmVmPScvRGVmYXVsdC5hc3B4JztyZXR1cm47In0sInAxIjp7IlRpdGxlIjoi6IOh5b+X5a6P77yIMTYxMjMxMTPvvInnmoTmr4/ml6XkuIDmiqUiLCJJRnJhbWVBdHRyaWJ1dGVzIjp7fX19",
		"F_TARGET":                   "p1_ctl01_btnSubmit",
	}
	d := url.Values{}
	for k, v := range postData {
		d.Set(k, v)
	}
	var data = strings.NewReader(d.Encode())
	return data
}
