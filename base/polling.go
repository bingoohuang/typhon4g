package base

import (
	"context"
	"time"

	"github.com/bingoohuang/gor"
)

// PollingService defines the polling service.
type PollingService struct {
	*Context
}

// Start starts the polling service loop.
func (p PollingService) Start(ctx context.Context) {
	d := p.RetryNetworkSleep

	for {
		// 固定启动时候，等待1秒，为了让应用准备好读取哪些文件，以方便polling时指定
		// 另外每次循环前固定等待1秒，则是防止空转
		time.Sleep(1 * time.Second) // nolint gomnd

		servers := p.GetConfigServers()
		ok, _ := gor.IterateSlice(servers, -1, func(addr string) bool {
			return p.Client.Polling(addr) == nil
		})

		select {
		case <-ctx.Done():
			return
		default: // required goon
		}

		if !ok {
			time.Sleep(d)
		}
	}
}
