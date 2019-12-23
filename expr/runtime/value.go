package runtime

import "github.com/dcaiafa/go-expr/expr/types"

type Func func(args []Value) Value

type Value struct {
	typ   types.Type
	other int
	num   float64
	str   string
}

func NewBoolValue(v bool) Value {
	vint := 0
	if v {
		vint = 1
	}
	return Value{typ: types.Bool, other: vint}
}

func NewNumberValue(v float64) Value {
	return Value{typ: types.Number, num: v}
}

func NewStringValue(v string) Value {
	return Value{typ: types.String, str: v}
}

func NewExternalFuncValue(typ types.Type, fnIndex int) Value {
	return Value{typ: typ, other: fnIndex}
}

func (v Value) Type() types.Type {
	return v.typ
}

func (v Value) Bool() bool {
	return v.other != 0
}

func (v Value) Number() float64 {
	return v.num
}

func (v Value) String() string {
	return v.str
}

func (v Value) ExternalFunc() int {
	return v.other
}
