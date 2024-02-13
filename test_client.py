from client_py.client import QueueClient

qc = QueueClient()


def push(key: str, val: bytes):
    return qc.push(key, val)


def pull():
    return qc.pull()


def subscribe(action):
    qc.subscribe(action)
