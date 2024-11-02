package cmd

import (
	"github.com/Oussamabh242/singularity/pkg/handlers"
	"github.com/Oussamabh242/singularity/pkg/parser"
	"github.com/Oussamabh242/singularity/pkg/queue"

	"fmt"
	"github.com/spf13/cobra"
	"log"
	"net"
	// "os"
)

var (
	maxSubscribers int
	port           string
	maxMessages    int
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts Singularity message Broker ",
	Long:  `This command starts a TCP server on the specified port. Which will make run the Message Broker.`,
	Run: func(cmd *cobra.Command, args []string) {
		startServer(port, maxSubscribers, maxMessages)
	},
}

func init() {
	// Add a string flag for the port
	serveCmd.Flags().StringVarP(&port, "port", "p", "1234", "Port to listen on")
	serveCmd.Flags().IntVarP(&maxSubscribers, "subscribers", "s", 10, "Maximum number of subscribers allowed per a single queue")
	serveCmd.Flags().IntVarP(&maxMessages, "messages", "m", 100, "Maximum number of messages a queue can hold before becoming blocked")
	rootCmd.AddCommand(serveCmd)
}

// parse packet type and assign it to a specified handler
func handleConnection(conn net.Conn, qs *queue.QStore) {

	b := make([]byte, 1024)
	length, err := conn.Read(b)
	if err != nil {
		fmt.Println("error reading message", err)
		conn.Close()
		return
	}
	thing := parser.Parse(b[:length])
	switch thing.PacketType {
	case parser.PING:
		handlers.HandlePing(conn)
		break
	case parser.PUBLISH:
		handlers.HandlePublish(conn, &thing, qs)
		break
	case parser.CREATEQUEUE:
		handlers.HandlerCreateQueue(conn, &thing, qs)
		break

	case parser.SUBSCRIBE:
		handlers.HandleSubscribe(conn, &thing, qs)
		break

	default:
		fmt.Println("UNKNOWN")
	}
}

func startServer(PORT string, maxSubscribers int, maxMessages int) {

	// var PORT string = os.Getenv("PORT")
	// if len(PORT) == 0 {
	// 	PORT = "1234"
	// }
	fmt.Println(PORT)
	ln, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		log.Panic("error starting the socket ", err)
	}

	qs := queue.NewQStore(maxSubscribers, maxMessages)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("error establishing connection : ", err)
		}
		go handleConnection(conn, &qs)

	}
}
