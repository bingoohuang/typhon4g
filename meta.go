package typhon4g

import (
	"time"

	"github.com/bingoohuang/goreflect"

	"github.com/sirupsen/logrus"
)

// MetaService defines the meta refreshing service.
type MetaService struct {
	C                        *Context
	ConfigServersAddrUpdater func(string)
}

// Start starts the meta refreshing loop
func (m MetaService) Start(stop chan bool) {
	d := SecondsDuration(m.C.MetaRefreshIntervalSeconds)
	timer := time.NewTimer(d)
	defer timer.Stop()

	for {
		select {
		case <-stop:
			return
		case <-timer.C:
			m.Try()
			timer.Reset(d)
		}
	}
}

// Try try to refresh meta.
func (m MetaService) Try() {
	var configServerUrls []string
	goreflect.IterateSlice(m.C.MetaServers, -1, func(url string) bool {
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

// TryURL tries to refresh meta by url.
func (m MetaService) TryURL(url string) ([]string, error) {
	var rsp struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
		Data    string `json:"data"`
	}
	if err := m.C.ReqOption.RestGet(url, &rsp); err != nil {
		return nil, err
	}

	m.ConfigServersAddrUpdater(rsp.Data)
	return m.C.createConfigServers(rsp.Data), nil
}
