package main

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/elliotmr/plug/example"
	"github.com/stretchr/testify/require"
)

func TestKVPlugin(t *testing.T) {
	c, err := example.Load("kv.exe")
	require.NoError(t, err)
	val, err := c.Get("test")
	require.Error(t, err)
	require.Nil(t, val)
	err = c.Put("test", []byte{1,2,3,4})
	require.NoError(t, err)
	val, err = c.Get("test")
	require.NoError(t, err)
	assert.EqualValues(t, []byte{1,2,3,4}, val)
}

func BenchmarkKVPlugin(b *testing.B) {
	c, err := example.Load("kv.exe")
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		err = c.Put("test", []byte{1,2,3,4})
		if err != nil {
			b.Fatal(err)
		}
	}
}