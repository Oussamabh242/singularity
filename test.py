import socket
import json
import time
import random
import sys 
import struct
HOST = "127.0.0.1"  # The server's hostname or IP address
PORT = 1234  # The port used by the server

class Packet :
    def __init__(self ,ptype , metadata ,  paylaod ) :
        self.ptype = ptype
        self.metadata = metadata
        self.meta_len = len(metadata) 
        self.paylaod = paylaod
        self.paylaod_len = len(paylaod)
        self.rem_len = self.meta_len +self.paylaod_len + 4
    def encode(self) : 
        
        marshalled_metadata =json.dumps( self.metadata).encode("utf-8")
        packet = bytes([self.ptype]) + self.rem_len.to_bytes(4, "big") + self.meta_len.to_bytes(2, "big") + marshalled_metadata + self.paylaod_len.to_bytes(2, "big") + self.paylaod.encode("utf-8")
        return packet
    @classmethod 
    def decode(cls , packet):
        #[3 , 0,0,0,4 , 0,2 , ---* , 4,0,---*]
        ptype = packet[0]
        rem_len = struct.unpack(">I" , packet[1:5])[0]
        meta_len = struct.unpack(">H",packet[5:7])[0]
        metadata = ""
        next_byte = 7
        if meta_len>0 : 
            metadata = packet[7:7+meta_len]
            next_byte = 7+meta_len
        paylaod_len = struct.unpack(">H" ,packet[next_byte:next_byte+2])[0]
        paylaod = ""
        if paylaod_len >0 : 
            paylaod = packet[next_byte+2:next_byte+2+paylaod_len]
        return cls(ptype ,metadata , paylaod)

        


def encoder(metadata :dict, message:str, type:int ):
    marshalled_metadata =json.dumps( metadata).encode("utf-8")
    
    packet = bytes([type]) + (len(marshalled_metadata) + len(message) + 2 + 2).to_bytes(4, "big") + len(marshalled_metadata).to_bytes(2, "big") + marshalled_metadata + len(message).to_bytes(2, "big") + message.encode("utf-8")

    return packet 

def decoder(packet) : 
   pass 


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
        # 5 , len(metadata+msg+4) , len(m) , m , 0
        
        packet = bytes([5]) + len(metadata).to_bytes(2, byteorder='big') + len(metadata).to_bytes(2, byteorder='big') +metadata + bytes([0])
        try :
            self.connect()
            self.sock.sendall(packet)
            ack = self.sock.recv(1024)
            # print(ack)
        finally :
            self.disconnect()
    def Publish(self , qname , msg : str ) :
        metadata =json.dumps( {
            "queue":qname
        }).encode("utf-8")
        msgByte = msg.encode("utf-8")

        packet = bytes([1]) + len(metadata).to_bytes(2, byteorder='big') + len(metadata).to_bytes(2, byteorder='big') +metadata + len(msgByte).to_bytes(2,byteorder="big" ) +msgByte
        
        try :
            self.connect()
            self.sock.sendall(packet)
            ack = self.sock.recv(1024)
            # print(ack)
        finally :
            self.disconnect()

    def Consume(self , qname, job ) :
        metadata ={
            "queue":qname
        }
        packet = encoder(metadata , "", 3) 
        
        self.connect()
        self.sock.sendall(packet)
        ack = self.sock.recv(4)         
        if int(ack[0]) == 4 :
            # length = int.from_bytes(byte_data, byteorder='big')
            while 1 : 
                type_len = self.sock.recv(5)
                length = struct.unpack(">I" , type_len[1:5])[0]
                msg = self.sock.recv(length)
                decoded = Packet.decode(type_len+msg) 
                job(decoded.paylaod)
                self.sock.sendall(b"ok")
     
def job(msg:str) :
    # wait = random.randint(0,5) 
    wait = 5
    # print("working on " , msg , wait)

    time.sleep(wait) 
    print(sys.argv[2] ,": job on " ,msg, " finished ." )



sing = Singularity(PORT ,HOST)
# # sing.CreateQueue("someQueue")
# # sing.Publish("someQueue" , "hello there 1")
# # sing.Publish("someQueue" , "hello there 2")
# # sing.Publish("someQueue" , "hello there 3")
#
print("client :" , sys.argv[2] ," starting")
sing.Consume(sys.argv[1], job)

#
# print("meta data only packet")
# content = bytes([5, 0, 0, 0, 21, 0, 17, 123, 34, 113, 117, 101, 117, 101, 34, 58, 32, 34, 116, 104, 105, 115, 34, 125, 0, 0])
# Packet.decode(content)
# print("----------------\nful packet")
# content = bytes([1,00,00, 18, 52, 0, 17, 123, 34, 113, 117, 101, 117, 101, 34, 58, 32, 34, 116, 104, 97, 116, 34, 125, 0, 12, 104, 101, 108, 108, 111, 32, 116, 104, 101, 114, 101, 48])
# Packet.decode(content)

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

# p = encoder({"queue":"this"},"" , 5)
# print(list(p[0]) , p[1],p[2])
