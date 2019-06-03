package typhon4g

import (
	"sync"
	"time"
)

// BaseConf defines the base structure of conf file
type BaseConf struct {
	raw       string
	confFile  string
	listeners []ConfFileChangeListener

	rawLock sync.RWMutex
	updater func(updated string)
}

// Raw gets the raw content of conf file
func (b *BaseConf) Raw() string {
	b.rawLock.RLock()
	defer b.rawLock.RUnlock()

	return b.raw
}

// UpdateRaw update the raw content of conf file
func (b *BaseConf) UpdateRaw(updated string) {
	b.rawLock.Lock()
	defer b.rawLock.Unlock()

	b.raw = updated
	if b.updater != nil {
		b.updater(updated)
	}
}

// Name gets the name of conf file
func (b *BaseConf) Name() string {
	return b.confFile
}

// Register registers the change listener of conf file
func (b *BaseConf) Register(listener ConfFileChangeListener) {
	b.listeners = append(b.listeners, listener)
}

// Unregister removes the register of the change listener of conf file
func (b *BaseConf) Unregister(listener ConfFileChangeListener) int {
	ls := make([]ConfFileChangeListener, 0, len(b.listeners))
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
func (b *BaseConf) UnregisterAll() {
	b.listeners = b.listeners[0:0]
}

// TriggerChange trigger the changes event
func (b *BaseConf) TriggerChange(old, new FileContent, changedTime time.Time) []ClientReportItem {
	oldRaw := b.Raw()
	b.UpdateRaw(new.Content)
	items := make([]ClientReportItem, 0)

	if len(b.listeners) == 0 {
		items = append(items, ClientReportItem{Msg: "No listeners", ConfFile: b.confFile,
			Time: time.Now().Format(time.RFC3339)})
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
			Time: time.Now().Format(time.RFC3339)})
	}

	return items
}
