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
	assert.Nil(t, q.Clause.RawClause)
}

func TestNewQueryWithWhere(t *testing.T) {
	q, err := sqlparser.NewQuery("SELECT * from students WHERE name =   'Jane' and lastname = 'Doe'")
	assert.Nil(t, err, "should be nil")
	assert.Equal(t, q.TableName, "students")
	assert.NotNil(t, q.Clause.RawClause)
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
