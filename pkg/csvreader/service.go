package csvreader

import (
	"encoding/csv"
	"fmt"
	"io"

	"github.com/xfiendx4life/gb_go_best_final/pkg/sqlparser"
)

func NewData() *Data {
	return &Data{
		headers:     make([]string, 0),
		currentLine: make([]string, 0),
		Table:       make(map[string][]string),
	}
}

func (r *Data) GetHeaders() []string {
	return r.headers
}

func (r *Data) ReadRow(source io.Reader) (row []string, err error) {
	r.reader = csv.NewReader(source)
	row, err = r.reader.Read()
	if err != nil {
		return nil, err
	}
	return row, nil
}

func (r *Data) ReadHeaders(source io.Reader) (headers []string, err error) {
	headers, err = r.ReadRow(source)
	if err != nil {
		return nil, fmt.Errorf("can't rad headers %s", err)
	}
	r.headers = headers
	return headers, nil
}

func (r *Data) composeRow(headers []string, row []string) (composedRow map[string]string) {
	composedRow = make(map[string]string)
	for i, header := range headers {
		composedRow[header] = row[i]
	}
	return composedRow
}

func (r *Data) ProceedQuery(query string, headers []string, row []string) (data *Data, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	q := sqlparser.NewQuery()
	postfix, err := q.ParseToPostfix(query)
	if err != nil {
		return nil, fmt.Errorf("cant' proceed query %s", err)
	}
	var isValid bool
	composed := r.composeRow(headers, row)
	isValid, err = q.SelectFromRow(postfix, composed)
	if err != nil {
		return nil, fmt.Errorf("can't proceed query %s", err)
	}
	if isValid {
		for k, v := range composed {
			r.Table[k] = append(r.Table[k], v)
		}
	}
	data = r
	return data, nil
}
