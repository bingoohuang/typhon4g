package typhon4g

import (
	"time"

	"github.com/bingoohuang/gou"
	"github.com/sirupsen/logrus"
)

type MetaService struct {
	C                        *TyphonContext
	ConfigServersAddrUpdater func(string)
}

func (m MetaService) Start() {
	d := SecondsDuration(m.C.MetaRefreshIntervalSeconds)
	timer := time.NewTimer(d)
	defer timer.Stop()

	for range timer.C {
		m.Try()
		timer.Reset(d)
	}
}

func (m MetaService) Try() {
	var configServerUrls []string
	gou.RandomIterateSlice(m.C.MetaServers, func(url string) bool {
		var err error
		configServerUrls, err = m.TryURL(url)
		if err != nil {
			logrus.Warnf("fail to TryURL %v", err)
			return false
		}

		if len(configServerUrls) == 0 {
			logrus.Warnf("fail to TryURL empty")
			return false
		}

		return true // break the iterate
	})

	if len(configServerUrls) > 0 {
		m.C.ConfigServers = configServerUrls
	}
}

type MetaRsp struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

func (m MetaService) TryURL(url string) ([]string, error) {
	var rsp MetaRsp
	if err := gou.RestGet(url, &rsp); err != nil {
		return nil, err
	}

	m.ConfigServersAddrUpdater(rsp.Data)
	return m.C.createConfigServers(rsp.Data), nil
}
