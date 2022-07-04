package oogway

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestInit(t *testing.T) {
	assert.NoError(t, os.RemoveAll("test"))
	assert.NoError(t, Init("test"))
}
