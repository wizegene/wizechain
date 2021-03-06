packag

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

// Worker contains all the logic required to aggregate the proofs
type Worker struct {
	ctx     context.Context
	cancel  context.CancelFunc
	Round   *Round
	Timeout time.Duration // Timeout in second for a proposal
}

// NewWorker returns a newly constructed worker
func NewWorker(r *Round) *Worker {
	return &Worker{
		ctx:     r.ctx,
		cancel:  r.cancel,
		Round:   r,
		Timeout: time.Duration(5) * time.Second,
	}
}

// Aggregate performs an aggregation
func (w *Worker) Aggregate(job *Job) (*Result, error) {

	if len(job.InputProofs) < 2 {
		log.Infof("Wtf happened got a poorly created proof : %v, label : %v",
			job.InputProofs,
			job.Label,
		)
		return nil, fmt.Errorf("Got an improper job")
	}

	data := w.Round.Participant.Aggregator.AggregateTrees(
		job.InputProofs[0], job.InputProofs[1], // TODO: Support multi-arity
	)

	return &Result{
		Result: data,
		Label:  job.Label,
		ID:     w.Round.Participant.ID,
	}, nil
}

// PublishProposal to alert the leader, we are ready
func (w *Worker) PublishProposal() error {
	proposal := &Proposal{
		ID:       w.Round.Participant.ID,
		Deadline: time.Now().Add(w.Timeout),
	}
	return w.Round.TopicProvider.PublishProposal(proposal)
}

// Start the Worker routine
func (w *Worker) Start() error {
	topic := w.Round.TopicProvider.JobTopic(
		w.ctx, w.Round.Participant.ID,
	)

	jobChan, err := topic.Chan()
	if err != nil {
		return err
	}

	go func() {
		defer w.cancel()

		for {

			err := w.PublishProposal()
			if err != nil {
				return
			}

			propCtx, propCancel := context.WithTimeout(
				context.Background(),
				w.Timeout,
			)

			select {
			case <-propCtx.Done():
				log.Info("Sent proposal expired")
				propCancel()
				continue

			case <-w.ctx.Done():
				propCancel()
				return

			case jobEncoded := <-jobChan:
				propCancel()
				job, err := MarshalledJob(jobEncoded).Decode()
				if err != nil {
					return
				}

				log.Infof("Got a job with label %v", job.Label)

				res, err := w.Aggregate(job)
				if err != nil {
					return
				}

				_ = w.Round.TopicProvider.PublishResult(res)
			}
		}
	}()

	return nil
}