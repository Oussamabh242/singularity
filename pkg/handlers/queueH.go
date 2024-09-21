// "github.com/Oussamabh242/singularity"
package handlers
   
import (
	"net"
   
	"github.com/Oussamabh242/singularity/pkg/parser"
	"github.com/Oussamabh242/singularity/pkg/queue"
)  
   
var ACKCREATE []byte = []byte{parser.ACKQREATE , 0 , 0 , 0}   

func HandlerCreateQueue(conn net.Conn ,parsed *parser.Packet, qs *queue.QStore){
  defer conn.Close()
  qs.CreateQueue(parsed.Metadata.Queue) 
  AckQueueCreate(conn)
  return
} 

func AckQueueCreate(conn net.Conn){
  conn.Write(ACKCREATE)
}
