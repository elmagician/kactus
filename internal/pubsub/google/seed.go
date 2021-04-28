package google

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"

	"cloud.google.com/go/pubsub"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"

	"github.com/elmagician/kactus/internal/picker"
)

type (
	Message struct {
		Topic    string                 `yaml:"topic"`
		Data     map[string]interface{} `yaml:"data"`
		Metadata map[string]string      `yaml:"metadata"`
	}

	Manifest struct {
		Send            []Message `yaml:"send"`
		Receive         []string  `yaml:"receive"`
		AutoAcknowledge bool      `yaml:"acknowledgeOnReceive"`
	}
)

func (m Manifest) Load(cli *Client) error {
	for _, subscription := range m.Receive {
		locSub := subscription

		go func() {
			err := cli.ReceiveOn(locSub, m.AutoAcknowledge)
			if err != nil && !errors.Is(err, context.Canceled) {
				log.Error("error while receiving", zap.Error(err))
			}
		}()
	}

	for _, message := range m.Send {
		bytesData, err := json.Marshal(message.Data)
		if err != nil {
			return err
		}

		if err := cli.sendMessage(message.Topic, &pubsub.Message{Data: bytesData, Attributes: message.Metadata}); err != nil {
			return err
		}
	}

	return nil
}

func (ps *Client) SeedFromFile(filePath string, store *picker.Store) error {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	strFile := string(file)

	if store != nil {
		log.Debug("injecting picked variables")

		strFile = store.InjectAll(strFile)
	}

	var manifest Manifest
	if err := yaml.Unmarshal([]byte(strFile), &manifest); err != nil {
		return err
	}

	return manifest.Load(ps)
}
