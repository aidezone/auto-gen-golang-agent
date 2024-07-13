package defines

type Chat struct {
	list []*Message
}

func NewChat() *Chat {
	return &Chat{}
}

func (s *Chat) AppendAsk(msg *string) []*Message {
	s.list = append(s.list, &Message{
		Actor: USER,
		Msg: msg,
	})
	return s.list
}

func (s *Chat) Ask(msg *string) []*Message {
	return append(s.list, &Message{
		Actor: USER,
		Msg: msg,
	})
}

func (s *Chat) AppendAnswer(msg *string) {
	s.list = append(s.list, &Message{
		Actor: ROBOT,
		Msg: msg,
	})
}