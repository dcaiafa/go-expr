package parser

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParser(t *testing.T) {
	_, err := Parse(`123+foorbar=="elephant"||_id < trob;butt&&less<=3`)
	require.NoError(t, err)
	_, err = Parse(`123+foorbar(hello() + 3425, "hi")*3`)
	require.NoError(t, err)
}
