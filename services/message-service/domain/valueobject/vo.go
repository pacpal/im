package valueobject

type MessageID struct {
	value string
}

func NewMessageID(value string) MessageID {
	return MessageID{value: value}
}

func (id MessageID) String() string {
	return id.value
}

func (id MessageID) IsValid() bool {
	return id.value != ""
}

type MessageType struct {
	value string
}

func NewMessageType(value string) MessageType {
	return MessageType{value: value}
}

func (t MessageType) String() string {
	return t.value
}

func (t MessageType) IsValid() bool {
	return t.value == "private" || t.value == "group"
}