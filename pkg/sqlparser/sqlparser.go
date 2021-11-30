package sqlparser

type WhereCondition struct {
	// тут похоже нужно делать дерево условий
	RawCondition []string
	tree         *Node
	postfix      []string
}

type Query struct {
	rawQuery  string
	TableName string
	Columns   []string // all columns if empty
	Condition WhereCondition
}

type Node struct {
	Data  string
	Left  *Node
	Right *Node
}

type Operation struct {
	name     string
	priority int
	binary   bool
}

type Operations map[string]Operation

func InitOperations() Operations {
	return Operations{
		"<":   Operation{"<", 1, true},
		">":   Operation{">", 1, true},
		"<=":  Operation{"<=", 1, true},
		">=":  Operation{">=", 1, true},
		"=":   Operation{"=", 1, true},
		"not": Operation{"not", 2, false},
		"and": Operation{"and", 3, true},
		"or":  Operation{"or", 4, true},
	}
}
