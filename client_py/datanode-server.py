import grpc
from concurrent import futures
from google.protobuf import empty_pb2
import queue_pb2_grpc
import queue_pb2


class DataNode(queue_pb2_grpc.QueueServicer):
    def __init__(self):
        self.queue_data = []

    def Push(self, request, context):
        print(f"received a push message:")
        key = request.key
        values = request.value
        self.queue_data.append((key, values))
        print(f"pushed. queue: {self.queue_data}")
        return empty_pb2.Empty()

    def Pull(self, request, context):
        try:
            key, values = self.queue_data.pop(0)
            print(key, values)
            print(f"server before pull: {self.queue_data}")
            response = queue_pb2.PullResponse(key=key, value=values)
            print(f"pulled. queue: {self.queue_data}")
            print(response)
            return response
        except:
            print("Queue is empty.")

    def AcknowledgePull(self, request, context):
        try:
            print('ack received!')
            return empty_pb2.Empty()
        except grpc.RpcError as e:
            print(f"Error in acknowledging. {e}")


def serve():
    port = "8888"
    host = "localhost"

    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    queue_pb2_grpc.add_QueueServicer_to_server(DataNode(), server)

    server.add_insecure_port('[::]:' + port)
    server.start()
    print('Server started, listening on ' + port)
    server.wait_for_termination()


if __name__ == '__main__':
    serve()