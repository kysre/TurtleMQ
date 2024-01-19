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

    @classmethod
    def get_stub(cls, host: str, port: str):
        if cls.stub is None:
            channel = grpc.insecure_channel(f"{host}:{port}")
            cls.stub = queue_pb2_grpc.QueueStub(channel)
        return cls.stub

    def push(self, key: str, value: List[bytes]):
        try:
            # logging.basicConfig(level=logging.INFO)
            stub = self.get_stub(HOST, PORT)
            # req = queue_pb2.PushRequest(key=key)
            # req.value.extend(value)

            # push_message = queue_pb2.PushRequest(key=key, value=value)

            stub.Push(queue_pb2.PushRequest(key=key, value=value))

        except grpc.RpcError as e:
            print(f"Error in pushing: {e}.")

    def pull(self):
        try:
            stub = self.get_stub(HOST, PORT)
            # print('here')
            print('reaches here')
            response = stub.Pull(_empty_pb2)
            print(f"got response: {response}")
            for item in response.value:
                print(f"KEY: {response.key}, VALUE: {item}")
        except grpc.RpcError as e:
            print(f"Error in pulling: {e}.")

    def subscribe(self):
        pass


def run():
    logging.basicConfig(level=logging.INFO)
    queue_clinet = QueueClient()

    # TEST 1
    key = 'k1'
    value = [b'hi.', b'goodbye.']
    # print(f"TRY TO PUSH MESSAGE: {key} : {value}")
    queue_clinet.push(key=key, value=value)

    # TEST 2
    key = 'k2'
    value = [b'second', b'1234']
    queue_clinet.push(key=key, value=value)

    queue_clinet.pull()
    # print('After pull')


if __name__ == "__main__":
    run()
