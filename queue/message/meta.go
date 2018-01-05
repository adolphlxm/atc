package message

type Meta struct {
	MessageId string
	Src       string
	Priority  MsgPriority
	Options   map[string]string
}

func (m *Meta) FormMessage(msg *Message) {
	if msg != nil {
		m.MessageId = msg.MessageId
		m.Options = msg.Options
		m.Priority = msg.Priority
	}
}
