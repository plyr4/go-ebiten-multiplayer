package ws

type Message interface {
	GetSubType() string
}

type Ping struct {
}

func (c Ping) GetSubType() string {
	return "ping"
}

type Join struct {
}

func (c Join) GetSubType() string {
	return "join"
}

type Update struct {
}

func (c Update) GetSubType() string {
	return "update"
}
