from unittest import TestCase
from client_py.client import QueueClient
from concurrent import futures
import asyncio


def sample_push_pull(client: QueueClient):
    asyncio.run(client.push("test key", [b'test value']))
    return asyncio.run(client.pull()).value


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
        self.fail()
