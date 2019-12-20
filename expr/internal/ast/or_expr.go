package ast

import (
	"fmt"

	"github.com/dcaiafa/go-expr/expr/internal/context"
	"github.com/dcaiafa/go-expr/expr/runtime"
	"github.com/dcaiafa/go-expr/expr/types"
)

type OrExpr struct {
	left  Expr
	right Expr
}

func NewOrExpr(left Expr, right Expr) *OrExpr {
	return &OrExpr{left: left, right: right}
}

func (e *OrExpr) Type() types.Type {
	return types.Bool
}

func (e *OrExpr) Print(p *context.GraphPrinter) {
	p.PrintNode("||", e.left, e.right)
}

func (e *OrExpr) RunPass(ctx *context.Context, pass context.Pass) error {
	err := e.left.RunPass(ctx, pass)
	if err != nil {
		return err
	}

	var skipRight *runtime.Label
	if pass == context.Emit {
		skipRight = ctx.Builder.NewLabel()
		ctx.Builder.EmitOp(runtime.Duplicate)
		ctx.Builder.EmitJump(runtime.JumpIfTrue, skipRight)
	}

	err = e.right.RunPass(ctx, pass)
	if err != nil {
		return err
	}

	if pass == context.Emit {
		ctx.Builder.EmitOp(runtime.Or)
		ctx.Builder.AssignLabel(skipRight)
	} else if pass == context.CheckTypes {
		err = e.checkTypes(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *OrExpr) checkTypes(ctx *context.Context) error {
	if e.left.Type() != types.Bool {
		return fmt.Errorf("left side of || is not bool")
	}
	if e.right.Type() != types.Bool {
		return fmt.Errorf("right side of || is not bool")
	}
	return nil
}
