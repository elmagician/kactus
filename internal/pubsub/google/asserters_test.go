package google_test

import (
	"context"
	"testing"
	"time"

	googlePubSub "cloud.google.com/go/pubsub"
	"cloud.google.com/go/pubsub/pstest"
	"github.com/cucumber/messages-go/v10"
	"github.com/elmagician/godog"
	"github.com/elmagician/kactus/internal/matchers"
	"github.com/elmagician/kactus/internal/picker"
	"github.com/elmagician/kactus/internal/pubsub/google"
	. "github.com/elmagician/kactus/internal/test"
	"github.com/elmagician/kactus/internal/types"
	. "github.com/smartystreets/goconvey/convey"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
)

func init() {
	google.NoLog()
	matchers.NoLog()
	types.NoLog()
}

func initTest(project string, opts ...pstest.ServerReactorOption) (*google.Client, *pstest.Server) {
	psTest := pstest.NewServer(opts...)

	conn, err := grpc.Dial(psTest.Addr, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	cli, err := googlePubSub.NewClient(context.Background(), project, option.WithGRPCConn(conn))
	if err != nil {
		panic(err)
	}

	subscriptionName := "foo"
	top, err := cli.CreateTopic(context.Background(), "someTopic")
	if err != nil {
		panic(err)
	}
	if _, err = cli.CreateSubscription(context.Background(), subscriptionName, googlePubSub.SubscriptionConfig{Topic: top}); err != nil {
		panic(err)
	}
	ps := google.Pubsub("foo", cli, picker.NewStore())
	go func() {
		if err := ps.ReceiveOn(subscriptionName, true); err != nil {
			panic(err)
		}
	}()

	return ps, psTest
}

func TestUnit_Client_AssertMessageReceived(t *testing.T) {
	Convey("When I try to assert message received", t, func() {
		ps, server := initTest("test")

		data := &godog.Table{
			Rows: []*messages.PickleStepArgument_PickleTable_PickleTableRow{
				{
					Cells: []*messages.PickleStepArgument_PickleTable_PickleTableRow_PickleTableCell{
						{Value: "field"},
						{Value: "matcher"},
						{Value: "value"},
					},
				}, {
					Cells: []*messages.PickleStepArgument_PickleTable_PickleTableRow_PickleTableCell{
						{Value: "foo"},
						{Value: "eq"},
						{Value: "bar"},
					},
				},
			},
		}
		Convey("should success if message received", func() {
			server.Publish("projects/test/topics/someTopic", []byte("{\"foo\":\"bar\"}"), nil)

			So(ps.AssertMessageReceived(data, time.Second), ShouldBeNil)
		})

		Convey("should fail", func() {
			Convey("if message not received", func() {
				So(ps.AssertMessageReceived(data, time.Second), ShouldBeError)
			})

			Convey("if matcher does not exist", func() {
				server.Publish("projects/test/topics/someTopic", []byte("{\"foo\":\"bar\"}"), map[string]string{"metadataKey": "metadataValue"})
				data.Rows = []*messages.PickleStepArgument_PickleTable_PickleTableRow{
					{
						Cells: []*messages.PickleStepArgument_PickleTable_PickleTableRow_PickleTableCell{
							{Value: "field"},
							{Value: "matcher"},
							{Value: "value"},
						},
					}, {
						Cells: []*messages.PickleStepArgument_PickleTable_PickleTableRow_PickleTableCell{
							{Value: "metadataKey"},
							{Value: "coucou"},
							{Value: "metadataValue"},
						},
					},
				}

				So(ps.AssertMessageReceived(data, time.Second), ShouldBeError)
			})
		})
	})
}

func TestUnit_Client_AssertMessageMetadata(t *testing.T) {
	Convey("When I try to assert message have meta data", t, func() {
		ps, server := initTest("test")

		data := &godog.Table{
			Rows: []*messages.PickleStepArgument_PickleTable_PickleTableRow{
				{
					Cells: []*messages.PickleStepArgument_PickleTable_PickleTableRow_PickleTableCell{
						{Value: "field"},
						{Value: "matcher"},
						{Value: "value"},
					},
				}, {
					Cells: []*messages.PickleStepArgument_PickleTable_PickleTableRow_PickleTableCell{
						{Value: "metadataKey"},
						{Value: "eq"},
						{Value: "metadataValue"},
					},
				},
			},
		}
		Convey("should success if message received", func() {
			server.Publish("projects/test/topics/someTopic", []byte("{\"foo\":\"bar\"}"), map[string]string{"metadataKey": "metadataValue"})

			So(ps.AssertMessageMetadata(data, time.Second), ShouldBeNil)
		})

		Convey("should fail", func() {
			Convey("if message not received", func() {
				So(ps.AssertMessageMetadata(data, time.Second), ShouldBeLikeError, google.ErrNoMatch)
			})

			Convey("if matcher does not exist", func() {
				server.Publish("projects/test/topics/someTopic", []byte("{\"foo\":\"bar\"}"), map[string]string{"metadataKey": "metadataValue"})
				data.Rows = []*messages.PickleStepArgument_PickleTable_PickleTableRow{
					{
						Cells: []*messages.PickleStepArgument_PickleTable_PickleTableRow_PickleTableCell{
							{Value: "field"},
							{Value: "matcher"},
							{Value: "value"},
						},
					}, {
						Cells: []*messages.PickleStepArgument_PickleTable_PickleTableRow_PickleTableCell{
							{Value: "metadataKey"},
							{Value: "coucou"},
							{Value: "metadataValue"},
						},
					},
				}

				So(ps.AssertMessageMetadata(data, time.Second), ShouldBeError)
			})
		})
	})
}
