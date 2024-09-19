package main

import (
	"fmt"
	"log"
	"net"
	"os"
  _ "github.com/Oussamabh242/singularity/pkg/parser"
)


func handleConnection(conn net.Conn)  {
  b := make([]byte , 1024)
  n, err := conn.Read(b) 
  if err != nil{
    log.Println("error reading from connection")
  }
  conn.Write([]byte(string(b[:n])+" hello there"))
  fmt.Println("closing conn")
}
 
 
func main()  {
  var PORT string = os.Getenv("PORT")
  if len(PORT)== 0 {
    PORT = "1234"
  }
  fmt.Println(PORT)
  ln , err := net.Listen("tcp" , ":"+PORT)
  if err != nil{
    log.Panic("error starting the socket ", err)
  }
  for {
    conn , err := ln.Accept()
    if err != nil{
      log.Println("error establishing connection : " , err)
    }
    go handleConnection(conn)
  }
}

