package p2p

import (
	"context"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-pubsub"
)

type Server struct {
	Host host.Host
	PubSub *pubsub.PubSub
	Bootstrap *Bootstrapping
}

func NewServerWithID(wlp WhiteListProvider, priv *identity.PrivateKey, addr string) (*Server, error) {
	listenAddr, err := ListenAddress(addr)
	if err != nil {
		return nil, err
	}

	privP2P, err := Identity(priv.Libp2p())

}



func (s *Server) Start() {
	_ = s.Bootstrap.Start()
}

func (s *Server) GetPeerInfo() peer.AddrInfo {
	return peer.AddrInfo{
		ID: s.Host.ID(),
		Addrs: s.Host.Addrs(),
	}
}

func (s *Server) GetTopic(ctx context.Context, topic string) {

}