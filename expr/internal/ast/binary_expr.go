package ast

import (
	"fmt"

	"github.com/dcaiafa/go-expr/expr/exprerrors"
	"github.com/dcaiafa/go-expr/expr/internal/context"
	"github.com/dcaiafa/go-expr/expr/runtime"
	"github.com/dcaiafa/go-expr/expr/types"
)

type BinaryOp int

const (
	Lt BinaryOp = iota
	Le
	Gt
	Ge
	Plus
	Minus
	Times
	Div
	Eq
)

func (o BinaryOp) String() string {
	switch o {
	case Lt:
		return "<"
	case Le:
		return "<="
	case Gt:
		return ">"
	case Ge:
		return ">="
	case Plus:
		return "+"
	case Minus:
		return "-"
	case Times:
		return "*"
	case Div:
		return "/"
	case Eq:
		return "=="
	default:
		return "???"
	}
}

type BinaryExpr struct {
	exprImpl
	left  Expr
	op    BinaryOp
	right Expr
}

func NewBinaryExpr(left Expr, op BinaryOp, right Expr) *BinaryExpr {
	return &BinaryExpr{
		left:  left,
		op:    op,
		right: right,
	}
}

func (e *BinaryExpr) Print(p *context.GraphPrinter) {
	p.PrintNode(e.op.String(), e.left, e.right)
}

func (e *BinaryExpr) RunPass(ctx *context.Context, pass context.Pass) error {
	err := e.left.RunPass(ctx, pass)
	if err != nil {
		return err
	}

	err = e.right.RunPass(ctx, pass)
	if err != nil {
		return err
	}

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
	}

	return nil
}

func (e *BinaryExpr) checkTypes(ctx *context.Context) error {
	if e.left.Type() == nil || e.right.Type() == nil {
		panic("sub-expressions have unevaluated types")
	}

	switch e.op {
	case Lt, Le, Gt, Ge:
		if e.left.Type() != types.Number || e.right.Type() != types.Number {
			return fmt.Errorf(
				"%w: binary expression: expected number",
				exprerrors.ErrSemantic)
		}
		e.typ = types.Bool

	case Plus, Minus, Times, Div:
		if e.left.Type() != types.Number || e.right.Type() != types.Number {
			return fmt.Errorf(
				"%w: binary expression: expected number",
				exprerrors.ErrSemantic)
		}
		e.typ = types.Number

	case Eq:
		if e.left.Type() != e.right.Type() {
			return fmt.Errorf(
				"%w: binary expression: both sides must have the same type",
				exprerrors.ErrSemantic)
		}
		if e.left.Type() != types.String &&
			e.left.Type() != types.Number &&
			e.left.Type() != types.Bool {
			return fmt.Errorf(
				"%w: binary expression: expected string, number or bool",
				exprerrors.ErrSemantic)
		}
		e.typ = types.Bool

	default:
		panic("invalid operator")
	}

	return nil
}

func (e *BinaryExpr) emit(ctx *context.Context) error {
	switch e.op {
	case Lt:
		ctx.Builder.EmitOp(runtime.CompareLT)
	case Le:
		ctx.Builder.EmitOp(runtime.CompareLE)
	case Gt:
		ctx.Builder.EmitOp(runtime.CompareGT)
	case Ge:
		ctx.Builder.EmitOp(runtime.CompareGE)

	case Plus:
		ctx.Builder.EmitOp(runtime.Add)
	case Minus:
		ctx.Builder.EmitOp(runtime.Subtract)
	case Times:
		ctx.Builder.EmitOp(runtime.Multiply)
	case Div:
		ctx.Builder.EmitOp(runtime.Divide)

	case Eq:
		if e.left.Type() == types.Number {
			ctx.Builder.EmitOp(runtime.CompareEq)
		} else if e.left.Type() == types.String {
			ctx.Builder.EmitOp(runtime.CompareEqString)
		} else if e.left.Type() == types.Bool {
			ctx.Builder.EmitOp(runtime.CompareEqBool)
		} else {
			panic("unexpected type with == operator")
		}
	}

	return nil
}
