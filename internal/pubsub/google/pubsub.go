package google

import (
	"sync"

	"cloud.google.com/go/pubsub"

	"github.com/elmagician/kactus/internal/picker"
)

// Client provides a structure to manage GCP instances.
type Client struct {
	name string
	cli  *pubsub.Client

	received          []*pubsub.Message
	nbMessageReceived int

	mu sync.RWMutex
}

// Pubsub initializes a Client instance for GCP pubsub. It will persists
// instance using provided name to be injectable through {{gcp.name}} in
// gherkin steps.
//
// Providing an empty picker does not impact initialization. It will just not
// picked the instance.
//
// You can always call Client.Persist to save client instance in a picker store.
func Pubsub(name string, cli *pubsub.Client, store *picker.Store) *Client {
	client := &Client{name: name, cli: cli}

	if store != nil {
		client.Persist(store)
	}

	return client
}

// Persist persists client instance through picker instance using gcp.name key.
func (ps *Client) Persist(store *picker.Store) {
	store.Pick(
		"gcp."+ps.name,
		picker.InstanceItem{Kind: picker.GCP, Instance: ps},
		picker.InstanceValue,
	)
}

// Reset instance.
func (ps *Client) Reset() {
	ps.received = []*pubsub.Message{}
	ps.nbMessageReceived = 0

	ResetLog()
}
