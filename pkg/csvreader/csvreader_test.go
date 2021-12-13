package csvreader_test

import (
	"fmt"
	"strings"
	"sync"
	"testing"

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
	r, err := r.ProceedQuery(query, q, row)
	require.Nil(t, err)
	fmt.Println(r.Table)
	require.Equal(t, "5", r.Table["a"][0])
}

func TestProceedConcurrentData(t *testing.T) {
	r := cs.NewData()
	query := "SELECT * FROM tablename where a > 4 and not c < 5"
	headers := []string{"a,c", "a,c"}
	rows := [][]string{{"5", "8"}, {"6", "10"}}
	var wg sync.WaitGroup
	q := sqlparser.NewQuery()
	wg.Add(2)
	for i, header := range headers {
		r.ReadHeaders(strings.NewReader(header))
		go func(row []string) {
			r.ProceedQuery(query, q, row)
			wg.Done()
		}(rows[i])
	}
	wg.Wait()
	assert.Equal(t, 2, len(r.Table["a"]))
}

// TODO: More tests
