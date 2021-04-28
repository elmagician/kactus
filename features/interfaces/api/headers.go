package api

import (
	"github.com/elmagician/godog"
)

// SetRequestHeaders replaces current headers with list of provided ones.
// It ALWAYS return nil.
func (cli *Client) SetRequestHeaders(headers *godog.Table) {
	cli.request.ResetHeader()
	cli.AddHeaders(headers)
}

// AddHeaders adds headers to known list.
// It ALWAYS return nil.
func (cli *Client) AddHeaders(headers *godog.Table) {
	var (
		settedHeaders string
		key, value    string
	)

	for _, row := range headers.Rows {
		key = row.Cells[0].Value
		value = row.Cells[0].Value
		settedHeaders += key + ", "

		cli.SetHeader(key, value) // nolint: errcheck
	}
}

// SetHeader adds a single key,value couple to request headers.
// It ALWAYS return nil.
func (cli *Client) SetHeader(key, values string) {
	if cli.request.Empty() {
		cli.InitRequest(true)
	}

	cli.request.AddHeader(key, values)
}
