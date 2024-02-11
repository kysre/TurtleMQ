import time

from client_py.client import QueueClient


def test_subscribe_function_1(key: str, value: bytes):
    print(f'{key}, {value}')


def test_subscribe_function_2(key: str, value: bytes):
    print(f'{key}, {value}')


if __name__ == '__main__':
    test_client_1 = QueueClient()
    test_client_2 = QueueClient()
    test_client_3 = QueueClient()

    test_client_2.subscribe(test_subscribe_function_1)
    test_client_3.subscribe(test_subscribe_function_2)

    for i in range(200):
        test_client_1.push(f'{i}', b'value')
        time.sleep(0.1)

    time.sleep(500)
