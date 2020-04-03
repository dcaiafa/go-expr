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
	value interface{}
}

func NewOrExpr(left Expr, right Expr) *OrExpr {
	return &OrExpr{left: left, right: right}
}

func (e *OrExpr) Type() types.Type {
	return types.Bool
}

func (e *OrExpr) Value() interface{} {
	return e.value
}

func (e *OrExpr) Print(p *context.GraphPrinter) {
	p.PrintNode("||", e.left, e.right)
}

func (e *OrExpr) RunPass(ctx *context.Context, pass context.Pass) error {
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

func (e *OrExpr) checkTypes(ctx *context.Context) error {
	err := e.left.RunPass(ctx, context.CheckTypes)
	if err != nil {
		return err
	}
	err = e.right.RunPass(ctx, context.CheckTypes)
	if err != nil {
		return err
	}
	if e.left.Type() != types.Bool {
		return fmt.Errorf("left side of || is not bool")
	}
	if e.right.Type() != types.Bool {
		return fmt.Errorf("right side of || is not bool")
	}
	return nil
}

func (e *OrExpr) fold(ctx *context.Context) error {
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

	e.value = e.left.Value().(bool) || e.right.Value().(bool)
	return nil
}

func (e *OrExpr) emit(ctx *context.Context) error {
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
	ctx.Builder.EmitJump(runtime.JumpIfTrue, skipRight)

	err = e.right.RunPass(ctx, context.Emit)
	if err != nil {
		return err
	}

	ctx.Builder.EmitOp(runtime.Or)
	ctx.Builder.AssignLabel(skipRight)

	return nil
}
