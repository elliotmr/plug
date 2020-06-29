package main

import (
	"github.com/elliotmr/plug/example"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"testing"
)

func TestSerializerPlugin(t *testing.T) {
	c, err := example.TestSerializer(&serializer{})
	require.NoError(t, err)
	ex := make(map[string]*anypb.Any)
	ex["question"], err = anypb.New(wrapperspb.String("what is the meaning of life?"))
	ex["answer"], err = anypb.New(wrapperspb.UInt32(42))
	out, err := c.Marshal(ex)
	require.NoError(t, err)
	require.Equal(
		t,
		`{"answer":{"value":42},"question":{"value":"what is the meaning of life?"}}`,
		string(out),
	)
}
