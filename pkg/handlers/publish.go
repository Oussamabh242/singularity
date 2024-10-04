package handlers

import (
	"fmt"
	"net"

	"github.com/Oussamabh242/singularity/pkg/messages"
	"github.com/Oussamabh242/singularity/pkg/parser"
	"github.com/Oussamabh242/singularity/pkg/queue"
)

var (
	ACK []byte = []byte{parser.ACKRECIVE, 0, 0, 0}
)

/*
 * Stores a messages recived from a publisher
 * respond with a packet of type ACKRECIVE to publisher
 */
func HandlePublish(conn net.Conn, parsed *parser.Packet, qs *queue.QStore, ms *messages.MsgStore) {
	defer conn.Close()
	_, err := qs.GetQueue(parsed.Metadata.Queue)
	if err != nil {
		conn.Write([]byte("Error Sending message , queue not found"))
		return
	}
	msg := messages.Message{
		Body:  []byte(parsed.Payload),
		Queue: parsed.Metadata.Queue,
		Topic: parsed.Metadata.Topic,
	}
	fmt.Println("writing msg : ", parsed.Payload)
	ms.Add(msg)
	AckPublish(conn)
	return
}

func AckPublish(conn net.Conn) {
	conn.Write(ACK)
}
