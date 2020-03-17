package typhon

// Polling polling the specified addr.
func (c *Client) Polling(addr string) error {
	pollingAddr := c.pollingAddr(addr)

	return c.readConfig(pollingAddr, "", nil, true)
}

func (c *Client) pollingAddr(addr string) string {
	return addr + "/client/notify/" + c.C.AppID
}
