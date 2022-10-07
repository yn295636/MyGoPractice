import asyncio
from datetime import datetime

import uvloop
from aiohttp import ClientSession, ClientTimeout
from shortuuid import random

from py_util.nsq_api_sample import API_HOST


async def post_redis(session: ClientSession, key: str, value: str):
    now = datetime.now().strftime('%x %X.%f')
    api_path = '/redis'
    resp = await session.post(f'{API_HOST}{api_path}', json={
        'key': key,
        'value': value
    })
    print(f'{now} post_redis status {resp.status}')


async def case():
    async with ClientSession(timeout=ClientTimeout(total=5)) as session:
        loop = asyncio.get_running_loop()
        tasks = []
        for _ in range(3):
            tasks.append(loop.create_task(post_redis(session, random(8), random(4))))
        await asyncio.gather(*tasks)
        tasks = []
        for _ in range(3):
            tasks.append(loop.create_task(post_redis(session, random(8), random(4))))
        await asyncio.gather(*tasks)


def main():
    asyncio.run(case())


if __name__ == '__main__':
    uvloop.install()
    main()
