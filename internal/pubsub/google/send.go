package google

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/pubsub"
	"github.com/cucumber/godog"
	"go.uber.org/zap"
)

type paramMessage struct {
	Data interface{}       `json:"data"`
	Meta map[string]string `json:"metadata"`
}

// Send send message to topic
func (ps *Client) Send(topic string, data *godog.DocString) error {
	var msg paramMessage
	if err := json.Unmarshal([]byte(data.Content), &msg); err != nil {
		return err
	}

	serializedData, err := json.Marshal(msg.Data)
	if err != nil {
		return err
	}

	return ps.sendMessage(topic, &pubsub.Message{Data: serializedData, Attributes: msg.Meta})
}

func (ps *Client) sendMessage(topic string, msg *pubsub.Message) error {
	log.Debug(
		"sending message",
		zap.String("data", string(msg.Data)),
		zap.Reflect("metadata", msg.Attributes),
		zap.String("topic", topic),
	)

	top := ps.cli.Topic(topic)
	defer top.Stop()

	id, err := top.
		Publish(context.Background(), msg).
		Get(context.Background())

	if err != nil {
		return fmt.Errorf("%w for message: %s", err, id)
	}

	return nil
}
