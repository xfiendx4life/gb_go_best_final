package sqlparser

type Querier interface {
	// Parses query from string to []commands,
	// which is postfix form of query
	ParseToPostfix(rawQuery string) ([]string, error) // to postfix form
	// returns cols names to be listed in result table
	GetResultCols() (cols []string)
	// returns tablename
	GetTableName() (tableName string)
	// proceed query on row
	SelectFromRow(postfix []string, row map[string]string) (result bool, err error)
}

type WhereCondition struct {
	// тут похоже нужно делать дерево условий
	RawCondition []string
	tree         *Node
	// postfix      []string
}

// we keep all the data in
type Query struct {
	// rawQuery  string
	tableName string
	columns   []string // all columns if *
	Condition WhereCondition
}

type Node struct {
	Data  string
	Left  *Node
	Right *Node
}
