from enum import Enum
import hashlib


class PartitionStatus(Enum):
    FREE = 1
    LOCKED = 2


class MessagesStatus(Enum):
    NOT_SENT = 1
    PENDING = 2
    SENT = 3


class PartitionsBusyError(Exception):
    pass


def hash_function(key: str, number_of_partitions: int) -> int:

    result = hashlib.md5(key.encode())
    hashed_bytes = result.digest()

    return int.from_bytes(hashed_bytes, byteorder='big') % number_of_partitions

