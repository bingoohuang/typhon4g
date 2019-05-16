package typhon4g

import (
	"github.com/bingoohuang/gou"
	"strings"
	"time"
)

type PollingService struct {
	ConfigService
}

func (p PollingService) startPolling() {
	d := time.Duration(p.c.RetryNetworkSleepSeconds) * time.Second

	for {
		urls := p.c.ConfigServerUrls
		ok := gou.IterateSlice(urls, gou.RandomIntN(uint64(len(urls))), func(url string) bool {
			pollUrl := strings.Replace(url, "/config/", "/notify/", 1)
			return p.tryUrl(pollUrl, "", &p.setting)
		})

		if !ok {
			time.Sleep(d)
		}
	}
}
