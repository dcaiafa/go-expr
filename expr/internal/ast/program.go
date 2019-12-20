package ast

import (
	"fmt"

	"github.com/dcaiafa/go-expr/expr/internal/context"
	"github.com/dcaiafa/go-expr/expr/runtime"
	"github.com/dcaiafa/go-expr/expr/types"
)

type Program struct {
	exprs []Expr
	typ   types.Type
}

func NewProgram(expr Expr) *Program {
	return &Program{exprs: []Expr{expr}}
}

func (p *Program) AddExpr(expr Expr) {
	p.exprs = append(p.exprs, expr)
}

func (p *Program) Type() types.Type {
	return p.typ
}

func (p *Program) Print(gp *context.GraphPrinter) {
	gp.PrintNode("Program", exprsAsPrinters(p.exprs)...)
}

func (p *Program) RunPass(ctx *context.Context, pass context.Pass) error {
	for _, expr := range p.exprs {
		err := expr.RunPass(ctx, pass)
		if err != nil {
			return err
		}

		if pass == context.CheckTypes {
			if p.typ == nil {
				p.typ = expr.Type()
			} else if !expr.Type().Equal(p.typ) {
				return fmt.Errorf("mistmatched expression types %v and %v",
					p.typ, expr.Type())
			}
		}

		if pass == context.Emit {
			ctx.Builder.EmitOp(runtime.Return)
			ctx.Builder.FinishExpr()
		}
	}
	return nil
}

func exprsAsPrinters(exprs []Expr) []context.Printer {
	printers := make([]context.Printer, len(exprs))
	for i, expr := range exprs {
		printers[i] = expr
	}
	return printers
}
