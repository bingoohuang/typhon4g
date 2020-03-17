package apollo

import (
	"sync"

	"github.com/bingoohuang/snow"
	"github.com/bingoohuang/typhon4g/base"
)

// Client defines the apollo client.
type Client struct {
	localIP string
	*base.Context
	fileRaw chan base.FileRawWait

	notifications sync.Map
	releaseKeys   sync.Map
}

// MakeClient makes a apollo client.
func MakeClient(c *base.Context, fileRaw chan base.FileRawWait) *Client {
	localIP := c.LocalIP
	if localIP == "" {
		localIP = snow.InferHostIPv4("")
	}

	return &Client{
		Context: c,
		localIP: localIP,
		fileRaw: fileRaw,
	}
}
