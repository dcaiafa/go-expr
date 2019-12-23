package expr

import (
	"testing"

	"github.com/dcaiafa/go-expr/expr/runtime"
	"github.com/dcaiafa/go-expr/expr/types"
	"github.com/stretchr/testify/require"
)

func runExpr(t *testing.T, input string, args ...interface{}) {
	compiler := NewCompiler()

	var values []runtime.Value
	for len(args) > 2 {
		symbol := args[0].(string)
		value := args[1]
		var typ types.Type
		switch n := value.(type) {
		case float64:
			typ = types.Number
			values = append(values, runtime.NewNumberValue(n))
		case int:
			typ = types.Number
			values = append(values, runtime.NewNumberValue(float64(n)))
		case string:
			typ = types.String
			values = append(values, runtime.NewStringValue(n))
		case bool:
			typ = types.Bool
			values = append(values, runtime.NewBoolValue(n))
		default:
			panic("invalid type")
		}
		compiler.RegisterInput(symbol, typ)
		args = args[2:]
	}

	expected := args[0]
	if n, ok := expected.(int); ok {
		expected = float64(n)
	}

	prog, err := compiler.Compile(input)
	require.NoError(t, err)

	run := runtime.NewRuntime(prog)
	res, err := run.Run(0, values)
	require.NoError(t, err)

	switch prog.ResultType {
	case types.Number:
		require.Equal(t, types.Number, res.Type())
		require.Equal(t, expected, res.Number())
	case types.Bool:
		require.Equal(t, types.Bool, res.Type())
		require.Equal(t, expected, res.Bool())
	case types.String:
		require.Equal(t, types.String, res.Type())
		require.Equal(t, expected, res.String())
	default:
		panic("invalid result/expectation")
	}
}

func TestExpr_BinaryExpr_Vars(t *testing.T) {
	run := func(name, input string, args ...interface{}) {
		t.Run(name, func(t *testing.T) {
			runExpr(t, input, args...)
		})
	}
	run("lt", "a < b", "a", 1, "b", 2, true)
	run("lt2", "a < b", "a", 2, "b", 2, false)
	run("le", "a <= b", "a", 1, "b", 2, true)
	run("le2", "a <= b", "a", 2, "b", 2, true)
	run("le3", "a <= b", "a", 3, "b", 2, false)
	run("gt", "a > b", "a", 2, "b", 1, true)
	run("gt2", "a > b", "a", 2, "b", 2, false)
	run("ge", "a >= b", "a", 2, "b", 1, true)
	run("ge2", "a >= b", "a", 2, "b", 2, true)
	run("ge3", "a >= b", "a", 2, "b", 3, false)
	run("plus", "a + b", "a", 2, "b", 3, 5)
	run("minus", "a - b", "a", 5, "b", 2, 3)
	run("mult", "a * b", "a", 5, "b", 2, 10)
	run("div", "a / b", "a", 10, "b", 5, 2)
	run("eq_num", "a == b", "a", 2, "b", 2, true)
	run("eq_num2", "a == b", "a", 3, "b", 2, false)
	run("eq_bool", "a == b", "a", true, "b", true, true)
	run("eq_bool2", "a == b", "a", true, "b", false, false)
	run("eq_str", "a == b", "a", "foo", "b", "foo", true)
	run("eq_str2", "a == b", "a", "foo", "b", "bar", false)
}

func TestExpr_BinaryExpr_Literals(t *testing.T) {
	run := func(name, input string, args ...interface{}) {
		t.Run(name, func(t *testing.T) {
			runExpr(t, input, args...)
		})
	}
	run("lt", "1 < 2", true)
	run("lt2", "2 < 2", false)
	run("le", "1 <= 2", true)
	run("le2", "2 <= 2", true)
	run("le3", "3 <= 2", false)
	run("gt", "2 > 1", true)
	run("gt2", "2 > 2", false)
	run("ge", "2 >= 1", true)
	run("ge2", "2 >= 2", true)
	run("ge3", "2 >= 3", false)
	run("plus", "2 + 3", 5)
	run("minus", "5 - 2", 3)
	run("mult", "5 * 2", 10)
	run("div", "10 / 5", 2)
	run("eq_num", "2 == 2", true)
	run("eq_num2", "3 == 2", false)
	run("eq_str", `"foo" == "foo"`, true)
	run("eq_str2", `"foo" == "bar"`, false)
}

func TestExpr_AndExpr(t *testing.T) {
	run := func(name, input string, args ...interface{}) {
		t.Run(name, func(t *testing.T) {
			runExpr(t, input, args...)
		})
	}

	run("1", "a && b", "a", true, "b", true, true)
	run("2", "a && b", "a", true, "b", false, false)
	run("3", "a && b", "a", false, "b", true, false)
	run("4", "a && b", "a", false, "b", false, false)
	run("5", "a && b && c", "a", true, "b", true, "c", false, false)
	run("6", "a && b && c", "a", false, "b", true, "c", true, false)
}

func TestExpr_OrExpr(t *testing.T) {
	run := func(name, input string, args ...interface{}) {
		t.Run(name, func(t *testing.T) {
			runExpr(t, input, args...)
		})
	}

	run("1", "a || b", "a", true, "b", true, true)
	run("2", "a || b", "a", true, "b", false, true)
	run("3", "a || b", "a", false, "b", true, true)
	run("4", "a || b", "a", false, "b", false, false)
	run("5", "a || b || c", "a", true, "b", false, "c", false, true)
	run("6", "a || b || c", "a", false, "b", false, "c", false, false)
}

func TestExpr_NegateExpr(t *testing.T) {
	run := func(name, input string, args ...interface{}) {
		t.Run(name, func(t *testing.T) {
			runExpr(t, input, args...)
		})
	}

	run("1", "!a", "a", true, false)
	run("2", "!a", "a", false, true)
}

func TestExpr_LiteralExpr(t *testing.T) {
	run := func(name, input string, args ...interface{}) {
		t.Run(name, func(t *testing.T) {
			runExpr(t, input, args...)
		})
	}

	run("num", "3.14", 3.14)
	run("string", `"foobar"`, "foobar")
}

func TestExpr_Precedence(t *testing.T) {
	runExpr(t, "2*3 * (2 + 3) - 5", 25)
}

func TestExpr_Func_Basic(t *testing.T) {
	compiler := NewCompiler()

	compiler.RegisterFunc(
		"div",
		func(args []runtime.Value) runtime.Value {
			a := args[0].Number()
			b := args[1].Number()
			return runtime.NewNumberValue(a / b)
		},
		types.Number, types.Number, types.Number,
	)

	prog, err := compiler.Compile("div(15,3)")
	require.NoError(t, err)

	r := runtime.NewRuntime(prog)
	res, err := r.Run(0, nil)
	require.NoError(t, err)
	require.Equal(t, types.Number, res.Type())
	require.Equal(t, float64(5), res.Number())
}

func TestComplex1(t *testing.T) {
	compiler := NewCompiler()

	compiler.RegisterFunc(
		"len",
		func(args []runtime.Value) runtime.Value {
			a := args[0].String()
			return runtime.NewNumberValue(float64(len(a)))
		},
		types.Number, types.String,
	)
	compiler.RegisterInput("a", types.String)
	compiler.RegisterInput("b", types.Number)
	compiler.RegisterConst("c", runtime.NewNumberValue(3))

	prog, err := compiler.Compile("len(a) + c == b")
	require.NoError(t, err)

	r := runtime.NewRuntime(prog)

	args := []runtime.Value{
		runtime.NewStringValue("hello"),
		runtime.NewNumberValue(8),
	}

	res, err := r.Run(0, args)
	require.NoError(t, err)
	require.True(t, res.Bool())
}

func Benchmark1(b *testing.B) {
	compiler := NewCompiler()

	compiler.RegisterFunc(
		"len",
		func(args []runtime.Value) runtime.Value {
			a := args[0].String()
			return runtime.NewNumberValue(float64(len(a)))
		},
		types.String, types.Number,
	)
	compiler.RegisterInput("a", types.String)
	compiler.RegisterInput("b", types.Number)
	compiler.RegisterConst("c", runtime.NewNumberValue(3))

	prog, err := compiler.Compile("len(a) + c == b")
	require.NoError(b, err)

	r := runtime.NewRuntime(prog)

	args := []runtime.Value{
		runtime.NewStringValue("hello"),
		runtime.NewNumberValue(8),
	}

	for i := 0; i < b.N; i++ {
		res, err := r.Run(0, args)
		if err != nil {
			panic(err)
		}
		if !res.Bool() {
			panic("unexpected result")
		}
	}
}
