package sqlparser_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xfiendx4life/gb_go_best_final/pkg/sqlparser"
)

func TestNewQueryNoWhere(t *testing.T) {
	q, err := sqlparser.NewQuery("SELECT * from students")
	assert.Nil(t, err, "should be nil")
	assert.Equal(t, q.TableName, "students")
	assert.Nil(t, q.Condition.RawCondition)
}

func TestNewQueryWithWhere(t *testing.T) {
	q, err := sqlparser.NewQuery("SELECT * from students WHERE name =   'Jane' and lastname = 'Doe'")
	assert.Nil(t, err, "should be nil")
	assert.Equal(t, q.TableName, "students")
	assert.NotNil(t, q.Condition.RawCondition)
}

func TestNewQueryWithoutSelect(t *testing.T) {
	q, err := sqlparser.NewQuery("* from students WHERE name =   'Jane' and lastname = 'Doe'")
	assert.Nil(t, q, "should be nil")
	assert.NotNil(t, err)
}

func TestNewQueryWithoutFrom(t *testing.T) {
	q, err := sqlparser.NewQuery("SELECT * students WHERE name =   'Jane' and lastname = 'Doe'")
	assert.Nil(t, q, "should be nil")
	assert.NotNil(t, err)
}

func TestParseQueryBinaryOnly(t *testing.T) {
	condition := []string{"a", ">", "b", "and", "c", "<", "d"}
	root := sqlparser.Node{}
	root.ParseQueryToTree(condition, sqlparser.InitOperations())
	assert.Equal(t, "and", root.Data)
	assert.Equal(t, ">", root.Left.Data)
	assert.Equal(t, "a", root.Left.Left.Data)
	assert.Equal(t, "b", root.Left.Right.Data)
	assert.Equal(t, "<", root.Right.Data)
	assert.Equal(t, "c", root.Right.Left.Data)
	assert.Equal(t, "d", root.Right.Right.Data)
}

func TestParseQuery(t *testing.T) {
	condition := []string{"a", ">", "b", "and", "not", "c", "<", "d"}
	root := sqlparser.Node{}
	root.ParseQueryToTree(condition, sqlparser.InitOperations())
	assert.Equal(t, "and", root.Data)
	assert.Equal(t, ">", root.Left.Data)
	assert.Equal(t, "a", root.Left.Left.Data)
	assert.Equal(t, "b", root.Left.Right.Data)
	assert.Equal(t, "not", root.Right.Data)
	assert.Equal(t, "<", root.Right.Right.Data)
	assert.Equal(t, "c", root.Right.Right.Left.Data)
	assert.Equal(t, "d", root.Right.Right.Right.Data)
}
