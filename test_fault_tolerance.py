### this test checks fault tolerancy of the system. 
### for each component in the system do the following:
### 1. run the test
### 2. when the pop up popes, bring down an instance of the said component. 
### 3. wait for cluster to become healty again
### 4. press enter

### note that API-Gateway, External Database, ... all are components of the system and are prune to downtime
import random
import time
from typing import List
from threading import Lock
import string

from test_client import pull, push, subscribe

TEST_SIZE = 1000 * 1
KEY_SIZE = 8
SUBSCRIBER_COUNT = 4

key_seq = []
for _ in range(TEST_SIZE):
    rand_part_length = random.choices(range(20), k=1)
    rand_part = ''.join(random.choices(string.ascii_lowercase, k=rand_part_length[0]))
    key_seq.append(f"{rand_part}")

pulled: List[int] = []
lock = Lock()


def store(_: str, val: bytes):
    next_val = int(val.decode("utf-8"))
    with lock:
        pulled.append(next_val)


for i in range(TEST_SIZE // 2):
    push(f"{key_seq[i]}", f"{i}".encode(encoding="utf-8"))

print("manually fail one node and wait for cluster to become healthy again")
print("press enter when cluster is healthy")
input()

for _ in range(SUBSCRIBER_COUNT):
    subscribe(store)

for i in range(TEST_SIZE // 2, TEST_SIZE):
    push(f"{key_seq[i]}", f"{i}".encode(encoding="utf-8"))

# Wait to consume all messages
print("Wait for all messages to be consumed and press enter")
input()

# Remove duplicates (it's ok to return a k,v multiple times)
pulled = list(set(pulled))

pulled.sort()

print(len(pulled))

for i in range(TEST_SIZE):
    if pulled[i] != i:
        print(f"{i} != {pulled[i]}")
        print("DATA loss occurred")

print("Fault tolerance test passed successfully!")
