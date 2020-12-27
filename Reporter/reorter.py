from .test import gen_url, datetime


class Reporter:
    def __init__(self, Requester=None) -> None:
        self.__req = Requester

    def setRequester(self, Requester):
        self.__req = Requester

    def PreviousReport(self, startDate: str):
        for each_date, t in gen_url(startDate):
            self.__req.report(each_date, t)

    def TodayReport(self):
        today = datetime.date.today()
        self.SunReport(today)
        self.MoonRepot(today)

    def SunReport(self, date: str):
        self.__req.report(date, 1)

    def MoonRepot(self, date: str):
        self.__req.report(date, 2)
