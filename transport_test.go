package plug

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestTransportError(t *testing.T) {
	recvr, sendl, err := os.Pipe()
	require.NoError(t, err)
	l := &transport{
		out: sendl,
		buf: make([]byte, os.Getpagesize()),
	}
	r := &transport{
		in: recvr,
		buf: make([]byte, os.Getpagesize()),
	}

	l.sendError(fmt.Errorf("test error"))
	srv, buf, err := r.recv()
	assert.Nil(t, buf)
	assert.Equal(t, serviceError, srv)
	assert.Error(t, err)
	assert.Equal(t, "remote error: test error", err.Error())

	sendDouble := &wrapperspb.DoubleValue{Value: 42.0}
	recvDouble := &wrapperspb.DoubleValue{}
	err = l.send(5, sendDouble)
	require.NoError(t, err)
	srv, buf, err = r.recv()
	require.NoError(t, err)
	assert.Equal(t, Service(5), srv)
	err = proto.Unmarshal(buf, recvDouble)
	require.NoError(t, err)
	assert.Equal(t, 42.0, recvDouble.Value)
}