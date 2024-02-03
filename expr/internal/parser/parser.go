package parser

import "github.com/dcaiafa/go-expr/expr/internal/ast"

//go:generate goyacc parser.y

func init() {
	// Enable extra details in the error message returned by the parser.
	yyErrorVerbose = true
}

func Parse(input string) (*ast.Program, error) {
	l := newLex(input)
	p := yyNewParser()
	p.Parse(l)
	return l.Program, l.err
}
