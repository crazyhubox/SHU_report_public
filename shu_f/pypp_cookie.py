import asyncio

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
