package election

import (
	"WizechainFoundation/wizechain/network/identity"
	"WizechainFoundation/wizechain/network/protocol"
	"context"

	log "github.com/sirupsen/logrus"
)

// Round encompasses the computation made in a single
type Round struct {
	ctx    context.Context
	cancel context.CancelFunc

	ID          identity.ID
	Batch       [][]byte
	Participant *Participant
	Members     protocol.IDList

	TopicProvider *TopicProvider
}

// NewRound construct a new round
func NewRound(ctx context.Context, par *Participant, batch [][]byte) *Round {
	ctx, cancel := context.WithCancel(ctx)
	r := &Round{
		ctx:    ctx,
		cancel: cancel,

		ID:          protocol.IDFromBatch(batch),
		Batch:       batch,
		Participant: par,

		Members: par.MemberProvider.GetMembers(),
	}
	return r.WithTopicProvider()
}

// WithTopicProvider add a topic provider in the Round object
func (r *Round) WithTopicProvider() *Round {
	r.TopicProvider = &TopicProvider{
		Network: r.Participant.Network,
		Round:   r,
	}
	return r
}

// Start run the Round
func (r *Round) Start() {
	defer r.Close()
	defer r.WaitForResult()

	log.Infof("Starting to compute the leader ID")
	if r.Participant.ID == r.GetLeaderID() {
		log.Infof("Starting a new leader")
		NewLeader(r).Start()
		return
	}

	log.Infof("Starting a new worker")
	NewWorker(r).Start()
}

// GetLeaderID returns the ID and position of the leader
func (r *Round) GetLeaderID() identity.ID {
	_, res := r.Members.SmallestHigherThan(r.ID)
	return res
}

// WaitForResult waits for the round aggregation result to be published on-chain
func (r *Round) WaitForResult() {
	r.Participant.BatchPubSub.NextBatchDone(r.ctx)
	log.Debugf("Found batch result, round is finished")
}

// Close terminate the round
func (r *Round) Close() {
	r.cancel()
}
