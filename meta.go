package typhon4g

import (
	"github.com/bingoohuang/gou"
	"github.com/sirupsen/logrus"
	"github.com/thoas/go-funk"
	"time"
)

type MetaService struct {
	c                       *TyphonContext
	configServerAddrUpdater func(string)
}

func (m MetaService) start() {
	d := time.Duration(m.c.MetaRefreshIntervalSeconds) * time.Second
	timer := time.NewTimer(d)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			m.try()
			timer.Reset(d)
		}
	}
}

func (m MetaService) try() {
	var configServerUrls []string
	urls := m.c.MetaServerUrls
	gou.IterateSlice(urls, gou.RandomIntN(uint64(len(urls))), func(url string) bool {
		var err error
		configServerUrls, err = m.tryUrl(url)
		if err != nil {
			logrus.Warnf("fail to tryUrl %v", err)
			return false
		}

		if len(configServerUrls) == 0 {
			logrus.Warnf("fail to tryUrl empty")
			return false
		}

		return true // break the iterate
	})

	if len(configServerUrls) > 0 {
		m.c.ConfigServerUrls = configServerUrls
	}
}

type MetaRsp struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

func (m MetaService) tryUrl(url string) ([]string, error) {
	var rsp MetaRsp
	if err := gou.RestGet(url, &rsp); err != nil {
		return nil, err
	}

	m.configServerAddrUpdater(rsp.Data)
	return CreateConfigServerUrls(m.c.AppID, rsp.Data), nil
}

func CreateConfigServerUrls(appID string, configServers string) []string {
	urls := gou.SplitN(configServers, ",", true, true)
	mf := func(url string) string { return url + "/client/config/" + appID }
	return funk.Map(urls, mf).([]string)
}
