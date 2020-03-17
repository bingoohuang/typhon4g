package typhon

import (
	"github.com/bingoohuang/gou/str"
	"github.com/bingoohuang/typhon4g/base"
	"github.com/thoas/go-funk"
)

// CreateMetaServers creates the meta servers addresses.
func (c *Client) CreateMetaServers() []string {
	f := func(meta string) string { return base.HTPPAddr(meta) + "/meta" }
	return funk.Map(str.SplitN(c.C.MetaServers, ",", true, true), f).([]string)
}
