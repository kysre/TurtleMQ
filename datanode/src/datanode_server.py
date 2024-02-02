import grpc
from concurrent import futures
from google.protobuf import empty_pb2
import protos.datanode_pb2_grpc as datanode_pb2_grpc
import protos.datanode_pb2 as datanode_pb2
from configs.configs import ConfigManager
from shared_partition import SharedPartitions
from loguru import logger


class DataNode(datanode_pb2_grpc.DataNodeServicer):
    def __init__(self, partition_count=1, home_path='datanode/server/'):

        self.shared_partition = SharedPartitions(partition_count, home_path=home_path)

    def Push(self, request, context):
        logger.info(f"received a push message: {request.message}")
        self.shared_partition.push(request.message)
        return empty_pb2.Empty()

    def Pull(self, request, context):
        try:
            message = self.shared_partition.pull()
            response = datanode_pb2.PullResponse(message=message)
            return response
        except grpc.RpcError as e:
            logger.exception(e)
        except Exception as e:
            logger.exception(e)

    def AcknowledgePull(self, request, context):
        try:
            key = request.key
            self.shared_partition.acknowledge(key)
            return empty_pb2.Empty()
        except grpc.RpcError as e:
            logger.exception(f"Error in acknowledging. {e}")


def serve():
    port = ConfigManager.get_prop('server_port')
    partitions_count = int(ConfigManager.get_prop('partition_count'))
    home_path = ConfigManager.get_prop('partition_home_path')

    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    datanode_pb2_grpc.add_DataNodeServicer_to_server(DataNode(partitions_count, home_path), server)

    server.add_insecure_port('[::]:' + port)
    server.start()
    logger.info('Server started, listening on ' + port)
    server.wait_for_termination()


if __name__ == '__main__':
    serve()
