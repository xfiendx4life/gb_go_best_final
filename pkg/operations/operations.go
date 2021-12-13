package operations

// type Command func(ops) (bool, error)

type Operation struct {
	Name     string
	Priority int
	Binary   bool
	BasicOp  func(_ basicOps) bool
	LogicOp  func(_ logicOps) bool
}

type Operations map[string]Operation

type basicOps interface {
	less() bool
	lessOrEqual() bool
	more() bool
	moreOrEqual() bool
	equal() bool
}

type logicOps interface {
	and() bool
	not() bool
	or() bool
}
