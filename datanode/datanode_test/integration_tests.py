import time
from unittest import TestCase
from datanode_client import QueueClient
from concurrent import futures
import asyncio


def sample_push_pull(client: QueueClient):
    asyncio.run(client.push("test key", b'test value'))
    return asyncio.run(client.pull()).value


class TestQueueClient(TestCase):
    def setUp(self):
        self.client = QueueClient()

    def test_push(self):
        asyncio.run(self.client.push('test key', b'test value'))

    def test_push_pull(self):
        self.assertEqual(sample_push_pull(self.client), b"test value")
        # self.fail()

    def test_multiple_push_pull(self):
        asyncio.run(self.client.push("test key 1", b'test value 1'))
        asyncio.run(self.client.push("test key 1", b'test value 2'))


        response = asyncio.run(self.client.pull_without_ack())
        value = response.value

        self.assertEqual(value, b'test value 1')

        response = asyncio.run(self.client.pull())

        self.assertEqual(response, 'error')

        time.sleep(10)

        response = asyncio.run(self.client.pull())
        value = response.value

        # The first pull request did not send the ack, so this one should also be [b'test value 1']
        self.assertEqual(value, b'test value 1')

        response = asyncio.run(self.client.pull())
        value = response.value

        self.assertEqual(value, b'test value 2')

    def test_ack(self):
        ack_res = asyncio.run(self.client.ack("test key"))
        self.assertTrue(ack_res)

    def test_multiple_partition(self):
        asyncio.run(self.client.push("test key 1", b'test value 1'))
        asyncio.run(self.client.push("dsajiojad", b'test value 2'))
        asyncio.run(self.client.push("testasflj", b'test value 3'))
        asyncio.run(self.client.push("test ;fjgijnas;fpj4", b'test value 4'))

        # ----------- run first pull ---------------
        response = asyncio.run(self.client.pull_without_ack())
        value_0 = response.value

        expected_values = [b'test value 1', b'test value 2', b'test value 3', b'test value 4']

        self.assertTrue(value_0 in expected_values)

        # ----------- run second pull ---------------
        response = asyncio.run(self.client.pull_without_ack())
        value_1 = response.value

        expected_values.remove(value_0)
        self.assertTrue(value_1 in expected_values)

        time.sleep(20)

        # ----------- run first pull ---------------
        response = asyncio.run(self.client.pull())
        value_0 = response.value

        expected_values = [b'test value 1', b'test value 2', b'test value 3', b'test value 4']

        self.assertTrue(value_0 in expected_values)

        # ----------- run second pull ---------------
        response = asyncio.run(self.client.pull())
        value_1 = response.value

        expected_values.remove(value_0)
        self.assertTrue(value_1 in expected_values)

        # ----------- run third pull ---------------
        response = asyncio.run(self.client.pull())
        value_2 = response.value

        expected_values.remove(value_1)
        self.assertTrue(value_2 in expected_values)

        # ----------- run fourth pull ---------------
        response = asyncio.run(self.client.pull())
        value_3 = response.value

        expected_values.remove(value_2)
        self.assertTrue(value_3 in expected_values)

        response = asyncio.run(self.client.pull())
        self.assertEqual(response, 'error')

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

        self.assertListEqual(res, [b'test value', b'test value', b'test value'])
