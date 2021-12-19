package sqlparser

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/xfiendx4life/gb_go_best_final/pkg/operations"
	"go.uber.org/zap"
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
func NewQuery() *Query {
	return &Query{
		columns: make([]string, 0),
	}
}

func deleteComas(data []string) []string {
	res := make([]string, len(data))
	for i, item := range data {
		if item[len(item)-1] == ',' {
			res[i] = item[:len(item)-1]
		} else {
			res[i] = item
		}
	}
	return res
}

// split raw query to condition and cols after validation
func (q *Query) SplitToConditionAndCols(rawQuery string) error {
	splitted, err := normalizeValidateQuery(rawQuery)
	if err != nil {
		return fmt.Errorf("can't parse query: %s", err)
	}
	ind := getIndexOf(splitted, "from")
	if ind == -1 {
		return fmt.Errorf("error parsing query: no table chosen")
	}
	q.columns = deleteComas(splitted[1:ind])
	q.tableName = strings.ToLower(splitted[ind+1])
	ind = getIndexOf(splitted, "where")
	if ind == -1 {
		q.Condition.RawCondition = nil
	} else {
		q.Condition.RawCondition = splitted[ind+1:]
	}
	return nil
}

// Parse raw query to postfix form
func (q *Query) ParseToPostfix(rawQuery string) ([]string, error) {
	err := q.SplitToConditionAndCols(rawQuery)
	if err != nil {
		return nil, fmt.Errorf("can't parse raw query %s", err)
	}
	root := Node{}
	ops := operations.InitOperations()
	root.ParseQueryToTree(q.Condition.RawCondition, ops)
	q.Condition.tree = &root
	postfix := make([]string, 0)
	root.postorder(&postfix)
	return postfix, nil
}

// find index of operation with the highest priority
func getHighestPriorityOperation(condition []string, ops operations.Operations) int {
	ind := -1
	prior := -1
	for i, el := range condition {
		if _, ok := ops[el]; ok && ops[el].Priority > prior {
			ind = i
			prior = ops[el].Priority
		}
	}
	return ind
}

// parsing condition of query to tree
func (root *Node) ParseQueryToTree(m []string, ops operations.Operations) {
	if len(m) == 1 {
		root.Data = m[0]
		return
	}
	ind := getHighestPriorityOperation(m, ops)
	root.Data = m[ind]
	if op := ops[m[ind]]; op.Binary {
		var Left Node
		root.Left = &Left
		root.Left.ParseQueryToTree(m[:ind], ops)
	}
	root.Right = &Node{}
	root.Right.ParseQueryToTree(m[ind+1:], ops)
}

func (root *Node) postorder(m *[]string) {
	if root != nil {
		root.Left.postorder(m)
		root.Right.postorder(m)
		*m = append(*m, root.Data)
	}
}

func (q *Query) GetResultCols() (cols []string) {
	cols = q.columns
	return
}

func (q *Query) GetTableName() (tableName string) {
	tableName = q.tableName
	return
}

// Checks if where clause returns true or false with current row
// if types are different SelectFromRow works as if they are strings
func (q *Query) SelectFromRow(ctx context.Context, postfix []string, row map[string]string, z *zap.SugaredLogger) (res bool, err error) {
	select {
	case <-ctx.Done():
		z.Errorf("done with context")
		return false, fmt.Errorf("stopped by context")
	default:
		stack := NewStack(0)
		ops := operations.InitOperations()
		z.Debugf("starting parsing query from postfix form")
		for _, item := range postfix {
			if o, ok := ops[item]; !ok {
				if value, ok := row[item]; ok {
					stack.Push(value)
				} else {
					stack.Push(item)
				}

			} else {
				var res string
				var b string
				b, err = stack.Pop()
				if err != nil {
					z.Errorf("stack error while making select %s", err)
					return false, fmt.Errorf("stack error while making select %s", err)
				}
				var a string
				// cheking if operation is binary
				if o.Binary {
					a, err = stack.Pop()
					if err != nil {
						z.Errorf("stack error while making select %s", err)
						return false, fmt.Errorf("stack error while making select %s", err)
					}
				}
				// check if operation is comapration
				if o.BasicOp != nil {
					op, terr := operations.OpsBuilder(a, b)
					if err != nil {
						z.Errorf("can't select row %s", terr)
						return false, fmt.Errorf("can't select row %s", terr)
					}
					res = strconv.FormatBool(o.BasicOp(op))
				} else { // check if operation is logic
					op, terr := operations.LogicBuilder(b, a)
					if err != nil {
						z.Errorf("error while making select %s", terr)
						return false, fmt.Errorf("error while making select %s", terr)
					}
					res = strconv.FormatBool(o.LogicOp(op))
				}
				stack.Push(res)
				// complete this shit
			}
		}
		if stack.IsEmpty() || stack.Len() > 1 {
			z.Errorf("not valid stack computation")
			return false, fmt.Errorf("not valid stack computation")
		}
		res, err = strconv.ParseBool(stack.data[0])
		if err != nil {
			z.Errorf("can't parse bool %s", err)
			return false, fmt.Errorf("can't parse bool %s", err)
		}
		return res, nil
	}
}
