package handlers

import (
	"net"

	"github.com/Oussamabh242/singularity/pkg/parser"
	"github.com/Oussamabh242/singularity/pkg/queue"
)

var ACKSUB []byte = []byte{parser.ACKSUBSCRIBE, 0, 0, 0}

/*
  - Adds the Subscriber to the queue's listenners
    if the queue is present

  - Respond by a Packet of type ACKSUB to initiate-
    the start of reciving messages
*/
func HandleSubscribe(conn net.Conn, parsed *parser.Packet, qs *queue.QStore) {
	q, err := qs.GetQueue(parsed.Metadata.Queue)
	if err != nil {
		conn.Write([]byte("error subscribing no queue"))
    return
	}
	q.Enqueue(conn)
	ackSubs(conn)
}

func ackSubs(conn net.Conn) {
	conn.Write(ACKSUB)
}
