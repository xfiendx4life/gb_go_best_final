package csvreader

import (
	"encoding/csv"
	"sync"
)

// rawData of csv file, where key is header of table and key is a slice of rows
type Data struct {
	headers     []string
	currentLine []string
	reader      *csv.Reader
	Table      map[string][]string
	mu          sync.Mutex
}
