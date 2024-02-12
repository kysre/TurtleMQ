import time
from unittest import TestCase
from datanode_client import QueueClient
from concurrent import futures
import asyncio
import random
import string
import hashlib


def hash_function(key: str, number_of_partitions: int) -> int:
    result = hashlib.md5(key.encode())
    hashed_bytes = result.digest()

    return int.from_bytes(hashed_bytes, byteorder='big') % number_of_partitions


def sample_push_pull(client: QueueClient):
    key, value = random_string(), random_string().encode()
    client.push(key, value)
    return client.pull().value, value


def random_string():
    return ''.join(random.choices(string.ascii_uppercase + string.digits, k=10))


class TestQueueClient(TestCase):
    def setUp(self):
        self.client = QueueClient()

    def test_push(self):
        self.client.push(random_string(), random_string().encode())

    def test_push_pull(self):
        real, expected = sample_push_pull(self.client)
        self.assertEqual(real, expected)

    def test_multiple_push_pull(self):
        key_0, value_0 = random_string(), random_string().encode()
        key_1, value_1 = key_0, random_string().encode()

        self.client.push(key_0, value_0)
        self.client.push(key_1, value_1)

        response = self.client.pull_without_ack()
        value = response.value

        self.assertEqual(value, value_0)

        response = self.client.pull()

        self.assertEqual(response, 'error')

        time.sleep(10)

        response = self.client.pull()
        value = response.value

        # The first pull request did not send the ack, so this one should also be [b'test value 1']
        self.assertEqual(value, value_0)

        response = self.client.pull()
        value = response.value

        self.assertEqual(value, value_1)


    def test_concurrent_push_without_order(self):
        for _ in range(50):

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

            real, expected = [r[0] for r in res], [r[1] for r in res]

            print(f'real: {real}, expected: {expected}')

            self.assertSetEqual(set(real), set(expected))

    def test_read_partition(self):
        client = self.client

        push_keys = [''.join(random.choices(string.ascii_uppercase + string.digits, k=10)) for _ in range(100)]
        push_values = [''.join(random.choices(string.ascii_uppercase + string.digits, k=10)).encode('utf-8') for _ in
                       range(100)]

        for push_key, push_value in zip(push_keys, push_values):
            client.push(push_key, push_value, is_replica=False)

        first_partition_messages = [(key, value)
                                    for key, value in zip(push_keys, push_values) if
                                    hash_function(key, 4) == 0]

        res = client.read_partition(partition_id=0, is_replica=False)
        res = [(r.key, r.value) for r in res]

        self.assertListEqual(res, first_partition_messages)

    def test_write_partition_and_replica(self):
        client = self.client

        push_keys = [''.join(random.choices(string.ascii_uppercase + string.digits, k=10)) for _ in range(100)]
        push_values = [''.join(random.choices(string.ascii_uppercase + string.digits, k=10)).encode('utf-8') for _ in
                       range(100)]

        for push_key, push_value in zip(push_keys, push_values):
            client.push(push_key, push_value, is_replica=False)

        first_partition_messages = [(key, value)
                                    for key, value in zip(push_keys, push_values) if
                                    hash_function(key, 4) == 0]

        res = client.read_partition(partition_id=0, is_replica=False)

        client.write_partition(partition_id=0, is_replica=True, partition_messages=res)

        res = client.read_partition(partition_id=0, is_replica=True)
        res = [(r.key, r.value) for r in res]

        self.assertListEqual(res, first_partition_messages)

    def test_purge_replica(self):
        client = self.client

        push_keys = [''.join(random.choices(string.ascii_uppercase + string.digits, k=10)) for _ in range(100)]
        push_values = [''.join(random.choices(string.ascii_uppercase + string.digits, k=10)).encode('utf-8') for _ in
                       range(100)]

        for push_key, push_value in zip(push_keys, push_values):
            client.push(push_key, push_value, is_replica=False)

        first_partition_messages = [(key, value)
                                    for key, value in zip(push_keys, push_values) if
                                    hash_function(key, 4) == 0]

        res_main = client.read_partition(partition_id=0, is_replica=False)

        client.write_partition(partition_id=0, is_replica=True, partition_messages=res_main)

        client.purge_replicas()

        res = client.read_partition(partition_id=0, is_replica=True)

        res = [(r.key, r.value) for r in res]

        self.assertListEqual(res, [])
        self.assertEqual(first_partition_messages, [(r.key, r.value) for r in res_main])

    def test_push_replica(self):
        client = self.client

        push_keys = [''.join(random.choices(string.ascii_uppercase + string.digits, k=10)) for _ in range(100)]
        push_values = [''.join(random.choices(string.ascii_uppercase + string.digits, k=10)).encode('utf-8') for _ in
                       range(100)]

        for push_key, push_value in zip(push_keys, push_values):
            client.push(push_key, push_value, is_replica=True)

        first_partition_messages = [(key, value)
                                    for key, value in zip(push_keys, push_values) if
                                    hash_function(key, 4) == 0]

        res = client.read_partition(partition_id=0, is_replica=True)
        res = [(r.key, r.value) for r in res]

        self.assertListEqual(res, first_partition_messages)

    def test_push_ack_replica(self):
        client = self.client
        client.purge_replicas()

        push_keys = [''.join(random.choices(string.ascii_uppercase + string.digits, k=10)) for _ in range(100)]
        push_values = [''.join(random.choices(string.ascii_uppercase + string.digits, k=10)).encode('utf-8') for _ in
                       range(100)]

        for push_key, push_value in zip(push_keys, push_values):
            client.push(push_key, push_value, is_replica=True)

        # remove messages from replica data
        for i in range(50):
            key, value = push_keys.pop(), push_values.pop()
            client.ack(key=key, is_replica=True)

        first_partition_messages = [(key, value)
                                    for key, value in zip(push_keys, push_values) if
                                    hash_function(key, 4) == 0]

        res = client.read_partition(partition_id=0, is_replica=True)
        res = [(r.key, r.value) for r in res]

        self.assertListEqual(res, first_partition_messages)
