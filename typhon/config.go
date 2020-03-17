package typhon

import (
	"os"
	"strings"

	"github.com/bingoohuang/gor"
	"github.com/bingoohuang/typhon4g/base"
	"github.com/sirupsen/logrus"
)

func (c *Client) createConfigServer(addr string) string {
	return addr + "/client/config/" + c.C.AppID
}

// ReadConfig tries to refresh conf defined by confFile or all (confFile is empty).
func (c *Client) ReadConfig(confFile string) <-chan bool {
	wait := make(chan bool)

	gor.IterateSlice(c.C.ConfigServersParsed, -1, func(addr string) bool {
		configAddr := c.createConfigServer(addr)
		return c.readConfig(configAddr, confFile, wait, false) == nil
	})

	return wait
}

type configRsp struct {
	Status  int            `json:"status"`
	Message string         `json:"message"`
	Data    []base.FileRaw `json:"data"`
}

// readConfig tries to refresh conf defined by confFile or all (confFile is empty) in specified URL.
func (c *Client) readConfig(url, confFile string, wait chan bool, isPoll bool) error {
	confFileCrc := c.createConfFileCrcs(confFile)
	if confFileCrc == "" {
		return nil
	}

	clientURL := url + "?confFileCrc=" + confFileCrc

	var rsp configRsp

	req := c.C.Req
	if isPoll {
		req = c.C.ReqPoll
	}

	err := req.RestGet(clientURL, &rsp)
	if err != nil {
		if isPoll && os.IsTimeout(err) {
			logrus.Infof("normal polling timeout %s", clientURL)
			return nil
		}

		logrus.Warnf("fail to ReadConfig %s, error %v", clientURL, err)

		return err
	}

	if len(rsp.Data) == 0 {
		return nil
	}

	for _, item := range rsp.Data {
		item.TriggerChange = true

		c.fileRaw <- base.FileRawWait{
			Raw:  item,
			Wait: wait,
		}
	}

	return nil
}

func (c *Client) createConfFileCrcs(confFile string) string {
	if confFile != "" {
		return confFile + ":0"
	}

	return c.CreateConfFileCrcs()
}

// CreateConfFileCrcs creates conf files and their crcs.
func (c *Client) CreateConfFileCrcs() string {
	confFileCrc := make([]string, 0)

	c.C.WalkFileContents(func(cf string, fc *base.FileContent) {
		confFileCrc = append(confFileCrc, fc.ConfFile+":"+fc.Crc)
	})

	return strings.Join(confFileCrc, ",")
}
