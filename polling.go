package typhon4g

import (
	"github.com/bingoohuang/gou"
	"strings"
	"time"
)

type PollingService struct {
	ConfigService
}

func (p PollingService) Start() {
	d := time.Duration(p.C.RetryNetworkSleepSeconds) * time.Second

	for {
		ok, _ := gou.RandomIterateSlice(p.C.ConfigServerUrls, func(url string) (bool, interface{}) {
			pollUrl := strings.Replace(url, "/config/", "/notify/", 1)
			return p.TryUrl(pollUrl, "", &p.Setting)
		})

		if !ok {
			time.Sleep(d)
		}
	}
}
