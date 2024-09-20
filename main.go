package main

import (
	"fmt"
	"github.com/Oussamabh242/singularity/pkg/parser"
	"log"
	"net"
	"os"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	b := make([]byte, 1024)
	conn.Read(b)
	thing := parser.Parse(b)
	if thing.PacketType == "PING" {
		conn.Write([]byte("PONG"))
		return
	}
	conn.Write([]byte("RESPONSE"))

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
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("error establishing connection : ", err)
		}
		go handleConnection(conn)
	}
}
