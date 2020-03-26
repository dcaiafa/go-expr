package runtime

import (
	"context"
	"fmt"

	"github.com/dcaiafa/go-expr/expr/types"
)

type Operation int

const (
	PushNumber Operation = iota
	PushString
	PushBool
	PushArray
	PushValue
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
	strings []string
	consts  []Value
	inputs  []types.Type
}

func (p *Program) ExprCount() int {
	return len(p.exprs)
}

type Runtime struct {
	program  *Program
	stack    []RawValue
	callArgs []Value
}

func NewRuntime(program *Program) *Runtime {
	return &Runtime{
		program: program,
		stack:   make([]RawValue, 0, 30),
	}
}

func (r *Runtime) Run(ctx context.Context, exprIndex int, inputs []Value) (Value, error) {
	r.stack = r.stack[:0]

	if len(inputs) != len(r.program.inputs) {
		return Value{}, fmt.Errorf(
			"program expects %d inputs but %d were provided",
			len(r.program.inputs), len(inputs))
	}
	for i, input := range inputs {
		if !input.Type().Equal(r.program.inputs[i]) {
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
			r.push(NewRawNumber(instr.vnum))
		case PushString:
			r.push(NewRawObject(r.program.strings[instr.extra]))
		case PushBool:
			r.push(NewRawBool(instr.extra != 0))
		case PushArray:
			elemCount := instr.extra
			arr := make([]RawValue, elemCount)
			for i := range arr {
				arr[len(arr)-i-1] = r.pop()
			}
			r.push(NewRawObject(arr))
		case LoadConst:
			r.push(r.program.consts[instr.extra].RawValue)
		case LoadInput:
			r.push(inputs[instr.extra].RawValue)
		case Duplicate:
			r.push(r.peek())
		case Add:
			right, left := r.pop(), r.pop()
			r.push(NewRawNumber(left.num + right.num))
		case Subtract:
			right, left := r.pop(), r.pop()
			r.push(NewRawNumber(left.num - right.num))
		case Multiply:
			right, left := r.pop(), r.pop()
			r.push(NewRawNumber(left.num * right.num))
		case Divide:
			right, left := r.pop(), r.pop()
			r.push(NewRawNumber(left.num / right.num))
		case Negate:
			v := r.pop()
			r.push(NewRawBool(!v.Bool()))
		case And:
			right, left := r.pop(), r.pop()
			r.push(NewRawBool(left.Bool() && right.Bool()))
		case Or:
			right, left := r.pop(), r.pop()
			r.push(NewRawBool(left.Bool() || right.Bool()))
		case CompareEqBool:
			right, left := r.pop(), r.pop()
			r.push(NewRawBool(left.Bool() == right.Bool()))
		case CompareEqString:
			right, left := r.pop(), r.pop()
			r.push(NewRawBool(left.String() == right.String()))
		case CompareEq:
			right, left := r.pop(), r.pop()
			r.push(NewRawBool(left.num == right.num))
		case CompareLT:
			right, left := r.pop(), r.pop()
			r.push(NewRawBool(left.num < right.num))
		case CompareLE:
			right, left := r.pop(), r.pop()
			r.push(NewRawBool(left.num <= right.num))
		case CompareGT:
			right, left := r.pop(), r.pop()
			r.push(NewRawBool(left.num > right.num))
		case CompareGE:
			right, left := r.pop(), r.pop()
			r.push(NewRawBool(left.num >= right.num))
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
				r.callArgs[i].RawValue = r.pop()
			}
			fn := r.pop().Object().(*Func)
			fnType := fn.Type.(*types.Function)
			for i := range r.callArgs {
				r.callArgs[i].typ = fnType.Params[i]
			}
			res := fn.Func(ctx, r.callArgs)
			if !res.Type().Equal(fnType.Ret) {
				return Value{}, fmt.Errorf("function returned %v expected %v",
					res.Type(), fnType.Ret)
			}
			r.push(res.RawValue)
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

	return Value{typ: r.program.ResultType, RawValue: r.stack[0]}, nil
}

func (r *Runtime) push(v RawValue) {
	r.stack = append(r.stack, v)
}

func (r *Runtime) pop() RawValue {
	v := r.stack[len(r.stack)-1]
	r.stack = r.stack[:len(r.stack)-1]
	return v
}

func (r *Runtime) peek() RawValue {
	return r.stack[len(r.stack)-1]
}
