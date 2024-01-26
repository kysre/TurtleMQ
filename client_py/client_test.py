from concurrent import futures

# import google
import grpc
from google.protobuf import empty_pb2 as _empty_pb2

from client_py import queue_pb2
from client_py import queue_pb2_grpc
from client_py.client import QueueClient

HOST, PORT = "localhost", "8888"


class QueueServicer(queue_pb2_grpc.QueueServicer):
    def __init__(self):
        self.queue_data = []

    def Push(self, request, context):
        # print(f"received a push message:")
        key = request.key
        values = request.value
        self.queue_data.append((key, values))
        # print(f"push. queue: {self.queue_data}")
        return f()

    def Pull(self, request, context):
        try:
            key, values = self.queue_data.pop(0)
            response = queue_pb2.PullResponse(key=key, value=values)
            # print(f"pull. queue: {self.queue_data}")
            return response
        except:
            print("Queue is empty.")

    def AcknowledgePull(self, request, context):
        try:
            print()
            return f()
        except grpc.RpcError as e:
            print(f"Error in acknowledging. {e}")


def f():
    return _empty_pb2.Empty()


def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    queue_client = QueueClient

    queue_pb2_grpc.add_QueueServicer_to_server(QueueServicer(), server)

    server.add_insecure_port(f'[::]:{PORT}')
    server.start()
    print(f"Server started. Listening on port {PORT}.")
    # server.wait_for_termination()

    # TEST 1
    key = 'k1'
    value = [b'hi.', b'goodbye.']
    queue_client.push(queue_client, key=key, value=value)
    # queue_client.push(key=key, value=value)

    # TEST 2
    key = 'k2'
    value = [b'second', b'1234']
    queue_client.push(queue_client, key=key, value=value)

    queue_client.pull(queue_client)

    server.wait_for_termination()


if __name__ == '__main__':
    serve()
