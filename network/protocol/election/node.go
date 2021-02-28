package election

import (
	"WizechainFoundation/wizechain/network/p2p"
	"WizechainFoundation/wizechain/network/protocol"
	"sync"
)

type Node struct {
	Tree                  *protocol.Tree
	Label                 int
	AggregateProof        []byte
	mu                    sync.RWMutex
	readiness             int
	arity                 int
	hookOnReadinessUpdate NodeHook
	hookOnRootProofUpdate NodeHook
	Topic                 p2p.Topic
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
		label:       n.Label,
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
	counter := protocol.MakeCounter() // Will be used to label the nodes with unique numbers
	return func(t *protocol.Tree) {

		arity := 0
		if t.Children != nil {
			arity = len(t.Children)
		}

		t.Node = &Node{
			Tree:      t,
			Label:     counter(), // Gives a unique label to each node
			arity:     arity,
			readiness: 0,
		}
	}
}

/ makeHookOnReadinessUpdateApplier returns a TreeFunc setting the NodeHook
func makeHookOnReadinessUpdateApplier(f NodeHook) protocol.TreeFunc {
	return func(t *protocol.Tree) {
		t.Node.(*Node).hookOnReadinessUpdate = f
	}
}

// makeNodeMapperByLabel returns a map of node indexed by their label
func makeNodeMapperByLabel() (protocol.TreeFunc, func() map[int]*Node) {
	nodeMap := make(map[int]*Node)

	mapNode := func(t *protocol.Tree) { nodeMap[t.Node.(*Node).Label] = t.Node.(*Node) }
	getNodeMap := func() map[int]*Node { return nodeMap }

	return mapNode, getNodeMap
}
