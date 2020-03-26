package ast

import (
	"fmt"

	"github.com/dcaiafa/go-expr/expr/internal/context"
	"github.com/dcaiafa/go-expr/expr/types"
)

type ArrayLiteralExpr struct {
	exprImpl
	elements []Expr
}

func NewArrayLiteralExpr(element Expr) *ArrayLiteralExpr {
	return &ArrayLiteralExpr{
		elements: []Expr{element},
	}
}

func (e *ArrayLiteralExpr) AddElement(element Expr) {
	e.elements = append(e.elements, element)
}

func (e *ArrayLiteralExpr) Print(p *context.GraphPrinter) {
	p.PrintNode("array_literal", exprsAsPrinters(e.elements)...)
}

func (e *ArrayLiteralExpr) RunPass(ctx *context.Context, pass context.Pass) error {
	switch pass {
	case context.CheckTypes:
		err := e.runPassOnElements(ctx, pass)
		if err != nil {
			return err
		}
		err = e.checkTypes()
		if err != nil {
			return err
		}

	case context.Emit:
		err := e.emit(ctx)
		if err != nil {
			return err
		}

	default:
		err := e.runPassOnElements(ctx, pass)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *ArrayLiteralExpr) runPassOnElements(ctx *context.Context, pass context.Pass) error {
	for _, ast := range e.elements {
		err := ast.RunPass(ctx, pass)
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *ArrayLiteralExpr) checkTypes() error {
	if len(e.elements) == 0 {
		e.typ = types.Void
		return nil
	}

	elemType := e.elements[0].Type()

	for _, element := range e.elements {
		if !element.Type().Equal(elemType) {
			return fmt.Errorf("all elements in array must have the same type")
		}
	}

	e.typ = &types.Array{ElementType: elemType}

	return nil
}

func (e *ArrayLiteralExpr) emit(ctx *context.Context) error {
	err := e.runPassOnElements(ctx, context.Emit)
	if err != nil {
		return err
	}
	ctx.Builder.EmitPushArray(len(e.elements))
	return nil
}
