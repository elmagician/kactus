package postgres

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/elmagician/godog"
	"go.uber.org/zap"

	"github.com/elmagician/kactus/internal/picker"
)

type query struct {
	query     string
	arguments []interface{}
}

// DB provides a structure to manage postgres databases.
type DB struct {
	name string
	db   *sql.DB
}

// Database initializes a DB instance for postgres.
// It will persist instance using provided name to
// be injectable through {{pg.name}} in gherkin steps.
// Providing an empty picker does not impact initialization.
// It will just not picked the instance.
// You can always call DB.Persist to save database instance
// in a picker store.
func Database(name string, db *sql.DB, store *picker.Store) *DB {
	database := &DB{name: name, db: db}

	if store != nil {
		database.Persist(store)
	}

	return database
}

// Persist persists database instance through picker instance using pg.name key.
func (db DB) Persist(store *picker.Store) {
	store.Pick(
		"pg."+db.name,
		picker.InstanceItem{Kind: picker.Postgres, Instance: db},
		picker.InstanceValue,
	)
}

func (db DB) PickField(store *picker.Store, table, field, key string, identifiedBy *godog.Table) (err error) {
	var (
		getValue     interface{}
		whereClauses wheres

		query = "SELECT " + field + " FROM " + table
	)

	keys := identifiedBy.Rows[0].Cells

	log.Debug("parsing data for assertion query", zap.Reflect("keys", keys))

	for rowNumber, row := range identifiedBy.Rows[1:] {
		whereClause := where{Operator: defaultCheckOperator}

		log.Debug("analyzing row", zap.Int("row number", rowNumber), zap.Reflect("data", row.Cells))

		for rangeID, cell := range row.Cells {
			switch keys[rangeID].Value {
			case "field":
				whereClause.Column = cell.Value
			case "operator", "op":
				whereClause.Operator = cell.Value
			case "value", "val":
				whereClause.ExpectedResult = cell.Value
			}
		}

		log.Debug("adding condition", zap.Reflect("where", whereClause))

		whereClauses = append(whereClauses, whereClause)
	}

	q := prepareQuery(query, whereClauses)

	log.Debug("tyring to retrieve data", zap.String("query", q.query), zap.Reflect("args", q.arguments))

	rows, err := db.db.Query(q.query, q.arguments...)
	if err != nil {
		return err
	}

	defer func() {
		if rowErr := rows.Err(); rowErr != nil && err == nil {
			err = rowErr
		}

		if rowErr := rows.Close(); rowErr != nil && err == nil {
			err = rows.Close()
		}
	}()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return err
		}

		return ErrUnmatched
	}

	if err := rows.Scan(&getValue); err != nil {
		return err
	}

	if byteVal, ok := getValue.([]byte); ok {
		getValue = string(byteVal)
	}

	store.Pick(key, getValue, picker.DisposableValue)

	return nil
}

func prepareQuery(baseQuery string, queryElements wheres) query {
	var (
		queryString = baseQuery
		args        []interface{}
		hasWhere    bool

		wheres   string
		argIndex = 0
	)

	for _, element := range queryElements {
		argIndex++
		wheres += fmt.Sprintf("%s %s $%v AND ", element.Column, element.Operator, argIndex)
		args = append(args, element.ExpectedResult)

		hasWhere = true
	}

	wheres = strings.TrimSuffix(wheres, " AND ")

	if hasWhere {
		queryString += " WHERE " + wheres
	}

	return query{query: queryString, arguments: args}
}
