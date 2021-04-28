package picker_test

import (
	"testing"

	"github.com/cucumber/messages-go/v10"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/elmagician/godog"

	"github.com/elmagician/kactus/internal/matchers"
	"github.com/elmagician/kactus/internal/picker"
	"github.com/elmagician/kactus/internal/types"
)

func init() {
	picker.NoLog()
	matchers.NoLog()
	types.NoLog()
}

func TestUnit_Store_NewReplacer(t *testing.T) {
	Convey("Given I wish to setup an auto replacer for Godog before step, ", t, func() {
		store := picker.NewStore()
		store.Pick("test", 124, picker.PersistentValue)
		store.Pick("pi", 3.14, picker.PersistentValue)
		store.Pick("testName", "storeInject", picker.DisposableValue)

		replacer := picker.NewReplacer(store)

		text := "Some {{testName}} should be tested"
		replacedText := "Some storeInject should be tested"

		Convey("It should replace provided value in step text", func() {
			msg := &messages.Pickle_PickleStep{Text: text}
			replacer(msg)
			So(msg.Text, ShouldEqual, replacedText)
		})

		Convey("It should replace provided value in step DocString arguments", func() {
			doc := &godog.DocString{Content: "This is {{pi}} for {{test}}"}
			msg := &messages.Pickle_PickleStep{
				Text: "Some {{testName}} should be tested",
				Argument: &messages.PickleStepArgument{
					Message: &messages.PickleStepArgument_DocString{DocString: doc},
				},
			}
			replacer(msg)
			So(msg.Text, ShouldEqual, replacedText)
			So(doc.Content, ShouldEqual, "This is 3.14 for 124")
		})

		Convey("It should replace provided value in step DataTable arguments", func() {
			table := &messages.PickleStepArgument_DataTable{
				DataTable: &messages.PickleStepArgument_PickleTable{
					Rows: []*messages.PickleStepArgument_PickleTable_PickleTableRow{
						{
							Cells: []*messages.PickleStepArgument_PickleTable_PickleTableRow_PickleTableCell{
								{Value: "{{pi}}"},
								{Value: "{{test}}"},
							},
						}, {
							Cells: []*messages.PickleStepArgument_PickleTable_PickleTableRow_PickleTableCell{
								{Value: "{{test}}"},
								{Value: "test"},
							},
						},
					},
				},
			}

			msg := &messages.Pickle_PickleStep{
				Text:     "Some {{testName}} should be tested",
				Argument: &messages.PickleStepArgument{Message: table},
			}
			replacer(msg)
			So(msg.Text, ShouldEqual, replacedText)
			So(table.DataTable.Rows[0].Cells[0].Value, ShouldEqual, "3.14")
			So(table.DataTable.Rows[0].Cells[1].Value, ShouldEqual, "124")
			So(table.DataTable.Rows[1].Cells[0].Value, ShouldEqual, "124")
			So(table.DataTable.Rows[1].Cells[1].Value, ShouldEqual, "test")
		})

	})
}
