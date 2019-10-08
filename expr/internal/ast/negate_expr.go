package ast

import (
	"fmt"

	"github.com/dcaiafa/go-expr/expr/exprerrors"
	"github.com/dcaiafa/go-expr/expr/internal/context"
	"github.com/dcaiafa/go-expr/expr/runtime"
	"github.com/dcaiafa/go-expr/expr/types"
)

type NegateExpr struct {
	expr Expr
}

func NewNegateExpr(expr Expr) *NegateExpr {
	return &NegateExpr{expr: expr}
}

func (e *NegateExpr) Type() types.Type {
	return types.Bool
}

func (e *NegateExpr) Print(p *context.GraphPrinter) {
	p.PrintNode("!", e.expr)
}

func (e *NegateExpr) RunPass(ctx *context.Context, pass context.Pass) error {
	err := e.expr.RunPass(ctx, pass)
	if err != nil {
		return err
	}

	switch pass {
	case context.CheckTypes:
		if e.expr.Type() != types.Bool {
			return fmt.Errorf("%w: ! operator requires bool operand",
				exprerrors.ErrSemantic)
		}

	case context.Emit:
		ctx.Builder.EmitOp(runtime.Negate)
	}

	return nil
}
