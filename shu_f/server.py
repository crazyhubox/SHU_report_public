from fastapi import FastAPI,Request
import uvicorn
from pyppeteer import launch
from cookie import Cookies


app = FastAPI()
page = None

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
    


@app.middleware("http")
async def launch_the_browser(request: Request, call_next):
    """Launch the browser before the request is handled by route."""
    global page
    if not page:
        page = await open_browser()
    response = await call_next(request)
    return response



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


@app.get("/cookies")
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
    return await Cookies(uid=id,password=password,page=page).Read()

@app.get("/test")
async def test_respone():
    await page.goto("https://baidu.com")
    if page:
        return page.url
    return "No page"





if __name__ == '__main__':
    # uvicorn.run("server:app", host="0.0.0.0", port=8989,reload=True)
    uvicorn.run(app, host="0.0.0.0", port=8989)
