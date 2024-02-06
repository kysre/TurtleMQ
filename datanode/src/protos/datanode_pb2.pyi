from google.protobuf import empty_pb2 as _empty_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class QueueMessage(_message.Message):
    __slots__ = ("key", "value")
    KEY_FIELD_NUMBER: _ClassVar[int]
    VALUE_FIELD_NUMBER: _ClassVar[int]
    key: str
    value: _containers.RepeatedScalarFieldContainer[bytes]
    def __init__(self, key: _Optional[str] = ..., value: _Optional[_Iterable[bytes]] = ...) -> None: ...

class PushRequest(_message.Message):
    __slots__ = ("message", "is_replica")
    MESSAGE_FIELD_NUMBER: _ClassVar[int]
    IS_REPLICA_FIELD_NUMBER: _ClassVar[int]
    message: QueueMessage
    is_replica: bool
    def __init__(self, message: _Optional[_Union[QueueMessage, _Mapping]] = ..., is_replica: bool = ...) -> None: ...

class PullResponse(_message.Message):
    __slots__ = ("message",)
    MESSAGE_FIELD_NUMBER: _ClassVar[int]
    message: QueueMessage
    def __init__(self, message: _Optional[_Union[QueueMessage, _Mapping]] = ...) -> None: ...

class AcknowledgePullRequest(_message.Message):
    __slots__ = ("key", "is_replica")
    KEY_FIELD_NUMBER: _ClassVar[int]
    IS_REPLICA_FIELD_NUMBER: _ClassVar[int]
    key: str
    is_replica: bool
    def __init__(self, key: _Optional[str] = ..., is_replica: bool = ...) -> None: ...

class ReadPartitionResponse(_message.Message):
    __slots__ = ("partition_messages",)
    PARTITION_MESSAGES_FIELD_NUMBER: _ClassVar[int]
    partition_messages: _containers.RepeatedCompositeFieldContainer[QueueMessage]
    def __init__(self, partition_messages: _Optional[_Iterable[_Union[QueueMessage, _Mapping]]] = ...) -> None: ...

class ReadPartitionRequest(_message.Message):
    __slots__ = ("partition_index", "is_replica")
    PARTITION_INDEX_FIELD_NUMBER: _ClassVar[int]
    IS_REPLICA_FIELD_NUMBER: _ClassVar[int]
    partition_index: int
    is_replica: bool
    def __init__(self, partition_index: _Optional[int] = ..., is_replica: bool = ...) -> None: ...

class WritePartitionRequest(_message.Message):
    __slots__ = ("partition_index", "is_replica", "partition_messages")
    PARTITION_INDEX_FIELD_NUMBER: _ClassVar[int]
    IS_REPLICA_FIELD_NUMBER: _ClassVar[int]
    PARTITION_MESSAGES_FIELD_NUMBER: _ClassVar[int]
    partition_index: int
    is_replica: bool
    partition_messages: _containers.RepeatedCompositeFieldContainer[QueueMessage]
    def __init__(self, partition_index: _Optional[int] = ..., is_replica: bool = ..., partition_messages: _Optional[_Iterable[_Union[QueueMessage, _Mapping]]] = ...) -> None: ...

class GetRemainingMessagesCountResponse(_message.Message):
    __slots__ = ("remaining_messages_count",)
    REMAINING_MESSAGES_COUNT_FIELD_NUMBER: _ClassVar[int]
    remaining_messages_count: int
    def __init__(self, remaining_messages_count: _Optional[int] = ...) -> None: ...
