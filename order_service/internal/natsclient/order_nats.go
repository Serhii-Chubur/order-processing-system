package natsclient

import "github.com/nats-io/nats.go"

type OrderNATS struct {
	Conn *nats.Conn
}

func NewNATS(url string) *OrderNATS {
	conn, err := nats.Connect(url)
	if err != nil {
		panic(err)
	}
	return &OrderNATS{Conn: conn}
}

func (n *OrderNATS) Subscribe(subject string, cb nats.MsgHandler) (*nats.Subscription, error) {
	return n.Conn.Subscribe(subject, cb)
}
