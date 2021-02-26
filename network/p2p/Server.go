package p2p

import (
	"WizechainFoundation/wizechain/network/identity"
	"context"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-pubsub"
	"time"
)

type Server struct {
	Host      host.Host
	PubSub    *pubsub.PubSub
	Bootstrap *Bootstrapping
}

func NewServerWithID(wlp WhiteListProvider, priv *identity.PrivKey, addr string) (*Server, error) {
	listenAddr, err := ListenAddress(addr)
	if err != nil {
		return nil, err
	}

	privateKey, _, _ := priv.Libp2p()
	privP2P, err := Identity(privateKey)
	if err != nil {
		return nil, err
	}

	hs, err := libp2p.New(context.Background(), listenAddr, privP2P)
	if err != nil {
		return nil, err
	}
	ps, err := pubsub.NewGossipSub(context.Background(), hs)

	bstr := NewBootstrapping(
		context.Background(), hs, 5, 9, wlp, time.Duration(2)*time.Second)

	return &Server{
		Host:      hs,
		PubSub:    ps,
		Bootstrap: bstr,
	}, nil

}

func (s *Server) Start() {
	_ = s.Bootstrap.Start()
}

func (s *Server) GetPeerInfo() peer.AddrInfo {
	return peer.AddrInfo{
		ID:    s.Host.ID(),
		Addrs: s.Host.Addrs(),
	}
}

func (s *Server) GetTopic(ctx context.Context, topic string) *Topic {
	res := Topic{
		ctx:  ctx,
		ps:   s.PubSub,
		Name: topic,
	}
	return &res
}

func (s *Server) Publish(ctx context.Context, topic string, msg []byte) error {
	t, _ := s.PubSub.Join(topic)
	return t.Publish(ctx, msg)

}
