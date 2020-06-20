package plug

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"os"
)

// GenPluginMethod is used internally to link plugin methods from
// protoc-gen-plug to the plugin runtime.
type GenPluginMethod func(proto.Message) (proto.Message, error)

// GenPlugin is a generated plugin instance from protoc-gen-plug.
type GenPlugin interface {
	Link(srv Service) (proto.Message, GenPluginMethod, error)
}

type plugin struct {
	g GenPlugin
	t *transport
}

// Run is the entrypoint for the plugin implementation. You must
// pass your implementation wrapped in your generated plugin to
// this function. Run will block forever.
func Run(genPlugin GenPlugin) error {
	fiIn, err := os.Stdin.Stat()
	if err != nil {
		return fmt.Errorf("unable to stat stdin: %w", err)
	}
	fiOut, err := os.Stdout.Stat()
	if err != nil {
		return fmt.Errorf("unable to stat stdout: %w", err)
	}
	if fiIn.Mode()&fiOut.Mode()&os.ModeNamedPipe == 0 {
		return fmt.Errorf("both stdin and stdout must be pipes")
	}
	s := &plugin{
		g: genPlugin,
		t: &transport{
			in:  os.Stdin,
			out: os.Stdout,
			buf: make([]byte, os.Getpagesize()),
		},
	}

	for {
		srv, buf, err := s.t.recv()
		if err != nil {
			s.t.sendError(err)
			continue
		}

		req, gsm, err := s.g.Link(srv)
		if err != nil {
			s.t.sendError(err)
			continue
		}

		err = proto.Unmarshal(buf, req)
		if err != nil {
			s.t.sendError(err)
			continue
		}

		resp, err := gsm(req)
		if err != nil {
			s.t.sendError(err)
			continue
		}

		err = s.t.send(srv, resp)
		if err != nil {
			s.t.sendError(err)
		}
	}
}
