package google

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/cucumber/godog"
	"go.uber.org/zap"

	"github.com/elmagician/kactus/internal/interfaces"
	match "github.com/elmagician/kactus/internal/matchers"
)

var (
	ErrNoMatch        = errors.New("no match")
	ErrUnknownMessage = errors.New("positions does not match a known message")
)

// AssertMessageReceived assert pubsub client received a message
// matching providing assertions from godog.Table.
//
// /!\ It does not apply filters to message metadata.
func (ps *Client) AssertMessageReceived(data *godog.Table, within time.Duration) error {
	return ps.analyzeOnReception(data, within, ps.AssertXMessageData)
}

// AssertMessageMetadata assert pubsub client received a message metadata
// matching providing assertions from godog.Table.
//
// /!\ It does not apply filters to message data.
func (ps *Client) AssertMessageMetadata(data *godog.Table, within time.Duration) error {
	return ps.analyzeOnReception(data, within, ps.AssertXMessageMetadata)
}

// AssertXMessageData assert pubsub client Xth message data matches provided
// conditions.
//
// /!\ It does not apply filters to message metadata.
func (ps *Client) AssertXMessageData(position int, expectedMessage *godog.Table) error {
	if position+1 > ps.nbMessageReceived {
		return fmt.Errorf("%w: position %d is outside of known message range %d", ErrUnknownMessage, position+1, ps.nbMessageReceived)
	}

	var (
		key, val      string
		matcher       string
		actualMessage map[string]interface{}
		msg           = ps.received[position].Data
		head          = expectedMessage.Rows[0].Cells
	)

	if err := json.Unmarshal(msg, &actualMessage); err != nil {
		return err
	}

	// nolint: dupl
	for i := 1; i < len(expectedMessage.Rows); i++ {
		for n, cell := range expectedMessage.Rows[i].Cells {
			switch head[n].Value {
			case "field":
				key = cell.Value
			case "matcher":
				matcher = cell.Value
			case "value":
				val = cell.Value
			default:
				return fmt.Errorf("unexpected column name %s", head[n].Value)
			}
		}

		if actualVal, exists := interfaces.GetFieldFromPath(actualMessage, key); exists {
			if err := match.Assert(matcher, actualVal, val); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("%w: key %v do not exist", ErrNoMatch, key)
		}

		key = ""
		val = ""
		matcher = ""
	}

	return nil
}

// AssertXMessageMetadata assert pubsub client Xth message metadata matches
// provided conditions.
//
// /!\ It does not apply filters to message data.
func (ps *Client) AssertXMessageMetadata(position int, expectedMetadata *godog.Table) error {
	if position+1 > len(ps.received) {
		return fmt.Errorf("%w: position %d is outside of known message range %d", ErrUnknownMessage, position+1, len(ps.received))
	}

	var (
		key, val string
		matcher  string
		metadata = ps.received[position].Attributes
		head     = expectedMetadata.Rows[0].Cells
	)

	// nolint: dupl
	for i := 1; i < len(expectedMetadata.Rows); i++ {
		for n, cell := range expectedMetadata.Rows[i].Cells {
			switch head[n].Value {
			case "field":
				key = cell.Value
			case "matcher":
				matcher = cell.Value
			case "value":
				val = cell.Value
			default:
				return fmt.Errorf("unexpected column name %s", head[n].Value)
			}
		}

		if actualVal, exists := interfaces.GetFieldFromPath(metadata, key); exists {
			if err := match.Assert(matcher, actualVal, val); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("%w: key %v do not exist", ErrNoMatch, key)
		}

		key = ""
		val = ""
		matcher = ""
	}

	return nil
}

func (ps *Client) analyzeOnReception(data *godog.Table, within time.Duration, singleAsserter func(int, *godog.Table) error) error {
	var (
		outdatedCheck = time.Now().Add(within)

		analyzeMessages = func(startAt, to int) (int, error) {
			var i int
			for i = startAt; i < to; i++ {
				err := singleAsserter(i, data)
				if err != nil && !errors.Is(err, ErrNoMatch) && !errors.Is(err, match.ErrUnmatched) {
					return 0, err
				}

				if err == nil {
					return i, nil
				}
			}

			return i, ErrNoMatch
		}

		nextPos int
		err     error
	)

	log.Debug("asserting messages received", zap.Int("nb message received", ps.nbMessageReceived))

	nextPos, err = analyzeMessages(0, ps.nbMessageReceived)
	if err == nil || err != ErrNoMatch {
		log.Debug("exiting search loop", zap.Error(err), zap.Int("message position", nextPos))

		return err
	}

	for time.Now().Before(outdatedCheck) {
		time.Sleep(1 * time.Second)

		log.Debug("asserting messages received", zap.Int("nb message received", ps.nbMessageReceived), zap.Int("from message", nextPos))

		if nextPos < ps.nbMessageReceived {
			nextPos, err = analyzeMessages(nextPos, ps.nbMessageReceived)
			if err == nil || err != ErrNoMatch {
				log.Debug("exiting search loop", zap.Error(err), zap.Int("message position", nextPos))

				return err
			}
		}
	}

	return ErrNoMatch
}
