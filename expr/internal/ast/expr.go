package ast

import "github.com/dcaiafa/go-expr/expr/types"

type Expr interface {
	AST
	Type() types.Type
}

type exprImpl struct {
	typ types.Type
}

func (i *exprImpl) Type() types.Type {
	return i.typ
}
