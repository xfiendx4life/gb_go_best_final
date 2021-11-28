package sqlparser

import (
	"fmt"
	"strings"
)

type WhereClause struct {
	// тут похоже нужно делать дерево условий
	RawClause []string
	tree      Tree
	postfix   []string
}

type Query struct {
	rawQuery  string
	TableName string
	Columns   []string // all columns if empty
	Clause    WhereClause
}

type Node struct {
	Data  string
	Left  *Node
	Right *Node
}

type Tree struct {
	root *Node
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

//Searching index of the first appereance of target string -1 if not in string
func getIndexOf(Data []string, target string) int {
	ind := -1
	for i, word := range Data {
		if strings.ToLower(word) == target {
			return i
		}
	}
	return ind
}

func normalize(query string) string {
	query = strings.ToLower(query)
	for strings.Contains(query, "  ") {
		query = strings.Replace(query, "  ", " ", -1)
	}
	return query
}

// Returns normalize and validated slice of words
func normalizeValidateQuery(query string) ([]string, error) {
	query = normalize(query)
	splitted := strings.Split(query, " ")
	if splitted[0] != "select" {
		return nil, fmt.Errorf("not valid select query: no select")
	}
	ind := getIndexOf(splitted, "from")
	if ind == -1 {
		return nil, fmt.Errorf("error parsing query: no table chosen")
	}
	return splitted, nil
}

func NewQuery(query string) (*Query, error) {
	q := Query{
		rawQuery: query,
		Columns:  make([]string, 0),
	}
	splitted, err := normalizeValidateQuery(query)
	if err != nil {
		return nil, fmt.Errorf("can't parse query: %s", err)
	}
	ind := getIndexOf(splitted, "from")
	if ind == -1 {
		return nil, fmt.Errorf("error parsing query: no table chosen")
	}
	q.Columns = splitted[1:ind]
	q.TableName = strings.ToLower(splitted[ind+1])
	ind = getIndexOf(splitted, "where")
	if ind == -1 {
		q.Clause.RawClause = nil
	} else {
		q.Clause.RawClause = splitted[ind:]
	}
	return &q, nil
}

// Parse raw query to postfix form
func (q *Query) ParseRawQuery(rawQuery string) []string {
	return []string{}
}

func getHighestPriorityOperation(clause []string, ops Operations) int {
	ind := -1
	prior := -1
	for i, el := range clause {
		if _, ok := ops[el]; ok && ops[el].priority > prior {
			ind = i
			prior = ops[el].priority
		}
	}
	return ind
}

func (root *Node) ParseQuery(m []string, ops Operations) {
	if len(m) == 1 {
		root.Data = m[0]
		return
	}
	ind := getHighestPriorityOperation(m, ops)
	root.Data = m[ind]
	if op := ops[m[ind]]; op.binary {
		var Left Node
		root.Left = &Left
		root.Left.ParseQuery(m[:ind], ops)
	}
	root.Right = &Node{}
	root.Right.ParseQuery(m[ind+1:], ops)
}
