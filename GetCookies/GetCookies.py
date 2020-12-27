
from pyppeteer import launch
from pyppeteer.browser import Browser
from pyppeteer.page import Page
import asyncio
from typing import Dict, List
import datetime
from os import system
from time import sleep

class GetCookies:
    """
    Get the cookies of the userInfo which user has given.
    [id_num]   : username
    [password] : password

    GetCookies().cookies() ==> cookies ==> PostRquest()
    """

    def __init__(self, id_num: str, password: str) -> None:
        self.user_info = {
            'user': id_num,
            'password': password
        }
        self.vstate = ''
        
        
    def setUserInfo(self,username:str,password:str):
        """
        Update the infomation from the new userInfo which user can use this api to give.
        Reset the self.vstate param
        

        :param username: username
        :type username: str
        :param password: password
        :type password: str
        """
        # for key, value in userInfo.items():
        self.user_info['user'] = username
        self.user_info['password'] = password
        self.vstate = ''

    
    def cookies(self) -> dict:
        return asyncio.get_event_loop().run_until_complete(self.__cookies())


    def viewstate(self) -> str:
        if not self.vstate:
            raise ValueError('The viewstate meta cannot be None.There must be some error.')
        return self.vstate


    async def __cookies(self):
        cookie_set = set()
        while True:
            page = await self.getPage()
            try:
                await self.login(page)
            except Exception as e:
                print('[ERROR]: Our ip addr is banned by the SHU.')
                system('pgrep chrome | xargs kill -s 9')
                print('[INFO]: Kill the all chrome')
                print('[WAIT]: 10mins.')
                sleep(600)
            else:
                break

        cookies_list = await page.cookies()
        
        for each_cookie in cookies_list:
            cookie_set.add(each_cookie['name'])
        # Check the get-cookies
        if '.ncov2019selfreport' not in cookie_set:
            raise ValueError('The userInfo is wrong.')

        cookies_dict = {}
        for each in cookies_list:
            cookies_dict[each['name']] = each['value']
        print('[INFO]: The cookies has been got.')
        return cookies_dict


    async def getPage(self):
        while True:
            try:
                browser = await launch({'headless': True, 'args': ['--disable-infobars', '--window-size=1920,1080', '--no-sandbox']})
            except Exception as e:
                print("[ERROR]: There are some unknown erros in the launching of the browser.")
                system('pgrep chrome | xargs kill -s 9')
                print(f"[ERROR]: {e}")
            else:
                print('[INFO]: The browser is launched successfully.')
                context = await browser.createIncognitoBrowserContext()
                page = await context.newPage()
                await page.setViewport({'width': 1920, 'height': 1080})   # 设置页面的大小
                return page

    async def login(self, page) -> None:
        await page.goto('https://selfreport.shu.edu.cn/Default.aspx')
        await page.waitFor('#username',{'timeout':500})
        await page.type('#username', self.user_info['user'])
        await page.type('#password', self.user_info['password'])
        await page.click('#submit')
        await page.waitFor('body')
        while True:
            try:
                await page.goto(f'https://selfreport.shu.edu.cn/XueSFX/HalfdayReport.aspx?day={getYesterday()}&t=1')
                await page.waitFor('input#__VIEWSTATE',{'timeout':500})
                break
            except Exception as e:
                print(e)
                pass
        self.vstate =  await page.querySelectorEval('#__VIEWSTATE','node => node.value')

def getYesterday(): 
    today=datetime.date.today() 
    oneday=datetime.timedelta(days=1) 
    yesterday=today-oneday  
    return yesterday

if __name__ == "__main__":
    from time import sleep
    for i in range(100):
        obj_get = GetCookies('id', 'password')
        cookies = obj_get.cookies()
        print(cookies)
        # sleep(3)
    # print(getYesterday())
