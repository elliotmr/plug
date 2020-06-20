package plug

import (
	"fmt"
	"os"

	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/proto"
)

/*
 * Notes
 * - send and recv share a buffer, if possible it will be reused without
 *   it will be grown as required.
 * - the only reason why the transport should return a service byte that
 *   doesn't match
 * - the transport must be locked externally
 * - transport should be allocated with a minimum of 11 bytes, probably
 *   it should have at least a page allocated to start
 */

type Service byte

type transport struct {
	in  *os.File
	out *os.File
	buf []byte
}

const (
	serviceError Service = 254
	serviceNone  Service = 255
)


func (t *transport) sendError(err error) {
	t.buf = t.buf[:0]
	t.buf = append(t.buf, byte(serviceError))
	t.buf = protowire.AppendString(t.buf, err.Error())
	_, err = t.out.Write(t.buf)
	if err != nil {
		panic(err)
	}
}

func (t *transport) send(s Service, msg proto.Message) error {
	var err error
	t.buf = t.buf[:0]
	t.buf = append(t.buf, byte(s))
	opts := proto.MarshalOptions{}
	t.buf = protowire.AppendVarint(t.buf, uint64(opts.Size(msg)))
	opts.UseCachedSize = true
	t.buf, err = opts.MarshalAppend(t.buf, msg)
	if err != nil {
		return fmt.Errorf("failed proto marshal: %w", err)
	}
	_, err = t.out.Write(t.buf)
	return err
}

func (t *transport) recv() (Service, []byte, error) {
	t.buf = t.buf[:cap(t.buf)]
	n, err := t.in.Read(t.buf)
	if err != nil {
		return serviceNone, nil, fmt.Errorf("failed to read in pipe: %w", err)
	}
	sz, m := protowire.ConsumeVarint(t.buf[1:])
	t.buf = t.buf[:n]

	// in case the buffer wasn't large enough, we need to allocate more
	// space and read again.
	extraNeeded := int(sz) + m + 1 - n
	if extraNeeded > 0 {
		extra := make([]byte, extraNeeded)
		o, err := t.in.Read(extra)
		if o != extraNeeded {
			return serviceNone, nil, fmt.Errorf("incorrect size for second read")
		}
		if err != nil {
			return serviceNone, nil, fmt.Errorf("failed to read extra bytes: %w", err)
		}
		t.buf = append(t.buf, extra...)
	}

	if Service(t.buf[0]) == serviceError {
		// sendError uses AppendString which uses the same varint length prefix
		return serviceError, nil, fmt.Errorf("remote error: %s", string(t.buf[1+m:]))
	}

	return Service(t.buf[0]), t.buf[m+1:], nil
}