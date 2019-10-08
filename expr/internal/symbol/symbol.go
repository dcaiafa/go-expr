package symbol

import (
	"github.com/dcaiafa/go-expr/expr/runtime"
	"github.com/dcaiafa/go-expr/expr/types"
)

type Symbol interface {
	Name() string
	Type() types.Type
	EmitAccess(builder *runtime.Builder)
}

type symbolImpl struct {
	name string
	typ  types.Type
}

func (s *symbolImpl) Name() string {
	return s.name
}

func (s *symbolImpl) Type() types.Type {
	return s.typ
}

type ConstSymbol struct {
	symbolImpl
	constIndex int
}

func NewConstSymbol(name string, typ types.Type, constIndex int) *ConstSymbol {
	return &ConstSymbol{
		symbolImpl: symbolImpl{
			name: name,
			typ:  typ,
		},
		constIndex: constIndex,
	}
}

func (s *ConstSymbol) EmitAccess(builder *runtime.Builder) {
	builder.EmitLoadConst(s.constIndex)
}

type InputSymbol struct {
	symbolImpl
	inputIndex int
}

func NewInputSymbol(name string, typ types.Type, inputIndex int) *InputSymbol {
	return &InputSymbol{
		symbolImpl: symbolImpl{
			name: name,
			typ:  typ,
		},
		inputIndex: inputIndex,
	}
}

func (s *InputSymbol) EmitAccess(builder *runtime.Builder) {
	builder.EmitLoadInput(s.inputIndex)
}
