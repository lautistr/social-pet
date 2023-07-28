package events

import (
	"bytes"
	"encoding/gob"

	"github.com/lautistr/social-pet/models"
	"github.com/nats-io/nats.go"
)

type NatsEventStore struct {
	conn            *nats.Conn
	postCreatedSub  *nats.Subscription
	postCreatedChan chan CreatedPostMessage
}

func NewNats(url string) (*NatsEventStore, error) {
	conn, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}

	return &NatsEventStore{
		conn: conn,
	}, nil
}

func (n *NatsEventStore) Close() {
	if n.conn != nil {
		n.conn.Close()
	}
	if n.postCreatedSub != nil {
		n.postCreatedSub.Unsubscribe()
	}
	close(n.postCreatedChan)
}

func (n *NatsEventStore) encodeMessage(m Message) ([]byte, error) {
	b := bytes.Buffer{}
	err := gob.NewEncoder(&b).Encode(m)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (n *NatsEventStore) PublishCreatedPost(post *models.Post) error {
	msg := CreatedPostMessage{
		Post: post,
	}

	data, err := n.encodeMessage(msg)
	if err != nil {
		return err
	}
	return n.conn.Publish(msg.Type(), data)
}

func (n *NatsEventStore) decodeMessage(data []byte, m interface{}) error {
	b := bytes.Buffer{}
	b.Write(data)
	return gob.NewDecoder(&b).Decode(m)
}

func (n *NatsEventStore) OnCreatedPost(f func(CreatedPostMessage)) (err error) {
	msg := CreatedPostMessage{}
	n.postCreatedSub, err = n.conn.Subscribe(msg.Type(), func(m *nats.Msg) {
		n.decodeMessage(m.Data, &msg)
		f(msg)
	})
	return
}

func (n *NatsEventStore) SubscribeCreatedPost() (<-chan CreatedPostMessage, error) {
	m := CreatedPostMessage{}
	n.postCreatedChan = make(chan CreatedPostMessage, 64)
	ch := make(chan *nats.Msg, 64)
	var err error
	n.postCreatedSub, err = n.conn.ChanSubscribe(m.Type(), ch)
	if err != nil {
		return nil, err
	}
	go func() {
		for {
			select {
			case msg := <-ch:
				n.decodeMessage(msg.Data, &m)
				n.postCreatedChan <- m
			}
		}
	}()
	return (<-chan CreatedPostMessage)(n.postCreatedChan), nil
}
