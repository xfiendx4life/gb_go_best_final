package csvreader

import (
	"context"
	"encoding/csv"
	"io"
	"sync"

	"github.com/xfiendx4life/gb_go_best_final/pkg/sqlparser"
	"go.uber.org/zap"
)

type Source interface {
	// Open file for reading
	Read(path string) (file []byte, err error)
}

type Table interface {
	// Reads headers from
	ReadHeaders(source io.Reader) (headers []string, err error)
	// Reads any line
	ReadRow(source io.Reader, sep rune) (row []string, err error)
	// parses whole command sequence
	// returns the same structure ready to parse
	//another query
	ProceedQuery(ctx context.Context, query sqlparser.Querier, row []string, postfix []string, z *zap.SugaredLogger) (Table, error)
	// proceeds select query to the whole table concurrently
	ProceedFullTable(ctx context.Context, source io.Reader, rawQuery string, z *zap.SugaredLogger, resChan chan Table, errChan chan error)
	GetTable() map[string][]string
	GetHeaders() []string
}

// rawData of csv file, where key is header of table and key is a slice of rows
type Data struct {
	headers     []string
	currentLine []string
	reader      *csv.Reader
	Table       map[string][]string
	mu          sync.Mutex
}
