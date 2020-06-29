package runtime

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

// Service is used to enumerate the available services for the plugin
// transport. It is used by the generated code from `protoc-gen-plug`
// and should not be needed by a plugin consumer or author.
type Service uint16

type transport struct {
	in  *os.File
	out *os.File
	buf []byte
}

const (
	serviceError Service = 0xFF00
	serviceNone  Service = 0xFF01
)

func (t *transport) sendError(err error) {
	t.buf = t.buf[:0]
	t.buf = append(t.buf, byte(serviceError>>8), byte(serviceError&0xFF))
	t.buf = protowire.AppendString(t.buf, err.Error())
	_, err = t.out.Write(t.buf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
}

func (t *transport) send(s Service, msg proto.Message) error {
	var err error
	t.buf = t.buf[:0]
	t.buf = append(t.buf, byte(s>>8), byte(s&0xFF))
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
	sz, headerSize := protowire.ConsumeVarint(t.buf[2:])
	t.buf = t.buf[:n]
	headerSize += 2 // service type

	// in case the buffer wasn't large enough, we need to allocate more
	// space and read again.
	extraNeeded := int(sz) + headerSize - n
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

	srv := Service(t.buf[0])<<8 | Service(t.buf[1])
	if srv == serviceError {
		// sendError uses AppendString which uses the same varint length prefix
		return serviceError, nil, fmt.Errorf("remote error: %s", string(t.buf[headerSize:]))
	}

	return srv, t.buf[headerSize:], nil
}
