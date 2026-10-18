package client

import (
	"context"
	"os/exec"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/xgodev/boost/model/errors"
)

// Plugin customizes the client after it is built and before it connects.
type Plugin func(context.Context, *Client) error

// Client wraps an *mcp.Client and the session it opened on the configured server.
type Client struct {
	instance *mcp.Client
	session  *mcp.ClientSession
	options  *Options
}

// NewClient connects a client from the default config path.
func NewClient(ctx context.Context, plugins ...Plugin) (*Client, error) {
	o, err := NewOptions()
	if err != nil {
		return nil, err
	}
	return NewClientWithOptions(ctx, o, plugins...)
}

// NewClientWithConfigPath connects a client from the given config path.
func NewClientWithConfigPath(ctx context.Context, path string, plugins ...Plugin) (*Client, error) {
	o, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewClientWithOptions(ctx, o, plugins...)
}

// NewClientWithOptions builds a client from o, runs the plugins and connects.
func NewClientWithOptions(ctx context.Context, o *Options, plugins ...Plugin) (*Client, error) {
	t, err := newTransport(o)
	if err != nil {
		return nil, err
	}

	cli := &Client{
		instance: mcp.NewClient(
			&mcp.Implementation{Name: o.Name, Version: o.Version},
			&mcp.ClientOptions{KeepAlive: o.KeepAlive},
		),
		options: o,
	}

	for _, plugin := range plugins {
		if err := plugin(ctx, cli); err != nil {
			return nil, err
		}
	}

	cli.session, err = cli.instance.Connect(ctx, t, nil)
	if err != nil {
		return nil, errors.Wrap(err, errors.Internalf("mcp client: connect via %s", o.Transport))
	}

	return cli, nil
}

func newTransport(o *Options) (mcp.Transport, error) {
	o.Transport = strings.ToUpper(o.Transport)
	switch o.Transport {
	case TransportHTTP:
		// TODO(maintainer-review): HTTPClient stays nil (http.DefaultClient). A
		// global client timeout would cut the long-lived SSE stream; per-call
		// deadlines belong on the ctx passed to CallTool.
		return &mcp.StreamableClientTransport{Endpoint: o.Endpoint, MaxRetries: o.MaxRetries}, nil
	case TransportCOMMAND:
		if o.Command == "" {
			return nil, errors.NotValidf("mcp client: COMMAND transport needs a command")
		}
		return &mcp.CommandTransport{Command: exec.Command(o.Command, o.Args...)}, nil
	default:
		return nil, errors.NotValidf("mcp client transport %q: want %s or %s", o.Transport, TransportHTTP, TransportCOMMAND)
	}
}

// Instance returns the wrapped client.
func (c *Client) Instance() *mcp.Client {
	return c.instance
}

// Session returns the connected session; call tools and list capabilities on it.
func (c *Client) Session() *mcp.ClientSession {
	return c.session
}

// Options returns the options the client was built with.
func (c *Client) Options() *Options {
	return c.options
}

// Close ends the session (and stops the subprocess on COMMAND).
func (c *Client) Close() error {
	return c.session.Close()
}
