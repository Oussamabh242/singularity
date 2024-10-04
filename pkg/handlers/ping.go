package handlers

import (
	"net"
)

/*
 * This is meant to test the connection
 * if We recive a PING package you acknowledge with
 * a PONG
 */
func HandlePing(conn net.Conn) {
	defer conn.Close()
	b := []byte("PONG")
	conn.Write(b)
	return
}
