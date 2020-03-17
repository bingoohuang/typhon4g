package base

import (
	"sync"
)

// Conf defines the base structure of conf file
type Conf struct {
	raw      string
	confFile string
	rawLock  sync.RWMutex
	updater  func(updated string)
}

// Raw gets the raw content of conf file
func (b *Conf) Raw() string {
	b.rawLock.RLock()
	defer b.rawLock.RUnlock()

	return b.raw
}

// UpdateRaw update the raw content of conf file
func (b *Conf) UpdateRaw(updated string) {
	b.rawLock.Lock()
	defer b.rawLock.Unlock()

	b.raw = updated
	if b.updater != nil {
		b.updater(updated)
	}
}

// Name gets the name of conf file
func (b *Conf) Name() string {
	return b.confFile
}
