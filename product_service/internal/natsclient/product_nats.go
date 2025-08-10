package natsclient

import "github.com/nats-io/nats.go"

type ProductNATS struct {
	Conn *nats.Conn
}

func NewNATS(url string) *ProductNATS {
	conn, err := nats.Connect(url)
	if err != nil {
		panic(err)
	}
	return &ProductNATS{Conn: conn}
}

func (n *ProductNATS) Publish(subject string, data []byte) error {
	return n.Conn.Publish(subject, data)
}
