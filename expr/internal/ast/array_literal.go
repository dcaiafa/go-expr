package ast

import (
	"fmt"
	"log"

	"github.com/dcaiafa/go-expr/expr/internal/context"
	"github.com/dcaiafa/go-expr/expr/runtime"
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
		for _, ast := range e.elements {
			err := ast.RunPass(ctx, pass)
			if err != nil {
				return err
			}
		}
		err := e.checkTypes()
		if err != nil {
			return err
		}

	case context.Emit:
		err := e.emit(ctx)
		if err != nil {
			return err
		}

	default:
		for _, ast := range e.elements {
			err := ast.RunPass(ctx, pass)
			if err != nil {
				return err
			}
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
	optimize := true
	for _, elem := range e.elements {
		if elem.Value() == nil {
			optimize = false
			break
		}
	}
	if optimize {
		array := make([]runtime.RawValue, len(e.elements))
		for i, elem := range e.elements {
			switch v := elem.Value().(type) {
			case float64:
				array[i] = runtime.NewRawNumber(v)
			case bool:
				array[i] = runtime.NewRawBool(v)
			case string:
				array[i] = runtime.NewRawObject(v)
			default:
				log.Fatal("invalid array literal folded value")
			}
		}
		arrayValue := runtime.NewObject(e.Type(), array)
		arrayConst := ctx.Builder.NewConst(arrayValue)
		ctx.Builder.EmitLoadConst(arrayConst)
		return nil
	}

	for _, ast := range e.elements {
		err := ast.RunPass(ctx, context.Emit)
		if err != nil {
			return err
		}
	}
	ctx.Builder.EmitPushArray(len(e.elements))
	return nil
}
