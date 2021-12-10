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
		Result:      make(map[string]string),
	}
}

func (r *Data) GetHeaders() []string {
	return r.headers
}

func (r *Data) ReadLine(source io.Reader) (err error) {
	r.reader = csv.NewReader(source)
	r.headers, err = r.reader.Read()
	if err != nil {
		return err
	}
	return nil
}

func (r *Data) ProceedQuery(query string) (data *Data, err error) {
	q := sqlparser.NewQuery()
	postfix, err := q.ParseToPostfix(query)
	if err != nil {
		return nil, fmt.Errorf("cant' proceed query %s", err)
	}
	data.Result, err = q.SelectFromRow(postfix)
	if err != nil {
		return nil, fmt.Errorf("can't proceed query %s", err)
	}
	return data, nil
}
