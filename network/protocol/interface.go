package protocol

import (
	"WizechainFoundation/wizechain/network"
	"WizechainFoundation/wizechain/network/identity"
	"github.com/libp2p/go-libp2p-core/peer"
)

type MemberProvider interface {
	GetMembers() []identity.ID
}

type DefaultMembersProvider struct {
	WLP network.WhiteListProvider
}

// GetMembers return the list of all the members
func (d *DefaultMembersProvider) GetMembers() []identity.ID {
	pis, _ := d.WLP.GetPeers()
	members := make([]identity.ID, len(pis))
	for index, b := range pis {
		pi := peer.AddrInfo{}
		_ = pi.UnmarshalJSON(b)
		members[index] = identity.ID(pi.ID)
	}
	return members
}
