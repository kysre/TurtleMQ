from enum import Enum


class PartitionStatus(Enum):
    FREE = 1
    LOCKED = 2


class MessagesStatus(Enum):
    NOT_SENT = 1
    PENDING = 2
    SENT = 3


class PartitionsBusyError(Exception):
    pass
