package sqlparser_test

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	ops "github.com/xfiendx4life/gb_go_best_final/pkg/operations"
	"github.com/xfiendx4life/gb_go_best_final/pkg/sqlparser"
)

func TestNewQueryNoWhere(t *testing.T) {
	q := sqlparser.NewQuery()
	err := q.SplitToConditionAndCols("SELECT * from students")
	assert.Nil(t, err, "should be nil")
	assert.Equal(t, q.GetTableName(), "students")
	assert.Nil(t, q.Condition.RawCondition)
}

func TestNewQueryWithWhere(t *testing.T) {
	q := sqlparser.NewQuery()
	err := q.SplitToConditionAndCols("SELECT * from students WHERE name =   'Jane' and lastname = 'Doe'")
	assert.Nil(t, err, "should be nil")
	assert.Equal(t, q.GetTableName(), "students")
	assert.NotNil(t, q.Condition.RawCondition)
}

func TestNewQueryWithoutSelect(t *testing.T) {
	q := sqlparser.NewQuery()
	err := q.SplitToConditionAndCols("* from students WHERE name =   'Jane' and lastname = 'Doe'")
	assert.NotNil(t, err)
}

func TestNewQueryWithoutFrom(t *testing.T) {
	q := sqlparser.NewQuery()
	err := q.SplitToConditionAndCols("SELECT * students WHERE name =   'Jane' and lastname = 'Doe'")
	assert.Nil(t, q.Condition.RawCondition, "should be nil")
	assert.NotNil(t, err)
}

func TestParseQueryBinaryOnly(t *testing.T) {
	condition := []string{"a", ">", "b", "and", "c", "<", "d"}
	root := sqlparser.Node{}
	root.ParseQueryToTree(condition, ops.InitOperations())
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
	root.ParseQueryToTree(condition, ops.InitOperations())
	assert.Equal(t, "and", root.Data)
	assert.Equal(t, ">", root.Left.Data)
	assert.Equal(t, "a", root.Left.Left.Data)
	assert.Equal(t, "b", root.Left.Right.Data)
	assert.Equal(t, "not", root.Right.Data)
	assert.Equal(t, "<", root.Right.Right.Data)
	assert.Equal(t, "c", root.Right.Right.Left.Data)
	assert.Equal(t, "d", root.Right.Right.Right.Data)
}

func TestParseRawQuery(t *testing.T) {
	q := sqlparser.NewQuery()
	res, err := q.ParseToPostfix("SELECT * FROM tablename where a > b and not c < d")
	assert.Nil(t, err)
	assert.NotEmpty(t, res)
	log.Print(res)
}

func TestStackIsEmpty(t *testing.T) {
	st := sqlparser.NewStack(0)
	assert.True(t, st.IsEmpty())
}

func TestStackPushAndPop(t *testing.T) {
	st := sqlparser.NewStack(0)
	st.Push("test")
	assert.False(t, st.IsEmpty())
	res, err := st.Pop()
	assert.Nil(t, err)
	assert.Equal(t, "test", res)
	assert.True(t, st.IsEmpty())
}

func TestStackPushAndLen(t *testing.T) {
	st := sqlparser.NewStack(0)
	st.Push("test")
	assert.False(t, st.IsEmpty())
	res := st.Len()
	assert.Equal(t, 1, res)
}

func TestSelectFromRow(t *testing.T) {
	q := sqlparser.NewQuery()
	postfix := []string{"a", "4", ">", "c", "6", "<", "not", "and"}
	row := map[string]string{
		"a": "5",
		"c": "8",
	}
	res, err := q.SelectFromRow(postfix, row)
	require.Nil(t, err)
	require.True(t, res)
}

//TODO: more tests before continue
