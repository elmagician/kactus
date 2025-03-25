package picker

import (
	"github.com/cucumber/godog"
	"go.uber.org/zap"
)

// NewReplacer provides function to replace picked variables in
// messages.Pickle_PickleStep Text and Arguments.
func NewReplacer(store *Store) func(step *godog.Step) {
	return func(step *godog.Step) {
		log.Debug("Replacing value in step.")
		log.Debug("Step text: " + step.Text)

		step.Text = store.InjectAll(step.Text)

		log.Debug("Updated step text: " + step.Text)

		args := step.Argument
		if args != nil {
			log.Debug("Step as arguments.")
			if args.DataTable != nil {

				log.Debug(" Replacing values in data table.")

				arg := args.DataTable
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
			}

			if args.DocString != nil {
				log.Debug("Replacing values in doc string.\n" + args.DocString.Content)

				arg := args.DocString
				arg.Content = store.InjectAll(args.DocString.Content)

				log.Debug("Updated content.\n" + arg.Content)
			}

		}
	}
}
