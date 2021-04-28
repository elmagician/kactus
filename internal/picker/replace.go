package picker

import (
	"github.com/cucumber/messages-go/v10"
	"go.uber.org/zap"
)

// NewReplacer provides function to replace picked variables in
// messages.Pickle_PickleStep Text and Arguments.
func NewReplacer(store *Store) func(step *messages.Pickle_PickleStep) {
	return func(step *messages.Pickle_PickleStep) {
		log.Debug("Replacing value in step.")
		log.Debug("Step text: " + step.Text)

		step.Text = store.InjectAll(step.Text)

		log.Debug("Updated step text: " + step.Text)

		args := step.Argument
		if args != nil {
			log.Debug("Step as arguments.")

			switch argTyped := args.Message.(type) {
			case *messages.PickleStepArgument_DataTable:
				log.Debug(" Replacing values in data table.")

				arg := argTyped.DataTable
				for i := 0; i < len(arg.Rows); i++ {
					for n, cell := range arg.Rows[i].Cells {
						log.Debug(
							"Replacing cell value: "+cell.Value,
							zap.Int("rows", i), zap.Int("column", n),
						)

						arg.Rows[i].Cells[n].Value = store.InjectAll(cell.Value)

						log.Debug(
							"Updated cell value: "+arg.Rows[i].Cells[n].Value,
							zap.Int("rows", i), zap.Int("column", n),
						)
					}
				}
			case *messages.PickleStepArgument_DocString:
				log.Debug("Replacing values in doc string.\n" + argTyped.DocString.GetContent())

				arg := argTyped.DocString
				arg.Content = store.InjectAll(arg.GetContent())

				log.Debug("Updated content.\n" + arg.Content)
			default:
				log.Error("unmanaged message type", zap.Reflect("arg", argTyped))
			}
		}
	}
}
