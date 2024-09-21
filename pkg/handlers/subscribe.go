package handlers

import (
	"net"

	"github.com/Oussamabh242/singularity/pkg/parser"
	"github.com/Oussamabh242/singularity/pkg/queue"
)

var ACKSUB []byte= []byte{parser.ACKSUBSCRIBE , 0, 0 , 0}

func HandleSubscribe(conn net.Conn , parsed *parser.Packet ,qs *queue.QStore)  {
  q , err := qs.GetQueue(parsed.Metadata.Queue)
  if err != nil{
    conn.Write([]byte("error subscribing no queue"))
  }
  client := queue.Listener{
    Conn : conn ,
    Free: true,
    Connected: true,
  }
  q.Listeners = append(q.Listeners,client )
  ackSubs(conn) 
}

func ackSubs(conn net.Conn){
  conn.Write(ACKSUB)
}
