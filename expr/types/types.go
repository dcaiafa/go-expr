package types

import "fmt"

// Type is the type of values and symbols in an expression.
type Type interface {
	fmt.Stringer
	Equal(t Type) bool
}

type basicKind int

const (
	numberKind basicKind = iota
	stringKind
	boolKind
)

type basic struct {
	kind basicKind
}

var _ Type = (*basic)(nil)

func (b *basic) String() string {
	switch b.kind {
	case numberKind:
		return "number"
	case stringKind:
		return "string"
	case boolKind:
		return "bool"
	default:
		return "invalid"
	}
}

func (b *basic) Equal(other Type) bool {
	return b == other
}

// Basic types.
var (
	Number = &basic{numberKind}
	String = &basic{stringKind}
	Bool   = &basic{boolKind}
)

// Function is type of function symbols and values.
type Function struct {
	Params []Type
	Ret    Type
}

var _ Type = (*Function)(nil)

// Equal determines whether 'other' is the same type as this type.
func (f *Function) Equal(other Type) bool {
	otherFn, ok := other.(*Function)
	if !ok {
		return false
	}
	if !f.Ret.Equal(otherFn.Ret) {
		return false
	}
	if len(f.Params) != len(otherFn.Params) {
		return false
	}
	for i, arg := range f.Params {
		if !arg.Equal(otherFn.Params[i]) {
			return false
		}
	}
	return true
}

func (f *Function) String() string {
	return "function"
}
