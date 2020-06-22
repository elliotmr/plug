package main

import (
	"testing"

	"github.com/elliotmr/plug/example"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKVPlugin(t *testing.T) {
	c, err := example.Test(&kvImpl{m: make(map[string][]byte)})
	require.NoError(t, err)
	val, err := c.Get("test")
	require.Error(t, err)
	require.Nil(t, val)
	err = c.Put("test", []byte{1, 2, 3, 4})
	require.NoError(t, err)
	val, err = c.Get("test")
	require.NoError(t, err)
	assert.EqualValues(t, []byte{1, 2, 3, 4}, val)
}