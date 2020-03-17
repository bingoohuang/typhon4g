package apollo

import "github.com/bingoohuang/typhon4g/base"

// PostConf posts the conf to the server with clientIps (blank/comma separated IP addresses or all)
// returns crc and error info.
func (c *Client) PostConf(confFile, raw, clientIps string) (string, error) {
	panic("not implemented")
}

// ListenerResults gets the listener results from the server.
func (c *Client) ListenerResults(confFile, crc string) ([]base.ClientReportRspItem, error) {
	panic("not implemented")
}
