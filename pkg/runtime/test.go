package runtime

import (
	"fmt"
	"os"
)

// TestPlugin creates a communication pipes just like launching
// the plugin normally, but instead of running from a sub-process
// it simply runs the plugin in a go-routine.
func Test(genPlugin GenPlugin, magic string, version uint32, serviceBase Service) (*Host, error) {
	rRecv, lSend, err := os.Pipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create send pipe: %w", err)
	}
	lRecv, rSend, err := os.Pipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create recv pipe: %w", err)
	}

	h := &Host{
		t: &transport{
			in:  lRecv,
			out: lSend,
			buf: make([]byte, os.Getpagesize()),
		},
	}

	p := &plugin{
		g: genPlugin,
		t: &transport{
			in:  rRecv,
			out: rSend,
			buf: make([]byte, os.Getpagesize()),
		},
	}

	go p.run(magic, version, serviceBase)
	return h, h.handshake(magic, version, serviceBase)
}
