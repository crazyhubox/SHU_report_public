from redis import Redis
Host = 'xxxxxx'
rdb = Redis()

def InitUserTable():
    users = [
        'id,password',
    ]

    for each_user in users:
        res = rdb.sadd('users',each_user)
        print(res)

def ReadUserInfo():
    user_dict = {}
    users = rdb.smembers('users')
    user_num = rdb.scard('users')
    for user_info in users:
        user_info = user_info.split(',')
        uid = user_info[0]
        password = user_info[1]
        user_dict[uid] = password    
    print(f'[INFO]: Read the user_info successfully! There are {user_num} users.')
    return user_dict

if __name__ == "__main__":
    print(ReadUserInfo())