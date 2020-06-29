package runtime

import (
	"fmt"
	"os"

	"github.com/elliotmr/plug/pkg/plugpb"
	"google.golang.org/protobuf/proto"
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

// Run is the entrypoint for the plugin implementation. It will
// be wrapped by the plug generated Run function.
func Run(genPlugin GenPlugin, magic string, version uint32, serviceBase Service) error {
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
	return p.run(magic, version, serviceBase)
}

func (p *plugin) run(magic string, version uint32, serviceBase Service) error {
	srv, buf, err := p.t.recv()
	if err != nil {
		err = fmt.Errorf("received error instead of handshake: %w", err)
		p.t.sendError(err)
		return err
	}
	if srv != serviceBase<<8|0xFF {
		err = fmt.Errorf("first message was not a handshake")
		p.t.sendError(err)
		return err
	}

	req := &plugpb.Handshake{}
	resp := &plugpb.Handshake{
		Version: version,
		Magic:   magic,
	}
	err = proto.Unmarshal(buf, req)
	if err != nil {
		err = fmt.Errorf("failed to unmarshal handshake request: %w", err)
		p.t.sendError(err)
		return err
	}

	if resp.Magic != req.Magic {
		err = fmt.Errorf("magic mismatch in handshake")
		p.t.sendError(err)
		return err
	}
	if resp.Version != req.Version {
		err = fmt.Errorf("version mismatch in handshake (host: %d) (plugin: %d)", req.Version, resp.Version)
		p.t.sendError(err)
		return err
	}

	err = p.t.send(srv, resp)
	if err != nil {
		p.t.sendError(err)
	}

	for {
		srv, buf, err = p.t.recv()
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
