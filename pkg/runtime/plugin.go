package runtime

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
	p := &plugin{
		g: genPlugin,
		t: &transport{
			in:  os.Stdin,
			out: os.Stdout,
			buf: make([]byte, os.Getpagesize()),
		},
	}
	p.run()
	return nil
}

func (p *plugin) run() {
	for {
		srv, buf, err := p.t.recv()
		if err != nil {
			p.t.sendError(err)
			continue
		}

		req, gsm, err := p.g.Link(srv)
		if err != nil {
			p.t.sendError(err)
			continue
		}

		err = proto.Unmarshal(buf, req)
		if err != nil {
			p.t.sendError(err)
			continue
		}

		resp, err := gsm(req)
		if err != nil {
			p.t.sendError(err)
			continue
		}

		err = p.t.send(srv, resp)
		if err != nil {
			p.t.sendError(err)
		}
	}
}