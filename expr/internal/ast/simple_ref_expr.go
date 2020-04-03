package ast

import (
	"github.com/dcaiafa/go-expr/expr/internal/context"
	"github.com/dcaiafa/go-expr/expr/internal/symbol"
	"github.com/dcaiafa/go-expr/expr/types"
)

type SimpleRefExpr struct {
	id  string
	sym symbol.Symbol
}

func NewSimpleRefExpr(id string) *SimpleRefExpr {
	return &SimpleRefExpr{
		id: id,
	}
}

func (e *SimpleRefExpr) Type() types.Type {
	return e.sym.Type()
}

func (e *SimpleRefExpr) Value() interface{} {
	return nil
}

func (e *SimpleRefExpr) Print(p *context.GraphPrinter) {
	p.PrintNode(e.id)
}

func (e *SimpleRefExpr) RunPass(ctx *context.Context, pass context.Pass) error {
	switch pass {
	case context.ResolveNames:
		err := e.resolveNames(ctx)
		if err != nil {
			return err
		}

	case context.Emit:
		err := e.emit(ctx)
		if err != nil {
			return err
		}

	default:
	}

	return nil
}

func (e *SimpleRefExpr) resolveNames(ctx *context.Context) error {
	var err error
	e.sym, err = ctx.GlobalScope.Get(e.id)
	return err
}

func (e *SimpleRefExpr) emit(ctx *context.Context) error {
	e.sym.EmitAccess(ctx.Builder)
	return nil
}
