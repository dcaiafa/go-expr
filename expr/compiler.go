package expr

import (
	"io"

	"github.com/dcaiafa/go-expr/expr/internal/context"
	"github.com/dcaiafa/go-expr/expr/internal/parser"
	"github.com/dcaiafa/go-expr/expr/internal/symbol"
	"github.com/dcaiafa/go-expr/expr/runtime"
	"github.com/dcaiafa/go-expr/expr/types"
)

// Compiler compiles an expression to produce a runtime.Program, to be
// executed in a runtime.Runtime.
type Compiler struct {
	ctx *context.Context
}

// NewCompiler creates a new Compiler.
func NewCompiler() *Compiler {
	return &Compiler{
		ctx: context.NewContext(),
	}
}

// RegisterInput creates an input parameter that can be used in the expression.
func (c *Compiler) RegisterInput(name string, typ types.Type) (int, error) {
	inputIndex := c.ctx.Builder.NewInput(typ)
	inputSymbol := symbol.NewInputSymbol(name, typ, inputIndex)
	err := c.ctx.GlobalScope.Add(inputSymbol)
	if err != nil {
		return 0, err
	}
	return inputIndex, nil
}

// RegisterConst creates a named constant that can be used in the expression.
func (c *Compiler) RegisterConst(name string, v runtime.Value) error {
	constIndex := c.ctx.Builder.NewConst(v)
	constSymbol := symbol.NewConstSymbol(name, v.Type(), constIndex)
	return c.ctx.GlobalScope.Add(constSymbol)
}

// RegisterFunc registers a function that can be used in the expression. The
// function implementation *must* return a value with the type specified at
// registration.
func (c *Compiler) RegisterFunc(
	name string,
	fn runtime.FuncFn,
	ret types.Type,
	args ...types.Type,
) error {
	fnType := &types.Function{
		Params: make([]types.Type, len(args)),
		Ret:    ret,
	}
	copy(fnType.Params, args)
	return c.registerFunc(name, fnType, fn)
}

func (c *Compiler) registerFunc(
	name string,
	typ *types.Function,
	fn runtime.FuncFn,
) error {
	v := runtime.NewObject(typ, &runtime.Func{
		Type: typ,
		Func: fn,
	})
	constIndex := c.ctx.Builder.NewConst(v)
	fnSymbol := symbol.NewConstSymbol(name, typ, constIndex)
	return c.ctx.GlobalScope.Add(fnSymbol)
}

// Compile compiles an expression into a program.
func (c *Compiler) Compile(expr string) (*runtime.Program, error) {
	progAST, err := parser.Parse(expr)
	if err != nil {
		return nil, err
	}

	for _, pass := range context.Passes {
		err = progAST.RunPass(c.ctx, pass)
		if err != nil {
			return nil, err
		}
	}

	prog := c.ctx.Builder.Build()
	prog.ResultType = progAST.Type()
	return prog, nil
}

// PrintAST parses an expression and prints the AST in Graphviz dot format.
func PrintAST(input string, out io.Writer) error {
	progAST, err := parser.Parse(input)
	if err != nil {
		return err
	}
	context.NewGraphPrinter(out).PrintGraph(progAST)
	return nil
}
