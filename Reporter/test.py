import calendar
import datetime
from typing import Union
from datetime import date

def get_daysFromMonth(year,month):
    monthRange = calendar.monthrange(year, month)
    daysCount = monthRange[1]
    return daysCount

def gen_url(start_date:Union[str,date]):
    """
    :type start_date: date or str
    """
    today = datetime.date.today()
    check_res = check_date(start_date)
    if not check_res:
        raise ValueError("The date value is illegal!")

    if isinstance(start_date,str):
        year, mon, day = start_date.split('-')
        year_int = int(year)
        mon_int = int(mon)
        day_int = int(day)
    elif isinstance(start_date,date):
        year_int = start_date.year
        mon_int = start_date.month
        day_int = start_date.day
    else:
        raise TypeError('The param start_date must be type str or date.')

    current_date = date(year_int,mon_int,day_int)
    while check_date(current_date):
        daysCount = get_daysFromMonth(year_int,mon_int)
        if mon_int == today.month:
            daysCount = today.day
        for i in range(day_int, 40):
            if i > daysCount:
                break
            mon_day = f'{mon_int}-{i}'
            for j in range(1, 3):
                # pageUrl = f'https://selfreport.shu.edu.cn/XueSFX/HalfdayReport.aspx?day={year_int}-{mon_day}&t={j}'
                # yield pageUrl
                yield (f'{year_int}-{mon_day}',j)
                report_time = '早报' if j == 1 else '晚报'
                # print(pageUrl,report_time)
        day_int, mon_int, year_int = gotoNextMonth(day_int, mon_int, year_int)
        current_date = date(year_int, mon_int, day_int)


def gotoNextMonth(day_int, mon_int, year_int):
    mon_int += 1
    if mon_int > 12:
        year_int += 1
        mon_int = mon_int - 12
    day_int = 1
    return day_int, mon_int, year_int


def check_date(start_date):
    """
    Check whether the start_date  is correct.
    """
    today = datetime.date.today()
    if isinstance(start_date, str):
        year, mon, day = start_date.split('-')
        year_int = int(year)
        mon_int = int(mon)
        day_int = int(day)
    elif isinstance(start_date, date):
        year_int = start_date.year
        mon_int = start_date.month
        day_int = start_date.day
    else:
        raise TypeError('The type of start_date is invailid.')
    daysCount = get_daysFromMonth(year_int,mon_int)

    if year_int > today.year:
        return False
    elif year_int == today.year and mon_int > today.month:
        return False
    elif year_int == today.year and mon_int == today.month and day_int > today.day:
        return False
    elif day_int > daysCount:
        return False
    elif mon_int > 12:
        return False
    else:
        return True


if __name__ == '__main__':
    test_date = date(2020,11,25)
    # res = check_date(test_date)
    # today = datetime.date.today()
    # if isinstance(test_date,date):
    #     pass
    # elif isinstance(today,date):
    #     print(1)
    # print(res)

    # date1 = date(2021,11,40)
    # date2 = date(2021,12,1)
    # if date1 > date2:
    #     print(1)
    # else:
    #     print(2)

    t = gen_url(test_date)
    for each in t:
        print(each)

