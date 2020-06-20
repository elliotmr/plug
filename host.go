package plug

import (
	"fmt"
	"os"
	"os/exec"
	"sync"

	"google.golang.org/protobuf/proto"
)

type Host struct {
	mu sync.Mutex
	t *transport
	cmd *exec.Cmd
}

func LaunchPlugin(filename string) (*Host, error) {
	rRecv, lSend, err := os.Pipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create send pipe: %w", err)
	}
	lRecv, rSend, err := os.Pipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create recv pipe: %w", err)
	}
	
	cmd := &exec.Cmd{
		Path:   filename,
		Stdin:  rRecv,
		Stdout: rSend,
	}
	
	h := &Host{
		t: &transport{
			in:  lRecv,
			out: lSend,
			buf: make([]byte, os.Getpagesize()),
		},
		cmd: cmd,
	}
	
	return h, cmd.Start()
}

func (c *Host) SendRecv(srv Service, req, resp proto.Message) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	err := c.t.send(srv, req)
	if err != nil {
		return fmt.Errorf("send error: %w", err)
	}
	srvResp, buf, err := c.t.recv()
	if err != nil {
		return fmt.Errorf("recv error: %w", err)
	}
	if srvResp != srv {
		return fmt.Errorf("recv service mismatch: %d != %d", srvResp, srv)
	}
	return proto.Unmarshal(buf, resp)
}
