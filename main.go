package main

import (
	"fmt"

	"github.com/Oussamabh242/singularity/pkg/handlers"
	"github.com/Oussamabh242/singularity/pkg/parser"
	"github.com/Oussamabh242/singularity/pkg/queue"

	"log"
	"net"
	"os"
)

func handleConnection(conn net.Conn , qs *queue.QStore) {
	b := make([]byte, 1024)
	conn.Read(b)
	thing := parser.Parse(b)
  switch thing.PacketType {
  case parser.PING :
    handlers.HandlePing(conn)
    break
  case parser.PUBLISH :
    handlers.HandlePublish(conn , &thing, qs)
    break
  case parser.CREATEQUEUE :
    handlers.HandlerCreateQueue(conn  , &thing , qs)
    break

  case parser.SUBSCRIBE :
    handlers.HandleSubscribe(conn , &thing , qs)
    break

  default :
    fmt.Println("UNKNOWN")
  }
  fmt.Println(thing.RLenght)
  // if thing.PacketType == parser.SUBSCRIBE {
  //   q , _ := qs.GetQueue(thing.Metadata.Queue)
  //   msg := <- q.Store
  //   fmt.Println(len(q.Listeners))
  //   time.Sleep(time.Second*1)
  //   q.Listeners[0].Conn.Write([]byte(msg))
  // }
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
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("error establishing connection : ", err)
		}
		go handleConnection(conn , &qs)
	}
}
