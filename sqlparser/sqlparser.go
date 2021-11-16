package sqlparser

import (
	"fmt"
	"strings"
)

type WhereClause struct {
	// тут похоже нужно делать дерево условий
	RawClause []string
}

type Query struct {
	rawQuery  string
	TableName string
	Columns   []string // all columns if empty
	Clause    WhereClause
}

//Searching index of the first appereance of target string -1 if not in string
func getIndexOf(data []string, target string) int {
	ind := -1
	for i, word := range data {
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
