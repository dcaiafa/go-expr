package ast

import (
	"fmt"

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
	Ne
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
	case Ne:
		return "!="
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
		err := e.runPassChildren(ctx, pass)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *BinaryExpr) runPassChildren(ctx *context.Context, pass context.Pass) error {
	err := e.left.RunPass(ctx, pass)
	if err != nil {
		return err
	}

	err = e.right.RunPass(ctx, pass)
	if err != nil {
		return err
	}

	return nil
}

func (e *BinaryExpr) checkTypes(ctx *context.Context) error {
	err := e.runPassChildren(ctx, context.CheckTypes)
	if err != nil {
		return err
	}

	if e.left.Type() == nil || e.right.Type() == nil {
		panic("sub-expressions have unevaluated types")
	}

	switch e.op {
	case Lt, Le, Gt, Ge:
		if e.left.Type() != types.Number || e.right.Type() != types.Number {
			return fmt.Errorf("operator requires number operands")
		}
		e.typ = types.Bool

	case Plus, Minus, Times, Div:
		if e.left.Type() != types.Number || e.right.Type() != types.Number {
			return fmt.Errorf("operator requires number operands")
		}
		e.typ = types.Number

	case Eq, Ne:
		if !e.left.Type().Equal(e.right.Type()) {
			return fmt.Errorf("invalid operation: mistmatched types %v and %v",
				e.left.Type(), e.right.Type())
		}
		if e.left.Type() != types.String &&
			e.left.Type() != types.Number &&
			e.left.Type() != types.Bool &&
			!e.left.Type().Equal(&types.Array{ElementType: types.Bool}) &&
			!e.left.Type().Equal(&types.Array{ElementType: types.Number}) &&
			!e.left.Type().Equal(&types.Array{ElementType: types.String}) {
			return fmt.Errorf("invalid operation: cannot compare type %v", e.left.Type())
		}
		e.typ = types.Bool

	default:
		panic("invalid operator")
	}

	return nil
}

func (e *BinaryExpr) fold(ctx *context.Context) error {
	err := e.runPassChildren(ctx, context.Fold)
	if err != nil {
		return err
	}

	if e.left.Value() == nil || e.right.Value() == nil {
		return nil
	}

	switch e.op {
	case Lt:
		e.value = e.left.Value().(float64) < e.right.Value().(float64)
	case Le:
		e.value = e.left.Value().(float64) <= e.right.Value().(float64)
	case Gt:
		e.value = e.left.Value().(float64) > e.right.Value().(float64)
	case Ge:
		e.value = e.left.Value().(float64) >= e.right.Value().(float64)

	case Plus:
		e.value = e.left.Value().(float64) + e.right.Value().(float64)
	case Minus:
		e.value = e.left.Value().(float64) - e.right.Value().(float64)
	case Times:
		e.value = e.left.Value().(float64) * e.right.Value().(float64)
	case Div:
		e.value = e.left.Value().(float64) / e.right.Value().(float64)

	case Eq, Ne:
		if e.left.Type() == types.Number {
			e.value = e.left.Value().(float64) == e.right.Value().(float64)
		} else if e.left.Type() == types.String {
			e.value = e.left.Value().(string) == e.right.Value().(string)
		} else if e.left.Type() == types.Bool {
			e.value = e.left.Value().(bool) == e.right.Value().(bool)
		} else {
			panic("unexpected type")
		}
		if e.op == Ne {
			e.value = !e.value.(bool)
		}
	}

	return nil
}

func (e *BinaryExpr) emit(ctx *context.Context) error {
	if e.value != nil {
		ctx.Builder.EmitPushBasicValue(e.value)
		return nil
	}

	err := e.runPassChildren(ctx, context.Emit)
	if err != nil {
		return err
	}

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

	case Eq, Ne:
		if e.left.Type() == types.Number {
			ctx.Builder.EmitOp(runtime.CompareEqNumber)
		} else if e.left.Type() == types.String {
			ctx.Builder.EmitOp(runtime.CompareEqString)
		} else if e.left.Type() == types.Bool {
			ctx.Builder.EmitOp(runtime.CompareEqBool)
		} else if e.left.Type().Equal(&types.Array{ElementType: types.Bool}) {
			ctx.Builder.EmitOp(runtime.CompareEqArrayBool)
		} else if e.left.Type().Equal(&types.Array{ElementType: types.Number}) {
			ctx.Builder.EmitOp(runtime.CompareEqArrayNumber)
		} else if e.left.Type().Equal(&types.Array{ElementType: types.String}) {
			ctx.Builder.EmitOp(runtime.CompareEqArrayString)
		} else {
			panic("unexpected type with == operator")
		}
		if e.op == Ne {
			ctx.Builder.EmitOp(runtime.Negate)
		}
	}

	return nil
}
