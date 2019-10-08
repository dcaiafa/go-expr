package context

type Pass int

const (
	Print Pass = iota

	ResolveNames
	CheckTypes
	Emit
)

type PassRunner interface {
	RunPass(ctx *Context, pass Pass) error
}
