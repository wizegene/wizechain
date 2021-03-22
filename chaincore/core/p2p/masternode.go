package p2p

import (
	"context"
	lru "github.com/karlseguin/ccache/v2"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/wizegene/wizechain/chaincore/core"
	"runtime"
	"time"
)

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
