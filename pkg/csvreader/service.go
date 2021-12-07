package csvreader

import (
	"encoding/csv"
	"io"
)

func NewRawData() *RawData {
	return &RawData{
		headers:     make([]string, 0),
		currentLine: make([]string, 0),
	}
}

func (r *RawData) GetHeaders() []string {
	return r.headers
}

func (r *RawData) ReadLine(source io.Reader) (err error) {
	r.reader = csv.NewReader(source)
	r.headers, err = r.reader.Read()
	if err != nil {
		return err
	}
	return nil
}
