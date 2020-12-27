from functools import partial


def test1(t,a):
    print(t,a)



def main():
    test_func = partial(test1,a=1,t=2)
    #partial return a function quote
    test_func()

if __name__ == '__main__':
    main()