package postgres

import (
	"io/ioutil"

	"go.uber.org/zap"

	"github.com/elmagician/kactus/internal/picker"
)

// SeedFromFile executes an SQL file content
// into known postgres database.
func (db *DB) SeedFromFile(filePath string, store *picker.Store) (err error) {
	var (
		ok      bool
		strFile string
	)

	log.Debug("Seeding file", zap.String("filePath", filePath))

	file, fileErr := ioutil.ReadFile(filePath)
	if fileErr != nil {
		return fileErr
	}

	strFile = string(file)

	if store != nil {
		log.Debug("injecting picked variables")

		strFile = store.InjectAll(strFile)
	}

	log.Debug("starting seed transaction", zap.String("file", strFile))

	tx, beginErr := db.db.Begin()
	if beginErr != nil {
		return beginErr
	}

	defer func() {
		if !ok {
			log.Debug("rolling back seed transaction")

			if errRb := tx.Rollback(); err == nil && errRb != nil {
				err = errRb
			}
		}
	}()

	if _, err := tx.Exec(strFile); err != nil {
		return err
	}

	ok = true

	return tx.Commit()
}
