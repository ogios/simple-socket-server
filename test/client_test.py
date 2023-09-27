import socket


# a = [255, 255, 12, 0]
# b = bytes(a)
a = "fetch\ndata"
b = a.encode("utf-8")

HOST = "127.0.0.1"
PORT = 15001

s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

print("connecting...")
s.connect((HOST, PORT))
print("sending...")
s.send(b)
print("closing...")
s.close()
