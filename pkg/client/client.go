package client

import (
	"github.com/xcoulon/converse-mcp/pkg/channel"

	"github.com/creachadair/jrpc2"
	"github.com/creachadair/jrpc2/jhttp"
)

type Client struct {
	*jrpc2.Client
}

func NewFromChannel(c channel.Channel) *Client {
	return &Client{
		Client: jrpc2.NewClient(c, nil),
	}
}

func NewFromURL(url string) *Client {
	c := jhttp.NewChannel(url, nil)
	return &Client{
		Client: jrpc2.NewClient(c, nil),
	}
}
