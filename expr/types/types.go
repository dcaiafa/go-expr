package types

type BasicKind int

const (
	NumberKind BasicKind = iota
	StringKind
	BoolKind
)

type Type interface {
	String() string
	Equal(t Type) bool
}

type Basic struct {
	kind BasicKind
}

func (b *Basic) String() string {
	switch b.kind {
	case NumberKind:
		return "number"
	case StringKind:
		return "string"
	case BoolKind:
		return "bool"
	default:
		return "invalid"
	}
}

func (t *Basic) Equal(other Type) bool {
	return t == other
}

var (
	Number = &Basic{NumberKind}
	String = &Basic{StringKind}
	Bool   = &Basic{BoolKind}
)

type Function struct {
	Args []Type
	Ret  Type
}

func (f *Function) Equal(other Type) bool {
	otherFn, ok := other.(*Function)
	if !ok {
		return false
	}
	if !f.Ret.Equal(otherFn.Ret) {
		return false
	}
	if len(f.Args) != len(otherFn.Args) {
		return false
	}
	for i, arg := range f.Args {
		if !arg.Equal(otherFn.Args[i]) {
			return false
		}
	}
	return true
}

func (f *Function) String() string {
	return "function"
}
