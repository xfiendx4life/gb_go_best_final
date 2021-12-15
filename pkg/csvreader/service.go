package csvreader

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"sync"

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

func (r *Data) GetTable() map[string][]string {
	return r.Table
}

func (r *Data) ReadRow(source io.Reader) (row []string, err error) {
	if r.reader == nil {
		r.reader = csv.NewReader(source)
	}
	// TODO: Read comma from config
	r.reader.Comma = ','
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

func (r *Data) ProceedQuery(ctx context.Context, query string, q sqlparser.Querier, row []string) (data Table, err error) {
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("done with context")
	default:
		r.mu.Lock()
		defer r.mu.Unlock()
		postfix, err := q.ParseToPostfix(query)
		if err != nil {
			return nil, fmt.Errorf("cant' proceed query %s", err)
		}
		var isValid bool
		composed := r.composeRow(r.headers, row)
		isValid, err = q.SelectFromRow(ctx, postfix, composed)
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
}

func (r *Data) ProceedFullTable(ctx context.Context, source io.Reader, rawQuery string) (table Table, err error) {
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("done with context")
	default:
		_, err = r.ReadHeaders(source)
		if err != nil {
			return nil, fmt.Errorf("can't read headers %s", err)
		}
		q := sqlparser.NewQuery()
		var wg sync.WaitGroup
		for {
			row, err := r.ReadRow(source)
			if err == io.EOF {
				break
			} else if err != nil {
				return nil, fmt.Errorf("can't parse table %s", err)
			}
			wg.Add(1)
			go func(row []string) {
				_, err = r.ProceedQuery(ctx, rawQuery, q, row)
				if err != nil {
					log.Printf("error while parsing row %s", err)
				}
				wg.Done()
			}(row)
		}
		wg.Wait()
		return r, nil
	}
}
