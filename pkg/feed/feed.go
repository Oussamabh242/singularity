package feed

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/Oussamabh242/singularity/pkg/messages"
	"github.com/Oussamabh242/singularity/pkg/queue"
)

func FeedMessages(qs *queue.QStore , ms *messages.MsgStore)  {
  i := 0
  for {
    i ++
    fmt.Println(i , "Waiting for a message")
    msg := ms.Get()
    ctx , cancel := context.WithTimeout(context.Background() ,time.Second*5)
    go func(ctx context.Context , cnacel context.CancelFunc) {
      defer cancel()
      AwaitForWork(ctx ,qs ,ms , msg)
    }(ctx ,cancel)
    
  }
  
}


func AwaitForWork(ctx context.Context ,qs *queue.QStore,ms *messages.MsgStore  ,msg messages.Message){
  q , _ := qs.GetQueue(msg.Queue) 
  fmt.Println("Waiting for a Listenner")
  select {
  case <-ctx.Done() :
    ms.Add(msg)
    fmt.Println("timed Out 1")
  case conn := <- q.Listeners :
    GoodConn := true
    var recv chan []byte = make(chan []byte , 100 ) 
    go WaitForJobDone(conn ,msg.Body , recv, &GoodConn ) 
    select {
    case <-ctx.Done() : 
      ms.Add(msg)
      if GoodConn == true {
        q.Enqueue(conn)
      }
      fmt.Println("timed Out 2")
    case recived := <-recv :
      fmt.Println("I recived this : " , string(recived))
      if GoodConn == true {
        q.Enqueue(conn)
      }
    }


  }
}


func MakePacket(msgBody []byte) []byte {
  packetType := byte(7)
  rLength := byte(2+len(msgBody))
  mLenght := byte(0) 
  bLength := byte(len(msgBody))
  arrByte := []byte{packetType  , rLength , mLenght , bLength }
  for _, val := range msgBody {
    arrByte = append(arrByte, val) 
  }
  return arrByte
}

func WaitForJobDone(conn net.Conn , msg []byte , recv chan []byte , goodConn *bool){
    fmt.Println("writing")
    conn.Write(MakePacket(msg)) 
    b := make([]byte , 40)
    fmt.Println("starting to read")
    n , err := conn.Read(b)
    if n==0 {
      *goodConn = false
      fmt.Println("no data to read")

      return 
    }
    if err != nil{
      *goodConn = false
      return
    }
    recv<-b
    fmt.Println("here")
}
