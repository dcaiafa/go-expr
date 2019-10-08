%{
package parser

import (
  "github.com/dcaiafa/go-expr/expr/internal/ast"
  "github.com/dcaiafa/go-expr/expr/types"
)

%}

%union {
  num float64
  str string
  ast ast.AST
  expr ast.Expr
}

%token LEXERR
%token ID
%token <num> NUMBER
%token <str> STRING
%token <str> ID

%type <ast> exprs opt_params params
%type <expr> expr binary_expr unary_expr term invocation number

%left OR
%left AND
%nonassoc '<' LE '>' GE EQ
%left '+' '-'
%left '*' '/'

%start program

%%

program: exprs { yylex.(*lex).Program = $1.(*ast.Program) }

exprs: exprs ';' expr  { $1.(*ast.Program).AddExpr($3.(ast.Expr)) }
     | expr            { $$ = ast.NewProgram($1.(ast.Expr)) }

expr: binary_expr 

binary_expr: unary_expr
           | binary_expr AND binary_expr  { $$ = ast.NewAndExpr($1, $3) }
           | binary_expr OR binary_expr   { $$ = ast.NewOrExpr($1, $3) }
           | binary_expr '<' binary_expr  { $$ = ast.NewBinaryExpr($1, ast.Lt, $3) }
           | binary_expr LE binary_expr   { $$ = ast.NewBinaryExpr($1, ast.Le, $3) }
           | binary_expr '>' binary_expr  { $$ = ast.NewBinaryExpr($1, ast.Gt, $3) }
           | binary_expr GE binary_expr   { $$ = ast.NewBinaryExpr($1, ast.Ge, $3) }
           | binary_expr EQ binary_expr   { $$ = ast.NewBinaryExpr($1, ast.Eq, $3) }
           | binary_expr '+' binary_expr  { $$ = ast.NewBinaryExpr($1, ast.Plus, $3) }
           | binary_expr '-' binary_expr  { $$ = ast.NewBinaryExpr($1, ast.Minus, $3) }
           | binary_expr '*' binary_expr  { $$ = ast.NewBinaryExpr($1, ast.Times, $3) }
           | binary_expr '/' binary_expr  { $$ = ast.NewBinaryExpr($1, ast.Div, $3) }

unary_expr: '!' term { $$ = ast.NewNegateExpr($2) }
          | term    

term: number
    | STRING       { $$ = ast.NewLiteralExpr(types.String, $1) }
    | ID           { $$ = ast.NewSimpleRefExpr($1) }
    | '(' expr ')' { $$ = $2 }
    | invocation

number: '-' NUMBER { $$ = ast.NewLiteralExpr(types.Number, -$2) } 
      | NUMBER     { $$ = ast.NewLiteralExpr(types.Number,  $1) } 

invocation: term '(' opt_params ')' { $$ = ast.NewCallExpr($1, $3.(*ast.Params)) }

opt_params: params
          |         { $$ = &ast.Params{} }

params: params ',' expr  { $1.(*ast.Params).AddParam($3.(ast.Expr)); $$ = $1 }
      | expr             { $$ = ast.NewParams($1) }
