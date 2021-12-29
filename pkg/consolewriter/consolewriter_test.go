package consolewriter_test

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xfiendx4life/gb_go_best_final/pkg/consolewriter"
	"github.com/xfiendx4life/gb_go_best_final/pkg/csvreader"
	"github.com/xfiendx4life/gb_go_best_final/pkg/sqlparser"
	"go.uber.org/zap"
)

type tableStub struct {
	table map[string][]string
}

func (t *tableStub) ReadHeaders(source io.Reader) (headers []string, err error) {
	return []string{}, nil
}

func (t *tableStub) ReadRow(source io.Reader, sep rune) (row []string, err error) {
	return
}

func (t *tableStub) ProceedQuery(ctx context.Context, query sqlparser.Querier, row []string, postfix []string, z *zap.SugaredLogger) (csvreader.Table, error) {
	return t, nil
}

func (t *tableStub) ProceedFullTable(ctx context.Context, source io.Reader, rawQuery string, z *zap.SugaredLogger, resChan chan csvreader.Table, errChan chan error) {

}
func (t *tableStub) GetTable() map[string][]string {
	return t.table
}

func (t *tableStub) GetHeaders() []string {
	res := make([]string, 0)
	for k := range t.table {
		res = append(res, k)
	}
	return res
}

func TestWriteOneCol(t *testing.T) {
	ts := tableStub{
		table: map[string][]string{"header 1": {"val1", "val2"}},
	}
	wr := new(strings.Builder)
	consolewriter.NewConsoleWriter().Write(wr, &ts)
	assert.Equal(t, "header 1 \nval1     \nval2     \n", wr.String())

}

func TestWrite(t *testing.T) {
	ts := tableStub{
		table: map[string][]string{"header 1": {"val1", "val2"}, "header two": {"valone", "valtwo"}},
	}
	wr := new(strings.Builder)
	consolewriter.NewConsoleWriter().Write(wr, &ts)
	assert.Equal(t, "header 1   header two \nval1       valone     \nval2       valtwo     \n", wr.String())

}
