package feed

import (
	"context"
	"fmt"
	"net"
	"time"
	"github.com/Oussamabh242/singularity/pkg/messages"
	"github.com/Oussamabh242/singularity/pkg/parser"
	"github.com/Oussamabh242/singularity/pkg/queue"
  "github.com/Oussamabh242/singularity/pkg/encoder"
)

/*
 * this a gourinte that will work in parallel
 * it keeps waiting for a listenner
 * if there is a listenner it waits for a message
 * it sends sends the message to the listenner (consumer)
 */
func FeedMessages(q *queue.Queue) {
	for {
    fmt.Println("listenning now are : " , len(q.Listeners))
		select {
		case conn := <-q.Listeners:

			select {
			case msg := <-q.Messages:

				go func(conn queue.Listener, msg messages.Message) {

					ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
					defer cancel()
					AwaitForWork(ctx, q, conn, msg)

				}(conn, msg)
			}

		}
	}
}

/*
  - AwaitForWork extract a listner and gets a message wait for the message
    to be done by one of the queue subscribers in case of reciving ack from
    listener * the listener is reEnqueued to the queue's Listeners
  - listen on the error channel if an error occurs the listener will not be
    reEnqueued
  - if context is timedOut both the listener and the message are reEnqueued

ARGS :

ctx : context : if done the message is reEnqueued
q   : @Queue  : reference to the queue that corresponds to the message
conn: q.Listener (net.Conn)
msg : Message
*/
func AwaitForWork(ctx context.Context, q *queue.Queue, conn net.Conn, msg messages.Message) {
	var recv = make(chan int)
	var errCh = make(chan struct{})

	go WaitForAck(conn, msg.Body, recv, errCh)

	select {
	case <-ctx.Done():
		q.Messages <- msg
		q.Enqueue(conn)
		fmt.Println("Timed out while waiting for job to finish")
	case <-recv:
		fmt.Println("Received acknowledgement")
		q.Enqueue(conn)
	case <-errCh:
		fmt.Println("Removing listener due to error")
		conn.Close()
		q.Messages <- msg
	}
}

/*
 * WaitForAck sends a packet  that contains the message to a listener
 * wait for ack from the listener
 * for reading error from the listener (conn closed)
 */
func WaitForAck(conn net.Conn, msg []byte, recv chan int, errCh chan struct{}) {
	defer close(errCh)

  packet := encoder.Encode(parser.JOB ,[]byte{} ,msg)
  fmt.Println("feeding ... " , packet)
	if _, err := conn.Write(packet); err != nil {
		errCh <- struct{}{}
		return
	}

	b := make([]byte, 40)
	n, err := conn.Read(b)
	fmt.Println("recived" ,n, err)
	if err != nil {
		errCh <- struct{}{}
		return
	}
	if n == 0 {
		errCh <- struct{}{}
		return
	}

	recv <- 1
}

/*
  - organize a packet of type JOB that contains the lenght
    of the message and the actual message
*/
func MakePacket(msgBody []byte) []byte {
	packetType := byte(parser.JOB)
	rLength := byte(2 + len(msgBody))
	mLength := byte(0)
	bLength := byte(len(msgBody))
	arrByte := []byte{packetType, rLength, mLength, bLength}
	arrByte = append(arrByte, msgBody...)
	return arrByte
}
