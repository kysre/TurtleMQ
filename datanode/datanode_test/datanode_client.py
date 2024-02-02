from typing import List

import grpc
import google.protobuf.empty_pb2 as empty_pb2

from datanode.src.protos import datanode_pb2_grpc
from protos import datanode_pb2

HOST, PORT = "localhost", "1234"


class QueueClient:
    stub = None
    HOST, PORT = HOST, PORT

    @classmethod
    def get_stub(cls, host: str, port: str):
        if cls.stub is None:
            channel = grpc.insecure_channel(f"{host}:{port}")
            cls.stub = datanode_pb2_grpc.DataNodeStub(channel)
        return cls.stub

    async def push(self, key: str, value: List[bytes]):
        try:
            stub = self.get_stub(HOST, PORT)
            message = datanode_pb2.QueueMessage(key=key, value=value)
            stub.Push(datanode_pb2.PushRequest(message=message))

        except grpc.RpcError as e:
            print(f"Error in pushing: {e}.")

    async def pull(self):
        try:
            stub = self.get_stub(HOST, PORT)
            response = stub.Pull(empty_pb2.Empty())
            message = response.message
            print(f"key and message: {message.key} - {message.value}")
            await self.ack(message.key)
            return message
        except grpc.RpcError as e:
            print(f"Error in pulling: {e}.")
            return 'error'

    async def pull_without_ack(self):
        try:
            stub = self.get_stub(HOST, PORT)
            response = stub.Pull(empty_pb2.Empty())
            message = response.message
            print(f"key and message: {message.key} - {message.value}")
            return message
        except grpc.RpcError as e:
            print(f"Error in pulling: {e}.")

    async def ack(self, acknowledgement: str):
        try:
            stub = self.get_stub(HOST, PORT)
            ack_request = datanode_pb2.AcknowledgePullRequest(key=acknowledgement)
            stub.AcknowledgePull(ack_request)
            return True
        except grpc.RpcError as e:
            print(f"Error in acknowledgement: {e}")
            return False

    async def subscribe(self):
        try:
            while True:
                await self.pull()
        except grpc.RpcError as e:
            print(f"Error in pulling: {e}.")
