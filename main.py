from GetCookies import GetCookies
from Requets import PostRequest
from Reporter import Reporter
import schedule
from time import sleep
from datetime import date,datetime
from getUserFromRedis import ReadUserInfo

def main(t):    
    today = date.today()
    accounts = ReadUserInfo()
    cookie_obj = GetCookies(...,...)
    req_obj = PostRequest()
    reporter = Reporter()
    
    for user, passw in accounts.items():
        print('[INFO]:',user)
        cookie_obj.setUserInfo(username=user,password=passw)

        cookies = cookie_obj.cookies()
        view_state = cookie_obj.viewstate()

        req_obj.setUserInfo(cookies,view_state)
        reporter.setRequester(req_obj)
        # reporter.PreviousReport('2020-12-20')
        if t:
            reporter.SunReport(today)
        else:
            reporter.MoonRepot(today)
        print('='*100)
        sleep(3)
    print('[Finished]:',datetime.now())


def run():
    schedule.every().day.at("07:30").do(main,t=1)
    schedule.every().day.at("20:00").do(main,t=0)

    while True:
        schedule.run_pending()
        sleep(1)


if __name__ == "__main__":
    # main(1)
    run()

    


