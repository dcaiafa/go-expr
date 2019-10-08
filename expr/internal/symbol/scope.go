package symbol

import (
	"fmt"

	"github.com/dcaiafa/go-expr/expr/exprerrors"
)

type Scope struct {
	symbols map[string]Symbol
}

func NewScope() *Scope {
	return &Scope{
		symbols: make(map[string]Symbol),
	}
}

func (s *Scope) Add(sym Symbol) error {
	if _, ok := s.symbols[sym.Name()]; ok {
		return fmt.Errorf("Scope already has a symbol named %v", sym.Name())
	}
	s.symbols[sym.Name()] = sym
	return nil
}

func (s *Scope) Has(name string) bool {
	return s.symbols[name] != nil
}

func (s *Scope) Get(name string) (Symbol, error) {
	sym := s.symbols[name]
	if sym == nil {
		return nil, exprerrors.NewErrUnknownSymbol(name)
	}
	return sym, nil
}
