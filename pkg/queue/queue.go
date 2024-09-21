package queue

import (
	"fmt"
	"net"
	"sync"
)

type Listener struct {
  Conn net.Conn 
  Free bool
  Connected bool
}

type Queue struct {
  Listeners []Listener
  Store chan string
}

func (q *Queue) Enqueue(msg string){
  q.Store<-msg
}

func (q *Queue) Dequeue() string{
  return <-q.Store
}

func (q *Queue) Channel() chan string {
  return q.Store
}

type QStore struct {
  Queues sync.Map 
}

func (qs *QStore) CreateQueue(name string)  {
  q := Queue{
    Store: make(chan string, 20),
    Listeners: []Listener{},
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
