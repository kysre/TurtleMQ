from typing import List

import grpc

from client_py import queue_pb2_grpc
from client_py import queue_pb2

HOST, PORT = "localhost", "8888"


class QueueClient:
    @classmethod
    def get_stub(cls, host: str, port: str):
        if cls.stub is None:
            channel = grpc.insecure_channel(f"{host}:{port}")
            cls.stub = queue_pb2_grpc.QueueStub(channel)
        return cls.stub

    def push(self, key: str, value: List[bytes]):
        try:
            stub = self.get_stub(HOST, PORT)
            req = queue_pb2.PushRequest(key=key)
            req.value.extend(value)
            stub.Push(req)
        except Exception:
            print("kir")

    def pull(self):
        pass

    def subscribe(self):
        pass


def run():
    queue_clinet = QueueClient()
    try:
        stub = queue_pb2_grpc.QueueStub(queue_clinet.stub)

    with grpc.insecure_channel(f"{HOST}:{PORT}") as channel:
        stub = queue_pb2_grpc.QueueStub(channel)

    queue_client = QueueClient()


if __name__ == "__main__":
    run()
