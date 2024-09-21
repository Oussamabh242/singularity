package handlers

import(
  "net"
)

func HandlePing(conn net.Conn){
  defer conn.Close()
  b := []byte("PONG")
  conn.Write(b)
  return 
} 
