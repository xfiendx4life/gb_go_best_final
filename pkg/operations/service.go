package operations

import (
	"fmt"
	"strconv"
)

// change this func for reading from config
func InitOperations() Operations {
	return Operations{
		"<":   Operation{"<", 1, true, basicOps.less, nil},
		">":   Operation{">", 1, true, basicOps.more, nil},
		"<=":  Operation{"<=", 1, true, basicOps.lessOrEqual, nil},
		">=":  Operation{">=", 1, true, basicOps.moreOrEqual, nil},
		"=":   Operation{"=", 1, true, basicOps.equal, nil},
		"not": Operation{"not", 2, false, nil, logicOps.not},
		"and": Operation{"and", 3, true, nil, logicOps.and},
		"or":  Operation{"or", 4, true, nil, logicOps.or},
	}
}

// checks if s is numeric and returns float64 if it is
func getNumeric(s string) (float64, bool) {
	if a, err := strconv.ParseFloat(s, 64); err == nil {
		return a, true
	}
	return 0, false
}

type floatOpers struct {
	a, b float64
}

func (fl *floatOpers) less() bool {
	return fl.a < fl.b
}

func (fl *floatOpers) lessOrEqual() bool {
	return fl.a <= fl.b
}

func (fl *floatOpers) more() bool {
	return fl.a > fl.b
}

func (fl *floatOpers) moreOrEqual() bool {
	return fl.a >= fl.b
}

func (fl *floatOpers) equal() bool {
	return fl.a == fl.b
}

// type dateOpers struct {
// 	a, b time
// }

func OpsBuilder(op1, op2 string) (basicOps, error) {
	if a, ok := getNumeric(op1); ok {
		if b, ok := getNumeric(op2); ok {
			return &floatOpers{a, b}, nil
		} else {
			return nil, fmt.Errorf("can't compare numeric and not numeric types")
		}
	} else {
		return &stringOpers{op1, op2}, nil
	}
}

type boolOps struct {
	a []bool
}

func (l *boolOps) and() bool {
	return l.a[0] && l.a[1]
}

func (l *boolOps) or() bool {
	return l.a[0] || l.a[1]
}

func (l *boolOps) not() bool {
	return !l.a[0]
}

func LogicBuilder(ops ...string) (logicOps, error) {
	var a, b bool
	var err error
	if a, err = strconv.ParseBool(ops[0]); err != nil {
		return nil, fmt.Errorf("can't parse bool %s", err)
	}
	if ops[1] != "" {
		if b, err = strconv.ParseBool(ops[1]); err != nil {
			return nil, fmt.Errorf("can't parse bool %s", err)
		}
	}
	return &boolOps{[]bool{a, b}}, nil
}

type stringOpers struct {
	a, b string
}

func (fl *stringOpers) less() bool {
	return fl.a < fl.b
}

func (fl *stringOpers) lessOrEqual() bool {
	return fl.a <= fl.b
}

func (fl *stringOpers) more() bool {
	return fl.a > fl.b
}

func (fl *stringOpers) moreOrEqual() bool {
	return fl.a >= fl.b
}

func (fl *stringOpers) equal() bool {
	return fl.a == fl.b
}
