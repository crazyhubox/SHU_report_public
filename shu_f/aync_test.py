

# base_url = 'https://newsso.shu.edu.cn'

    # def coo(cookies:str):
    #     c = {}
    #     coo_temp = cookies.split(";")[0].split(":")        
    #     coo_temp = coo_temp[1].split("=")
    #     print(coo_temp)
    #     c[coo_temp[0].strip()] = coo_temp[1]
    #     return c


    # async def POST1(session:ClientSession,data):
    #     async with session.post('https://newsso.shu.edu.cn/login/eyJ0aW1lc3RhbXAiOjE2MTg5OTUyNzg0ODkwMzA1MjYsInJlc3BvbnNlVHlwZSI6ImNvZGUiLCJjbGllbnRJZCI6IldVSFdmcm50bldZSFpmelE1UXZYVUNWeSIsInNjb3BlIjoiMSIsInJlZGlyZWN0VXJpIjoiaHR0cHM6Ly9zZWxmcmVwb3J0LnNodS5lZHUuY24vTG9naW5TU08uYXNweD9SZXR1cm5Vcmw9JTJmRGVmYXVsdC5hc3B4Iiwic3RhdGUiOiIifQ==',data=data,allow_redirects=False) as response:
    #         session.cookie_jar.filter_cookies(base_url)
    #         await response.text()
    #         next_url = response.headers['Location']
    #         cookies = response.cookies
    #         cookies = coo(str(cookies))
    #         print(cookies)
    #         return next_url,cookies

    # async def GET2(session:ClientSession,url,cookies):
    #     print("="*100)
    #     print(url)
    #     async with session.get(url,allow_redirects=False,cookies=cookies) as response:
    #         session.cookie_jar.filter_cookies("https://selfreport.shu.edu.cn")
    #         await response.text()
    #         next_url = response.headers['Location']

    #         return next_url

    
    # async def GET3(session:ClientSession,url,cookies):
    #     print("="*100)
    #     print(url)
    #     async with session.get(url,allow_redirects=False,cookies=cookies) as response:
    #         await response.text()
    #         cookies = response.headers['Set-Cookie']
    #         print(response.cookies)
    #         print(session.cookie_jar.filter_cookies("https://selfreport.shu.edu.cn"))
    #         return cookies

 # async with ClientSession() as session:
            
        #     url_temp,c1 = await POST1(session,data)

        #     url = urljoin(base=base_url,url=url_temp)
        #     url_temp = await GET2(session,url,cookies=c1)
        #     print(url_temp)
        #     cookies = await GET3(session,url_temp,cookies=None)
            
        # return cookies