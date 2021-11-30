package sqlparser

import (
	"fmt"
	"strings"
)

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

// clear double spaces
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

// Func to create query and separate cols from condition
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
		q.Condition.RawCondition = nil
	} else {
		q.Condition.RawCondition = splitted[ind:]
	}
	return &q, nil
}

// Parse raw query to postfix form
func (q *Query) ParseRawQuery(rawQuery string) []string {
	return []string{}
}

// find index of operation with the highest priority
func getHighestPriorityOperation(condition []string, ops Operations) int {
	ind := -1
	prior := -1
	for i, el := range condition {
		if _, ok := ops[el]; ok && ops[el].priority > prior {
			ind = i
			prior = ops[el].priority
		}
	}
	return ind
}

// parsing condition of query to tree
func (root *Node) ParseQueryToTree(m []string, ops Operations) {
	if len(m) == 1 {
		root.Data = m[0]
		return
	}
	ind := getHighestPriorityOperation(m, ops)
	root.Data = m[ind]
	if op := ops[m[ind]]; op.binary {
		var Left Node
		root.Left = &Left
		root.Left.ParseQueryToTree(m[:ind], ops)
	}
	root.Right = &Node{}
	root.Right.ParseQueryToTree(m[ind+1:], ops)
}
