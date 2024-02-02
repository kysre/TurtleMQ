import json
import threading
import time
import os
import random
from protos import datanode_pb2
from configs.utils import MessagesStatus, PartitionsBusyError, PartitionStatus
from configs.configs import ConfigManager
import threading
from loguru import logger


class Message:
    encoding_method = ConfigManager.get_prop('encoding_method')

    def __init__(self, message: datanode_pb2.QueueMessage, file_address):
        self.key = message.key
        self.message_status = MessagesStatus.NOT_SENT
        self.file_address = file_address
        self.index = self.write_file_system(message.value)
        logger.info(self.index)

    def get_message(self):
        line = None
        with open(self.file_address, 'r') as file_system:
            for _ in range(self.index + 1):
                line = file_system.readline()
        assert line is not None
        message = json.loads(line)
        value = message['value']
        return datanode_pb2.QueueMessage(key=self.key, value=[b.encode(Message.encoding_method) for b in value])

    def pend_message(self):
        self.message_status = MessagesStatus.PENDING
        self.write_file_system([])

    def sent_message(self):
        self.message_status = MessagesStatus.SENT
        self.write_file_system([])

    def set_message_free(self):
        self.message_status = MessagesStatus.NOT_SENT
        self.write_file_system([])

    def write_file_system(self, value):
        count_lines = 0
        with open(self.file_address, 'r') as file_system:
            for line in file_system:
                count_lines += 1

        with open(self.file_address, 'a') as file_system:
            value = [b.decode(Message.encoding_method) for b in value]
            file_system.write(json.dumps(dict(key=self.key, value=value, status=self.message_status.value)) + '\n')

        return count_lines


class Partition:
    def __init__(self, path):
        self.messages = []
        self.dir = path

        self.last_pull = {'message': None,
                          'time': None}

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

            self.last_pull['time'] = time.time()
            self.last_pull['message'] = response

        return response.get_message()

    def remove_message(self):
        if len(self.messages) != 0:
            message = self.messages.pop(0)
            message.sent_message()

    def add_message(self, message: datanode_pb2.QueueMessage):
        self.messages.append(Message(message, self.dir))


class SharedPartitions:
    pull_timeout = int(ConfigManager.get_prop('pull_timeout'))

    def __init__(self, partition_count, home_path, cleaner=True):
        self.partition_lock = threading.Lock()

        assert type(partition_count) is int
        self.partition_count = partition_count

        self.partitions_status = partition_count * [PartitionStatus.FREE]
        self.partitions = [Partition(path=f'{home_path}/partition_{i + 1}.tmq') for i in range(partition_count)]

        if cleaner:
            self.cleaner = threading.Thread(target=clean_partitions, args=[self])
            self.cleaner.start()

    def pull(self):
        start_time = time.time()
        while True:
            try:
                partition_index = self.get_free_partition()
                break
            except PartitionsBusyError:
                time.sleep(1)
                if time.time() - start_time > self.pull_timeout:
                    raise TimeoutError

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


def clean_partitions(shared_partition: SharedPartitions):
    logger.info('cleaner is working.')
    pending_timeout = int(ConfigManager.get_prop('pending_timeout'))
    cleaner_period = int(ConfigManager.get_prop('cleaner_period'))

    while True:
        for i, partition_status in enumerate(shared_partition.partitions_status):
            partition = shared_partition.partitions[i]

            # Do nothing if partition is already free
            if partition_status == PartitionStatus.FREE:
                continue

            last_pull_status = partition.last_pull

            assert last_pull_status['time'] is not None

            if time.time() - last_pull_status['time'] > pending_timeout:
                with shared_partition.partition_lock:
                    shared_partition.partitions_status[i] = PartitionStatus.FREE
                last_pull_status['message'].set_message_free()
                logger.info(f'Partition {i} is now free!')

        time.sleep(cleaner_period)
