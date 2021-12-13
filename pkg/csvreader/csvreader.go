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
	ProceedQuery(raw_query string, query sqlparser.Query, row []string) (Table, error)
}

// rawData of csv file, where key is header of table and key is a slice of rows
type Data struct {
	headers     []string
	currentLine []string
	reader      *csv.Reader
	Table       map[string][]string
	mu          sync.Mutex
}
