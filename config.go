package typhon4g

import (
	"github.com/bingoohuang/gou"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

type ConfigService struct {
	c        *TyphonContext
	setting  gou.UrlHttpSettings
	updateFn func([]FileContent)
}

type ConfigRsp struct {
	Status  int           `json:"status"`
	Message string        `json:"message"`
	Data    []FileContent `json:"data"`
}

func (c ConfigService) start() {
	c.try("")

	d := time.Duration(c.c.ConfigRefreshIntervalSeconds) * time.Second
	timer := time.NewTimer(d)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			c.try("")
			timer.Reset(d)
		}
	}
}

func (c ConfigService) try(confFile string) bool {
	urls := c.c.ConfigServerUrls
	return gou.IterateSlice(urls, gou.RandomIntN(uint64(len(urls))), func(url string) bool {
		return c.tryUrl(url, confFile, &c.setting)
	})
}

func (c ConfigService) tryUrl(url, confFile string, setting *gou.UrlHttpSettings) bool {
	confFileCrc := ""
	if confFile != "" {
		confFileCrc = confFile + ":0"
	} else {
		confFileCrc = c.createConfFileCrcs()
	}
	if confFileCrc == "" {
		return true
	}

	clientUrl := url + "?confFileCrc=" + confFileCrc
	var rsp ConfigRsp
	err := gou.RestGetV2(clientUrl, &rsp, setting)
	if err != nil {
		logrus.Warnf("fail to RefreshConfig %s, error %v", clientUrl, err)
		return false
	}

	if len(rsp.Data) > 0 {
		c.updateFn(rsp.Data)
		c.c.saveCaches(rsp.Data)
		return true
	}

	return false
}

func (c ConfigService) createConfFileCrcs() string {
	confFileCrc := make([]string, 0)
	c.c.iterateCache(func(cf string, fc *FileContent) {
		confFileCrc = append(confFileCrc, fc.ConfFile+":"+fc.Crc)
	})

	return strings.Join(confFileCrc, ",")
}
