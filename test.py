import socket

HOST = "127.0.0.1"  # The server's hostname or IP address
PORT = 1234  # The port used by the server
x = 1024
print(x.to_bytes(2 , 'big'))

content = bytes([1, 136, 0, 128, 42,10, 123, 10, 34, 113, 117, 101, 117, 101, 34, 32, 58, 32, 34,
  101, 120, 97, 109, 112, 108, 101, 34, 44, 10, 34, 116, 111, 112,
  105, 99, 34, 32, 58, 32, 34, 115, 101, 120, 34, 10, 125, 10 ,11, 72, 101, 108, 108, 111, 32, 87, 111, 114, 108, 100])

with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
    s.connect((HOST, PORT))
    s.sendall(content)
    d = s.recv(1024)
    print(d)
    if int(d[0]) == 4 :
        print("subscribing")
        while 1 :
            data = s.recv(1024)
            print(data)
            break
