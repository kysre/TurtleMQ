# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: leader.proto
# Protobuf Python Version: 4.25.0
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import symbol_database as _symbol_database
from google.protobuf.internal import builder as _builder
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from google.protobuf import empty_pb2 as google_dot_protobuf_dot_empty__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x0cleader.proto\x12\x06leader\x1a\x1bgoogle/protobuf/empty.proto\"%\n\x12\x41\x64\x64\x44\x61taNodeRequest\x12\x0f\n\x07\x61\x64\x64ress\x18\x01 \x01(\t2\x88\x01\n\x06Leader\x12;\n\tIsHealthy\x12\x16.google.protobuf.Empty\x1a\x16.google.protobuf.Empty\x12\x41\n\x0b\x41\x64\x64\x44\x61taNode\x12\x1a.leader.AddDataNodeRequest\x1a\x16.google.protobuf.EmptyB-Z+github.com/kysre/TurtleMQ/leader/pkg/leaderb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'leader_pb2', _globals)
if _descriptor._USE_C_DESCRIPTORS == False:
  _globals['DESCRIPTOR']._options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z+github.com/kysre/TurtleMQ/leader/pkg/leader'
  _globals['_ADDDATANODEREQUEST']._serialized_start=53
  _globals['_ADDDATANODEREQUEST']._serialized_end=90
  _globals['_LEADER']._serialized_start=93
  _globals['_LEADER']._serialized_end=229
# @@protoc_insertion_point(module_scope)
