from typing import List

import grpc
import google.protobuf.empty_pb2 as empty_pb2

from client_py import queue_pb2_grpc
from client_py import queue_pb2

HOST, PORT = "localhost", "8888"


class QueueClient:
    stub = None

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
            print(f"values: {value}")
            # for item in value:
            #     print(type(item), item)
            #     print(type(req.value))
            #     req.value.extend(item)
            req.value.extend(value)
            stub.Push(req)
        except grpc.RpcError as e:
            print(f"Error in pushing: {e}.")

    def pull(self):
        try:
            stub = self.get_stub(HOST, PORT)
            # print('here')
            # response = queue_pb2.PullResponse(empty_pb2)
            response = stub.Pull(empty_pb2)
            print(f"got response: {response}")
            for item in response.value:
                print(f"KEY: {response.key}, VALUE: {item}")
        except grpc.RpcError as e:
            print(f"Error. Couldn't pull: {e}.")

    def subscribe(self):
        pass


def run():
    queue_clinet = QueueClient()
    # TEST 1
    key = 'k1'
    value = [b'hi.', b'goodbye.']
    queue_clinet.push(key=key, value=value)
    queue_clinet.pull()


if __name__ == "__main__":
    # print(type(bytearray([4,6])))
    # print(bytearray([4,6]))
    run()
