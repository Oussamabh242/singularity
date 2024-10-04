package queue

import (
	"fmt"
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
}

/*
  - concurrent map to hold the Queues
  - the map is returns a reference to a Queue
    when provided with its name of type map[string]*Queue
*/
type QStore struct {
	Queues sync.Map
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
func NewQStore() QStore {
	return QStore{}
}

// add a queue inside the QueueStore with max 20 Listener
func (qs *QStore) CreateQueue(name string) {
	q := Queue{
		Listeners: make(chan Listener, 20),
	}
	qs.Queues.Store(name, &q)
}

// retrieve a reference to a Queue based on its name or an error
func (qs *QStore) GetQueue(name string) (*Queue, error) {
	val, ok := qs.Queues.Load(name)
	if !ok {
		return nil, fmt.Errorf("No queue found")
	}
	return val.(*Queue), nil
}
