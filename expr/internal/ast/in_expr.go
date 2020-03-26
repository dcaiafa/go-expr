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

func (e *InExpr) Print(p *context.GraphPrinter) {
	p.PrintNode("in", e.left, e.right)
}

func (e *InExpr) RunPass(ctx *context.Context, pass context.Pass) error {
	if pass == context.Emit {
		if e.left.Type().Equal(types.String) {
			ctx.Builder.EmitLoadConst(runtime.InternalInStringArray)
		} else if e.left.Type().Equal(types.Number) {
			ctx.Builder.EmitLoadConst(runtime.InternalInNumberArray)
		} else {
			log.Fatal("unexpected left type")
		}
	}

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
		if !e.left.Type().Equal(types.String) &&
			!e.left.Type().Equal(types.Number) {
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

	case context.Emit:
		ctx.Builder.EmitCall(2)
	}

	return nil
}
