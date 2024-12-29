// "github.com/Oussamabh242/singularity"
package handlers

import (
	"fmt"
	"net"
	"github.com/Oussamabh242/singularity/pkg/feed"
	"github.com/Oussamabh242/singularity/pkg/parser"
	"github.com/Oussamabh242/singularity/pkg/queue"
  "github.com/Oussamabh242/singularity/pkg/encoder"
)

var ACKCREATE []byte = []byte{parser.ACKQREATE, 0, 0, 0,0}

/*
 * Create a queue inside the queue Store
 * respond with a packet of type ACKCREATE
 */
func HandlerCreateQueue(conn net.Conn, parsed *parser.Packet, qs *queue.QStore) {
	defer conn.Close()
	queueName := parsed.Metadata.Queue
	_, ok := qs.GetQueue(queueName)
	if ok == nil {
		AckQueueCreate(conn)
    return
	}

	if len(queueName) == 0 {
		conn.Write([]byte("error creating a queue"))
		return
	}
	q := qs.CreateQueue(parsed.Metadata.Queue)
	fmt.Println("new queue created :",parsed.Metadata.Queue, parsed.Metadata)
	go feed.FeedMessages(q)
	AckQueueCreate(conn)
	return
}

func AckQueueCreate(conn net.Conn) {
	conn.Write(encoder.Encode(parser.ACKQREATE ,[]byte{} , []byte{}))
}
