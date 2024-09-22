package messages


type Message struct {
  Body  []byte 
  Queue string
  Topic string
}


type MsgStore struct {
  Store chan Message
}

func (Ms MsgStore) Add(msg Message) {
  Ms.Store <- msg
}

func (Ms MsgStore) Get() Message {
  return <-Ms.Store
}

func NewMessageStore() MsgStore {
  return MsgStore{
    Store: make(chan Message, 20),
  }
}
