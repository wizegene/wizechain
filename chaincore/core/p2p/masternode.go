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
	discovery "github.com/libp2p/go-libp2p-discovery"
	"github.com/wizegene/wizechain/chaincore/core/p2p/handlers"
	"os"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/wizegene/wizechain/chaincore/core"
	"runtime"
	"sync"
	"time"
)

var logger = log.Logger("rendezvous")

type MasterNode struct {
	ID                       peer.ID
	Addrs                    []ma.Multiaddr
	rawHost                  host.Host
	Context                  context.Context
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
	maxServerLoad            float32
	currentServerLoad        chan float32
	maxCPUs                  int
	maxMem                   float32
	localTime                time.Time
	timeSinceStart           time.Duration
}

func NewMasterNode() *MasterNode {
	return &MasterNode{}

}

func (m *MasterNode) create(mastername, networkid, host, port string, maxChildrens,
	maxTotalPeers int, servicesAvailable []string, acceptConnsOnStart bool) {

	m.configure(mastername, networkid, host, port, maxChildrens,
		maxTotalPeers, servicesAvailable, acceptConnsOnStart)

	h, err := libp2p.New(m.Context, libp2p.ListenAddrStrings("/ip4/"+m.Host+"/tcp/"+m.Port))
	if err != nil {
		panic(err)
	}
	defer h.Close()
	m.HostID = h.ID().String()
	m.ID = h.ID()
	m.Addrs = h.Addrs()
	m.rawHost = h

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
	m.currentServerLoad = make(chan float32, 0)
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

func (m *MasterNode) bootstrapNode(addrs []string) {

	m.rawHost.SetStreamHandler("sync", handleStream)

	if err := m.DHT.Bootstrap(m.Context); err != nil {
		panic(err)
	}

	var wg sync.WaitGroup

	for _, peerAddr := range addrs {
		addr, _ := ma.NewMultiaddr(peerAddr)
		peerinfo, _ := peer.AddrInfoFromP2pAddr(addr)
		wg.Add(1)
		go func() {

			defer wg.Done()
			if err := m.rawHost.Connect(m.Context, *peerinfo); err != nil {
				logger.Debug(err)
			} else {
				logger.Info("Connection established with bootstrap node:", *peerinfo)
			}

		}()
	}
	wg.Wait()

	logger.Info("announcing the masternode...")
	routingDiscovery := discovery.NewRoutingDiscovery(m.DHT)
	discovery.Advertise(m.Context, routingDiscovery, m.NetworkID)
	logger.Info("masternode announced, continuing...")
	peerChan, err := routingDiscovery.FindPeers(m.Context, m.NetworkID)
	if err != nil {
		panic(err)
	}

	for peer := range peerChan {
		if peer.ID == m.rawHost.ID() {
			continue
		}

		if m.maxTotalPeersWithinGroup <= len(m.ChildrenNodes) {
			logger.Warn("max number of peers within group reached", len(m.ChildrenNodes))
			continue
		}

		if m.MasterServices["fullnode"] == false {
			continue
		}

		logger.Debug("found other peer:", peer)
		logger.Debug("connecting to peer:", peer)

		logger.Info("will try to synchronize the wizechain network")

		stream, err := m.rawHost.NewStream(m.Context, peer.ID, "sync")
		if err != nil {
			logger.Warn("connection failed:", err)
			continue
		}

		rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))
		readData(rw)
		writeData(rw)

	}

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
