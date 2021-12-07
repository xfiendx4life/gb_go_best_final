package domain

import "io"

type Table interface {
	// Reads headers from
	ReadLine(source io.Reader) (err error)
	// parses whole command sequence
	ParseCommands(commands []Command) Table
}

type Source interface {
	// Open file for reading
	Read(path string) (file []byte, err error)
}

type Command interface {
	// Executes one command and returns resulting table
	Execute(table Table) Table
}

type Query interface {
	// Parses query from string to []commands,
	// which is postfix form of query
	ParseToPostfix(rawQuery string) ([]string, error) // to postfix form
	// returns cols names to be listed in result table
	GetResultCols() (cols []string)
	// returns tablename
	GetTableName() (tableName string, err error)
}
