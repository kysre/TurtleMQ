from typing import List
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
            ack_message = 'acknowledged!'
            stub.AcknowledgePull(ack_message)
        except grpc.RpcError as e:
            print(f"Error in pulling: {e}.")

    def subscribe(self):
        pass


def f():
    return _empty_pb2.Empty()


def run():
    logging.basicConfig(level=logging.INFO)
    queue_clinet = QueueClient()

    # TEST 1
    key = 'k1'
    value = [b'hi.', b'goodbye.']
    queue_clinet.push(key=key, value=value)

    # TEST 2
    key = 'k2'
    value = [b'second', b'1234']
    queue_clinet.push(key=key, value=value)

    queue_clinet.pull()


if __name__ == "__main__":
    run()
