package election

import (
	net "WizechainFoundation/wizechain/network"
	"WizechainFoundation/wizechain/network/identity"
	"WizechainFoundation/wizechain/network/protocol"
	"context"
)

type Participant struct {
	ctx context.Context

	ID             identity.ID
	MemberProvider protocol.MemberProvider

	Network     net.PubSub
	BatchPubSub BatchPubSub
	Aggregator  aggregator.Aggregator
}

// NewParticipant ..
func NewParticipant(
	ctx context.Context,
	id identity.ID,
	aggregator aggregator.Aggregator,
	memberProvider protocol.MemberProvider,
	batchPubSub BatchPubSub,
	network net.PubSub,

) *Participant {
	return &Participant{
		ctx: ctx,

		MemberProvider: memberProvider,
		ID:             id,

		Network:     network,
		BatchPubSub: batchPubSub,
		Aggregator:  aggregator,
	}
}

// Start the main routine of the participant
func (par *Participant) Start() {
	go par.background()
}

func (par *Participant) background() {
	for {
		batch, err := par.BatchPubSub.NextNewBatch(par.ctx)
		if err != nil {
			return
		}
		NewRound(par.ctx, par, batch).Start()
	}
}
