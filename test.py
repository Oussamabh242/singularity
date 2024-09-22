import socket
import json
import time
HOST = "127.0.0.1"  # The server's hostname or IP address
PORT = 1234  # The port used by the server


class Singularity :
    def __init__(self , port , host) -> None:
        self.port = port 
        self.host = host
        self.sock = socket.socket(socket.AF_INET , socket.SOCK_STREAM)
    def connect(self) :
        self.sock.connect((self.host , self.port))
    def disconnect(self) :
        self.sock.close()
    def CreateQueue(self , qname ) :
        metadata =json.dumps( {
            "queue":qname
        }).encode("utf-8")
        # 5 , len(n)+2 , len(n) , n , 0
        packet= bytes([5 ,len(metadata)+2 , len(metadata)  ]+list(metadata)+[0])
        try :
            self.connect()
            self.sock.sendall(packet)
            ack = self.sock.recv(1024)
            print(ack)
        finally :
            self.disconnect()
    def Publish(self , qname , msg : str ) :
        metadata =json.dumps( {
            "queue":qname
        }).encode("utf-8")
        msgByte = list(msg.encode("utf-8"))

        packet= bytes([1 ,len(metadata)+2 , len(metadata)  ]+list(metadata)+[len(msgByte)]+msgByte)

        try :
            self.connect()
            self.sock.sendall(packet)
            ack = self.sock.recv(1024)
            print(ack)
        finally :
            self.disconnect()

    def Consume(self , qname ) :
        metadata =json.dumps( {
            "queue":qname
        }).encode("utf-8")
        packet= bytes([3 ,len(metadata)+2 , len(metadata)  ]+list(metadata)+[0])

        
        self.connect()
        self.sock.sendall(packet)
        ack = self.sock.recv(50)         
        if int(ack[0]) == 4 :
            print("subscribing")
            while 1 : 
                job = self.sock.recv(50)
                print("working on" ,str(job))
                time.sleep(4)
                self.sock.sendall(b"ok")
     

sing = Singularity(PORT ,HOST)
# sing.CreateQueue("someQueue")
# sing.Publish("someQueue" , "hello there 1")
# sing.Publish("someQueue" , "hello there 2")
# sing.Publish("someQueue" , "hello there 3")
# sing.Consume("someQueue")


content = bytes([3, 136, 0, 128, 42,10, 123, 10, 34, 113, 117, 101, 117, 101, 34, 32, 58, 32, 34,
  101, 120, 97, 109, 112, 108, 101, 34, 44, 10, 34, 116, 111, 112,
  105, 99, 34, 32, 58, 32, 34, 115, 101, 120, 34, 10, 125, 10 ,11, 72, 101, 108, 108, 111, 32, 87, 111, 114, 108, 100])

# with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
#     s.connect((HOST, PORT))
#
#     s.sendall(content)
#     d = s.recv(1024)
#     print(d)
#     if int(d[0]) == 4 :
#         print("subscribing")
#         while 1 :
#             data = s.recv(1024)
#             print(data)
