package operations_test

import (
	"testing"

	//	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
	"github.com/xfiendx4life/gb_go_best_final/pkg/operations"
)

func TestOpsBuilderBasic(t *testing.T) {
	_, err := operations.OpsBuilder("5", "10")
	require.Nil(t, err)
}

func TestCorrectBasic(t *testing.T) {
	op := operations.InitOperations()
	a := "5"
	b := "6"
	o, err := operations.OpsBuilder(a, b)
	require.Nil(t, err)
	require.True(t, op["<"].BasicOp(o))
	require.True(t, op["<="].BasicOp(o))
	require.False(t, op[">"].BasicOp(o))
	require.False(t, op[">="].BasicOp(o))
	require.False(t, op["="].BasicOp(o))
}

func TestInCorrectBasic(t *testing.T) {
	a := "5"
	b := "hello"
	_, err := operations.OpsBuilder(a, b)
	require.NotNil(t, err)
}

func TestCorrectStringBasic(t *testing.T) {
	op := operations.InitOperations()
	a := "test1"
	b := "test2"
	o, err := operations.OpsBuilder(a, b)
	require.Nil(t, err)
	require.True(t, op["<"].BasicOp(o))
	require.True(t, op["<="].BasicOp(o))
	require.False(t, op[">"].BasicOp(o))
	require.False(t, op[">="].BasicOp(o))
	require.False(t, op["="].BasicOp(o))
}

func TestCorrectLogic(t *testing.T) {
	op := operations.InitOperations()
	a := "true"
	b := "false"
	o, err := operations.LogicBuilder(a, b)
	require.Nil(t, err)
	require.False(t, op["and"].LogicOp(o))
	require.True(t, op["or"].LogicOp(o))
	require.False(t, op["not"].LogicOp(o))
}

func TestInCorrectLogic(t *testing.T) {
	a := "true"
	b := "String"
	_, err := operations.LogicBuilder(a, b)
	require.NotNil(t, err)
	a = "string"
	b = "true"
	_, err = operations.LogicBuilder(a, b)
	require.NotNil(t, err)
}
