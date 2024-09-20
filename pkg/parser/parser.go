package parser

import (
	"encoding/json"
)

const (
	PUBLISH     = 1 // Publish
	SUBSCRIBE   = 2 // Subscribe to Queue
	CREATEQUEUE = 3
	RECEIVE     = 4
	ACKMSG      = 5
	REJECTMSG   = 6
	PING        = 7

	BYTE_MASK = 0x7F
)

type Metadata struct {
	Queue       string `json:"queue"`
	Topic       string `json:"Topic"`
	ContentType string `json:"content-type"`
}

type Packet struct {
	PacketType  string
	RLenght     uint
	MetadataLen uint
	Metadata    Metadata
	PayloadLen  uint
	Payload     string
}

func Parse(packet []byte) Packet {
	totalsize := len(packet)
	parsed := Packet{}
	parsed.PacketType = packetTypeToString(packet[0])
	rLenght, nextByte := getLength(packet, 1) // at most 4 bytes starting from 2nd bit
	parsed.RLenght = rLenght

	if rLenght == 0 || nextByte >= totalsize {
		return parsed
	}

	mLength, nextByte := getLength(packet, nextByte) // at most 2 bytes
	parsed.MetadataLen = mLength
	if mLength > 0 {
		var metadata Metadata
		metadata, nextByte = parseMetadata(packet, nextByte, mLength)
		parsed.Metadata = metadata
	}
	if nextByte >= totalsize {
		return parsed
	}
	pLength, nextByte := getLength(packet, nextByte)
	parsed.PayloadLen = pLength
	if pLength > 0 {
		payload, _ := getString(packet, nextByte, pLength)
		parsed.Payload = payload
	}

	return parsed

}

func parseMetadata(packet []byte, start int, length uint) (Metadata, int) {
	md := Metadata{}
	json.Unmarshal(packet[start:start+int(length)], &md)
	return md, start + int(length)
}

// get the remaining Length using the variable Length encoding
// checking the Most Singnificant Bit (MSB)
// 0-> no related following bit | 1 -> there is
func getLength(sequence []byte, start int) (uint, int) {
	nextByte := start + len(sequence)
	size := ""
	for i := start; i < start+4; i++ {
		if i > len(sequence) {
			nextByte = -1
			break
		}
		msb := (sequence[i] >> 7)
		length := sequence[i] & byte(BYTE_MASK)
		size += byteToBinary(length)
		if msb == 0 {
			nextByte = i + 1
			break
		}
	}
	if nextByte > len(sequence) {
		nextByte = -1
	}
	return binaryToUint(size), nextByte
}

// get the content that based on its lenght
func getString(packet []byte, start int, length uint) (string, int) {
	str := []byte{}
	for i := start; i < start+int(length); i++ {
		str = append(str, packet[i])
	}
	return string(str), start + int(length)

}

// determining the packet type
func packetTypeToString(packetType byte) string {
	switch packetType {
	case PUBLISH:
		return "PUBLISH"
	case SUBSCRIBE:
		return "SUBSCRIBE"
	case ACKMSG:
		return "ACKMSG"
	case CREATEQUEUE:
		return "CREATEQUEUE"
	case RECEIVE:
		return "RECEIVE"
	case REJECTMSG:
		return "REJECTMSG"
	case PING:
		return "PING"
	default:
		return "UNKNOWN"
	}
}
