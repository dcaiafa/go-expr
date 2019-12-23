package runtime

import (
	"fmt"

	"github.com/dcaiafa/go-expr/expr/types"
)

type Operation int

const (
	PushNumber Operation = iota
	PushString
	PushBool
	LoadConst
	LoadInput
	Duplicate
	Add
	Subtract
	Multiply
	Divide
	Negate
	And
	Or
	CompareEqBool
	CompareEqString
	CompareEq
	CompareLT
	CompareLE
	CompareGT
	CompareGE
	Jump
	JumpIfTrue
	JumpIfFalse
	Call
	Return
)

type Instruction struct {
	op    Operation
	extra int
	vnum  float64
}

type Label struct {
	index int
	addr  int
}

type Expr []Instruction

type Program struct {
	ResultType types.Type

	exprs   []Expr
	funcs   []Func
	strings []string
	consts  []Value
	inputs  []types.Type
}

func (p *Program) ExprCount() int {
	return len(p.exprs)
}

type Runtime struct {
	program  *Program
	stack    []Value
	callArgs []Value
}

func NewRuntime(program *Program) *Runtime {
	return &Runtime{
		program: program,
		stack:   make([]Value, 0, 30),
	}
}

func (r *Runtime) Run(exprIndex int, inputs []Value) (Value, error) {
	r.stack = r.stack[:0]

	if len(inputs) != len(r.program.inputs) {
		return Value{}, fmt.Errorf(
			"program expects %d inputs but %d were provided",
			len(r.program.inputs), len(inputs))
	}
	for i, input := range inputs {
		if input.Type() != r.program.inputs[i] {
			return Value{}, fmt.Errorf(
				"program expects input index %d type %v but %v was provided",
				i, r.program.inputs[i], input.Type())
		}
	}

	exprInstr := r.program.exprs[exprIndex]

Loop:
	for n := 0; n < len(exprInstr); {
		instr := exprInstr[n]
		switch instr.op {
		case PushNumber:
			r.push(NewNumberValue(instr.vnum))
		case PushString:
			r.push(NewStringValue(r.program.strings[instr.extra]))
		case PushBool:
			v := false
			if instr.extra != 0 {
				v = true
			}
			r.push(NewBoolValue(v))
		case LoadConst:
			r.push(r.program.consts[instr.extra])
		case LoadInput:
			r.push(inputs[instr.extra])
		case Duplicate:
			r.push(r.peek())
		case Add, Subtract, Multiply, Divide:
			right := r.pop()
			left := r.pop()
			var res float64
			switch instr.op {
			case Add:
				res = left.num + right.num
			case Subtract:
				res = left.num - right.num
			case Multiply:
				res = left.num * right.num
			case Divide:
				res = left.num / right.num
			}
			r.push(NewNumberValue(res))
		case Negate:
			v := r.pop()
			r.push(NewBoolValue(!v.Bool()))
		case And:
			right := r.pop()
			left := r.pop()
			r.push(NewBoolValue(left.Bool() && right.Bool()))
		case Or:
			right := r.pop()
			left := r.pop()
			r.push(NewBoolValue(left.Bool() || right.Bool()))
		case CompareEqBool:
			right := r.pop()
			left := r.pop()
			r.push(NewBoolValue(left.Bool() == right.Bool()))
		case CompareEqString:
			right := r.pop()
			left := r.pop()
			r.push(NewBoolValue(left.String() == right.String()))
		case CompareEq, CompareLT, CompareLE, CompareGT, CompareGE:
			right := r.pop()
			left := r.pop()
			var res bool
			switch instr.op {
			case CompareEq:
				res = left.num == right.num
			case CompareLT:
				res = left.num < right.num
			case CompareLE:
				res = left.num <= right.num
			case CompareGT:
				res = left.num > right.num
			case CompareGE:
				res = left.num >= right.num
			}
			r.push(NewBoolValue(res))
		case Jump:
			n = instr.extra
			continue
		case JumpIfTrue:
			boolValue := r.pop()
			if boolValue.Bool() {
				n = instr.extra
				continue
			}
		case JumpIfFalse:
			boolValue := r.pop()
			if !boolValue.Bool() {
				n = instr.extra
				continue
			}
		case Call:
			argCount := instr.extra
			if cap(r.callArgs) < argCount {
				r.callArgs = make([]Value, argCount)
			} else {
				r.callArgs = r.callArgs[:argCount]
			}
			for i := argCount - 1; i >= 0; i-- {
				r.callArgs[i] = r.pop()
			}
			fnValue := r.pop()
			fn := r.program.funcs[fnValue.ExternalFunc()]
			res := fn(r.callArgs)
			r.push(res)

		case Return:
			break Loop
		default:
			panic("invalid op")
		}

		n++
	}
	if len(r.stack) != 1 {
		return Value{}, fmt.Errorf("invalid program: after execution stack len: %d",
			len(r.stack))
	}
	return r.stack[0], nil
}

func (r *Runtime) push(v Value) {
	r.stack = append(r.stack, v)
}

func (r *Runtime) pop() Value {
	v := r.stack[len(r.stack)-1]
	r.stack = r.stack[:len(r.stack)-1]
	return v
}

func (r *Runtime) peek() Value {
	return r.stack[len(r.stack)-1]
}
