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
func (c *Client) ReadConfig(namespace string, wait bool) <-chan bool {
	if _, ok := c.notifications.Load(namespace); !ok {
		c.notifications.Store(namespace, int64(0))
	}

	var waitCh chan bool

	if wait {
		waitCh = make(chan bool)
	}

	c.readConfig(namespace, waitCh)

	return waitCh
}

func (c *Client) readConfig(namespace string, wait chan bool) {
	releaseKey, _ := c.releaseKeys.LoadOrStore(namespace, "")

	servers := c.GetConfigServers()
	gor.IterateSlice(servers, -1, func(addr string) bool {
		configAddr := c.configAddr(addr, namespace, releaseKey.(string))

		logrus.Infof("config address %s", configAddr)

		var result configResult
		if err := c.Req.RestGet(configAddr, &result); err != nil {
			return false
		}

		props, _ := properties.LoadMap(result.Configurations)

		c.releaseKeys.Store(namespace, result.ReleaseKey)

		c.fileRaw <- base.FileRawWait{
			FileRaw: base.FileRaw{
				AppID:    c.AppID,
				ConfFile: namespace,
				Content:  props.String(),
				Crc:      "",
			},
			Wait: wait,
		}

		return true
	})
}

func (c *Client) configAddr(addr, namespace, releaseKey string) string {
	return fmt.Sprintf("%s/configs/%s/%s/%s?releaseKey=%s&ip=%s",
		base.HTPPAddr(addr),
		url.QueryEscape(c.AppID),
		url.QueryEscape(c.Cluster),
		url.QueryEscape(namespace),
		releaseKey,
		c.localIP)
}
