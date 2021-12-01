package domain

type Table interface {
	// Reads data from csv
	Read(csv_table []byte) (n int, err error)
	// parses whole command sequence
	ParseCommands(commands []Command) Table 
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