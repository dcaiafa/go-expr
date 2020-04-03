package ast

import (
	"fmt"
	"log"

	"github.com/dcaiafa/go-expr/expr/internal/context"
	"github.com/dcaiafa/go-expr/expr/runtime"
	"github.com/dcaiafa/go-expr/expr/types"
)

type InExpr struct {
	left  Expr
	right Expr
}

func NewInExpr(left, right Expr) *InExpr {
	return &InExpr{left: left, right: right}
}

func (e *InExpr) Type() types.Type {
	return types.Bool
}

func (e *InExpr) Value() interface{} {
	return nil
}

func (e *InExpr) Print(p *context.GraphPrinter) {
	p.PrintNode("in", e.left, e.right)
}

func (e *InExpr) RunPass(ctx *context.Context, pass context.Pass) error {
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

func (e *InExpr) checkTypes(ctx *context.Context) error {
	err := e.left.RunPass(ctx, context.CheckTypes)
	if err != nil {
		return err
	}
	err = e.right.RunPass(ctx, context.CheckTypes)
	if err != nil {
		return err
	}

	if !e.left.Type().Equal(types.String) && !e.left.Type().Equal(types.Number) {
		return fmt.Errorf(
			"only number and string supported by 'in' expression, "+
				"but left side is %v", e.left.Type())
	}

	arrayType, ok := e.right.Type().(*types.Array)
	if !ok || !e.left.Type().Equal(arrayType.ElementType) {
		return fmt.Errorf(
			"right side of 'in' expression should be array of %v, but it is %v",
			e.left.Type(), e.right.Type())
	}

	return nil
}

func (e *InExpr) emit(ctx *context.Context) error {
	err := e.left.RunPass(ctx, context.Emit)
	if err != nil {
		return err
	}
	err = e.right.RunPass(ctx, context.Emit)
	if err != nil {
		return err
	}

	if e.left.Type() == types.Number {
		ctx.Builder.EmitOp(runtime.InArrayNumber)
	} else if e.left.Type() == types.String {
		ctx.Builder.EmitOp(runtime.InArrayString)
	} else {
		log.Fatal("unexpected left type")
	}

	return nil
}
