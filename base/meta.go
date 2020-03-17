package base

import (
	"context"
	"time"

	"github.com/bingoohuang/gor"

	"github.com/sirupsen/logrus"
)

// MetaService defines the meta refreshing service.
type MetaService struct {
	*Context
	ConfigServersUpdater func(addr []string)
}

// Start starts the meta refreshing loop
func (m MetaService) Start(ctx context.Context) {
	if len(m.MetaServersParsed) == 0 { // 没有配置metaServers
		return
	}

	timer := time.NewTimer(m.MetaRefreshInterval)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			m.Try()
		}
	}
}

// Try try to refresh meta.
func (m MetaService) Try() {
	gor.IterateSlice(m.MetaServersParsed, -1, func(url string) bool {
		configServers, err := m.Client.MetaGet(url)
		if err != nil {
			logrus.Warnf("fail to MetaGet %v", err)
			return false
		}

		if m.UpdateConfigServers(configServers) {
			m.ConfigServersUpdater(configServers)
		}

		return true // break the iterate
	})
}
