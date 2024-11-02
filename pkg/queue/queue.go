package queue

import (
	"fmt"
	"github.com/Oussamabh242/singularity/pkg/messages"
	"net"
	"sync"
)

type Listener net.Conn

/*
  - Queue is defined by a channel of Listeners for the ease of
    concurrent use
*/
type Queue struct {
	Listeners chan Listener
	Messages  chan messages.Message
}

/*
  - concurrent map to hold the Queues
  - the map is returns a reference to a Queue
    when provided with its name of type map[string]*Queue
*/
type QStore struct {
	Queues        sync.Map
	maxListenners int
	maxMessages   int
}

// adds a Connection to the Listeners
func (q *Queue) Enqueue(conn net.Conn) {
	fmt.Println("enqueuing")
	q.Listeners <- conn
}

// Dequeues a Connection
func (q *Queue) Dequeue() Listener {
	return <-q.Listeners
}

func (q *Queue) Channel() chan Listener {
	return q.Listeners
}

// initialize a QueueStore
func NewQStore(maxSubscribers int, maxMessages int) QStore {
	fmt.Println("here ", maxSubscribers, maxMessages)
	return QStore{
		maxListenners: maxSubscribers,
		maxMessages:   maxMessages,
	}

}

// add a queue inside the QueueStore with max 20 Listener
func (qs *QStore) CreateQueue(name string) *Queue {
	q := Queue{
		Listeners: make(chan Listener, qs.maxListenners),
		Messages:  make(chan messages.Message, qs.maxMessages),
	}
	fmt.Println("cap", cap(q.Listeners), cap(q.Messages))
	qs.Queues.Store(name, &q)
	return &q
}

// retrieve a reference to a Queue based on its name or an error
func (qs *QStore) GetQueue(name string) (*Queue, error) {
	val, ok := qs.Queues.Load(name)
	if !ok {
		return nil, fmt.Errorf("No queue found")
	}
	return val.(*Queue), nil
}
