from fastapi import FastAPI
import uvicorn
import asyncio
from pyppeteer import launch

app = FastAPI()

async def open_browser():
    """
    Open the browser obj

    :return : Page
    :rtype  : pyppeteer.Page
    """
    browser = await launch({'headless': True, 'args': ['--disable-infobars', '--window-size=1920,1080', '--no-sandbox']})
    # 打开一个页面
    page = await browser.newPage()
    await page.setViewport({'width': 1920, 'height': 1080})   # 设置页面的大小
    return page
    

async def get_cookies(page,uid,password):
    """
    Get the cookies.
    
    :param page     : Page obj
    :type page      : pyppeteer.Page
    :param uid      : user_id
    :type uid       : str
    :param password : passwod
    :type password  : str
    :return         : cookies
    :rtype          : Dict[name,value]
    """
    # 打开链接
    await page.goto('https://newsso.shu.edu.cn/login/eyJ0aW1lc3RhbXAiOjE2MTg2MjU2NjkyNjAyMDUzODksInJlc3BvbnNlVHlwZSI6ImNvZGUiLCJjbGllbnRJZCI6IldVSFdmcm50bldZSFpmelE1UXZYVUNWeSIsInNjb3BlIjoiMSIsInJlZGlyZWN0VXJpIjoiaHR0cHM6Ly9zZWxmcmVwb3J0LnNodS5lZHUuY24vTG9naW5TU08uYXNweD9SZXR1cm5Vcmw9JTJmRGVmYXVsdC5hc3B4Iiwic3RhdGUiOiIifQ==') 

    await page.type("#username",uid)
    await page.type("#password",password)
    await asyncio.wait([
        page.click("#submit"),
        page.waitForNavigation(),
    ])
    cookies = await page.cookies()
    for each_cookie in cookies:
        pass
        if each_cookie['name'] == '.ncov2019selfreport':
            return each_cookie


# @app.on_event('startup')
# async def start():
#     """
#     Open the browser at the start of the server.
#     Create the global object page.
#     """
#     global page
#     page = await open_browser()

@app.get("/init")
async def openBrowser():
    """
    Use the Page to visit the shu_report URL.
    GET the user cookies for report.

    :return :  cookies of user
    :rtype  :  dict
    """
    global page
    page = await open_browser()
    return "Browser has inited finished."


@app.get("/cookies/")
async def root(id:str,password:str):
    """
    Use the Page to visit the shu_report URL.
    GET the user cookies for report.

    :param user_id      : user_id
    :type user_id       : str
    :param password     : password
    :type password      : str
    :return             : cookies of user
    :rtype              : dict
    """
    print(id,password)
    return await get_cookies(page,id,password)

if __name__ == '__main__':
    # uvicorn.run("server:app", host="0.0.0.0", port=8989,reload=True)
    uvicorn.run(app, host="0.0.0.0", port=8989)
