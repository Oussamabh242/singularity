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
package encoder

import (
	"encoding/binary"
)

func Encode(packetType uint8, metadata []byte, message []byte) []byte {
	rLength := 2 + len(metadata) + 2 + len(message)
	packet := make([]byte, 0)

	rLengthByte := make([]byte, 4)
	binary.BigEndian.PutUint32(rLengthByte, uint32(rLength))

	metadataByte := make([]byte, 2)
	binary.BigEndian.PutUint16(metadataByte, uint16(len(metadata)))

	BodyByte := make([]byte, 2)
	binary.BigEndian.PutUint16(BodyByte, uint16(len(message)))

	packet = append(packet, packetType)
	packet = append(packet, rLengthByte...)
	packet = append(packet, metadataByte...)
	packet = append(packet, metadata...)
	packet = append(packet, BodyByte...)
	packet = append(packet, message...)
	return packet
}
