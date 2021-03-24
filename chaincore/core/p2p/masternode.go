package p2p

import (
	"bufio"
	"context"
	"fmt"
	"github.com/ipfs/go-log"
	lru "github.com/karlseguin/ccache/v2"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	noise "github.com/libp2p/go-libp2p-noise"
	"github.com/libp2p/go-libp2p/p2p/discovery"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/rogpeppe/fastuuid"
	"github.com/wizegene/wizechain/chaincore/core"
	"github.com/wizegene/wizechain/chaincore/core/p2p/handlers"
	"os"
	"runtime"

	"time"
)

var logger = log.Logger("rendezvous")

const DiscoveryInterval = time.Minute

type MasterNode struct {
	ID                       peer.ID
	Addrs                    []ma.Multiaddr
	Bootstrapper             *Bootstrapping
	rawHost                  host.Host
	Context                  context.Context
	PubSub                   *pubsub.PubSub
	NetworkID                string
	MasterName               string
	MasterServices           map[string]bool
	Host                     string
	Port                     string
	HostID                   string
	HostAddrs                []string
	DHT                      *dht.IpfsDHT
	ChildrenNodes            []*peer.AddrInfo
	Neighbors                []*MasterNode
	NodeState                *core.State
	workpool                 map[int][]byte
	cache                    *lru.Cache
	whitelist                []*peer.ID
	blacklist                []*peer.ID
	connectedMasterNodes     map[string]bool
	acceptConnsOnStart       bool
	status                   chan bool
	maxChildrens             int
	maxTotalPeersWithinGroup int
	maxServerLoad            float64
	currentServerLoad        float64
	maxCPUs                  int
	maxMem                   float32
	localTime                time.Time
	timeSinceStart           time.Duration
}

func NewMasterNode() *MasterNode {

	name := fastuuid.MustNewGenerator().Hex128()
	mn := &MasterNode{}
	rawHost, err := libp2p.New(context.Background(), libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/9000"))
	libp2p.ChainOptions(libp2p.DefaultSecurity, libp2p.Security(noise.ID, noise.New))
	if err != nil {
		panic(err)
	}

	mn.rawHost = rawHost
	mn.configure(name, "wize_1", "0.0.0.0", "5565", 2, 2, nil, true)
	mn.currentServerLoad = GetLoadAverage()
	return mn

}

func (m *MasterNode) newGossipPubSub() {
	ps, err := pubsub.NewGossipSub(m.Context, m.rawHost)
	if err != nil {
		panic(err)
	}
	m.PubSub = ps
}

type discoveryNotifee struct {
	h host.Host
}

func (d *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	fmt.Printf("discovered new peer %s\n", pi.ID.Pretty())
	err := d.h.Connect(context.Background(), pi)
	if err != nil {
		fmt.Printf("error connecting to peer %s: %s\n", pi.ID.Pretty(), err)
	}
}

func (m *MasterNode) discoverServices() (err error) {

	var disc discovery.Service

	for service, _ := range m.MasterServices {

		disc, err = discovery.NewMdnsService(m.Context, m.rawHost, DiscoveryInterval, service)
		if err != nil {
			return
		}
		n := discoveryNotifee{h: m.rawHost}
		disc.RegisterNotifee(&n)
	}
	return nil

}

func (m *MasterNode) configure(mastername, networkid, host, port string, maxChildrens,
	maxTotalPeers int, servicesAvailable []string, acceptConnsOnStart bool) {

	m.Context = context.Background()
	m.NetworkID = networkid
	m.MasterName = mastername
	services := make(map[string]bool, len(servicesAvailable))
	for _, service := range servicesAvailable {
		services[service] = true
	}

	m.Host = host
	m.Port = port
	m.ChildrenNodes = make([]*peer.AddrInfo, 0)
	m.Neighbors = make([]*MasterNode, 0)
	m.NodeState = &core.State{}
	m.workpool = make(map[int][]byte, 0)
	m.cache = lru.New(lru.Configure().MaxSize(10000).ItemsToPrune(100))
	m.whitelist = make([]*peer.ID, 0)
	m.blacklist = make([]*peer.ID, 0)
	m.connectedMasterNodes = make(map[string]bool, 0)
	m.acceptConnsOnStart = acceptConnsOnStart
	m.status = make(chan bool, 0)
	m.maxChildrens = maxChildrens
	m.maxTotalPeersWithinGroup = maxTotalPeers
	m.maxCPUs = runtime.NumCPU() - 1
	m.maxMem = 4096
	m.maxServerLoad = 10.0
	m.localTime = new(time.Time).Local()

	kdht, err := dht.New(m.Context, m.rawHost)
	if err != nil {
		panic(err)
	}

	m.DHT = kdht

}

func (m *MasterNode) connectToOtherMasterNode(peerAddr string) {
	peerMA, err := ma.NewMultiaddr(peerAddr)
	if err != nil {
		panic(err)
	}
	peerAddrInfo, err := peer.AddrInfoFromP2pAddr(peerMA)
	if err := m.rawHost.Connect(m.Context, *peerAddrInfo); err != nil {
		panic(err)
	}
	m.connectedMasterNodes[peerAddr] = true

}

func readData(rw *bufio.ReadWriter) {
	for {
		str, err := rw.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from buffer")
			panic(err)
		}

		if str == "" {
			return
		}
		if str != "\n" {
			// Green console colour: 	\x1b[32m
			// Reset console colour: 	\x1b[0m
			fmt.Printf("\x1b[32m%s\x1b[0m> ", str)
		}

	}
}

func writeData(rw *bufio.ReadWriter) {
	stdReader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		sendData, err := stdReader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from stdin")
			panic(err)
		}

		_, err = rw.WriteString(fmt.Sprintf("%s\n", sendData))
		if err != nil {
			fmt.Println("Error writing to buffer")
			panic(err)
		}
		err = rw.Flush()
		if err != nil {
			fmt.Println("Error flushing buffer")
			panic(err)
		}
	}
}

func handleStream(stream network.Stream) {
	logger.Info("Got a new stream!")

	if stream.Protocol() == "sync" {
		handlers.SyncHandler(stream)
	}
	// Create a buffer stream for non blocking read and write.

	// 'stream' will stay open until you close it (or the other side closes it).
}
