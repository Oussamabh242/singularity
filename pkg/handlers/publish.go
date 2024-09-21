package handlers

import (
	"net"

	"github.com/Oussamabh242/singularity/pkg/parser"
	"github.com/Oussamabh242/singularity/pkg/queue"
)

var (
  ACK []byte = []byte{parser.ACKRECIVE , 0 , 0 ,0}
)

func HandlePublish(conn net.Conn , parsed *parser.Packet , qs *queue.QStore){
  defer conn.Close()
  queue , err := qs.GetQueue(parsed.Metadata.Queue) ;  
  if err != nil {
    conn.Write([]byte("Error Sending message , queue not found"))
    return
  }
  queue.Store <- parsed.Payload
  AckPublish(conn)
  return
}


func AckPublish(conn net.Conn){
  conn.Write(ACK) 
}
