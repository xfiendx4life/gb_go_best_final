package csvreader_test

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"

	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xfiendx4life/gb_go_best_final/internal/logger"
	cs "github.com/xfiendx4life/gb_go_best_final/pkg/csvreader"
	"github.com/xfiendx4life/gb_go_best_final/pkg/sqlparser"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func newLogger() *zap.SugaredLogger {
	level := zapcore.DebugLevel
	return logger.InitLogger(&level, "")
}

func TestReadLine(t *testing.T) {
	r := cs.NewData()
	rdr := strings.NewReader("iso_code,continent,location,date,total_cases")
	row, err := r.ReadRow(rdr, ',')
	assert.Nil(t, err)
	assert.Equal(t, "iso_code", row[0])
}

func TestReadHeaders(t *testing.T) {
	r := cs.NewData()
	rdr := strings.NewReader("iso_code,continent,location,date,total_cases")
	_, err := r.ReadHeaders(rdr)
	assert.Nil(t, err)
	assert.Equal(t, "iso_code", r.GetHeaders()[0])
}

func TestProceedData(t *testing.T) {
	r := cs.NewData()
	query := []string{"a", "4", ">", "c", "5", "<", "not", "and"}
	r.ReadHeaders(strings.NewReader("a,c"))
	row := []string{"5", "8"}
	q := sqlparser.NewQuery()
	ctx := context.Background()
	r1, err := r.ProceedQuery(ctx, q, row, query, newLogger())
	require.Nil(t, err)
	require.Equal(t, "5", r1.GetTable()["a"][0])
}

func TestProceedConcurrentData(t *testing.T) {
	r := cs.NewData()
	query := []string{"a", "4", ">", "c", "5", "<", "not", "and"}
	headers := []string{"a,c", "a,c"}
	rows := [][]string{{"5", "8"}, {"6", "10"}}
	var wg sync.WaitGroup
	ctx := context.Background()
	q := sqlparser.NewQuery()
	wg.Add(2)
	for i, header := range headers {
		r.ReadHeaders(strings.NewReader(header))
		go func(row []string) {
			r.ProceedQuery(ctx, q, row, query, newLogger())
			wg.Done()
		}(rows[i])
	}
	wg.Wait()
	assert.Equal(t, 2, len(r.GetTable()["a"]))
}

func TestProceedFullTable(t *testing.T) {
	source := `a,c
5,8
5,6
`
	query := "SELECT * FROM tablename where a > 4 and not c < 5"
	r := cs.NewData()
	ctx := context.Background()
	resChan := make(chan cs.Table)
	errChan := make(chan error)
	go r.ProceedFullTable(ctx, strings.NewReader(source), query, newLogger(), resChan, errChan)
	select {
	case <-ctx.Done():
		assert.Nil(t, 0)
	case err := <-errChan:
		assert.Nil(t, err)
	case tab := <-resChan:
		assert.Equal(t, 2, len(tab.GetTable()["a"]))
		assert.Equal(t, []string{"5", "5"}, tab.GetTable()["a"])
	}

}

func TestProceedFullTableOneOfTwo(t *testing.T) {
	source := `a,c
5,8
5,4
`
	query := "SELECT a FROM tablename where a >= 4 and not c < 5"
	r := cs.NewData()
	ctx := context.Background()
	resChan := make(chan cs.Table)
	errChan := make(chan error)
	go r.ProceedFullTable(ctx, strings.NewReader(source), query, newLogger(), resChan, errChan)
	select {
	case <-ctx.Done():
		assert.Nil(t, 0)
	case err := <-errChan:
		assert.Nil(t, err)
	case tab := <-resChan:
		assert.Equal(t, 1, len(tab.GetTable()["a"]))
		assert.Nil(t, tab.GetTable()["c"])
	}

}

func TestProceedFullTableWithContext(t *testing.T) {
	source := `a,c
5,8
5,4
`
	query := "SELECT * FROM tablename where a >= 4 and not c < 5"
	r := cs.NewData()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*1))
	time.Sleep(1 * time.Second)
	resChan := make(chan cs.Table)
	errChan := make(chan error)
	go r.ProceedFullTable(ctx, strings.NewReader(source), query, newLogger(), resChan, errChan)
	assert.NotNil(t, <-ctx.Done())
	cancel()
}

// TODO: More tests
func TestProcessQueryError(t *testing.T) {
	q := []string{"some", "bull", "shit"}
	_, err := cs.NewData().ProceedQuery(context.Background(), sqlparser.NewQuery(), []string{"1", "2", "3"}, q, newLogger())
	assert.NotNil(t, err)
}

func TestGetNotAllCols(t *testing.T) {
	r := cs.NewData()
	query := []string{"a", "4", ">", "c", "5", "<", "not", "and"}
	headers := []string{"a,c", "a,c"}
	rows := [][]string{{"5", "8"}, {"6", "2"}}
	var wg sync.WaitGroup
	ctx := context.Background()
	q := sqlparser.NewQuery()
	wg.Add(2)
	for i, header := range headers {
		r.ReadHeaders(strings.NewReader(header))
		go func(row []string) {
			r.ProceedQuery(ctx, q, row, query, newLogger())
			wg.Done()
		}(rows[i])
	}
	wg.Wait()
	fmt.Println(r.GetTable())
	assert.Equal(t, "5", r.GetTable()["a"][0])
	assert.Equal(t, "8", r.GetTable()["c"][0])
}

func TestProceedFullTableONoWhere(t *testing.T) {
	source := `a,c
5,8
5,4
`
	query := "SELECT * FROM tablename"
	r := cs.NewData()
	ctx := context.Background()
	resChan := make(chan cs.Table)
	errChan := make(chan error)
	go r.ProceedFullTable(ctx, strings.NewReader(source), query, newLogger(), resChan, errChan)
	select {
	case <-ctx.Done():
		assert.Nil(t, 0)
	case err := <-errChan:
		assert.Nil(t, err)
	case tab := <-resChan:
		assert.Equal(t, 2, len(tab.GetTable()["a"]))
		assert.Equal(t, 2, len(tab.GetTable()["c"]))
	}

}
