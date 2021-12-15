package csvreader_test

import (
	"context"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	cs "github.com/xfiendx4life/gb_go_best_final/pkg/csvreader"
	"github.com/xfiendx4life/gb_go_best_final/pkg/sqlparser"
)

func TestReadLine(t *testing.T) {
	r := cs.NewData()
	rdr := strings.NewReader("iso_code,continent,location,date,total_cases")
	row, err := r.ReadRow(rdr)
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
	query := "SELECT * FROM tablename where a > 4 and not c < 5"
	r.ReadHeaders(strings.NewReader("a,c"))
	row := []string{"5", "8"}
	q := sqlparser.NewQuery()
	ctx := context.Background()
	r1, err := r.ProceedQuery(ctx, query, q, row)
	require.Nil(t, err)
	require.Equal(t, "5", r1.GetTable()["a"][0])
}

func TestProceedConcurrentData(t *testing.T) {
	r := cs.NewData()
	query := "SELECT * FROM tablename where a > 4 and not c < 5"
	headers := []string{"a,c", "a,c"}
	rows := [][]string{{"5", "8"}, {"6", "10"}}
	var wg sync.WaitGroup
	ctx := context.Background()
	q := sqlparser.NewQuery()
	wg.Add(2)
	for i, header := range headers {
		r.ReadHeaders(strings.NewReader(header))
		go func(row []string) {
			r.ProceedQuery(ctx, query, q, row)
			wg.Done()
		}(rows[i])
	}
	wg.Wait()
	assert.Equal(t, 2, len(r.Table["a"]))
}

func TestProceedFullTable(t *testing.T) {
	source := `a,c
5,8
5,6
`
	query := "SELECT * FROM tablename where a > 4 and not c < 5"
	r := cs.NewData()
	ctx := context.Background()
	tab, err := r.ProceedFullTable(ctx, strings.NewReader(source), query)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(tab.GetTable()["a"]))
	assert.Equal(t, []string{"5", "5"}, tab.GetTable()["a"])
}

// TODO: More tests
func TestProceedFullTableOneOfTwo(t *testing.T) {
	source := `a,c
5,8
5,4
`
	query := "SELECT * FROM tablename where a >= 4 and not c < 5"
	r := cs.NewData()
	ctx := context.Background()
	tab, err := r.ProceedFullTable(ctx, strings.NewReader(source), query)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(tab.GetTable()["a"]))
	assert.Equal(t, []string{"8"}, tab.GetTable()["c"])
}

func TestProceedFullTableWithContext(t *testing.T) {
	source := `a,c
5,8
5,4
`
	query := "SELECT * FROM tablename where a >= 4 and not c < 5"
	r := cs.NewData()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*1))
	go func() {
		cancel()
	}()
	time.Sleep(1 * time.Second)
	tab, err := r.ProceedFullTable(ctx, strings.NewReader(source), query)
	assert.NotNil(t, err)
	assert.Nil(t, tab)
}
