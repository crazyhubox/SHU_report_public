from typing import Dict
import requests
from re import compile

# from GetCookies import GetCookies
# Give the cookies of the user to the class 'PostRequest'.
# Give the 'date' and the 'flag' of the report url and the report method will post the selfReport.
class PostRequest:

    res_find = compile(r'F\.alert\((.+?)\);')

    def __init__(self, cookie: Dict[str, str] = None, viewstate: str = None):
        self._cookie = cookie
        self.viewstate = viewstate

    def setUserInfo(self, cookie: Dict[str, str], viewstate: str):
        self._cookie = cookie
        self.viewstate = viewstate

    def report(self, date: str, state: str):
        # cookies = {
        #     # 'ASP.NET_SessionId': 'tidem3dffsr5mnj3puze0ajw',
        #     '.ncov2019selfreport': '060D77C972DEF6C8C096F58FC8019A72A767B3EAC03499B3318A23C43CEA3A41119ED331922A644AA2BC367861D0F5847A45E1332C372668E67D918350361FADE689289BD820F0A399E73E4EFE9007292FEE2E080EE01E69F6DCCF482CDAA354BD975A7F6D7ACDDD1428F16ECC496EB1',
        # }
        cookies = self._cookie

        headers = {
            'Connection': 'keep-alive',
            'Accept': 'text/plain, */*; q=0.01',
            'X-Requested-With': 'XMLHttpRequest',
            'X-FineUI-Ajax': 'true',
            'User-Agent': 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36 Edg/87.0.664.60',
            'Content-Type': 'application/x-www-form-urlencoded; charset=UTF-8',
            'Origin': 'https://selfreport.shu.edu.cn',
            'Sec-Fetch-Site': 'same-origin',
            'Sec-Fetch-Mode': 'cors',
            'Sec-Fetch-Dest': 'empty',
            'Referer': f'https://selfreport.shu.edu.cn/XueSFX/HalfdayReport.aspx?day={date}&t={state}',
            'Accept-Language': 'zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6',
        }

        params = (
            ('day', date),
            ('t', state),
        )

        data = {
            '__EVENTTARGET': 'p1$ctl00$btnSubmit',
            '__EVENTARGUMENT': '',
            '__VIEWSTATE': self.viewstate,
            '__VIEWSTATEGENERATOR': 'DC4D08A3',
            'p1$ChengNuo': 'p1_ChengNuo',
            'p1$BaoSRQ': date,
            'p1$DangQSTZK': '\u826F\u597D',
            'p1$TiWen': '36',
            'p1$ZaiXiao': '\u5B9D\u5C71',
            'p1$ddlSheng$Value': '\u4E0A\u6D77',
            'p1$ddlSheng': '\u4E0A\u6D77',
            'p1$ddlShi$Value': '\u4E0A\u6D77\u5E02',
            'p1$ddlShi': '\u4E0A\u6D77\u5E02',
            'p1$ddlXian$Value': '\u5B9D\u5C71\u533A',
            'p1$ddlXian': '\u5B9D\u5C71\u533A',
            'p1$FengXDQDL': '\u5426',
            'p1$TongZWDLH': '\u5426',
            'p1$XiangXDZ': '\u65B0\u6821\u533A\u6821\u5185S2',
            'p1$QueZHZJC$Value': '\u5426',
            'p1$QueZHZJC': '\u5426',
            'p1$DangRGL': '\u5426',
            'p1$GeLDZ': '',
            'p1$CengFWH': '\u5426',
            'p1$CengFWH_RiQi': '',
            'p1$CengFWH_BeiZhu': '',
            'p1$JieChu': '\u5426',
            'p1$JieChu_RiQi': '',
            'p1$JieChu_BeiZhu': '',
            'p1$TuJWH': '\u5426',
            'p1$TuJWH_RiQi': '',
            'p1$TuJWH_BeiZhu': '',
            'p1$JiaRen_BeiZhu': '',
            'p1$SuiSM': '\u7EFF\u8272',
            'p1$LvMa14Days': '\u662F',
            'p1$Address2': '',
            'F_TARGET': 'p1_ctl00_btnSubmit',
            'p1_GeLSM_Collapsed': 'false',
            'p1_Collapsed': 'false',
            'F_STATE': 'eyJwMV9CYW9TUlEiOnsiVGV4dCI6IjIwMjAtMTItMTUifSwicDFfRGFuZ1FTVFpLIjp7IkZfSXRlbXMiOltbIuiJr+WlvSIsIuiJr+WlvSIsMV0sWyLkuI3pgIIiLCLkuI3pgIIiLDFdXSwiU2VsZWN0ZWRWYWx1ZSI6IuiJr+WlvSJ9LCJwMV9aaGVuZ1podWFuZyI6eyJIaWRkZW4iOnRydWUsIkZfSXRlbXMiOltbIuaEn+WGkiIsIuaEn+WGkiIsMV0sWyLlkrPll70iLCLlkrPll70iLDFdLFsi5Y+R54OtIiwi5Y+R54OtIiwxXV0sIlNlbGVjdGVkVmFsdWVBcnJheSI6W119LCJwMV9UaVdlbiI6eyJUZXh0IjoiMzYuMCJ9LCJwMV9aYWlYaWFvIjp7IlNlbGVjdGVkVmFsdWUiOiLlrp3lsbEiLCJGX0l0ZW1zIjpbWyLkuI3lnKjmoKEiLCLkuI3lnKjmoKEiLDFdLFsi5a6d5bGxIiwi5a6d5bGx5qCh5Yy6IiwxXSxbIuW7tumVvyIsIuW7tumVv+agoeWMuiIsMV0sWyLlmInlrpoiLCLlmInlrprmoKHljLoiLDFdLFsi5paw6Ze46LevIiwi5paw6Ze46Lev5qCh5Yy6IiwxXV19LCJwMV9kZGxTaGVuZyI6eyJGX0l0ZW1zIjpbWyItMSIsIumAieaLqeecgeS7vSIsMSwiIiwiIl0sWyLljJfkuqwiLCLljJfkuqwiLDEsIiIsIiJdLFsi5aSp5rSlIiwi5aSp5rSlIiwxLCIiLCIiXSxbIuS4iua1tyIsIuS4iua1tyIsMSwiIiwiIl0sWyLph43luoYiLCLph43luoYiLDEsIiIsIiJdLFsi5rKz5YyXIiwi5rKz5YyXIiwxLCIiLCIiXSxbIuWxseilvyIsIuWxseilvyIsMSwiIiwiIl0sWyLovr3lroEiLCLovr3lroEiLDEsIiIsIiJdLFsi5ZCJ5p6XIiwi5ZCJ5p6XIiwxLCIiLCIiXSxbIum7kem+meaxnyIsIum7kem+meaxnyIsMSwiIiwiIl0sWyLmsZ/oi48iLCLmsZ/oi48iLDEsIiIsIiJdLFsi5rWZ5rGfIiwi5rWZ5rGfIiwxLCIiLCIiXSxbIuWuieW+vSIsIuWuieW+vSIsMSwiIiwiIl0sWyLnpo/lu7oiLCLnpo/lu7oiLDEsIiIsIiJdLFsi5rGf6KW/Iiwi5rGf6KW/IiwxLCIiLCIiXSxbIuWxseS4nCIsIuWxseS4nCIsMSwiIiwiIl0sWyLmsrPljZciLCLmsrPljZciLDEsIiIsIiJdLFsi5rmW5YyXIiwi5rmW5YyXIiwxLCIiLCIiXSxbIua5luWNlyIsIua5luWNlyIsMSwiIiwiIl0sWyLlub/kuJwiLCLlub/kuJwiLDEsIiIsIiJdLFsi5rW35Y2XIiwi5rW35Y2XIiwxLCIiLCIiXSxbIuWbm+W3nSIsIuWbm+W3nSIsMSwiIiwiIl0sWyLotLXlt54iLCLotLXlt54iLDEsIiIsIiJdLFsi5LqR5Y2XIiwi5LqR5Y2XIiwxLCIiLCIiXSxbIumZleilvyIsIumZleilvyIsMSwiIiwiIl0sWyLnlJjogoMiLCLnlJjogoMiLDEsIiIsIiJdLFsi6Z2S5rW3Iiwi6Z2S5rW3IiwxLCIiLCIiXSxbIuWGheiSmeWPpCIsIuWGheiSmeWPpCIsMSwiIiwiIl0sWyLlub/opb8iLCLlub/opb8iLDEsIiIsIiJdLFsi6KW/6JePIiwi6KW/6JePIiwxLCIiLCIiXSxbIuWugeWkjyIsIuWugeWkjyIsMSwiIiwiIl0sWyLmlrDnloYiLCLmlrDnloYiLDEsIiIsIiJdLFsi6aaZ5rivIiwi6aaZ5rivIiwxLCIiLCIiXSxbIua+s+mXqCIsIua+s+mXqCIsMSwiIiwiIl0sWyLlj7Dmub4iLCLlj7Dmub4iLDEsIiIsIiJdXSwiU2VsZWN0ZWRWYWx1ZUFycmF5IjpbIuS4iua1tyJdfSwicDFfZGRsU2hpIjp7IkVuYWJsZWQiOnRydWUsIkZfSXRlbXMiOltbIi0xIiwi6YCJ5oup5biCIiwxLCIiLCIiXSxbIuS4iua1t+W4giIsIuS4iua1t+W4giIsMSwiIiwiIl1dLCJTZWxlY3RlZFZhbHVlQXJyYXkiOlsi5LiK5rW35biCIl19LCJwMV9kZGxYaWFuIjp7IkVuYWJsZWQiOnRydWUsIkZfSXRlbXMiOltbIi0xIiwi6YCJ5oup5Y6/5Yy6IiwxLCIiLCIiXSxbIum7hOa1puWMuiIsIum7hOa1puWMuiIsMSwiIiwiIl0sWyLljaLmub7ljLoiLCLljaLmub7ljLoiLDEsIiIsIiJdLFsi5b6Q5rGH5Yy6Iiwi5b6Q5rGH5Yy6IiwxLCIiLCIiXSxbIumVv+WugeWMuiIsIumVv+WugeWMuiIsMSwiIiwiIl0sWyLpnZnlronljLoiLCLpnZnlronljLoiLDEsIiIsIiJdLFsi5pmu6ZmA5Yy6Iiwi5pmu6ZmA5Yy6IiwxLCIiLCIiXSxbIuiZueWPo+WMuiIsIuiZueWPo+WMuiIsMSwiIiwiIl0sWyLmnajmtabljLoiLCLmnajmtabljLoiLDEsIiIsIiJdLFsi5a6d5bGx5Yy6Iiwi5a6d5bGx5Yy6IiwxLCIiLCIiXSxbIumXteihjOWMuiIsIumXteihjOWMuiIsMSwiIiwiIl0sWyLlmInlrprljLoiLCLlmInlrprljLoiLDEsIiIsIiJdLFsi5p2+5rGf5Yy6Iiwi5p2+5rGf5Yy6IiwxLCIiLCIiXSxbIumHkeWxseWMuiIsIumHkeWxseWMuiIsMSwiIiwiIl0sWyLpnZLmtabljLoiLCLpnZLmtabljLoiLDEsIiIsIiJdLFsi5aWJ6LSk5Yy6Iiwi5aWJ6LSk5Yy6IiwxLCIiLCIiXSxbIua1puS4nOaWsOWMuiIsIua1puS4nOaWsOWMuiIsMSwiIiwiIl0sWyLltIfmmI7ljLoiLCLltIfmmI7ljLoiLDEsIiIsIiJdXSwiU2VsZWN0ZWRWYWx1ZUFycmF5IjpbIuWuneWxseWMuiJdfSwicDFfRmVuZ1hEUURMIjp7IkxhYmVsIjoiMTLmnIgwMeaXpeiHszEy5pyIMTXml6XmmK/lkKblnKjkuK3pq5jpo47pmanlnLDljLrpgJfnlZk8c3BhbiBzdHlsZT0nY29sb3I6cmVkOyc+77yI5YaF6JKZ5Y+k5ruh5rSy6YeM5Lic5bGx6KGX6YGT44CB5YyX5Yy66KGX6YGT5Yqe5LqL5aSE77yM5omO6LWJ6K+65bCU5Yy656ys5LiJ44CB56ys5Zub44CB56ys5LqU6KGX6YGT5Yqe5LqL5aSE77yM5oiQ6YO95biC6YOr6YO95Yy66YOr562S6KGX6YGT5aSq5bmz5p2R44CB6I+g6JCd56S+5Yy65Lit6ZOB5aWl57u05bCU5LqM5pyf44CB5LiJ5pyf44CB6YOr6YO95Yy65ZSQ5piM6ZWH5rC45a6J5p2ROOe7hOOAgeaIkOWNjuWMuuW0lOWutuW6l+WNjumDveS6keaZr+WPsOWwj+WMuu+8jOm7kem+meaxn+eJoeS4ueaxn+W4guS4nOWugeW4guS4reW/g+ekvuWMuuOAgee7peiKrOays+W4gumdkuS6keWwj+WMuu+8jOaWsOeWhuWQkOmygeeVquW4gue6ouebvuWwj+WMuu+8iTwvc3Bhbj4iLCJTZWxlY3RlZFZhbHVlIjoi5ZCmIiwiRl9JdGVtcyI6W1si5pivIiwi5pivIiwxXSxbIuWQpiIsIuWQpiIsMV1dfSwicDFfVG9uZ1pXRExIIjp7IkxhYmVsIjoi5LiK5rW35ZCM5L2P5Lq65ZGY5piv5ZCm5pyJMTLmnIgwMeaXpeiHszEy5pyIMTXml6XmnaXoh6rkuK3pq5jpo47pmanlnLDljLrnmoTkuro8c3BhbiBzdHlsZT0nY29sb3I6cmVkOyc+77yI5YaF6JKZ5Y+k5ruh5rSy6YeM5Lic5bGx6KGX6YGT44CB5YyX5Yy66KGX6YGT5Yqe5LqL5aSE77yM5omO6LWJ6K+65bCU5Yy656ys5LiJ44CB56ys5Zub44CB56ys5LqU6KGX6YGT5Yqe5LqL5aSE77yM5oiQ6YO95biC6YOr6YO95Yy66YOr562S6KGX6YGT5aSq5bmz5p2R44CB6I+g6JCd56S+5Yy65Lit6ZOB5aWl57u05bCU5LqM5pyf44CB5LiJ5pyf44CB6YOr6YO95Yy65ZSQ5piM6ZWH5rC45a6J5p2ROOe7hOOAgeaIkOWNjuWMuuW0lOWutuW6l+WNjumDveS6keaZr+WPsOWwj+WMuu+8jOm7kem+meaxn+eJoeS4ueaxn+W4guS4nOWugeW4guS4reW/g+ekvuWMuuOAgee7peiKrOays+W4gumdkuS6keWwj+WMuu+8jOaWsOeWhuWQkOmygeeVquW4gue6ouebvuWwj+WMuu+8iTwvc3Bhbj4iLCJTZWxlY3RlZFZhbHVlIjoi5ZCmIiwiRl9JdGVtcyI6W1si5pivIiwi5pivIiwxXSxbIuWQpiIsIuWQpiIsMV1dfSwicDFfWGlhbmdYRFoiOnsiVGV4dCI6IuaWsOagoeWMuuagoeWGhVMyIn0sInAxX1F1ZVpIWkpDIjp7IkZfSXRlbXMiOltbIuaYryIsIuaYryIsMSwiIiwiIl0sWyLlkKYiLCLlkKYiLDEsIiIsIiJdXSwiU2VsZWN0ZWRWYWx1ZUFycmF5IjpbIuWQpiJdfSwicDFfRGFuZ1JHTCI6eyJTZWxlY3RlZFZhbHVlIjoi5ZCmIiwiRl9JdGVtcyI6W1si5pivIiwi5pivIiwxXSxbIuWQpiIsIuWQpiIsMV1dfSwicDFfR2VMU00iOnsiSGlkZGVuIjp0cnVlLCJJRnJhbWVBdHRyaWJ1dGVzIjp7fX0sInAxX0dlTEZTIjp7IlJlcXVpcmVkIjpmYWxzZSwiSGlkZGVuIjp0cnVlLCJGX0l0ZW1zIjpbWyLlsYXlrrbpmpTnprsiLCLlsYXlrrbpmpTnprsiLDFdLFsi6ZuG5Lit6ZqU56a7Iiwi6ZuG5Lit6ZqU56a7IiwxXV0sIlNlbGVjdGVkVmFsdWUiOm51bGx9LCJwMV9HZUxEWiI6eyJIaWRkZW4iOnRydWV9LCJwMV9DZW5nRldIIjp7IkxhYmVsIjoiMTLmnIgwMeaXpeiHszEy5pyIMTXml6XmmK/lkKblnKjkuK3pq5jpo47pmanlnLDljLrpgJfnlZnov4c8c3BhbiBzdHlsZT0nY29sb3I6cmVkOyc+77yI5YaF6JKZ5Y+k5ruh5rSy6YeM5Lic5bGx6KGX6YGT44CB5YyX5Yy66KGX6YGT5Yqe5LqL5aSE77yM5omO6LWJ6K+65bCU5Yy656ys5LiJ44CB56ys5Zub44CB56ys5LqU6KGX6YGT5Yqe5LqL5aSE77yM5oiQ6YO95biC6YOr6YO95Yy66YOr562S6KGX6YGT5aSq5bmz5p2R44CB6I+g6JCd56S+5Yy65Lit6ZOB5aWl57u05bCU5LqM5pyf44CB5LiJ5pyf44CB6YOr6YO95Yy65ZSQ5piM6ZWH5rC45a6J5p2ROOe7hOOAgeaIkOWNjuWMuuW0lOWutuW6l+WNjumDveS6keaZr+WPsOWwj+WMuu+8jOm7kem+meaxn+eJoeS4ueaxn+W4guS4nOWugeW4guS4reW/g+ekvuWMuuOAgee7peiKrOays+W4gumdkuS6keWwj+WMuu+8jOaWsOeWhuWQkOmygeeVquW4gue6ouebvuWwj+WMuu+8iTwvc3Bhbj4iLCJGX0l0ZW1zIjpbWyLmmK8iLCLmmK8iLDFdLFsi5ZCmIiwi5ZCmIiwxXV0sIlNlbGVjdGVkVmFsdWUiOiLlkKYifSwicDFfQ2VuZ0ZXSF9SaVFpIjp7IkhpZGRlbiI6dHJ1ZX0sInAxX0NlbmdGV0hfQmVpWmh1Ijp7IkhpZGRlbiI6dHJ1ZX0sInAxX0ppZUNodSI6eyJMYWJlbCI6IjEy5pyIMDHml6Xoh7MxMuaciDE15pel5piv5ZCm5LiO5p2l6Ieq5Lit6auY6aOO6Zmp5Zyw5Yy65Y+R54Ot5Lq65ZGY5a+G5YiH5o6l6KemPHNwYW4gc3R5bGU9J2NvbG9yOnJlZDsnPu+8iOWGheiSmeWPpOa7oea0sumHjOS4nOWxseihl+mBk+OAgeWMl+WMuuihl+mBk+WKnuS6i+WkhO+8jOaJjui1ieivuuWwlOWMuuesrOS4ieOAgeesrOWbm+OAgeesrOS6lOihl+mBk+WKnuS6i+WkhO+8jOaIkOmDveW4gumDq+mDveWMuumDq+etkuihl+mBk+WkquW5s+adkeOAgeiPoOiQneekvuWMuuS4remTgeWlpee7tOWwlOS6jOacn+OAgeS4ieacn+OAgemDq+mDveWMuuWUkOaYjOmVh+awuOWuieadkTjnu4TjgIHmiJDljY7ljLrltJTlrrblupfljY7pg73kupHmma/lj7DlsI/ljLrvvIzpu5HpvpnmsZ/niaHkuLnmsZ/luILkuJzlroHluILkuK3lv4PnpL7ljLrjgIHnu6XoiqzmsrPluILpnZLkupHlsI/ljLrvvIzmlrDnloblkJDpsoHnlarluILnuqLnm77lsI/ljLrvvIk8L3NwYW4+IiwiU2VsZWN0ZWRWYWx1ZSI6IuWQpiIsIkZfSXRlbXMiOltbIuaYryIsIuaYryIsMV0sWyLlkKYiLCLlkKYiLDFdXX0sInAxX0ppZUNodV9SaVFpIjp7IkhpZGRlbiI6dHJ1ZX0sInAxX0ppZUNodV9CZWlaaHUiOnsiSGlkZGVuIjp0cnVlfSwicDFfVHVKV0giOnsiTGFiZWwiOiIxMuaciDAx5pel6IezMTLmnIgxNeaXpeaYr+WQpuS5mOWdkOWFrOWFseS6pOmAmumAlOW+hOS4remrmOmjjumZqeWcsOWMujxzcGFuIHN0eWxlPSdjb2xvcjpyZWQ7Jz7vvIjlhoXokpnlj6Tmu6HmtLLph4zkuJzlsbHooZfpgZPjgIHljJfljLrooZfpgZPlip7kuovlpITvvIzmiY7otYnor7rlsJTljLrnrKzkuInjgIHnrKzlm5vjgIHnrKzkupTooZfpgZPlip7kuovlpITvvIzmiJDpg73luILpg6vpg73ljLrpg6vnrZLooZfpgZPlpKrlubPmnZHjgIHoj6DokJ3npL7ljLrkuK3pk4HlpaXnu7TlsJTkuozmnJ/jgIHkuInmnJ/jgIHpg6vpg73ljLrllJDmmIzplYfmsLjlronmnZE457uE44CB5oiQ5Y2O5Yy65bSU5a625bqX5Y2O6YO95LqR5pmv5Y+w5bCP5Yy677yM6buR6b6Z5rGf54mh5Li55rGf5biC5Lic5a6B5biC5Lit5b+D56S+5Yy644CB57ul6Iqs5rKz5biC6Z2S5LqR5bCP5Yy677yM5paw55aG5ZCQ6bKB55Wq5biC57qi55u+5bCP5Yy677yJPC9zcGFuPiIsIlNlbGVjdGVkVmFsdWUiOiLlkKYiLCJGX0l0ZW1zIjpbWyLmmK8iLCLmmK8iLDFdLFsi5ZCmIiwi5ZCmIiwxXV19LCJwMV9UdUpXSF9SaVFpIjp7IkhpZGRlbiI6dHJ1ZX0sInAxX1R1SldIX0JlaVpodSI6eyJIaWRkZW4iOnRydWV9LCJwMV9KaWFSZW4iOnsiTGFiZWwiOiIxMuaciDAx5pel6IezMTLmnIgxNeaXpeWutuS6uuaYr+WQpuacieWPkeeDreetieeXh+eKtiJ9LCJwMV9KaWFSZW5fQmVpWmh1Ijp7IkhpZGRlbiI6dHJ1ZX0sInAxX1N1aVNNIjp7IlNlbGVjdGVkVmFsdWUiOiLnu7/oibIiLCJGX0l0ZW1zIjpbWyLnuqLoibIiLCLnuqLoibIiLDFdLFsi6buE6ImyIiwi6buE6ImyIiwxXSxbIue7v+iJsiIsIue7v+iJsiIsMV1dfSwicDFfTHZNYTE0RGF5cyI6eyJTZWxlY3RlZFZhbHVlIjoi5pivIiwiRl9JdGVtcyI6W1si5pivIiwi5pivIiwxXSxbIuWQpiIsIuWQpiIsMV1dfSwicDEiOnsiVGl0bGUiOiLmr4/ml6XkuKTmiqXvvIjkuIvljYjvvIkiLCJJRnJhbWVBdHRyaWJ1dGVzIjp7fX19'
        }

        response = requests.post('https://selfreport.shu.edu.cn/XueSFX/HalfdayReport.aspx',
                                 headers=headers, params=params, cookies=cookies, data=data)
        print('Referer:', headers['Referer'],end=f'[{response.status_code}]')
        res = self.res_find.search(response.text)
        if res:
            print(res[1])


def main(): pass


if __name__ == "__main__":
    pass