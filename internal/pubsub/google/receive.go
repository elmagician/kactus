package google

import (
	"context"

	"cloud.google.com/go/pubsub"
	"go.uber.org/zap"
)

// ReceiveOn set up pubsub message reception from subscription.
// Messages will be stored in a slice until Client.Reset is called.
func (ps *Client) ReceiveOn(subscription string, acknowledgeOnReception bool) error {
	log.Debug("initializing reception", zap.String("subscription", subscription))

	return ps.cli.
		Subscription(subscription).
		Receive(context.Background(), ps.receive(acknowledgeOnReception))
}

func (ps *Client) receive(acknowledgeOnReception bool) func(ctx context.Context, msg *pubsub.Message) {
	return func(ctx context.Context, msg *pubsub.Message) {
		ps.mu.Lock()
		defer ps.mu.Unlock()

		log.Debug("message received", zap.String("data", string(msg.Data)), zap.Reflect("metadata", msg.Attributes))

		ps.received = append(ps.received, msg)
		ps.nbMessageReceived++

		if acknowledgeOnReception {
			msg.Ack()
		}
	}
}
