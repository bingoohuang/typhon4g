package typhon4g

import (
	"sync"
	"time"
)

type BaseConf struct {
	raw       string
	confFile  string
	listeners []ConfFileChangeListener

	rawLock sync.RWMutex
	updater func(updated string)
}

func (b BaseConf) Raw() string {
	lock := b.rawLock
	lock.RLock()
	defer lock.RUnlock()

	return b.raw
}

func (b *BaseConf) UpdateRaw(updated string) {
	lock := b.rawLock
	lock.Lock()
	defer lock.Unlock()

	b.raw = updated
	if b.updater != nil {
		b.updater(updated)
	}
}

func (b BaseConf) Name() string {
	return b.confFile
}

func (b *BaseConf) Register(listener ConfFileChangeListener) {
	b.listeners = append(b.listeners, listener)
}

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

func (b *BaseConf) UnregisterAll() {
	b.listeners = b.listeners[0:0]
}

func (b *BaseConf) TriggerChange(old, new *FileContent, changedTime time.Time) []ClientReportItem {
	v := b.Raw()
	b.UpdateRaw(new.Content)

	if len(b.listeners) == 0 {
		return nil
	}

	items := make([]ClientReportItem, 0)
	for _, l := range b.listeners {
		msg, ok := l.OnChange(ConfFileChangeEvent{
			ConfFile:       b.confFile,
			ConfFileFormat: Properties,
			Old:            v,
			Current:        new.Content,
			ChangedTime:    time.Now(),
		})

		items = append(items, ClientReportItem{Msg: msg, Ok: ok, ConfFile: b.confFile, Crc: new.Crc,
			Time: time.Now().Format(time.RFC3339)})
	}

	return items
}
