import asyncio
from datetime import datetime

import uvloop
from aiohttp import ClientSession, ClientTimeout
from shortuuid import random

from py_util.nsq_api_sample import API_HOST


async def post_greet(session: ClientSession, name: str):
    now = datetime.now().strftime('%x %X.%f')
    api_path = '/greet'
    resp = await session.post(f'{API_HOST}{api_path}', json={
        'name': name,
    })
    print(f'{now} post_greet status {resp.status}')


async def case():
    async with ClientSession(timeout=ClientTimeout(total=5)) as session:
        while True:
            await post_greet(session, random(8))
            await asyncio.sleep(0.05)


def main():
    asyncio.run(case())


if __name__ == '__main__':
    uvloop.install()
    main()
