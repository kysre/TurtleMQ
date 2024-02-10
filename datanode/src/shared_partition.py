import json
import threading
import time
import os
import shutil
import random
from protos import datanode_pb2
from configs.utils import MessagesStatus, PartitionsBusyError, PartitionStatus
from configs.configs import ConfigManager
import threading
from loguru import logger
from configs.utils import hash_function


def clear_path(path):
    for filename in os.listdir(path):
        file_path = os.path.join(path, filename)
        try:
            if os.path.isfile(file_path) or os.path.islink(file_path):
                os.unlink(file_path)
            elif os.path.isdir(file_path):
                shutil.rmtree(file_path)
        except Exception as e:
            print('Failed to delete %s. Reason: %s' % (file_path, e))
    pass


class Message:
    encoding_method = ConfigManager.get_prop('encoding_method')

    def __init__(self, message: datanode_pb2.QueueMessage, file_address, read_write_lock):
        self.key = message.key

        self.read_write_lock = read_write_lock

        self.message_status = MessagesStatus.NOT_SENT
        self.file_address = file_address
        self.index = self.write_file_system(message.value)

    def get_message(self):
        line = None
        with open(self.file_address, 'r') as file_system:
            for _ in range(self.index + 1):
                line = file_system.readline()
        assert line is not None

        message = json.loads(line)

        key, decoded_value = message['key'], message['value']

        encoded_value = decoded_value.encode(Message.encoding_method)
        messageProto = datanode_pb2.QueueMessage(key=key, value=encoded_value)
        return messageProto

    def pend_message(self):
        self.message_status = MessagesStatus.PENDING
        self.write_file_system(b'')

    def sent_message(self):
        self.message_status = MessagesStatus.SENT
        self.write_file_system(b'')

    def set_message_free(self):
        self.message_status = MessagesStatus.NOT_SENT
        self.write_file_system(b'')

    def write_file_system(self, value):
        with self.read_write_lock:
            count_lines = 0
            with open(self.file_address, 'r') as file_system:
                for line in file_system:
                    count_lines += 1

            with open(self.file_address, 'a') as file_system:
                value = value.decode(Message.encoding_method)
                file_system.write(json.dumps(dict(key=self.key, value=value, status=self.message_status.value)) + '\n')

            return count_lines


class Partition:
    def __init__(self, path):
        self.messages = []
        self.dir = path

        self.messages_lock = threading.Lock()

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
        self.messages.append(Message(message, self.dir, self.messages_lock))

    def get_remaining_message_count(self) -> int:
        return len(self.messages)


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
            for available_partition in available_partitions:
                if len(self.partitions[available_partition].messages) > 0:
                    return available_partition
            raise PartitionsBusyError('There is no a free - non empty partition!')

    def push(self, message, partition_index=None):
        if partition_index is None:
            partition_index = self.get_dest_partition(message.key)
        destination_partition = self.partitions[partition_index]
        destination_partition.add_message(message)

    def acknowledge(self, key):
        partition_index = self.get_dest_partition(key)
        with self.partition_lock:
            self.partitions_status[partition_index] = PartitionStatus.FREE
            self.partitions[partition_index].remove_message()

    def get_dest_partition(self, key):
        return hash_function(key, self.partition_count)

    def get_remaining_messages_count(self) -> int:
        remaining_count = 0
        for partition in self.partitions:
            remaining_count += partition.get_remaining_message_count()
        return remaining_count


def clean_partitions(shared_partition: SharedPartitions):
    logger.info('cleaner is working.')
    pending_timeout = int(ConfigManager.get_prop('pending_timeout'))
    cleaner_period = int(ConfigManager.get_prop('cleaner_period'))

    while True:
        with shared_partition.partition_lock:
            for i, partition_status in enumerate(shared_partition.partitions_status):
                partition = shared_partition.partitions[i]

                # Do nothing if partition is already free
                if partition_status == PartitionStatus.FREE:
                    continue

                last_pull_status = partition.last_pull

                assert last_pull_status['time'] is not None

                if time.time() - last_pull_status['time'] > pending_timeout:
                    shared_partition.partitions_status[i] = PartitionStatus.FREE
                    last_pull_status['message'].set_message_free()
                    logger.info(f'Partition {i} is now free!')

        time.sleep(cleaner_period)
