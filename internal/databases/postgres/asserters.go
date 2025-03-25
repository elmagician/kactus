package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/cucumber/godog"
	"go.uber.org/zap"
)

const (
	defaultCheckOperator = "="
	defaultTypeMatcher   = ""
	separatorToken       = "((" // separator used to parse String to Where element
)

// matches sql where condition as string. They are express like:
// DATA ((::TYPE))((OPERATOR))
// Space is mandatory between DATA and OPTIONS due to Golang not doing backwards regex for perf.
// cf: https://regex101.com/r/lIe4Kr/1/ for examples
var sqlModifierToken = regexp.MustCompile(`(?P<data>.+) (?:\(\()?(?P<typeMatcher>:[^()]+)?(?:\)\)\(\()?(?P<operator>[^()]*)(?:\)\))?`)

// ErrUnmatched indicates that some expected assertions are not met.
var ErrUnmatched = errors.New("some conditions are not met")

type (
	wheres []where

	where struct {
		Column         string
		ExpectedResult string
		TypeMatcher    string
		Operator       string
	}

	assertQuery struct {
		query             string
		arguments         []interface{}
		nbMatchesExpected int
	}
)

// AssertData asserts provided data from godog table
// matches provided conditions in database.
func (db DB) AssertData(data *godog.Table) error {
	var queries []assertQuery

	for table, elements := range parseDataToQuery(data) {
		queries = append(queries, prepareAssertionQuery(table, elements))
	}

	return runAssertQuery(db.db, queries)
}

func parseDataToQuery(data *godog.Table) map[string][]wheres {
	var (
		table  string
		tables = make(map[string][]wheres)
	)

	keys := data.Rows[0].Cells

	log.Debug("parsing data for assertion query", zap.Reflect("keys", keys))

	for rowNumber, row := range data.Rows[1:] {
		var (
			checker    where
			oneElement []where
		)

		log.Debug("analyzing row", zap.Int("row number", rowNumber), zap.Reflect("data", row.Cells))

		for rangeID, cell := range row.Cells {
			if keys[rangeID].Value == "table" {
				table = cell.Value

				if _, exits := tables[table]; !exits {
					tables[table] = make([]wheres, 0) // equivalent to []wheres{}
				}

				continue
			}

			checker.Column = keys[rangeID].Value

			if !strings.Contains(cell.Value, separatorToken) {
				checker.Operator = defaultCheckOperator
				checker.TypeMatcher = defaultTypeMatcher
				checker.ExpectedResult = cell.Value

				oneElement = append(oneElement, checker)

				continue
			}

			detailed := sqlModifierToken.FindStringSubmatch(cell.Value)

			// As group are named, they will always be provided in the same order.
			// When using FindStringSubmatch, first return value is a full match
			// it is ignored for us.
			// Main block is the first subgroup
			// Type 2d
			// Operator the 3d
			checker.ExpectedResult = detailed[1]
			checker.TypeMatcher = detailed[2]
			checker.Operator = detailed[3]

			if checker.Operator == "" {
				checker.Operator = defaultCheckOperator
			}

			if checker.TypeMatcher == "" {
				checker.TypeMatcher = defaultTypeMatcher
			}

			checker.Operator = strings.Replace(checker.Operator, "_", " ", -1)
			checker.TypeMatcher = strings.Replace(checker.TypeMatcher, "_", " ", -1)

			oneElement = append(oneElement, checker)
		}

		log.Debug(
			"adding data for query preparation",
			zap.String("table", table), zap.Reflect("checks", oneElement),
		)

		tables[table] = append(tables[table], oneElement)
		table = ""
	}

	return tables
}

func prepareAssertionQuery(table string, queryElements []wheres) assertQuery {
	var (
		queryString string
		args        []interface{}

		wheres   string
		argIndex = 0

		nbMatches = 0
	)

	for _, orBlockElement := range queryElements {
		wheres += " ("

		for _, element := range orBlockElement {
			if element.ExpectedResult == "NULL" {
				wheres += fmt.Sprintf("%s %s NULL AND ", element.Column, element.Operator)
			} else {
				argIndex++
				wheres += fmt.Sprintf("%s%s %s $%v%s AND ", element.Column, element.TypeMatcher, element.Operator, argIndex, element.TypeMatcher)
				args = append(args, element.ExpectedResult)
			}
		}

		wheres = strings.TrimSuffix(wheres, " AND ")
		wheres = strings.TrimSuffix(wheres, " (")
		wheres += ") OR"

		if argIndex > 0 {
			nbMatches++
		}
	}

	wheres = strings.TrimSuffix(wheres, " OR")

	queryString = "SELECT COUNT(*) FROM " + table + " WHERE " + wheres

	return assertQuery{query: queryString, arguments: args, nbMatchesExpected: nbMatches}
}

// Ignore sqlclosecheck as it wish to use a defer and defer should not
// be used in a loop.
// nolint: sqlclosecheck
func runAssertQuery(db *sql.DB, queries []assertQuery) error {
	for _, query := range queries {
		var nbMatch int

		log.Debug(
			"running query",
			zap.String("query", query.query), zap.Reflect("arguments", query.arguments),
			zap.Int("expected match", query.nbMatchesExpected),
		)

		rows, err := db.Query(query.query, query.arguments...)
		if err != nil {
			return err
		}

		if !rows.Next() {
			_ = rows.Close()
			return rows.Err()
		}

		if err := rows.Scan(&nbMatch); err != nil {
			_ = rows.Close() // close row error is not relevant here for test as main issue is on scan
			return err
		}

		if nbMatch != query.nbMatchesExpected {
			log.Debug(
				"incorrect match",
				zap.Int("matched", nbMatch), zap.Int("expected", query.nbMatchesExpected),
			)

			_ = rows.Close() // close row error is not relevant here for test as main issue is that not all conditions where met

			return fmt.Errorf("%w, query: %s, arguments: %v", ErrUnmatched, query.query, query.arguments)
		}

		if err := rows.Close(); err != nil { // nolint: sqlclosecheck
			return err
		}
	}

	return nil
}
