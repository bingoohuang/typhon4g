package typhon4g

import (
	"strings"
	"time"

	"github.com/bingoohuang/gou"
)

// PollingService defines the polling service.
type PollingService struct {
	ConfigService
}

// Start starts the polling service loop.
func (p PollingService) Start(stop chan bool) {
	d := SecondsDuration(p.C.RetryNetworkSleepSeconds)

	for {
		ok, _ := gou.RandomIterateSlice(p.C.ConfigServers, func(url string) (bool, interface{}) {
			pollURL := strings.Replace(url, "/config/", "/notify/", 1)
			return p.TryURL(pollURL, "", p.Setting)
		})

		select {
		case <-stop:
			return
		default:
			// required goon
		}

		if !ok {
			time.Sleep(d)
		}
	}
}
