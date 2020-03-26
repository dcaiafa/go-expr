package runtime

import (
	"context"

	"github.com/dcaiafa/go-expr/expr/types"
)

type FuncFn func(ctx context.Context, args []Value) Value

type Func struct {
	Type types.Type
	Func FuncFn
}

type RawValue struct {
	num float64
	obj interface{}
}

func NewRawNumber(v float64) RawValue {
	return RawValue{num: v}
}

func NewRawBool(v bool) RawValue {
	num := float64(0)
	if v {
		num = 1
	}
	return NewRawNumber(num)
}

func NewRawObject(v interface{}) RawValue {
	return RawValue{obj: v}
}

func (v RawValue) Bool() bool {
	return v.num != 0
}

func (v RawValue) Number() float64 {
	return v.num
}

func (v RawValue) String() string {
	return v.obj.(string)
}

func (v RawValue) Object() interface{} {
	return v.obj
}

type Value struct {
	RawValue
	typ types.Type
}

func (v Value) Type() types.Type {
	return v.typ
}

func NewNumber(v float64) Value {
	return Value{typ: types.Number, RawValue: NewRawNumber(v)}
}

func NewBool(v bool) Value {
	num := float64(0)
	if v {
		num = 1
	}
	return Value{typ: types.Bool, RawValue: NewRawNumber(num)}
}

func NewObject(typ types.Type, v interface{}) Value {
	return Value{typ: typ, RawValue: NewRawObject(v)}
}

func NewString(v string) Value {
	return NewObject(types.String, v)
}
