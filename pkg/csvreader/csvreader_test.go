package csvreader_test

import (
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	cs "github.com/xfiendx4life/gb_go_best_final/pkg/csvreader"
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
	headers := []string{"a", "c"}
	row := []string{"5", "8"}
	r, err := r.ProceedQuery(query, headers, row)
	require.Nil(t, err)
	fmt.Println(r.Table)
	require.Equal(t, "5", r.Table["a"][0])
}

func TestProceedConcurrentData(t *testing.T) {
	r := cs.NewData()
	query := "SELECT * FROM tablename where a > 4 and not c < 5"
	headers := [][]string{{"a", "c"}, {"a", "c"}}
	rows := [][]string{{"5", "8"}, {"6", "10"}}
	var wg sync.WaitGroup
	wg.Add(2)
	for i, header := range headers {
		go func(header, row []string) {
			r.ProceedQuery(query, header, row)
			wg.Done()
		}(header, rows[i])
	}
	wg.Wait()
	assert.Equal(t, 2, len(r.Table["a"]))
}

// TODO: More tests
