package csvreader

import "fmt"

// rawData of csv file, where key is header of table and key is a slice of rows
type RawData struct {
	data map[string][]string
}

func NewRawData() *RawData {
	return &RawData{
		data: make(map[string][]string),
	}
}

func (r *RawData) Read(source []byte) (n int, err error) {
	// read data from csv to map[string]string
	return 1, fmt.Errorf("")
}
