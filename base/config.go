package base

import (
	"context"
	"time"
)

// ConfigService defines the structure used for config refresh/write service
type ConfigService struct {
	C *Context
}

// Start starts the refreshing loop of config service.
func (c *ConfigService) Start(ctx context.Context) {
	timer := time.NewTicker(c.C.ConfigRefreshInterval)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			c.C.Client.ReadConfig("")
		}
	}
}
