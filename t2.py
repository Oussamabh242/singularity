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
        
        packet = bytes([5]) + len(metadata).to_bytes(2, byteorder='big') + len(metadata).to_bytes(2, byteorder='big') +metadata + bytes([0])
        self.connect()
        self.sock.sendall(packet)
        ack = self.sock.recv(1024)
        print(ack)
        
    def Publish(self , qname , msg : str ) :
        metadata =json.dumps( {
            "queue":qname
        }).encode("utf-8")
        msgByte = msg.encode("utf-8")

        packet = bytes([1]) + (4660).to_bytes(2, byteorder='big') + len(metadata).to_bytes(2, byteorder='big') +metadata + len(msgByte).to_bytes(2,byteorder="big" ) +msgByte
        
        self.connect()
        self.sock.sendall(packet)
        ack = self.sock.recv(1024)
        print(ack)
    def Consume(self , qname ) :
        metadata =json.dumps( {
            "queue":qname
        }).encode("utf-8")
        packet = bytes([3]) + len(metadata).to_bytes(2, byteorder='big') + len(metadata).to_bytes(2, byteorder='big') +metadata + bytes([0])

        
        self.connect()
        self.sock.sendall(packet)
        ack = self.sock.recv(50)         
        if int(ack[0]) == 4 :
            print("subscribing")
            while 1 : 
                job = str(self.sock.recv(50))
                job = job[4:]
                print("working on" ,str(job))
                time.sleep(1)
                self.sock.sendall(b"ok")
     

import sys
sing = Singularity(PORT ,HOST)
sing.CreateQueue(sys.argv[1])
# sing.Publish("someQueue" , "hello there 1")
# sing.Publish("someQueue" , "hello there 2")
# sing.Publish("someQueue" , "hello there 3")

n = int(sys.argv[2])
for i in range(n) :
    sing = Singularity(PORT ,HOST)
    sing.Publish(sys.argv[1], "hello there"+str(i))
# sing.Consume("test")

