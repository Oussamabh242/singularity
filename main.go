package main

import (
	"fmt"
	"github.com/Oussamabh242/singularity/pkg/feed"
	"github.com/Oussamabh242/singularity/pkg/handlers"
	"github.com/Oussamabh242/singularity/pkg/messages"
	"github.com/Oussamabh242/singularity/pkg/parser"
	"github.com/Oussamabh242/singularity/pkg/queue"

	"log"
	"net"
	"os"
)

func handleConnection(conn net.Conn , qs *queue.QStore , ms *messages.MsgStore) {
	b := make([]byte, 1024)
	conn.Read(b)
	thing := parser.Parse(b)
  switch thing.PacketType {
  case parser.PING :
    handlers.HandlePing(conn)
    break
  case parser.PUBLISH :
    handlers.HandlePublish(conn , &thing, qs , ms)
    break
  case parser.CREATEQUEUE :
    handlers.HandlerCreateQueue(conn  , &thing , qs)
    break

  case parser.SUBSCRIBE :
    handlers.HandleSubscribe(conn , &thing , qs )
    break

  default :
    fmt.Println("UNKNOWN")
  }


}

func main() {
	var PORT string = os.Getenv("PORT")
	if len(PORT) == 0 {
		PORT = "1234"
	}
	fmt.Println(PORT)
	ln, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		log.Panic("error starting the socket ", err)
	}
  
  qs := queue.NewQStore()
  qs.CreateQueue("q")
  mStore := messages.NewMessageStore()
  go feed.FeedMessages(&qs , &mStore)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("error establishing connection : ", err)
		}
		go handleConnection(conn , &qs , &mStore)

	}
}
