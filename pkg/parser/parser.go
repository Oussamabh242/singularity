package parser

import (
	"fmt"
)

const (
    PUBLISH     = 0x30 // Publish     
    SUBSCRIBE   = 0x82 // Subscribe to topic
    UNSUBSCRIBE = 0xA2 // Unsubscribe from topics
    CONNECT     = 0x10 // Initiate connection
    CONNACK     = 0x20 // Acknowledgment of connection
    DISCONNECT  = 0xE0 // Disconnect from the broker
    PINGREQ     = 0xC0 // Ping request
    PINGRESP    = 0xD0 // Ping response
    PUBACK      = 0x40 // Acknowledgment for publish (QoS 1)
    PUBREC      = 0x50 // Publish received (QoS 2)
    PUBREL      = 0x62 // Publish release (QoS 2)
    PUBCOMP     = 0x70 // Publish complete (QoS 2)

    BYTE_MASK  = 0x7F  // Get the first 7 bit and ignore the first one ()
)

type Packet struct {
  PacketType string
  RLenght    uint
  TopicLen   uint
  TopicName  string
  PayloadLen uint 
  Payload    string
}

func Parse(packet  []byte) Packet  {

  pType := packetTypeToString(packet[0])
  rLenght , nextByte := getLength(packet[1 : 5] , 1 ) // at most 4 bytes starting from 2nd bit 
  fmt.Println("1" , pType , rLenght ,nextByte)
  tLength , nextByte := getLength(packet[nextByte:nextByte+2] , nextByte) // at most 2 bytes
  fmt.Println("2" , tLength , nextByte)

  tName , nextByte := getString(packet , nextByte , tLength) 
  pLength , nextByte := getLength(packet[nextByte : nextByte+4] , nextByte) 
  payload , _ := getString(packet , nextByte , pLength) 
  return Packet{
    PacketType: pType ,
    RLenght: rLenght,
    TopicLen: tLength,
    TopicName: tName,
    PayloadLen: pLength,
    Payload: payload,
  }


}

// get the remaining Length using the variable Length encoding
// checking the Most Singnificant Bit (MSB)
// 0-> no related following bit | 1 -> there is
func getLength(sequence []byte , start int) (uint , int)   {
  nextByte := start +len(sequence)
  size := ""
  for i , val := range sequence {
    msb := (val>>7)  
    length := val & byte(BYTE_MASK)
    fmt.Println(val , length , byteToBinary(length))
    size+=byteToBinary(length)
    if msb == 0 {
      nextByte = start +i +1
      break
    }
  }  
  return binaryToUint(size) ,nextByte 
}






//get the content that based on its lenght
func getString(packet []byte , start int , length uint  ) (string ,int)  {
  str := []byte{}
  for i:= start ; i< start+int(length) ; i++ {
    str = append(str, packet[i])
  }
  return string(str) , start+int(length)
  
}
// determining the packet type 
func packetTypeToString(packetType byte) string {
    switch packetType {
    case PUBLISH:
        return "PUBLISH"
    case SUBSCRIBE:
        return "SUBSCRIBE"
    case UNSUBSCRIBE:
        return "UNSUBSCRIBE"
    case CONNECT:
        return "CONNECT"
    case CONNACK:
        return "CONNACK"
    case DISCONNECT:
        return "DISCONNECT"
    case PINGREQ:
        return "PINGREQ"
    case PINGRESP:
        return "PINGRESP"
    case PUBACK:
        return "PUBACK"
    case PUBREC:
        return "PUBREC"
    case PUBREL:
        return "PUBREL"
    case PUBCOMP:
        return "PUBCOMP"
    default:
        return "UNKNOWN"
    }
}
