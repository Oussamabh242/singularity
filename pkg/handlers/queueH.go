// "github.com/Oussamabh242/singularity"
package handlers

import (
	"fmt"
	"net"

	"github.com/Oussamabh242/singularity/pkg/parser"
	"github.com/Oussamabh242/singularity/pkg/queue"
)  
   
var ACKCREATE []byte = []byte{parser.ACKQREATE , 0 , 0 , 0}   

func HandlerCreateQueue(conn net.Conn ,parsed *parser.Packet, qs *queue.QStore){
  defer conn.Close()
  queueName := parsed.Metadata.Queue
  _ , ok := qs.GetQueue(queueName) 
  if ok == nil {
    return
  }

  if len(queueName)==0 {
    conn.Write([]byte("error creating a queue"))
    return 
  } 
  qs.CreateQueue(parsed.Metadata.Queue) 
  q , _ := qs.Queues.Load(queueName)
  fmt.Println("new queue created :" , q)
  AckQueueCreate(conn)
  return
} 

func AckQueueCreate(conn net.Conn){
  conn.Write(ACKCREATE)
}
