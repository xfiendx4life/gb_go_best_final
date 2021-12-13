package csvreader

import (
	"encoding/csv"
	"io"
	"sync"

	"github.com/xfiendx4life/gb_go_best_final/pkg/sqlparser"
)

type Source interface {
	// Open file for reading
	Read(path string) (file []byte, err error)
}

type Table interface {
	// Reads headers from
	ReadHeaders(source io.Reader) (headers []string, err error)
	// Reads any line
	ReadRow(source io.Reader) (row []string, err error)
	// parses whole command sequence
	// returns the same structure ready to parse
	//another query
	ProceedQuery(rawQuery string, query sqlparser.Querier, row []string) (Table, error)
	// proceeds select query to the whole table concurrently
	ProceedFullTable(source io.Reader, rawQuery string) (table Table, err error)
	GetTable() map[string][]string
}

// rawData of csv file, where key is header of table and key is a slice of rows
type Data struct {
	headers     []string
	currentLine []string
	reader      *csv.Reader
	Table       map[string][]string
	mu          sync.Mutex
}
