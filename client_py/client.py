from typing import List
import asyncio
import logging

import grpc
# import google.protobuf.empty_pb2 as empty_pb2
from google.protobuf import empty_pb2 as _empty_pb2

from client_py import queue_pb2_grpc
from client_py import queue_pb2

HOST, PORT = "localhost", "8888"


class QueueClient:
    stub = None
    HOST, PORT = "localhost", "8888"

    @classmethod
    def get_stub(cls, host: str, port: str):
        if cls.stub is None:
            channel = grpc.insecure_channel(f"{host}:{port}")
            cls.stub = queue_pb2_grpc.QueueStub(channel)
        return cls.stub

    def push(self, key: str, value: List[bytes]):
        try:
            stub = self.get_stub(HOST, PORT)

            stub.Push(queue_pb2.PushRequest(key=key, value=value))

        except grpc.RpcError as e:
            print(f"Error in pushing: {e}.")

    def pull(self):
        try:
            stub = self.get_stub(HOST, PORT)
            response = stub.Pull(f())
            print(f"key and message: {response.key} - {response.value}")
            return response
        except grpc.RpcError as e:
            print(f"Error in pulling: {e}.")

    def ack(self, acknowledgement: str):
        try:
            stub = self.get_stub(HOST, PORT)
            ack_request = queue_pb2.AcknowledgePullRequest(key=acknowledgement)
            stub.AcknowledgePull(ack_request)
            return True
        except grpc.RpcError as e:
            print(f"Error in acknowledgement: {e}")
            return False

    async def non_blocking_pull(self):
        try:
            stub = self.get_stub(HOST, PORT)
            response = stub.Pull(f())
            print(f"key and message: {response.key} - {response.value}")
            await self.non_blocking_ack(response.key)
            return response
        except grpc.RpcError as e:
            print(f"Error in pulling: {e}.")

    async def non_blocking_ack(self, acknowledgement: str):
        try:
            stub = self.get_stub(HOST, PORT)
            ack_request = queue_pb2.AcknowledgePullRequest(key=acknowledgement)
            stub.AcknowledgePull(ack_request)
            return True
        except grpc.RpcError as e:
            print(f"Error in acknowledgement: {e}")
            return False

    async def subscribe(self):
        try:
            while True:
                await self.non_blocking_pull()
        except grpc.RpcError as e:
            print(f"Error in pulling: {e}.")
        pass


def f():
    return _empty_pb2.Empty()
