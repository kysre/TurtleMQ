import json
import threading
import time
import os
from enum import Enum
import random
from loguru import logger
from protos import datanode_pb2
from utils import MessagesStatus, PartitionsBusyError, PartitionStatus


class Message:
    def __init__(self, message: datanode_pb2.QueueMessage, file_address):
        self.key = message.key
        self.value = list(message.value)
        self.message_status = MessagesStatus.NOT_SENT
        self.file_address = file_address
        self.write_file_system()

    def pend_message(self):
        self.message_status = MessagesStatus.PENDING
        self.write_file_system()

    def sent_message(self):
        self.message_status = MessagesStatus.SENT
        self.write_file_system()

    def write_file_system(self):
        with open(self.file_address, 'a') as file_system:
            value = [b.decode('utf-8') for b in self.value]
            file_system.write(json.dumps(dict(key=self.key, value=value, status=self.message_status.value)) + '\n')


class Partition:
    def __init__(self, path):
        self.partition_status = PartitionStatus.FREE
        self.messages = []
        self.dir = path

        os.makedirs(os.path.dirname(path), exist_ok=True)
        with open(path, 'a') as file:
            file.write('')

    def is_empty(self):
        for message in self.messages:
            if message.message_status == MessagesStatus.NOT_SENT:
                return False
        return True

    def get_messages(self):
        response = None
        if len(self.messages) != 0:
            response = self.messages[0]
            response.pend_message()
        return response

    def remove_message(self):
        if len(self.messages) != 0:
            message = self.messages.pop(0)
            message.sent_message()

    def add_message(self, message: datanode_pb2.QueueMessage):
        self.messages.append(Message(message, self.dir))


class SharedPartitions:
    def __init__(self, partition_count, home_path):
        self.partition_lock = threading.Lock()

        assert type(partition_count) is int
        self.partition_count = partition_count

        self.partitions_status = partition_count * [PartitionStatus.FREE]
        self.partitions = [Partition(path=f'{home_path}/partition_{i + 1}.tmq') for i in range(partition_count)]

    def pull(self):
        while True:
            try:
                partition_index = self.get_free_partition()
                break
            except PartitionsBusyError:
                time.sleep(1)

        with self.partition_lock:
            self.partitions_status[partition_index] = PartitionStatus.LOCKED

        return self.partitions[partition_index].get_messages()

    def get_free_partition(self):
        with self.partition_lock:
            available_partitions = []
            for i, status in enumerate(self.partitions_status):
                if status == PartitionStatus.FREE:
                    available_partitions.append(i)
            if len(available_partitions) == 0:
                raise PartitionsBusyError('There is no available partitions!')
            return random.choice(available_partitions)

    def push(self, message):
        partition_index = self.get_dest_partition(message.key)
        destination_partition = self.partitions[partition_index]
        destination_partition.add_message(message)

    def acknowledge(self, key):
        partition_index = self.get_dest_partition(key)
        with self.partition_lock:
            self.partitions_status[partition_index] = PartitionStatus.FREE
            self.partitions[partition_index].remove_message()

    def get_dest_partition(self, key):
        return 0
