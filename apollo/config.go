package apollo

import (
	"fmt"
	"net/url"

	"github.com/sirupsen/logrus"

	"github.com/bingoohuang/gor"
	"github.com/bingoohuang/properties"
	"github.com/bingoohuang/typhon4g/base"
)

// configResult of query config
type configResult struct {
	NamespaceName  string            `json:"namespaceName"`
	Configurations map[string]string `json:"configurations"`
	ReleaseKey     string            `json:"releaseKey"`
}

// ReadConfig reads the config related to namespace.
func (c *Client) ReadConfig(namespace string) <-chan bool {
	if _, ok := c.notifications.Load(namespace); !ok {
		c.notifications.Store(namespace, int64(0))
	}

	wait := make(chan bool)

	c.readConfig(namespace, wait)

	return wait
}

func (c *Client) readConfig(namespace string, wait chan bool) {
	releaseKey, _ := c.releaseKeys.LoadOrStore(namespace, "")

	servers := c.C.GetConfigServers()
	gor.IterateSlice(servers, -1, func(addr string) bool {
		configAddr := c.configAddr(addr, namespace, releaseKey.(string))

		logrus.Infof("config address %s", configAddr)

		var result configResult
		if err := c.C.Req.RestGet(configAddr, &result); err != nil {
			return false
		}

		props, _ := properties.LoadMap(result.Configurations)

		c.releaseKeys.Store(namespace, result.ReleaseKey)

		c.fileRaw <- base.FileRawWait{
			Raw: base.FileRaw{
				AppID:         c.C.AppID,
				ConfFile:      namespace,
				Content:       props.String(),
				Crc:           "",
				TriggerChange: true,
			},
			Wait: wait,
		}

		return true
	})
}

func (c *Client) configAddr(addr, namespace, releaseKey string) string {
	return fmt.Sprintf("%s/configs/%s/%s/%s?releaseKey=%s&ip=%s",
		base.HTPPAddr(addr),
		url.QueryEscape(c.C.AppID),
		url.QueryEscape(c.C.Cluster),
		url.QueryEscape(namespace),
		releaseKey,
		c.localIP)
}
