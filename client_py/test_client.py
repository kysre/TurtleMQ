from unittest import TestCase
from client_py.client import QueueClient


class TestQueueClient(TestCase):
    def setup(self):
        self.client = QueueClient()

    def test_push(self):
        self.client.push("test key", [b'test value'])
        self.assertEqual(self.client.pull().value, [b'test value'])
        # self.fail()

    def test_pull(self):
        self.client.push("test_key", [b"test_value"])
        self.assertEqual(self.client.pull().value, [b"test_value"])
        # self.fail()

    def test_subscribe(self):
        self.fail()
