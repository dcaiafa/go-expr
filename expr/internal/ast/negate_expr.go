package ast

import (
	"fmt"

	"github.com/dcaiafa/go-expr/expr/internal/context"
	"github.com/dcaiafa/go-expr/expr/runtime"
	"github.com/dcaiafa/go-expr/expr/types"
)

type NegateExpr struct {
	expr  Expr
	value interface{}
}

func NewNegateExpr(expr Expr) *NegateExpr {
	return &NegateExpr{expr: expr}
}

func (e *NegateExpr) Type() types.Type {
	return types.Bool
}

func (e *NegateExpr) Value() interface{} {
	return e.value
}

func (e *NegateExpr) Print(p *context.GraphPrinter) {
	p.PrintNode("!", e.expr)
}

func (e *NegateExpr) RunPass(ctx *context.Context, pass context.Pass) error {
	switch pass {
	case context.CheckTypes:
		err := e.checkTypes(ctx)
		if err != nil {
			return err
		}

	case context.Fold:
		err := e.fold(ctx)
		if err != nil {
			return err
		}

	case context.Emit:
		err := e.emit(ctx)
		if err != nil {
			return err
		}

	default:
		err := e.expr.RunPass(ctx, pass)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *NegateExpr) checkTypes(ctx *context.Context) error {
	err := e.expr.RunPass(ctx, context.CheckTypes)
	if err != nil {
		return err
	}

	if e.expr.Type() != types.Bool {
		return fmt.Errorf("operator ! requires bool operand")
	}

	return nil
}

func (e *NegateExpr) fold(ctx *context.Context) error {
	err := e.expr.RunPass(ctx, context.Fold)
	if err != nil {
		return err
	}

	if e.expr.Value() == nil {
		return nil
	}

	e.value = !e.expr.Value().(bool)

	return nil
}

func (e *NegateExpr) emit(ctx *context.Context) error {
	if e.value != nil {
		ctx.Builder.EmitPushBool(e.value.(bool))
		return nil
	}

	err := e.expr.RunPass(ctx, context.Emit)
	if err != nil {
		return err
	}

	ctx.Builder.EmitOp(runtime.Negate)

	return nil
}
