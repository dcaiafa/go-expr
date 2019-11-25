package expr

import (
	"fmt"
	"io"

	"github.com/dcaiafa/go-expr/expr/internal/context"
	"github.com/dcaiafa/go-expr/expr/internal/parser"
	"github.com/dcaiafa/go-expr/expr/internal/symbol"
	"github.com/dcaiafa/go-expr/expr/runtime"
	"github.com/dcaiafa/go-expr/expr/types"
)

type Compiler struct {
	ctx *context.Context
}

func NewCompiler() *Compiler {
	return &Compiler{
		ctx: context.NewContext(),
	}
}

func (c *Compiler) RegisterInput(name string, t runtime.ValueType) (int, error) {
	typ, err := runtimeTypeToSymbolType(t)
	if err != nil {
		return 0, err
	}
	inputIndex := c.ctx.Builder.NewInput(t)
	inputSymbol := symbol.NewInputSymbol(name, typ, inputIndex)
	err = c.ctx.GlobalScope.Add(inputSymbol)
	if err != nil {
		return 0, err
	}
	return inputIndex, nil
}

func (c *Compiler) RegisterConst(name string, v runtime.Value) error {
	typ, err := runtimeTypeToSymbolType(v.Type())
	if err != nil {
		return err
	}

	constIndex := c.ctx.Builder.RegisterConst(v)
	constSymbol := symbol.NewConstSymbol(name, typ, constIndex)
	return c.ctx.GlobalScope.Add(constSymbol)
}

func (c *Compiler) RegisterFunc(
	name string,
	fn runtime.Func,
	ret runtime.ValueType,
	args ...runtime.ValueType,
) error {
	var err error

	fnType := &types.Function{
		Args: make([]types.Type, len(args)),
	}

	for i, arg := range args {
		fnType.Args[i], err = runtimeTypeToSymbolType(arg)
		if err != nil {
			return err
		}
	}

	fnType.Ret, err = runtimeTypeToSymbolType(ret)
	if err != nil {
		return err
	}

	return c.registerFunc(name, fnType, fn)
}

func (c *Compiler) registerFunc(
	name string,
	typ *types.Function,
	fn runtime.Func,
) error {
	fnIndex := c.ctx.Builder.RegisterExternalFunc(fn)
	v := runtime.NewExternalFuncValue(fnIndex)
	constIndex := c.ctx.Builder.RegisterConst(v)
	fnSymbol := symbol.NewConstSymbol(name, typ, constIndex)
	return c.ctx.GlobalScope.Add(fnSymbol)
}

func (c *Compiler) Compile(input string) (*runtime.Program, error) {
	progAST, err := parser.Parse(input)
	if err != nil {
		return nil, err
	}
	err = progAST.RunPass(c.ctx, context.ResolveNames)
	if err != nil {
		return nil, err
	}
	err = progAST.RunPass(c.ctx, context.CheckTypes)
	if err != nil {
		return nil, err
	}
	err = progAST.RunPass(c.ctx, context.Emit)
	if err != nil {
		return nil, err
	}

	prog := c.ctx.Builder.Build()
	prog.ResultType = progAST.Type()

	return prog, nil
}

func PrintAST(input string, out io.Writer) error {
	progAST, err := parser.Parse(input)
	if err != nil {
		return err
	}
	context.NewGraphPrinter(out).PrintGraph(progAST)
	return nil
}

func runtimeTypeToSymbolType(t runtime.ValueType) (types.Type, error) {
	switch t {
	case runtime.Bool:
		return types.Bool, nil
	case runtime.Number:
		return types.Number, nil
	case runtime.String:
		return types.String, nil
	default:
		return nil, fmt.Errorf("unsupported type %v", t)
	}
}
