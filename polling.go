package typhon4g

import (
	"strings"
	"time"

	"github.com/bingoohuang/gou"
)

type PollingService struct {
	ConfigService
}

func (p PollingService) Start() {
	d := time.Duration(p.C.RetryNetworkSleepSeconds) * time.Second

	for {
		ok, _ := gou.RandomIterateSlice(p.C.ConfigServers, func(url string) (bool, interface{}) {
			pollURL := strings.Replace(url, "/config/", "/notify/", 1)
			return p.TryURL(pollURL, "", &p.Setting)
		})

		if !ok {
			time.Sleep(d)
		}
	}
}
