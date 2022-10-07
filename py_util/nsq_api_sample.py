import asyncio
from random import randint

from aiohttp import ClientSession
from shortuuid import random

API_HOST = 'http://192.168.2.23:8081'


class UserInfo:
    username = ""
    gender = 0
    age = 0
    external_uid = 0

    def __init__(self, username, gender=0, age=0, external_uid=0):
        self.username = username
        self.gender = gender
        self.age = age
        self.external_uid = external_uid


class UserPhoneInfo:
    uid = 0
    phone = ''

    def __init__(self, uid, phone):
        self.uid = uid
        self.phone = phone


async def post_new_user(session: ClientSession, user: UserInfo):
    api_path = '/db/users'
    resp = await session.post(f'{API_HOST}{api_path}', json={
        'username': user.username,
        'gender': user.gender,
        'age': user.age,
        'external_uid': user.external_uid
    })
    if resp.status != 200:
        print(f'post_new_user failed {await resp.text()}')
        return 0
    payload = await resp.json()
    return payload.get('uid')


async def post_user_phone(session: ClientSession, user_phone: UserPhoneInfo):
    api_path = f'/db/users/{user_phone.uid}/phones'
    resp = await session.post(f'{API_HOST}{api_path}', json={
        'phone': user_phone.phone
    })
    if resp.status != 200:
        print(f'post_user_phone failed {await resp.text()}')
        return False
    return True


async def case():
    count = 20
    users = []
    phones = []

    for x in range(count):
        username = random(length=6)
        users.append(UserInfo(username, age=randint(18, 80), external_uid=randint(0, 9999999)))
        phones.append(str(randint(100000, 999999)))

    async with ClientSession() as session:
        uids = []
        for user in users:
            uid = await post_new_user(session, user)
            if uid != 0:
                uids.append(uid)
        print(f'Created {len(uids)} users')
        await asyncio.sleep(10)

        tasks = []
        loop = asyncio.get_running_loop()
        for i, uid in enumerate(uids):
            user_phone = UserPhoneInfo(uid, phones[i])
            tasks.append(loop.create_task(post_user_phone(session, user_phone)))
            await asyncio.sleep(0.01)

        await asyncio.gather(*tasks)
        failures = 0
        for t in tasks:
            if not t.result():
                failures += 1
        print(f'Failed tasks {failures}')
        return


def main():
    asyncio.run(case())


if __name__ == '__main__':
    main()
