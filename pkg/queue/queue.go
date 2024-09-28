package queue

import (
	"fmt"
	"net"
	"sync"
)

type Listener net.Conn

type Queue struct {
  Listeners chan Listener
  mux       sync.Mutex
}

func (q *Queue) Enqueue(conn net.Conn){
  fmt.Println("enqueuing")
  q.Listeners <- conn
}

func (q *Queue) Dequeue() Listener {
  return <-q.Listeners
}

func (q *Queue) Channel() chan Listener {
  return q.Listeners
}

type QStore struct {
  Queues sync.Map 
}

func (qs *QStore) CreateQueue(name string)  {
  q := Queue{
    Listeners: make(chan Listener , 20),
  }
  qs.Queues.Store(name , &q)
}

func (qs *QStore) GetQueue(name string) (*Queue , error) {
  val , ok := qs.Queues.Load(name)
  if !ok {
    return nil , fmt.Errorf("No queue found")
  }
  return val.(*Queue) , nil
}


func NewQStore() QStore {
  return QStore{}
}
