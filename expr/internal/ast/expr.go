package ast

import "github.com/dcaiafa/go-expr/expr/types"

type Expr interface {
	AST
	Type() types.Type
	Value() interface{}
}

type exprImpl struct {
	typ   types.Type
	value interface{}
}

func (i *exprImpl) Type() types.Type {
	return i.typ
}

func (i *exprImpl) Value() interface{} {
	return i.value
}
