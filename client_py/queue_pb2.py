# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: client_py/queue.proto
# Protobuf Python Version: 4.25.0
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import symbol_database as _symbol_database
from google.protobuf.internal import builder as _builder
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from google.protobuf import empty_pb2 as google_dot_protobuf_dot_empty__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x15\x63lient-py/queue.proto\x12\x06\x63lient\x1a\x1bgoogle/protobuf/empty.proto\")\n\x0bPushRequest\x12\x0b\n\x03key\x18\x01 \x01(\t\x12\r\n\x05value\x18\x02 \x03(\x0c\"*\n\x0cPullResponse\x12\x0b\n\x03key\x18\x01 \x01(\t\x12\r\n\x05value\x18\x02 \x03(\x0c\x32r\n\x05Queue\x12\x33\n\x04Push\x12\x13.client.PushRequest\x1a\x16.google.protobuf.Empty\x12\x34\n\x04Pull\x12\x16.google.protobuf.Empty\x1a\x14.client.PullResponseb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'client_py.queue_pb2', _globals)
if _descriptor._USE_C_DESCRIPTORS == False:
  DESCRIPTOR._options = None
  _globals['_PUSHREQUEST']._serialized_start=62
  _globals['_PUSHREQUEST']._serialized_end=103
  _globals['_PULLRESPONSE']._serialized_start=105
  _globals['_PULLRESPONSE']._serialized_end=147
  _globals['_QUEUE']._serialized_start=149
  _globals['_QUEUE']._serialized_end=263
# @@protoc_insertion_point(module_scope)
