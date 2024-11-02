package messages

/*
 * a message will contain :
 * the body : the actual message
 * Queue : the queue that is sent to
 * Topic
 */
type Message struct {
	Body  []byte
	Queue string
	Topic string
	// todo: content Type
}

// A channel of messages
type MsgStore struct {
	Store chan Message
}

// initialize a messageStore with max 100 messages
func NewMessageStore() MsgStore {
	return MsgStore{
		Store: make(chan Message, 100),
	}
}

// Enqueues a message
func (Ms MsgStore) Add(msg Message) {
	Ms.Store <- msg
}

// Dequeues a message
func (Ms MsgStore) Get() Message {
	return <-Ms.Store
}
