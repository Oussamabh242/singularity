import socket

HOST = "127.0.0.1"  # The server's hostname or IP address
PORT = 1234  # The port used by the server


content = bytes([7, 10, 0 , 0])

with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
    s.connect((HOST, PORT))
    s.sendall(content)
    msg = s.recv(1024)
    print(msg)

