package ast

import (
	"fmt"

	"github.com/dcaiafa/go-expr/expr/exprerrors"
	"github.com/dcaiafa/go-expr/expr/internal/context"
	"github.com/dcaiafa/go-expr/expr/runtime"
	"github.com/dcaiafa/go-expr/expr/types"
)

type CallExpr struct {
	exprImpl
	receiver Expr
	params   *Params
}

func NewCallExpr(receiver Expr, params *Params) *CallExpr {
	return &CallExpr{
		receiver: receiver,
		params:   params,
	}
}

func (e *CallExpr) Print(p *context.GraphPrinter) {
	p.PrintNode("call", e.receiver, e.params)
}

func (e *CallExpr) RunPass(ctx *context.Context, pass context.Pass) error {
	err := e.receiver.RunPass(ctx, pass)
	if err != nil {
		return err
	}
	err = e.params.RunPass(ctx, pass)
	if err != nil {
		return err
	}

	switch pass {
	case context.CheckTypes:
		err = e.checkTypes()
		if err != nil {
			return err
		}
	case context.Emit:
		e.emit(ctx.Builder)
	}

	return nil
}

func (e *CallExpr) checkTypes() error {
	fn, ok := e.receiver.Type().(*types.Function)
	if !ok {
		return fmt.Errorf("%w: receiver in call expression is not a function",
			exprerrors.ErrSemantic)
	}

	if len(fn.Args) != len(e.params.params) {
		return fmt.Errorf(
			"%w: function expected %d parameters but %d were provided",
			exprerrors.ErrSemantic, len(fn.Args), len(e.params.params))
	}

	for i, arg := range fn.Args {
		if !arg.Equal(e.params.params[i].Type()) {
			return fmt.Errorf(
				"%w: parameter %d expected type is %v but %v was provided",
				exprerrors.ErrSemantic, i, e.params.params[i].Type(), arg)
		}
	}

	e.typ = fn.Ret

	return nil
}

func (e *CallExpr) emit(builder *runtime.Builder) {
	builder.EmitCall(len(e.params.params))
}

type Params struct {
	params []Expr
}

func NewParams(param Expr) *Params {
	return &Params{
		params: []Expr{param},
	}
}

func (p *Params) AddParam(param Expr) {
	p.params = append(p.params, param)
}

func (p *Params) Print(gp *context.GraphPrinter) {
	gp.PrintNode("params", exprsAsPrinters(p.params)...)
}

func (p *Params) RunPass(ctx *context.Context, pass context.Pass) error {
	for _, ast := range p.params {
		err := ast.RunPass(ctx, pass)
		if err != nil {
			return err
		}
	}
	return nil
}
