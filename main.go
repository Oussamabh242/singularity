// package main
//
// import (
//
// 	"github.com/Oussamabh242/singularity/pkg/handlers"
// 	"github.com/Oussamabh242/singularity/pkg/parser"
// 	"github.com/Oussamabh242/singularity/pkg/queue"
//   "github.com/Oussamabh242/singularity/cmd"
//
//
// 	"fmt"
// 	"log"
// 	"net"
// 	"os"
// )
//
// // parse packet type and assign it to a specified handler
// func handleConnection(conn net.Conn, qs *queue.QStore) {
//
// 	b := make([]byte, 1024)
// 	length, err := conn.Read(b)
// 	if err != nil {
// 		fmt.Println("error reading message", err)
// 		conn.Close()
// 		return
// 	}
// 	thing := parser.Parse(b[:length])
// 	switch thing.PacketType {
// 	case parser.PING:
// 		handlers.HandlePing(conn)
// 		break
// 	case parser.PUBLISH:
// 		handlers.HandlePublish(conn, &thing, qs)
// 		break
// 	case parser.CREATEQUEUE:
// 		handlers.HandlerCreateQueue(conn, &thing, qs)
// 		break
//
// 	case parser.SUBSCRIBE:
// 		handlers.HandleSubscribe(conn, &thing, qs)
// 		break
//
// 	default:
// 		fmt.Println("UNKNOWN")
// 	}
// }
//
// func main() {
//
//   cmd.Execute()
// 	var PORT string = os.Getenv("PORT")
// 	if len(PORT) == 0 {
// 		PORT = "1234"
// 	}
// 	fmt.Println(PORT)
// 	ln, err := net.Listen("tcp", ":"+PORT)
// 	if err != nil {
// 		log.Panic("error starting the socket ", err)
// 	}
//
// 	qs := queue.NewQStore()
// 	for {
// 		conn, err := ln.Accept()
// 		if err != nil {
// 			log.Println("error establishing connection : ", err)
// 		}
// 		go handleConnection(conn, &qs)
//
// 	}
// }

package main

import "github.com/Oussamabh242/singularity/cmd"

func main() {
	cmd.Execute()
}
