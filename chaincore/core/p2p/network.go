package p2p

type Network struct {
	ID                         string
	nType                      string
	Configuration              map[string]string
	MinMasterNodesNeededToBoot int
	MaxConcurrentMasterNodes   int
	CurrentSpeedMBPS           float64
	UpNodes                    []*MasterNode
	DownNodes                  []*MasterNode
	workQueue                  map[string][]byte
	proofChan                  chan []byte
	sigChan                    chan bool
	statusCode                 int
	minFees                    float64
	maxFees                    float64
}

func InitNetwork() {

	net := new(Network)
	net.Bootstrap()

}

func (n *Network) Bootstrap() {

}
