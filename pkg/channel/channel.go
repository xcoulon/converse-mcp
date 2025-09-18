package channel

import (
	"os"

	jrpc2channel "github.com/creachadair/jrpc2/channel"
)

type Channel = jrpc2channel.Channel

var StdioChannel Channel = jrpc2channel.Line(os.Stdin, os.Stdout)

func Direct() (Channel, Channel) {
	return jrpc2channel.Direct()
}
