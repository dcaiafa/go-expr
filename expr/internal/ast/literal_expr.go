package ast

import (
	"fmt"

	"github.com/dcaiafa/go-expr/expr/internal/context"
	"github.com/dcaiafa/go-expr/expr/types"
)

type LiteralExpr struct {
	exprImpl
}

func NewLiteralExpr(typ types.Type, value interface{}) *LiteralExpr {
	return &LiteralExpr{
		exprImpl: exprImpl{
			typ:   typ,
			value: value,
		},
	}
}

func (e *LiteralExpr) Print(p *context.GraphPrinter) {
	p.PrintNode(fmt.Sprintf("%v", e.value))
}

func (e *LiteralExpr) RunPass(ctx *context.Context, pass context.Pass) error {
	switch pass {
	case context.Emit:
		err := e.emit(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *LiteralExpr) emit(ctx *context.Context) error {
	ctx.Builder.EmitPushBasicValue(e.value)
	return nil
}
