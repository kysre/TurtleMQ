import time
from typing import List
import asyncio
import logging
from threading import Thread
import grpc
# import google.protobuf.empty_pb2 as empty_pb2
from google.protobuf import empty_pb2 as _empty_pb2
from concurrent.futures import ThreadPoolExecutor
from client_py import queue_pb2_grpc
from client_py import queue_pb2


class QueueClient:
    stub = None
    HOST, PORT = "localhost", "8000"
    SUBSCRIBE_WORKERS = 3
    SUBSCRIBE_SLEEP_TIMEOUT = 2

    @classmethod
    def get_stub(cls, host: str, port: str):
        if cls.stub is None:
            channel = grpc.insecure_channel(f"{host}:{port}")
            cls.stub = queue_pb2_grpc.QueueStub(channel)
        return cls.stub

    def push(self, key: str, value: bytes):
        try:
            stub = self.get_stub(self.HOST, self.PORT)

            stub.Push(queue_pb2.PushRequest(key=key, value=value))

        except grpc.RpcError as e:
            print(f"Error in pushing: {e}.")

    def pull(self) -> (str, bytes):
        try:
            stub = self.get_stub(self.HOST, self.PORT)
            response = stub.Pull(_empty_pb2.Empty())
            self.ack(response.key)
            return response.key, response.value
        except grpc.RpcError as e:
            print(f"Error in pulling: {e}.")

    def ack(self, acknowledgement: str):
        try:
            stub = self.get_stub(self.HOST, self.PORT)
            stub.AcknowledgePull(queue_pb2.AcknowledgePullRequest(key=acknowledgement))
            return None
        except grpc.RpcError as e:
            print(f"Error in acknowledgement: {e}")
            return False

    def subscribe(self, f):
        thread = Thread(target=self.run_subscribe, args=(f,))
        thread.start()

    def run_subscribe(self, f):
        try:
            futures = []
            while True:
                with ThreadPoolExecutor(max_workers=QueueClient.SUBSCRIBE_WORKERS) as executer:
                    for _ in range(QueueClient.SUBSCRIBE_WORKERS):
                        pull_response = self.pull()
                        if pull_response is not None and pull_response is not False:
                            futures.append(executer.submit(f, pull_response[0], pull_response[1]))
                        time.sleep(QueueClient.SUBSCRIBE_SLEEP_TIMEOUT)

                    _ = [future.result() for future in futures]

        except grpc.RpcError as e:
            print(f"Error in pulling: {e}.")
