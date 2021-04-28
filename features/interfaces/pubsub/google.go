package pubsub

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/elmagician/godog"

	"github.com/elmagician/kactus/features/interfaces/picker"
	internalPicker "github.com/elmagician/kactus/internal/picker"
	"github.com/elmagician/kactus/internal/pubsub/google"
)

var (
	// ErrUnknown is raised when trying to operate on an unpicked gcp client.
	ErrUnknown = errors.New("unknown client")

	// ErrInvalidInstance is raised when trying to  operate on a non gcp instance.
	ErrInvalidInstance = errors.New("expected GCP instance")
)

var (
	instance *Google
	once     sync.Once
)

type (
	// Google manages instance for kactus Google and provides tools to
	// manipulate and asserts on known  postgres databases.
	Google struct {
		picker *internalPicker.Store
	}

	// GoogleInfo provides a structure to  register a new
	// google pubsub client instance to kactus.
	//
	// Key will be used to to pick instance.
	// `gcp.` will be append to provided key when picking.
	//
	// Client indicates *pubsub.Client instance to store under provided key.
	GoogleInfo struct {
		Key    string
		Client *pubsub.Client
	}
)

// NewGoogle initializes a postgres kactus manager for
// provided google pubsub client instances.
func NewGoogle(pickerInstance *picker.Picker, clients ...GoogleInfo) *Google {
	once.Do(func() {
		instance = &Google{
			picker: pickerInstance.This(),
		}

		// initialize clients
		for _, cli := range clients {
			google.Pubsub(cli.Key, cli.Client, pickerInstance.This())
		}
	})

	return instance
}

// HasMessage asserts client received a message matching provided data.
func (gcp *Google) HasMessage(instance string, within int, data *godog.Table) error {
	gcpInstance, err := gcp.getGCPInstance(instance)
	if err != nil {
		return err
	}

	return gcpInstance.AssertMessageReceived(data, time.Duration(within)*time.Second)
}

// HasMessageWithMetadata asserts client received a message matching provided meta data.
func (gcp *Google) HasMessageWithMetadata(instance string, within int, data *godog.Table) error {
	gcpInstance, err := gcp.getGCPInstance(instance)
	if err != nil {
		return err
	}

	return gcpInstance.AssertMessageMetadata(data, time.Duration(within)*time.Second)
}

func (gcp *Google) SendMessage(instance, topic string, data *godog.DocString) error {
	gcpInstance, err := gcp.getGCPInstance(instance)
	if err != nil {
		return err
	}

	splitTopic := strings.Split(topic, ".")
	if len(splitTopic) < 2 { // nolint: gomnd
		return gcpInstance.Send(topic, data)
	}

	return gcpInstance.Send(splitTopic[1], data)
}

// Reset resets google instance.
func (gcp *Google) Reset() {
	google.ResetLog()
}

func (gcp *Google) getGCPInstance(instance string) (*google.Client, error) {
	kind, localInstance, exists := gcp.picker.GetInstance(instance)
	if !exists {
		return nil, fmt.Errorf("%w: %s", ErrUnknown, instance)
	}

	if kind != internalPicker.GCP {
		return nil, fmt.Errorf("%w", ErrInvalidInstance)
	}

	gcpInstance, ok := localInstance.(*google.Client)
	if !ok {
		return nil, fmt.Errorf("%w, got: %T", ErrInvalidInstance, gcp)
	}

	return gcpInstance, nil
}
