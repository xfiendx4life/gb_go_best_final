package csvreader

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/xfiendx4life/gb_go_best_final/pkg/sqlparser"
	"go.uber.org/zap"
)

func NewData() Table {
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

func (r *Data) ReadRow(source io.Reader, sep rune) (row []string, err error) {
	if r.reader == nil {
		r.reader = csv.NewReader(source)
	}
	r.reader.Comma = sep
	row, err = r.reader.Read()
	if err != nil {
		return nil, err
	}
	return row, nil
}

func (r *Data) ReadHeaders(source io.Reader) (headers []string, err error) {
	headers, err = r.ReadRow(source, ',')
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

func checkSlice(target string, data []string) bool {
	if len(data) == 0 || data[0] == "*" {
		return true
	}
	for _, item := range data {
		if strings.HasPrefix(item, target) {
			return true
		}
	}
	return false
}

func (r *Data) ProceedQuery(ctx context.Context, q sqlparser.Querier, row []string, postfix []string, z *zap.SugaredLogger) (data Table, err error) {
	select {
	case <-ctx.Done():
		z.Errorf("ProceedQuery ended with context")
		return nil, fmt.Errorf("done with context")
	default:
		r.mu.Lock()
		defer r.mu.Unlock()
		var isValid bool
		composed := r.composeRow(r.headers, row)
		if len(postfix) != 0 {
			isValid, err = q.SelectFromRow(ctx, postfix, composed, z)
			if err != nil {
				z.Errorf("can't proceed query %s", err)
				return nil, fmt.Errorf("can't proceed query %s", err)
			}
		} else {
			isValid = true
		}
		if isValid {
			z.Debugf("row %v is valid for query", row)
			for k, v := range composed {
				if checkSlice(k, q.GetResultCols()) {
					r.Table[k] = append(r.Table[k], v)
				}
			}
		}
		data = r
		return data, nil
	}
}

func (r *Data) ProceedFullTable(ctx context.Context, source io.Reader, rawQuery string, z *zap.SugaredLogger, resChan chan Table, errChan chan error) {
	select {
	case <-ctx.Done():
		z.Errorf("ProceedTable ended with context")
		errChan <- fmt.Errorf("done with context")
		return
	default:
		_, err := r.ReadHeaders(source)
		if err != nil {
			errChan <- fmt.Errorf("can't read headers %s", err)
		}
		z.Debugf("headers of the table are %v", r.headers)
		q := sqlparser.NewQuery()
		postfix, err := q.ParseToPostfix(rawQuery)
		if err != nil {
			z.Errorf("can't parse to postfix %s", err)
			errChan <- fmt.Errorf("cant' proceed query %s", err)
		}
		z.Debugf("query parsed to postfix form %v", postfix)
		var wg sync.WaitGroup
		for {
			row, err := r.ReadRow(source, ',')
			if err == io.EOF {
				break
			} else if err != nil {
				z.Errorf("error while parsing table %s", err)
				errChan <- fmt.Errorf("can't parse table %s", err)
				return
			}
			wg.Add(1)
			go func(row []string) {
				_, err = r.ProceedQuery(ctx, q, row, postfix, z)
				if err != nil {
					z.Warnf("error while parsing row %s", err)
				}
				wg.Done()
			}(row)
		}
		wg.Wait()
		resChan <- r
	}
}
