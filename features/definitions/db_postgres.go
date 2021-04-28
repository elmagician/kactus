package definitions

import (
	"github.com/elmagician/godog"

	"github.com/elmagician/kactus/features/interfaces/database"
)

const pgKey = "[a-zA-Z_0-9 -]+"

// InstallPostgres adds postgres steps to expose basic pg features to known steps.
//
// Provided steps:
// - (?:I )?expect(?:ing)? data to be in (pg.[^:]+): => check data exists in database
//   I expect data to be in pg.example:
//      | table | id | email          | name      |
//      | users | 0  | test@gmail.com | something |
//      | users | 1  | lol@gmail.com  | lilacou   |
func InstallPostgres(s *godog.ScenarioContext, db *database.Postgres) {
	s.Step("^(?:I )?expect(?:ing)? data to be in (pg\\."+pgKey+"):$", db.AssertData)
	s.Step("^(?:I )?do not expect data to be in (pg\\."+pgKey+"):$", db.AssertDataDoesNotExists)
	s.Step("^(?:I )?pick(?:ing)? row ("+pgKey+") from (pg\\."+pgKey+")\\.("+pgKey+") as ([^:]+):$", db.StoreField)
}
