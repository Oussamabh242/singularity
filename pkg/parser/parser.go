package parser

/*
  ***
    the job of this package is to parse or transform incomming array of bytes
    (packet) into a struct to deal with it later
  ***

  ** Packet Structure **

      1 Byte      2 Byte          2 Byte
  +---------------+----------------+----------------+
  |packet type    |remainingLenght |metadata length |
  +---------------|----------------|----------------|
  | Metadata      | message length |  message       |
  +---------------+----------------+----------------+
  0->(2^16)-3 Byte     2 Byte       0->(2^16)-3 Byte

*/

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
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
	Payload     string
}

func Parse(packet []byte) Packet {
	// Parses the incomming bytes into a struct
	totalsize := len(packet)
	parsed := Packet{}
	parsed.PacketType = int(packet[0])
	rLenght := intify(packet[1:3]) // at most 2 bytes
	parsed.RLenght = uint(rLenght)
	nextByte := 3
	if rLenght == 0 || nextByte >= totalsize {
		return parsed
	}

	mLength := intify(packet[nextByte : nextByte+2]) // at most 2 bytes
	parsed.MetadataLen = mLength
	nextByte = nextByte + 2
	if mLength > 0 {
		var metadata Metadata
		metadata, nextByte = parseMetadata(packet, nextByte, mLength)
		parsed.Metadata = metadata
	}
	if nextByte >= totalsize {
		return parsed
	}
	pLength := intify(packet[nextByte : nextByte+2])
	parsed.PayloadLen = pLength
	nextByte = nextByte + 2
	if pLength > 0 {
		payload, _ := getString(packet, nextByte, pLength)
		parsed.Payload = payload
	}
	fmt.Println(parsed)
	return parsed

}

/*
 * transforming the bytes designed to the metadata which is a json
 * to a struct
 */
func parseMetadata(packet []byte, start int, length uint) (Metadata, int) {
	md := Metadata{}
	json.Unmarshal(packet[start:start+int(length)], &md)
	return md, start + int(length)
}

/*
*
  - get the remaining Length using the variable Length encoding
  - all the lengths gonna be maxed to two bytes which will make
    the max lenght 2^16
  - the length is expected to be encoded in two bytes eg:
    encoding 1024 should give us first byte : 4 , second 0 ( 00000100 00000000)
*/
func getLength(sequence []byte, start int) (uint, int) {

	arr := []byte{sequence[start], sequence[start+1]}
	var num uint16
	err := binary.Read(bytes.NewReader(arr), binary.BigEndian, &num)
	if err != nil {
		fmt.Println(err)
		return 0, start + 2
	}
	return uint(num), start + 2
}

/*
 * get the content that based on its lenght
 *
 */
func getString(packet []byte, start int, length uint) (string, int) {
	str := []byte{}
	for i := start; i < start+int(length); i++ {
		str = append(str, packet[i])
	}
	return string(str), start + int(length)

}
