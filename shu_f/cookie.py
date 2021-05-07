from abc import ABCMeta, abstractmethod
from pyppeteer import launch
import requests
from pypp_cookie import get_cookies

class CookiesGetter(metaclass=ABCMeta):

    def __init__(self,uid:str,password:str,page=None) -> None:
        self.page =  page
        self.username  = uid
        self.password = password

    @abstractmethod
    def Read(self):
        pass

class Cookies(CookiesGetter):

    def post_(self):
        data = {
            'username': self.username,
            'password': self.password
        }
        response = requests.post('https://newsso.shu.edu.cn/login/eyJ0aW1lc3RhbXAiOjE2MTg5OTUyNzg0ODkwMzA1MjYsInJlc3BvbnNlVHlwZSI6ImNvZGUiLCJjbGllbnRJZCI6IldVSFdmcm50bldZSFpmelE1UXZYVUNWeSIsInNjb3BlIjoiMSIsInJlZGlyZWN0VXJpIjoiaHR0cHM6Ly9zZWxmcmVwb3J0LnNodS5lZHUuY24vTG9naW5TU08uYXNweD9SZXR1cm5Vcmw9JTJmRGVmYXVsdC5hc3B4Iiwic3RhdGUiOiIifQ==', data=data)
        for each in response.history:
            for k, v in each.cookies.items():
                if k == '.ncov2019selfreport':
                    return f'{k}={v}'
        return None

    async def Read(self):
        await self.get_key_()
        if not isinstance(self.password,str) or len(self.password) < 20:
            return 'erropw'
        cookies = self.post_()
        if not cookies:
            return 'nocookie'
        return cookies

    @staticmethod
    def read_js(path) -> str:
        with open(path, 'r', encoding='utf-8') as f:
            return f.read()

    async def get_key_(self):
        # page = await open_browser()
        key_file_path = '/Users/tomjack/Desktop/code/Python/SHU_report_public/js_test/jiami.js'
        key_js = self.read_js(key_file_path)
        key_js += f'\ntest("{self.password}")'

        key_value = await self.page.evaluate(key_js)
        self.password = key_value

    async def open_browser(self):
        """
        Open the browser obj

        :return : Page
        :rtype  : pyppeteer.Page
        """
        if self.page:
            return 
        browser = await launch({'headless': True, 'args': ['--disable-infobars', '--window-size=1920,1080', '--no-sandbox']})
        # 打开一个页面
        page = await browser.newPage()
        await page.setViewport({'width': 1920, 'height': 1080})   # 设置页面的大小
        self.page = page
        return page

    async def test_single(self):
        print(self.password)
        page = await self.open_browser()
        password = await self.get_key_(page=page,password=self.password)
        data = {
            'username': self.username,
            'password': password
        }
        if not isinstance(password,str):
            return 
        return self.post_(data=data)

class CookiesPpeteer(CookiesGetter):

    async def Read(self):
        cookies = await get_cookies(self.page,self.username,self.password)
        if not cookies:
            return 'nocookie'
        return cookies



if __name__ == '__main__':
    cookies: str = get_cookies(None,"", "")
    print(cookies)
    # print(getUrl())
