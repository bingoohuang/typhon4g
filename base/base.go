package base

import (
	"sync"
	"time"

	"github.com/bingoohuang/now"
)

// Conf defines the base structure of conf file
type Conf struct {
	raw       string
	confFile  string
	listeners []ChangeListener

	rawLock sync.RWMutex
	updater func(updated string)
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

// Register registers the change listener of conf file
func (b *Conf) Register(listener ChangeListener) {
	b.listeners = append(b.listeners, listener)
}

// Unregister removes the register of the change listener of conf file
func (b *Conf) Unregister(listener ChangeListener) int {
	ls := make([]ChangeListener, 0, len(b.listeners))
	count := 0

	for _, l := range b.listeners {
		if l != listener {
			ls = append(ls, l)
		} else {
			count++
		}
	}

	b.listeners = ls

	return count
}

// UnregisterAll  removes all registers of the change listener of conf file
func (b *Conf) UnregisterAll() {
	b.listeners = b.listeners[0:0]
}

// TriggerChange trigger the changes event
func (b *Conf) TriggerChange(old, new FileRaw, changedTime time.Time) []ClientReportItem {
	oldRaw := b.Raw()
	b.UpdateRaw(new.Content)

	items := make([]ClientReportItem, 0)

	if !new.TriggerChange {
		return items
	}

	if len(b.listeners) == 0 {
		items = append(items, ClientReportItem{Msg: "No listeners", ConfFile: b.confFile,
			Time: now.MakeNow().P})
		return items
	}

	for _, l := range b.listeners {
		msg, ok := l.OnChange(ConfFileChangeEvent{
			ConfFile:       b.confFile,
			ConfFileFormat: PropertiesFmt,
			Old:            oldRaw,
			Current:        new.Content,
			ChangedTime:    time.Now(),
		})

		items = append(items, ClientReportItem{Msg: msg, Ok: ok, ConfFile: b.confFile, Crc: new.Crc,
			Time: now.MakeNow().P})
	}

	return items
}
