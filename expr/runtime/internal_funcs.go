package runtime

import "context"

const (
	InternalInStringArray int = iota
	InternalInNumberArray
)

func internalInStringArray(ctx context.Context, args []Value) Value {
	left := args[0].String()
	right := args[1].Object().([]RawValue)
	for _, elem := range right {
		if left == elem.String() {
			return NewBool(true)
		}
	}
	return NewBool(false)
}

func internalInNumberArray(ctx context.Context, args []Value) Value {
	left := args[0].Number()
	right := args[1].Object().([]RawValue)
	for _, elem := range right {
		if left == elem.Number() {
			return NewBool(true)
		}
	}
	return NewBool(false)
}
