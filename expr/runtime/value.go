package runtime

type Func func(args []Value) Value

type ValueType byte

const (
	Bool ValueType = iota
	Number
	String
	ExternalFunc
)

type Value struct {
	typ   ValueType
	other int
	num   float64
	str   string
}

func NewBoolValue(v bool) Value {
	vint := 0
	if v {
		vint = 1
	}
	return Value{typ: Bool, other: vint}
}

func NewNumberValue(v float64) Value {
	return Value{typ: Number, num: v}
}

func NewStringValue(v string) Value {
	return Value{typ: String, str: v}
}

func NewExternalFuncValue(fnIndex int) Value {
	return Value{typ: ExternalFunc, other: fnIndex}
}

func (v Value) Type() ValueType {
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
