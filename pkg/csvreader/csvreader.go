package csvreader

import "encoding/csv"

// rawData of csv file, where key is header of table and key is a slice of rows
type RawData struct {
	headers     []string
	currentLine []string
	reader      *csv.Reader
}
