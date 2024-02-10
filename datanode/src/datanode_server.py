import grpc
from concurrent import futures
from google.protobuf import empty_pb2
import protos.datanode_pb2_grpc as datanode_pb2_grpc
import protos.datanode_pb2 as datanode_pb2
import leader_protos.leader_pb2_grpc as leader_pb2_grpc
import leader_protos.leader_pb2 as leader_pb2
from configs.configs import ConfigManager
from shared_partition import SharedPartitions
from shared_partition import clear_path
from loguru import logger


class DataNode(datanode_pb2_grpc.DataNodeServicer):
    def __init__(self, partition_count=1, home_path='datanode/server/'):
        self.home_path = home_path
        self.partition_count = partition_count
        self.shared_partition = SharedPartitions(partition_count, home_path=home_path + '/main/')
        self.replica = SharedPartitions(partition_count, home_path=home_path + '/replica/')


    def Push(self, request, context):
        logger.info(f"received a push message: {request.message}")
        if request.is_replica:
            self.replica.push(request.message)
        else:
            self.shared_partition.push(request.message)
        return empty_pb2.Empty()

    def Pull(self, request, context):
        try:
            logger.info(f"received a pull message: {request}")
            message = self.shared_partition.pull()
            response = datanode_pb2.PullResponse(message=message)
            return response
        except grpc.RpcError as e:
            logger.exception(e)
        except Exception as e:
            logger.exception(e)

    def WritePartition(self, request, context):
        try:
            logger.info(f"received partition write message for partition: {request.partition_index}")
            partition_messages = request.partition_messages
            partition_index = request.partition_index
            if request.is_replica:
                for message in partition_messages:
                    push_to_partition(partition_index, self.replica, message)
            else:
                for message in partition_messages:
                    push_to_partition(partition_index, self.shared_partition, message)
            return empty_pb2.Empty()
        except grpc.RpcError as e:
            logger.exception(e)
        except Exception as e:
            logger.exception(e)

    def ReadPartition(self, request, context):
        try:
            logger.info(f"received partition read message for partition: {request.partition_index}")
            partition_index = request.partition_index
            if request.is_replica:
                return datanode_pb2.ReadPartitionResponse(
                    partition_messages=self.replica.read_partition_non_blocking(partition_index))
            else:
                return datanode_pb2.ReadPartitionResponse(
                    partition_messages=self.shared_partition.read_partition_non_blocking(partition_index))
        except grpc.RpcError as e:
            logger.exception(e)
        except Exception as e:
            logger.exception(e)

    def PurgeReplicaData(self, request, context):
        logger.info('received purge replica request.')
        clear_path(f'{self.home_path}/replica')
        self.replica = SharedPartitions(partition_count=self.partition_count,
                                        home_path=self.home_path + '/replica/')
        return empty_pb2.Empty()

    def AcknowledgePull(self, request, context):
        try:
            key = request.key
            logger.info(f"received an acknowledge message: {key}")
            if request.is_replica:
                self.replica.acknowledge(key)
            else:
                self.shared_partition.acknowledge(key)
            return empty_pb2.Empty()
        except grpc.RpcError as e:
            logger.exception(f"Error in acknowledging. {e}")

    def IsHealthy(self, request, context):
        try:
            return empty_pb2.Empty()
        except grpc.RpcError as e:
            logger.exception(f"Error in acknowledging. {e}")

    def GetRemainingMessagesCount(self, request, context):
        try:
            count = self.shared_partition.get_remaining_messages_count()
            res = datanode_pb2.GetRemainingMessagesCountResponse(remaining_messages_count=count)
            return res
        except grpc.RpcError as e:
            logger.exception(f"Error in getting remaining messages count: {e}")


def push_to_partition(partition_index: int,
                      shared_partition: SharedPartitions,
                      partition_message):
    shared_partition.push(partition_message, partition_index)


def serve():
    port = ConfigManager.get_prop('server_port')
    partitions_count = int(ConfigManager.get_prop('partition_count'))
    home_path = ConfigManager.get_prop('partition_home_path')

    # remove data-storage
    clear_path(home_path)

    datanode_name = ConfigManager.get_prop('datanode_name')

    server = grpc.server(futures.ThreadPoolExecutor(max_workers=50))
    datanode_pb2_grpc.add_DataNodeServicer_to_server(DataNode(partitions_count, home_path), server)

    server.add_insecure_port('[::]:' + port)
    server.start()
    logger.info('Server started, listening on ' + port)

    # notify leader
    try:
        leader_host, leader_port = ConfigManager.get_prop('leader_host'), ConfigManager.get_prop('leader_port')
        channel = grpc.insecure_channel(f"{leader_host}:{leader_port}")
        stub = leader_pb2_grpc.LeaderStub(channel)
        add_request = leader_pb2.AddDataNodeRequest(address=f'{datanode_name}:{port}')
        stub.AddDataNode(add_request)
    except grpc.RpcError as e:
        logger.exception(f"Error in notifying leader: {e}.")

    server.wait_for_termination()


if __name__ == '__main__':
    serve()
