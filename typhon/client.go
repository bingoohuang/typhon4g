package typhon

import (
	"github.com/bingoohuang/gou/str"
	"github.com/bingoohuang/typhon4g/base"
)

// MetaRsp defines the meta response structure of typhon meta service.
type MetaRsp struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

// Client defines the typhon client.
type Client struct {
	C       *base.Context
	fileRaw chan base.FileRawWait
}

// MetaGet gets the config servers address from the meta server.
func (c *Client) MetaGet(addr string) ([]string, error) {
	var rsp MetaRsp

	if err := c.C.Req.RestGet(addr, &rsp); err != nil {
		return nil, err
	}

	return str.SplitN(rsp.Data, ",", true, true), nil
}

// MakeClient makes a apollo client.
func MakeClient(c *base.Context, fileRaw chan base.FileRawWait) *Client {
	return &Client{
		C:       c,
		fileRaw: fileRaw,
	}
}
