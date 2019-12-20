package ast

import (
	"fmt"

	"github.com/dcaiafa/go-expr/expr/internal/context"
	"github.com/dcaiafa/go-expr/expr/runtime"
	"github.com/dcaiafa/go-expr/expr/types"
)

type AndExpr struct {
	left  Expr
	right Expr
}

func NewAndExpr(left Expr, right Expr) *AndExpr {
	return &AndExpr{left: left, right: right}
}

func (e *AndExpr) Type() types.Type {
	return types.Bool
}

func (e *AndExpr) Print(p *context.GraphPrinter) {
	p.PrintNode("&&", e.left, e.right)
}

func (e *AndExpr) RunPass(ctx *context.Context, pass context.Pass) error {
	err := e.left.RunPass(ctx, pass)
	if err != nil {
		return err
	}

	var skipRight *runtime.Label
	if pass == context.Emit {
		skipRight = ctx.Builder.NewLabel()
		ctx.Builder.EmitOp(runtime.Duplicate)
		ctx.Builder.EmitJump(runtime.JumpIfFalse, skipRight)
	}

	err = e.right.RunPass(ctx, pass)
	if err != nil {
		return err
	}

	if pass == context.Emit {
		ctx.Builder.EmitOp(runtime.And)
		ctx.Builder.AssignLabel(skipRight)
	} else if pass == context.CheckTypes {
		err = e.checkTypes(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *AndExpr) checkTypes(ctx *context.Context) error {
	if e.left.Type() != types.Bool {
		return fmt.Errorf("left side of && is not bool")
	}
	if e.right.Type() != types.Bool {
		return fmt.Errorf("right side of && is not bool")
	}
	return nil
}
