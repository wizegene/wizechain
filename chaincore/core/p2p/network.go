package p2p

import (
	csms "github.com/libp2p/go-conn-security-multistream"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/peerstore"
	"github.com/libp2p/go-libp2p-core/sec/insecure"
	swarm "github.com/libp2p/go-libp2p-swarm"
	yamux "github.com/libp2p/go-libp2p-yamux"
	msmux "github.com/libp2p/go-stream-muxer-multistream"
)

type Network struct {
	ID                         string
	nType                      string
	Configuration              map[string]string
	MinMasterNodesNeededToBoot int
	MaxConcurrentMasterNodes   int
	CurrentSpeedMBPS           float64
	MasterNodes                MappedNodes
	workQueue                  map[string][]byte
	proofChan                  chan []byte
	sigChan                    chan bool
	statusCode                 int
	minFees                    float64
	maxFees                    float64
}

type MappedNodes struct {
	Active   map[string]*MasterNode
	Inactive map[string]*MasterNode
	Invalid  map[string]*MasterNode
}

var n *swarm.Swarm

func InitNetwork() {

	net := new(Network)

	swarmID := n.LocalPeer()
	pk := n.Peerstore().PrivKey(swarmID)
	secMuxer := new(csms.SSMuxer)
	secMuxer.AddTransport(insecure.ID, insecure.NewWithIdentity(swarmID, pk))
	stMuxer := msmux.NewBlankTransport()
	stMuxer.AddTransport("/yamux/1.0.0", yamux.DefaultTransport)

	net.Bootstrap()

}

func addPeerToSwarm(ps peerstore.Peerstore) {
	for _, p := range ps.Peers() {
		pid, _ := p.Marshal()
		pub, _ := p.ExtractPublicKey()
		_ = ps.AddPubKey(peer.ID(pid), pub)

	}
}

func (n *Network) Bootstrap() {

}
