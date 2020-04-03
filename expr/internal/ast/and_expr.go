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
	value interface{}
}

func NewAndExpr(left Expr, right Expr) *AndExpr {
	return &AndExpr{left: left, right: right}
}

func (e *AndExpr) Type() types.Type {
	return types.Bool
}

func (e *AndExpr) Value() interface{} {
	return e.value
}

func (e *AndExpr) Print(p *context.GraphPrinter) {
	p.PrintNode("&&", e.left, e.right)
}

func (e *AndExpr) RunPass(ctx *context.Context, pass context.Pass) error {
	switch pass {
	case context.CheckTypes:
		err := e.checkTypes(ctx)
		if err != nil {
			return err
		}

	case context.Emit:
		err := e.emit(ctx)
		if err != nil {
			return err
		}

	case context.Fold:
		err := e.fold(ctx)
		if err != nil {
			return err
		}

	default:
		err := e.left.RunPass(ctx, pass)
		if err != nil {
			return err
		}
		err = e.right.RunPass(ctx, pass)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *AndExpr) checkTypes(ctx *context.Context) error {
	err := e.left.RunPass(ctx, context.CheckTypes)
	if err != nil {
		return err
	}
	err = e.right.RunPass(ctx, context.CheckTypes)
	if err != nil {
		return err
	}

	if e.left.Type() != types.Bool {
		return fmt.Errorf("left side of && is not bool")
	}
	if e.right.Type() != types.Bool {
		return fmt.Errorf("right side of && is not bool")
	}
	return nil
}

func (e *AndExpr) fold(ctx *context.Context) error {
	err := e.left.RunPass(ctx, context.Fold)
	if err != nil {
		return err
	}
	err = e.right.RunPass(ctx, context.Fold)
	if err != nil {
		return err
	}

	if e.left.Value() == nil || e.right.Value() == nil {
		return nil
	}

	e.value = e.left.Value().(bool) && e.right.Value().(bool)
	return nil
}

func (e *AndExpr) emit(ctx *context.Context) error {
	if e.value != nil {
		ctx.Builder.EmitPushBool(e.value.(bool))
		return nil
	}

	err := e.left.RunPass(ctx, context.Emit)
	if err != nil {
		return err
	}

	skipRight := ctx.Builder.NewLabel()
	ctx.Builder.EmitOp(runtime.Duplicate)
	ctx.Builder.EmitJump(runtime.JumpIfFalse, skipRight)

	err = e.right.RunPass(ctx, context.Emit)
	if err != nil {
		return err
	}

	ctx.Builder.EmitOp(runtime.And)
	ctx.Builder.AssignLabel(skipRight)
	return nil
}
