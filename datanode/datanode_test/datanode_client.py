from typing import List

import grpc
import google.protobuf.empty_pb2 as empty_pb2

from datanode.src.protos import datanode_pb2_grpc
from protos import datanode_pb2
from datanode.src.configs.configs import ConfigManager


class QueueClient:
    stub = None
    HOST, PORT = 'localhost', ConfigManager.get_prop('server_port')

    @classmethod
    def get_stub(cls):
        if cls.stub is None:
            channel = grpc.insecure_channel(f"{QueueClient.HOST}:{QueueClient.PORT}")
            cls.stub = datanode_pb2_grpc.DataNodeStub(channel)
        return cls.stub

    def push(self, key: str, value, is_replica=False):
        try:
            stub = self.get_stub()
            message = datanode_pb2.QueueMessage(key=key, value=value)
            stub.Push(datanode_pb2.PushRequest(message=message, is_replica=is_replica))
        except grpc.RpcError as e:
            print(f"Error in pushing: {e}.")

    def pull(self):
        try:
            stub = self.get_stub()
            response = stub.Pull(empty_pb2.Empty())
            message = response.message
            print(f"key and message: {message.key} - {message.value}")
            self.ack(message.key, is_replica=False)
            return message
        except grpc.RpcError as e:
            raise TimeoutError

    def pull_without_ack(self):
        try:
            stub = self.get_stub()
            response = stub.Pull(empty_pb2.Empty())
            message = response.message
            return message
        except grpc.RpcError as e:
            raise TimeoutError

    def ack(self, key: str, is_replica: bool):
        try:
            stub = self.get_stub()
            ack_request = datanode_pb2.AcknowledgePullRequest(key=key,
                                                              is_replica=is_replica)
            stub.AcknowledgePull(ack_request)
            return True
        except grpc.RpcError as e:
            print(f"Error in acknowledgement: {e}")
            return False

    def read_partition(self, partition_id: int, is_replica: bool):
        try:
            stub = self.get_stub()
            read_request = datanode_pb2.ReadPartitionRequest(partition_index=partition_id,
                                                             is_replica=is_replica)
            response = stub.ReadPartition(read_request)
            return response.partition_messages
        except grpc.RpcError as e:
            print(f"Error in reading: {e}.")

    def write_partition(self, partition_id: int,
                        is_replica: bool,
                        partition_messages: List[datanode_pb2.QueueMessage]):

        try:
            stub = self.get_stub()
            write_request = datanode_pb2.WritePartitionRequest(partition_index=partition_id,
                                                               is_replica=is_replica,
                                                               partition_messages=partition_messages)
            stub.WritePartition(write_request)
        except grpc.RpcError as e:
            print(f"Error in writing: {e}.")

    def purge_replicas(self):
        try:
            stub = self.get_stub()
            stub.PurgeReplicaData(empty_pb2.Empty())
        except grpc.RpcError as e:
            print(f"Error in PurgeReplicas: {e}.")

    def purge_main(self):
        try:
            stub = self.get_stub()
            stub.PurgeMainData(empty_pb2.Empty())
        except grpc.RpcError as e:
            print(f"Error in PurgeMain: {e}.")
