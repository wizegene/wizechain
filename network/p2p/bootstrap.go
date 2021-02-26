package p2p

import (
	"context"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/libp2p/go-libp2p-core/network"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

type Bootstrapping struct {
	ctx                 context.Context
	Host                host.Host
	mu                  sync.RWMutex
	peerChan            chan [][]byte
	sortedPeers         []peer.ID
	whitelist           map[peer.ID][]ma.Multiaddr
	minConns            int
	target              int
	maxConns            int
	Wlp                 WhiteListProvider
	bootstrappingPeriod time.Duration
}

func NewBootstrapping(
	ctx context.Context,
	host host.Host,
	minConns int,
	maxConns int,
	Wlp WhiteListProvider,
	bootstrappingPeriod time.Duration) *Bootstrapping {

	pchan, _ := Wlp.GetPeersChan()
	return &Bootstrapping{
		ctx:                 ctx,
		Host:                host,
		Wlp:                 Wlp,
		peerChan:            pchan,
		minConns:            minConns,
		maxConns:            maxConns,
		bootstrappingPeriod: bootstrappingPeriod,
	}

}

func (b *Bootstrapping) Start() error {
	initialPeers, err := b.Wlp.GetPeers()
	if err != nil {
		return err
	}

	b.SetNewWhiteList(initialPeers)
	go b.background()
	return nil
}

func (b *Bootstrapping) background() error {

	err := b.RunBootstrap()
	if err != nil {
		return err
	}

	ticker := time.NewTicker(b.bootstrappingPeriod)
	for {
		select {

		case <-b.ctx.Done():

		case <-ticker.C:
			err := b.RunBootstrap()
			if err != nil {
				return err
			}
		case newPeers := <-b.peerChan:
			err := b.SetNewWhiteList(newPeers)
			if err != nil {
				return err
			}
			b.TrimUnlistedPeers()

		}
	}

}

func (b *Bootstrapping) RunBootstrap() (err error) {

	nConns := len(b.Host.Network().Peers())
	if nConns < b.minConns {
		log.Debugf("Got %v connexion, adding new one to reach the ceil %v", nConns, b.minConns)
		return b.AddConnections()
	}
	if nConns > b.maxConns {
		log.Debugf("Got %v connexion, pruning some to reach the floor %v", nConns, b.minConns)
		return b.PruneConnections()
	}

	log.Debugf("Got %v connexion, keeping it like this", nConns)
	return nil

}

func (b *Bootstrapping) AddConnections() error {
	b.mu.RLock()
	defer b.mu.RUnlock()
	self := b.Host.ID()
	nConns := len(b.Host.Network().Peers())

ForEachPeerConnectLoop:
		for _, id := range b.sortedPeers {
			canConnect := b.Host.Network().Connectedness(id)
			if canConnect == network.CanConnect || canConnect == network.NotConnected {
				if id == self {
					continue ForEachPeerConnectLoop
				}

				addrInfo := peer.AddrInfo{
					ID: id,
					Addrs: b.whitelist[id],
				}

				if err := b.Host.Connect(context.Background(), addrInfo); err != nil {
					continue ForEachPeerConnectLoop
				}

				log.Debugf("Found new peer : %v, from : %v", id, b.Host.ID())
				nConns++
				if nConns == b.target {
					return nil
				}
			}
		}
		return nil
}

func (b *Bootstrapping) PruneConnections() error {

	b.mu.RLock()
	defer b.mu.RUnlock()
	nConns := len(b.Host.Network().Peers())

	ForEachPeerDisconnectLoop:
		for i := len(b.sortedPeers) - 1; i > 0; i-- {
			id := b.sortedPeers[i]
			connectedNess := b.Host.Network().Connectedness(id)
			if connectedNess == network.Connected {
				err := b.Host.Network().ClosePeer(id)
				if err != nil {
					continue ForEachPeerDisconnectLoop
				}

				nConns--
				if nConns == b.target {
					return nil
				}
			}
		}


	return nil
}

func (b *Bootstrapping) SetNewWhiteList(newPeers [][]byte) error {

	newWhiteList := make(map[peer.ID][]ma.Multiaddr)
	for _, marshalled := range newPeers {
		var unmarshalled peer.AddrInfo
		err := unmarshalled.UnmarshalJSON(marshalled)
		if err != nil {
			return err
		}
		newWhiteList[unmarshalled.ID] = unmarshalled.Addrs
	}

	bxd := byXORDistance{
		slice: make([]peer.ID, len(newPeers)),
		reference: b.Host.ID(),
	}

	index := 0
	for id := range newWhiteList {
		bxd.slice[index] = id
		index++
	}
	bxd.Sort()
	b.mu.Lock()
	b.whitelist = newWhiteList
	b.sortedPeers = bxd.slice
	b.mu.Unlock()
	return nil


}

func (b *Bootstrapping) TrimUnlistedPeers() {
	for _, connectedPeer := range b.Host.Network().Peers() {
		if !b.authorized(connectedPeer) {
			log.Infof("Found a non whitelisted connection")
			b.Host.Network().ClosePeer(connectedPeer)
		}
	}
}

func (b *Bootstrapping) Permissionize() {
	onlyWhitelisted := func(conn network.Conn) {
		if !b.authorized(conn.RemotePeer()) {
			conn.Close()
		}
	}
	b.Host.Network().SetConnHandler(onlyWhitelisted)
}

func (b *Bootstrapping) authorized(p peer.ID) bool {

	_, listed := b.whitelist[p]
	return listed

}

type WhiteListProvider interface {
	GetPeersChan() (chan [][]byte, error)
	GetPeers() ([][]byte, error)
}
