package glue

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGlue(t *testing.T) {
	assert.Equal(t, GetGlue("gluetest"), "Hello world!\n")
}
