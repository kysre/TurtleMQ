# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: datanode.proto
# Protobuf Python Version: 4.25.0
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import symbol_database as _symbol_database
from google.protobuf.internal import builder as _builder
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from google.protobuf import empty_pb2 as google_dot_protobuf_dot_empty__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x0e\x64\x61tanode.proto\x12\x08\x64\x61tanode\x1a\x1bgoogle/protobuf/empty.proto\"*\n\x0cQueueMessage\x12\x0b\n\x03key\x18\x01 \x01(\t\x12\r\n\x05value\x18\x02 \x01(\x0c\"J\n\x0bPushRequest\x12\'\n\x07message\x18\x01 \x01(\x0b\x32\x16.datanode.QueueMessage\x12\x12\n\nis_replica\x18\x02 \x01(\x08\"7\n\x0cPullResponse\x12\'\n\x07message\x18\x01 \x01(\x0b\x32\x16.datanode.QueueMessage\"9\n\x16\x41\x63knowledgePullRequest\x12\x0b\n\x03key\x18\x01 \x01(\t\x12\x12\n\nis_replica\x18\x02 \x01(\x08\"K\n\x15ReadPartitionResponse\x12\x32\n\x12partition_messages\x18\x01 \x03(\x0b\x32\x16.datanode.QueueMessage\"C\n\x14ReadPartitionRequest\x12\x17\n\x0fpartition_index\x18\x01 \x01(\x05\x12\x12\n\nis_replica\x18\x02 \x01(\x08\"x\n\x15WritePartitionRequest\x12\x17\n\x0fpartition_index\x18\x01 \x01(\x05\x12\x12\n\nis_replica\x18\x02 \x01(\x08\x12\x32\n\x12partition_messages\x18\x03 \x03(\x0b\x32\x16.datanode.QueueMessage\"E\n!GetRemainingMessagesCountResponse\x12 \n\x18remaining_messages_count\x18\x01 \x01(\x05\x32\x87\x05\n\x08\x44\x61taNode\x12\x35\n\x04Push\x12\x15.datanode.PushRequest\x1a\x16.google.protobuf.Empty\x12\x36\n\x04Pull\x12\x16.google.protobuf.Empty\x1a\x16.datanode.PullResponse\x12K\n\x0f\x41\x63knowledgePull\x12 .datanode.AcknowledgePullRequest\x1a\x16.google.protobuf.Empty\x12P\n\rReadPartition\x12\x1e.datanode.ReadPartitionRequest\x1a\x1f.datanode.ReadPartitionResponse\x12I\n\x0eWritePartition\x12\x1f.datanode.WritePartitionRequest\x1a\x16.google.protobuf.Empty\x12\x42\n\x10PurgeReplicaData\x12\x16.google.protobuf.Empty\x1a\x16.google.protobuf.Empty\x12?\n\rPurgeMainData\x12\x16.google.protobuf.Empty\x1a\x16.google.protobuf.Empty\x12;\n\tIsHealthy\x12\x16.google.protobuf.Empty\x1a\x16.google.protobuf.Empty\x12`\n\x19GetRemainingMessagesCount\x12\x16.google.protobuf.Empty\x1a+.datanode.GetRemainingMessagesCountResponseb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'datanode_pb2', _globals)
if _descriptor._USE_C_DESCRIPTORS == False:
  DESCRIPTOR._options = None
  _globals['_QUEUEMESSAGE']._serialized_start=57
  _globals['_QUEUEMESSAGE']._serialized_end=99
  _globals['_PUSHREQUEST']._serialized_start=101
  _globals['_PUSHREQUEST']._serialized_end=175
  _globals['_PULLRESPONSE']._serialized_start=177
  _globals['_PULLRESPONSE']._serialized_end=232
  _globals['_ACKNOWLEDGEPULLREQUEST']._serialized_start=234
  _globals['_ACKNOWLEDGEPULLREQUEST']._serialized_end=291
  _globals['_READPARTITIONRESPONSE']._serialized_start=293
  _globals['_READPARTITIONRESPONSE']._serialized_end=368
  _globals['_READPARTITIONREQUEST']._serialized_start=370
  _globals['_READPARTITIONREQUEST']._serialized_end=437
  _globals['_WRITEPARTITIONREQUEST']._serialized_start=439
  _globals['_WRITEPARTITIONREQUEST']._serialized_end=559
  _globals['_GETREMAININGMESSAGESCOUNTRESPONSE']._serialized_start=561
  _globals['_GETREMAININGMESSAGESCOUNTRESPONSE']._serialized_end=630
  _globals['_DATANODE']._serialized_start=633
  _globals['_DATANODE']._serialized_end=1280
# @@protoc_insertion_point(module_scope)
