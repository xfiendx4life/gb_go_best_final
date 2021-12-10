package sqlparser

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



