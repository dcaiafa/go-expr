package parser

import "github.com/dcaiafa/go-expr/expr/internal/ast"

//go:generate goyacc parser.y

func Parse(input string) (*ast.Program, error) {
	//yyDebug = 10
	yyErrorVerbose = true
	l := newLex(input)
	p := yyNewParser()
	p.Parse(l)
	return l.Program, l.err
}
