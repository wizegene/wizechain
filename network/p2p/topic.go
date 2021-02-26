package p2p

import (
	"context"
	"fmt"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"sync"
)

type Topic struct {
	ps *pubsub.PubSub
	ctx context.Context
	Subs *pubsub.Subscription
	Name string
	cancelChan context.CancelFunc
	onceChan sync.Once
}

func (t *Topic) Join(topic string) (*pubsub.Topic, error) {
	return t.ps.Join(topic)
}

func (t *Topic) Publish(topic string, msg []byte) error {

	top, err := t.Join(topic)
	if err != nil {
		return err
	}
	err = top.Publish(t.ctx, msg)
	if err != nil {
		return err
	}

	return nil

}

func (t *Topic) Chan() (<-chan []byte, error) {
	var out chan []byte
	var err error

	t.onceChan.Do(func() {
		s, err := t.Join(t.Name)
		if err != nil {
			return
		}
		subs, err := s.Subscribe()
		if err != nil {
			return
		}
		t.Subs = subs
		out = make(chan []byte, 20)
		go t.background(out)

	})


	if err != nil {
		return nil, err
	}

	if out == nil {
		return nil, fmt.Errorf("topic channel can only be called once")
	}

	return out, nil

}

func (t *Topic) Close() {
	if t.cancelChan != nil {
		t.cancelChan()
	}
	t.Subs.Cancel()
}

func (t *Topic) background(out chan []byte) {
	defer close(out)
	ctx, can := context.WithCancel(context.Background())
	t.cancelChan = can

	for {
		fromSub := make(chan *pubsub.Message)
		errorChan := make(chan error)

		go func() {
			defer close(errorChan)
			defer close(fromSub)

			msg, err := t.Subs.Next(ctx)
			if err != nil {
				errorChan <- err
				return
			}

			fromSub <- msg
		}()

		select {
		case <-t.ctx.Done():
			return
			case <-errorChan:
				t.cancelChan()
				return
				case msg := <-fromSub:
					out <- msg.GetData()
		}
	}
}
