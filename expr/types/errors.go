package types

import (
	"errors"
	"fmt"
)

var ErrUnknownSymbol = errors.New("symbol not found")

func NewErrUnknownSymbol(name string) error {
	return fmt.Errorf("%w: %v", ErrUnknownSymbol, name)
}
