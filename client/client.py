class Client:
    global datanode

    def __init__(self):
        pass

    def push(self, key: str, value: []):
        datanode.append((key, value))

    def pull(self):
        if datanode:
            key, value = datanode.pop(0)
            return key, value
        else:
            return None, None
        pass

    def subscribe(self, func):
        pass


client = Client()

message0 = 'hi. first message.'  # k0
message1 = 'second message.'  # k1
message2 = 'testing... 2.'  # k0
message3 = 'testing done.'  # k2

b_mess0 = list(bytes(message0, 'utf-8'))
b_mess1 = list(bytes(message1, 'utf-8'))
b_mess2 = list(bytes(message2, 'utf-8'))
b_mess3 = list(bytes(message3, 'utf-8'))

test0 = ('k0', b_mess0)
test1 = ('k1', b_mess1)
test2 = ('k0', b_mess2)
test3 = ('k2', b_mess3)

datanode = [('k3', list(bytes('message in queue', 'utf-8'))), ('k3', list(bytes('message in queue2', 'utf-8')))]
# datanode = []

key, value = client.pull()
if key is None:
    pass
else:
    for e in value:
        print(chr(e), end='')
    print()

client.push(test0[0], test0[1])
print(datanode)
