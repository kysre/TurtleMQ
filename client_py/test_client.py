import json
import time
from unittest import TestCase
from client_py.client import QueueClient
from concurrent import futures
import asyncio
import os


def sample_push_pull(client: QueueClient):
    asyncio.run(client.push("test key", [b'test value']))
    return asyncio.run(client.pull()).value


def store_key_value(key, value):
    print(key, value)
    print(os.path)
    path = '/tmp/test/'
    os.makedirs(os.path.dirname(path), exist_ok=True)
    with open(f'{path}/test.json', 'w') as file:
        file.write(json.dumps(dict(key=key, value=value.decode())))


class TestQueueClient(TestCase):
    def setUp(self):
        self.client = QueueClient()

    def test_push_pull(self):
        self.client.push("test key", [b'test value'])
        self.assertEqual(self.client.pull().value, [b"test value"])
        # self.fail()

    def test_ack(self):
        ack_res = self.client.ack("test key")
        # ack_res = asyncio.run(self.client.ack("test key"))
        self.assertTrue(ack_res)

    def test_concurrent_push_without_order(self):

        _futures = []

        num_threads = 3

        res = []

        with futures.ThreadPoolExecutor(max_workers=num_threads) as executor:

            for i in range(num_threads):
                _futures.append(executor.submit(sample_push_pull, self.client))

            for pulled_value in futures.as_completed(_futures):
                try:
                    res.append(pulled_value.result())
                except Exception as e:
                    print(f'Exception: {e}')

        self.assertListEqual(res, [[b'test value'], [b'test value'], [b'test value']])

    def test_subscribe(self):
        client_2 = QueueClient()
        self.client.subscribe(f=store_key_value)

        client_2.push("test key", b'test value')

        time.sleep(5)

        with open('/tmp/test/test.json') as file:
            json_string = file.read()
            json_dict = json.loads(json_string)

            print(json_dict['key'])

            self.assertEqual(json_dict['key'], 'test key')