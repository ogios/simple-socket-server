import socket

HOST = "127.0.0.1"
PORT = 15001

a = "push\ndata"
b = a.encode("utf-8")
s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

s.connect((HOST, PORT))
s.send(b)
s.close()
