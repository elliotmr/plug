package main

import (
	"testing"

	"github.com/elliotmr/plug/example"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestSerializerPlugin(t *testing.T) {
	c, err := example.TestSerializer(&serializer{})
	require.NoError(t, err)
	ex := make(map[string]*anypb.Any)
	ex["question"], err = anypb.New(wrapperspb.String("what is the meaning of life?"))
	require.NoError(t, err)
	ex["answer"], err = anypb.New(wrapperspb.UInt32(42))
	require.NoError(t, err)
	out, err := c.Marshal(ex)
	require.NoError(t, err)
	require.Equal(
		t,
		`{"answer":{"value":42},"question":{"value":"what is the meaning of life?"}}`,
		string(out),
	)
}
