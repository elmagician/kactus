package definitions

import (
	"github.com/elmagician/godog"

	"github.com/elmagician/kactus/features/interfaces/pubsub"
)

// InstallGooglePubsub adds google pubsub steps to expose basic google pubsub
// features to known steps.
//
// Provided steps:
// - (?:I )?expect(?:ing)? message to be received by (gcp\.[^:]+) within ([0-9]+) seconds::
//          => check if a message matching provided data was received within X seconds
//   I expect message to be received by gcp.main within 1 seconds:
//      | field  | matcher  | value |
//      | user   | contains | legit |
//      | reader | =        | me    |
//
// - (?:I )?expect(?:ing)? message to be received by (gcp.[^:]+) within ([0-9]+) seconds having metadata:
//          => check if a message matching provided metadata was received
//   I expect message to be received by gcp.main within 1000 seconds:
//      | field  | matcher  | value |
//      | user   | contains | legit |
//      | reader | =        | me    |
//
// - (?:I )?send(?:ing)? message to (gcp\.[^:]+) in (topic\..+):
//          => send a given message to gcp instance topic. Message has to contain
//              a data key witch is a free json representing main message and
//              a metadata key witch is a string: string object representing metadata
//   I expect message to be sent by gcp.main in topic.testTopic:
//      """
//      {
//          "data": {
//              "foo": "bar",
//              "my_value": 1234
//          },
//          "metadata":{
//              "foo": "bar"
//          }
//      }
//      """
func InstallGooglePubsub(s *godog.ScenarioContext, ps *pubsub.Google) {
	s.Step("^(?:I )?expect(?:ing)? message to be received by (gcp\\.[^:]+) within ([0-9]+) seconds:$", ps.HasMessage)
	s.Step(
		"^(?:I )?expect(?:ing)? message to be received by (gcp\\.[^:]+) within ([0-9]+) seconds having metadata:$",
		ps.HasMessageWithMetadata,
	)
	s.Step("^(?:I )?send(?:ing)? message to (gcp\\.[^:]+) in (topic\\..+):$", ps.SendMessage)
}
