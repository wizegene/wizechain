package p2p

import (
	"context"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	ma "github.com/multiformats/go-multiaddr"
	"sync"
	"time"
)

type Bootstrapping struct {
	ctx                 context.Context
	Host                host.Host
	mu                  sync.RWMutex
	whiteList           map[peer.ID][]ma.Multiaddr
	peerChan            chan [][]byte
	sortedPeers         []peer.ID
	minConns            int
	maxConns            int
	target              int
	bootstrappingPeriod time.Duration
}
