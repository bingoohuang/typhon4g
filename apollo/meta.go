package apollo

import (
	"github.com/bingoohuang/gou/str"
	"github.com/bingoohuang/typhon4g/base"
	"github.com/thoas/go-funk"
)

// CreateMetaServers creates the meta servers addresses.
func (c *Client) CreateMetaServers() []string {
	f := func(meta string) string { return base.HTPPAddr(meta) + "/services/config" }
	return funk.Map(str.SplitN(c.C.MetaServers, ",", true, true), f).([]string)
}

// MetaGet gets the config servers address from the meta server.
func (c *Client) MetaGet(url string) ([]string, error) {
	var metaRsps []MetaRsp

	if err := c.C.Req.RestGet(url, &metaRsps); err != nil {
		return nil, err
	}

	configServes := make([]string, len(metaRsps))

	for i, item := range metaRsps {
		configServes[i] = item.HomepageURL
	}

	return configServes, nil
}

// MetaRsp defines the meta response structure of apollo meta service.
type MetaRsp struct {
	AppName     string `json:"appName"`
	InstanceID  string `json:"instanceId"`
	HomepageURL string `json:"homepageUrl"`
}
