from asyncio.runners import run
from server import open_browser
from cookie import Cookies
import  asyncio
from time import time

now = lambda: time()


def count_time(func):
    def new_f():
        start = now()
        res = func()
        print("time:",now() - start)
        return res
    return new_f


@count_time
def main():
    async def run():
        page = await open_browser()
        cooer = Cookies("16123113","130E2d898",page)    
        for _ in range(1000):
            await cooer.get_key_()
            print(cooer.password)
        return cooer.password

    return asyncio.get_event_loop().run_until_complete(run())


if __name__ == '__main__':
    key = main()
    print(key)