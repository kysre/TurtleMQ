### this test checks if the cluster performance scales with the size of the cluster
### run the test and press enter untill you see a thourghput cap in the monitoring of the cluster
### manually scale up cluster 
### press enter once more and see if you have a increase in throughput rate of the cluster

from test_client import subscribe, pull, push
from client_py.client import QueueClient
import string
import random
from threading import Thread

TEST_SIZE = 1000 * 1000
KEY_SIZE = 8
SUBSCRIBER_COUNT = 4


def to_infinity():
    index = 1
    while True:
        yield index
        index += 1


def push_key(key: str):
    qc = QueueClient()
    for j in to_infinity():
        rand_part_length = random.choices(range(20), k=1)
        rand_part = ''.join(random.choices(string.ascii_lowercase, k=rand_part_length[0]))
        qc.push(f"{key}-{rand_part}", f"{j}".encode("utf-8"))
        print(f"{key}-{rand_part}", f"{j}".encode("utf-8"))


subscribe(lambda key, val: ...)
for i in to_infinity():
    t = Thread(target=push_key, args=(i,))
    t.start()
    print("did it cap?")
    print("if not, press enter to increase throughput")
    print("if capped, manually scale up the cluster and press enter to see if you can increase the throughput")
    input()
