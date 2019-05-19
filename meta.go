package typhon4g

import (
	"github.com/bingoohuang/gou"
	"github.com/sirupsen/logrus"
	"time"
)

type MetaService struct {
	C                        *TyphonContext
	ConfigServersAddrUpdater func(string)
}

func (m MetaService) Start() {
	d := SecondsDuration(m.C.MetaRefreshIntervalSeconds)
	timer := time.NewTimer(d)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			m.Try()
			timer.Reset(d)
		}
	}
}

func (m MetaService) Try() {
	var configServerUrls []string
	gou.RandomIterateSlice(m.C.MetaServerUrls, func(url string) bool {
		var err error
		configServerUrls, err = m.TryUrl(url)
		if err != nil {
			logrus.Warnf("fail to TryUrl %v", err)
			return false
		}

		if len(configServerUrls) == 0 {
			logrus.Warnf("fail to TryUrl empty")
			return false
		}

		return true // break the iterate
	})

	if len(configServerUrls) > 0 {
		m.C.ConfigServerUrls = configServerUrls
	}
}

type MetaRsp struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

func (m MetaService) TryUrl(url string) ([]string, error) {
	var rsp MetaRsp
	if err := gou.RestGet(url, &rsp); err != nil {
		return nil, err
	}

	m.ConfigServersAddrUpdater(rsp.Data)
	return m.C.CreateConfigServerUrls(rsp.Data), nil
}
