
package feed

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/Oussamabh242/singularity/pkg/messages"
	"github.com/Oussamabh242/singularity/pkg/queue"
)

func FeedMessages(qs *queue.QStore, ms *messages.MsgStore) {
	for {
		msg := ms.Get()
		q, _ := qs.GetQueue(msg.Queue)
		fmt.Println(len(q.Listeners))
		if len(q.Listeners) > 0 {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			go func(ctx context.Context, cancel context.CancelFunc, msg messages.Message) {
				defer cancel()
				AwaitForWork(ctx, q, ms, msg)
			}(ctx, cancel, msg)
		} else {
			ms.Add(msg)

		}
	}
}

func AwaitForWork(ctx context.Context, q *queue.Queue, ms *messages.MsgStore, msg messages.Message) {
	var recv = make(chan int)
	var errCh = make(chan struct{})

	select {
	case <-ctx.Done():
		ms.Add(msg)
		fmt.Println("Timed out waiting for work")
	case conn := <-q.Listeners:
		go WaitForJobDone(conn, msg.Body, recv, errCh)

		select {
		case <-ctx.Done():
			ms.Add(msg)
			q.Enqueue(conn)
			fmt.Println("Timed out while waiting for job to finish")
		case <-recv:
			fmt.Println("Received acknowledgement")
			q.Enqueue(conn)
		case <-errCh:
			fmt.Println("Removing listener due to error")
			conn.Close()
			ms.Add(msg)
		}
	}
}

func MakePacket(msgBody []byte) []byte {
	packetType := byte(7)
	rLength := byte(2 + len(msgBody))
	mLength := byte(0)
	bLength := byte(len(msgBody))
	arrByte := []byte{packetType, rLength, mLength, bLength}
	arrByte = append(arrByte, msgBody...)
	return arrByte
}

func WaitForJobDone(conn net.Conn, msg []byte, recv chan int, errCh chan struct{}) {
	defer close(errCh)

	if _, err := conn.Write(MakePacket(msg)); err != nil {
		errCh <- struct{}{}
		return
	}

	b := make([]byte, 40)
	n, err := conn.Read(b)
  fmt.Println(n , err)
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

