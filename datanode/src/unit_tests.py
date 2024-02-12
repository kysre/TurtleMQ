import threading
from unittest import TestCase
from shared_partition import Message, Partition
from configs.utils import MessagesStatus
from protos import datanode_pb2
import random
import string
import os


def random_string():
    return ''.join(random.choices(string.ascii_uppercase + string.digits, k=10))


class MessagesTests(TestCase):
    def test_get_message(self):
        path = 'test.msq'
        with open(path, 'a') as file:
            file.write('')
        queue_message = datanode_pb2.QueueMessage(key=random_string(), value=random_string().encode('utf-8'))
        message = Message(queue_message, path, threading.Lock())
        self.assertEqual(message.get_message(), queue_message)
        os.remove(path)


class PartitionsTests(TestCase):

    def setUp(self):
        self.path = './test.msq'
        self.test_partition = Partition(self.path)

    def test_add_messages(self):
        queue_message = datanode_pb2.QueueMessage(key=random_string(), value=random_string().encode('utf-8'))
        self.test_partition.add_message(queue_message)
        self.assertEqual(self.test_partition.messages[0].get_message(), queue_message)

    def test_get_messages(self):
        queue_message = datanode_pb2.QueueMessage(key=random_string(), value=random_string().encode('utf-8'))
        self.test_partition.add_message(queue_message)
        self.assertEqual(self.test_partition.get_messages(), queue_message)

    def test_pending_state(self):
        queue_message = datanode_pb2.QueueMessage(key=random_string(), value=random_string().encode('utf-8'))
        self.test_partition.add_message(queue_message)
        self.test_partition.get_messages()
        self.assertEqual(self.test_partition.messages[0].message_status,
                         MessagesStatus.PENDING)

    def test_removing_messages(self):
        key, value = random_string(), random_string().encode('utf-8')
        queue_message = datanode_pb2.QueueMessage(key=key, value=value)
        self.test_partition.add_message(queue_message)
        message = self.test_partition.messages[0]
        self.test_partition.get_messages()
        self.test_partition.remove_message(key)

        self.assertEqual(message.message_status, MessagesStatus.SENT)
        self.assertListEqual(self.test_partition.messages, [])

    def test_message_generation(self):
        messages = [datanode_pb2.QueueMessage(key=random_string(),
                                              value=random_string().encode('utf-8'))
                    for _ in range(10)]
        for message in messages:
            self.test_partition.add_message(message)

        self.assertListEqual(list(self.test_partition.message_generator()), messages)

    def tearDown(self):
        os.remove(self.path)
