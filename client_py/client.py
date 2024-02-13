import time
from typing import List
import asyncio
import logging
from threading import Thread
import grpc
from google.protobuf import empty_pb2 as _empty_pb2
from concurrent.futures import ThreadPoolExecutor
from client_py import queue_pb2_grpc
from client_py import queue_pb2


class QueueClient:
    stub = None
    replica_stub = None
    HOST = "64.226.122.208"
    PORT, REPLICA_PORT = "8000", "8001"
    SUBSCRIBE_WORKERS = 1
    SUBSCRIBE_SLEEP_TIMEOUT = 0.01

    @classmethod
    def get_stub(cls, host: str, port: str):
        if cls.stub is None:
            channel = grpc.insecure_channel(f"{host}:{port}")
            cls.stub = queue_pb2_grpc.QueueStub(channel)
        return cls.stub

    @classmethod
    def get_replica_stub(cls, host: str, port: str):
        if cls.replica_stub is None:
            channel = grpc.insecure_channel(f"{host}:{port}")
            cls.replica_stub = queue_pb2_grpc.QueueStub(channel)
        return cls.replica_stub

    def push(self, key: str, value: bytes):
        try:
            stub = self.get_stub(self.HOST, self.PORT)
            stub.Push(queue_pb2.PushRequest(key=key, value=value))

        except grpc.RpcError as e:
            try:
                stub = self.get_replica_stub(self.HOST, self.REPLICA_PORT)
                stub.Push(queue_pb2.PushRequest(key=key, value=value))
            except grpc.RpcError as e:
                pass

    def pull(self) -> (str, bytes):
        try:
            stub = self.get_stub(self.HOST, self.PORT)
            response = stub.Pull(_empty_pb2.Empty())
            self.ack(response.key)
            return response.key, response.value
        except grpc.RpcError as e:
            try:
                stub = self.get_replica_stub(self.HOST, self.REPLICA_PORT)
                response = stub.Pull(_empty_pb2.Empty())
                self.ack(response.key)
                return response.key, response.value
            except grpc.RpcError as e:
                return '', None

    def ack(self, acknowledgement: str):
        try:
            stub = self.get_stub(self.HOST, self.PORT)
            stub.AcknowledgePull(queue_pb2.AcknowledgePullRequest(key=acknowledgement))
            return None
        except grpc.RpcError as e:
            try:
                stub = self.get_replica_stub(self.HOST, self.REPLICA_PORT)
                stub.AcknowledgePull(queue_pb2.AcknowledgePullRequest(key=acknowledgement))
                return None
            except grpc.RpcError as e:
                return None

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
                        if pull_response is not None and pull_response is not False and pull_response[0] != '':
                            futures.append(executer.submit(f, pull_response[0], pull_response[1]))
                        else:
                            time.sleep(QueueClient.SUBSCRIBE_SLEEP_TIMEOUT)

                    _ = [future.result() for future in futures]

        except grpc.RpcError as e:
            print(f"Error in pulling: {e}.")
