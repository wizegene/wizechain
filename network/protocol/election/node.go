package election

import (
	"WizechainFoundation/wizechain/network/p2p"
	"WizechainFoundation/wizechain/network/protocol"
	"sync"
)

type Node struct {
	Tree *protocol.Tree
	Label int
	AggregateProof []byte
	mu sync.RWMutex
	readiness int
	arity int
	hookOnReadinessUpdate NodeHook
	hookOnRootProofUpdate NodeHook
	Topic p2p.Topic
}

func (n *Node) SetAggregateProof(aggregateProof []byte) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.AggregateProof = aggregateProof

	if n.Tree.Parent != nil {
		n.Tree.Parent.Node.(*Node).IncrementReadiness()
		return
	}
	n.hookOnRootProofUpdate(n)
}

func (n *Node) IncrementReadiness() {
	n.readiness++
	n.hookOnReadinessUpdate(n)
}

func (n *Node) IsReady() bool {
	return n.readiness == n.arity
}

func (n *Node) Job() *Job {
	inputProofs := make([][]byte, n.arity)

	for index, children := range n.Tree.Children {
		inputProofs[index] = children.Node.(*Node).AggregateProof
	}

	return &Job{
		InputProofs: inputProofs,
		label: n.Label,
	}
}

type NodeHook func(n *Node)

func InitializeNodes(t *protocol.Tree, f, g NodeHook) {
	initializeNodes := makeNodeInitializer()
	applyOnReadinessUpdateHook := makeHookOnReadinessUpdateApplier(f)
	t.Walk(initializeNodes)
	t.Walk(applyOnReadinessUpdateHook)
	t.Node.(*Node).hookOnRootProofUpdate = g
}

func makeNodeInitializer() protocol.TreeFunc {

}

func makeNodeMapperByLabel() (protocol.TreeFunc, func() map[int]*Node) {

}