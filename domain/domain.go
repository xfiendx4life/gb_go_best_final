package domain

import "io"

type Table interface {
	// Reads headers from
	ReadHeaders(source io.Reader) (headers []string, err error)
	// Reads any line
	ReadRow(source io.Reader) (row []string, err error)
	// parses whole command sequence
	// returns the same structure ready to parse
	//another query
	ProceedQuery(query string, headers []string) (Table, error)
}

type Source interface {
	// Open file for reading
	Read(path string) (file []byte, err error)
}

type Query interface {
	// Parses query from string to []commands,
	// which is postfix form of query
	ParseToPostfix(rawQuery string) ([]string, error) // to postfix form
	// returns cols names to be listed in result table
	GetResultCols() (cols []string)
	// returns tablename
	GetTableName() (tableName string, err error)
	// proceed query on row
	SelectFromRow(postfixQuery []string) (result bool, err error)
}

type Config interface {
	// get path of table
	GetPath() (path string)
	// read config from file with path
	ReadConfig(path string) (err error)
}
