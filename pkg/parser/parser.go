package parser

/*
  ***
    the job of this package is to parse or transform incomming array of bytes
    (packet) into a struct to deal with it later
  ***

  ** Packet Structure **

      1 Byte          4 Byte          2 Byte
  +---------------+----------------+----------------+
  |packet type    |remainingLength |metadata length |
  +---------------|----------------|----------------|
  | Metadata      | message length |  message       |
  +---------------+----------------+----------------+
  0->(2^16)-3 Byte     2 Byte       0->(2^16)-3 Byte

*/

import (
	"encoding/json"
	"fmt"
	"log"
)

const (
	PUBLISH      = 1 // Publish
	ACKRECIVE    = 2 // Acknowledge Publishing
	SUBSCRIBE    = 3 // Subscribe to Queue
	ACKSUBSCRIBE = 4
	CREATEQUEUE  = 5
	ACKQREATE    = 6
	JOB          = 7
	ACKJOB       = 8
	PING         = 100

	BYTE_MASK = 0x7F

  // Parts indexes 
  PACKET_HEAD_START = 1
  PACKET_HEAD_END   = 5// 5 bytes (1 => type  , 4 => remaining length)
  META_LENGTH_START = 5 
  META_LENGTH_END   = 7 

)

// the Metadata is conserned with everything that is not the messgae
// queue name , topic , message Content Type , message for the consumer
// Acknowledgement

type Metadata struct {
	Queue       string `json:"queue"`
	Topic       string `json:"Topic"`
	ContentType string `json:"content-type"`
}

type Packet struct {
	PacketType  int
	RLenght     uint
	MetadataLen uint
	Metadata    Metadata
	PayloadLen  uint
	Payload     []byte
}



func Parse(packet []byte) Packet {
  
	// Parses the incomming bytes into a struct
  totalsize := len(packet)

  fmt.Println("recived: " , packet)
	parsed := Packet{}
	parsed.PacketType = int(packet[0])
	rLength := Intify[uint32](packet[PACKET_HEAD_START:PACKET_HEAD_END]) // at most 4 bytes
	parsed.RLenght = uint(rLength)
  fmt.Println("--- reaining Length is " ,rLength ) // next Byte 5
	if rLength == 0 || len(packet)<PACKET_HEAD_END {
		return parsed
	}
	mLength := Intify[uint16](packet[META_LENGTH_START: META_LENGTH_END]) // at most 2 bytes
	parsed.MetadataLen = uint(mLength)
  fmt.Println("-- metadata Length is ", mLength)
  if mLength ==0 {
    return parsed
  }
  md, err := extractMetadata(packet[META_LENGTH_END: META_LENGTH_END+mLength])
  if err != nil{
    log.Print("Error Decoding Metadata" , err)
  }
  parsed.Metadata = md
  fmt.Println(md)
  nextByte := META_LENGTH_END+uint(mLength) 
      
	if nextByte>= uint(totalsize) {
		return parsed
	}
	pLength := Intify[uint16](packet[nextByte : nextByte+2])
	parsed.PayloadLen = uint(pLength)
  fmt.Println("-- Payload Length is " , pLength)
	nextByte = nextByte + 2
	if pLength == 0 {
		return parsed	
  }
  parsed.Payload = packet[nextByte : nextByte+uint(pLength)]
  fmt.Println("-- msg is " , string(parsed.Payload))
	fmt.Println(parsed)
	return parsed

}
// func Parse(packet []byte) Packet {
// 	// Parses the incomming bytes into a struct
//   fmt.Println("recived: " , packet)
// 	totalsize := len(packet)
// 	parsed := Packet{}
// 	parsed.PacketType = int(packet[0])
// 	rLength := Intify[uint32](packet[1:5]) // at most 4 bytes
// 	parsed.RLenght = uint(rLength)
//   fmt.Println("--- reaining Length is " ,rLength )
// 	nextByte := 5
// 	if rLength == 0 || nextByte >= totalsize {
// 		return parsed
// 	}
//
// 	mLength := Intify[uint16](packet[nextByte : nextByte+2]) // at most 2 bytes
// 	parsed.MetadataLen = uint(mLength)
//   fmt.Println("-- metadata Length is ", mLength ,packet[nextByte : nextByte+2])
// 	nextByte = nextByte + 2
// 	if mLength > 0 {
// 		var metadata Metadata
// 		metadata, nextByte = parseMetadata(packet, nextByte, uint(mLength))
// 		parsed.Metadata = metadata
// 	}
// 	if nextByte >= totalsize {
// 		return parsed
// 	}
// 	pLength := Intify[uint16](packet[nextByte : nextByte+2])
// 	parsed.PayloadLen = uint(pLength)
// 	nextByte = nextByte + 2
// 	if pLength > 0 {
// 		payload, _ := getString(packet, nextByte, uint(pLength))
// 		parsed.Payload = payload
// 	}
// 	fmt.Println(parsed)
// 	return parsed
//
// }


/*
 * transforming the bytes designed to the metadata which is a json
 * to a struct
 */
func extractMetadata(arr []byte ) (Metadata , error)  {
  md := Metadata{}
  err := json.Unmarshal(arr, &md)
  if err != nil{
    return md , err
  }
	return md , nil  
}

