package base

import (
	"context"
	"fmt"
	"time"

	"github.com/bingoohuang/now"

	"github.com/bingoohuang/gou/str"
)

// Runner defines the typhon-client typhon service.
type Runner struct {
	*Context
	SnapshotService *SnapshotService
	ConfigService   *ConfigService
	MetaService     *MetaService
	PollingService  *PollingService
	cancelContext   context.Context
	cancelFunc      context.CancelFunc

	listeners map[string][]ChangeListener
}

// Start start the typhon.
func (r *Runner) Start() {
	r.initConfigServerUrls()

	r.MetaService.ConfigServersUpdater = func(addr []string) { r.SnapshotService.SaveMeta(addr) }
	r.MetaService.Try()

	r.cancelContext, r.cancelFunc = context.WithCancel(context.Background())
	r.listeners = make(map[string][]ChangeListener)

	go r.ConsumeChan()
	go r.MetaService.Start(r.cancelContext)
	go r.ConfigService.Start(r.cancelContext)
	go r.PollingService.Start(r.cancelContext)
}

// ConsumeChan consumes the updating config changes from the channel.
func (r *Runner) ConsumeChan() {
	for raw := range r.FileRawChan {
		r.updateCache(raw)
	}
}

func (r *Runner) updateCache(raw FileRawWait) {
	r.cacheLock.Lock()
	defer r.cacheLock.Unlock()

	changed := false

	if old, ok := r.Cache[raw.ConfFile]; !ok {
		old = &FileContent{
			FileRaw: FileRaw{
				AppID:    raw.AppID,
				ConfFile: raw.ConfFile,
				Content:  "",
				Crc:      "",
			},
		}
		old.init()

		r.Cache[raw.ConfFile] = old
		changed = true
	} else if old.FileRaw.Content != raw.Content {
		changed = true
	}

	old := r.Cache[raw.ConfFile]

	if changed {
		oldFileRaw := old.FileRaw
		old.conf.UpdateRaw(raw.Content)
		old.FileRaw = raw.FileRaw

		if !raw.TriggerChangeIgnore {
			r.triggerChange(old.conf.Name(), oldFileRaw, raw.FileRaw)
		}
	}

	if !raw.SnapshotIgnore {
		r.SnapshotService.Save(raw.ConfFile, raw.Content)
	}

	if raw.Wait != nil {
		raw.Wait <- true
	}
}

// Stop stops the runner.
func (r *Runner) Stop() {
	r.cancelFunc()
	close(r.FileRawChan)
}

func (r *Runner) initConfigServerUrls() {
	if len(r.MetaServersParsed) > 0 {
		return
	}

	if cfs := r.SnapshotService.LoadMeta(); cfs != "" {
		r.UpdateConfigServers(str.SplitN(cfs, ",", true, true))
	}
}

// Properties gets the properties conf file.
func (r *Runner) Properties(confFile string) (*PropertiesConfFile, error) {
	c, e := r.ConfFile(confFile)
	if c == nil {
		return nil, e
	}

	return c.(*PropertiesConfFile), e
}

// ConfFile gets the conf file.
func (r *Runner) ConfFile(confFile string) (ConfFile, error) {
	cf := r.LoadConfFile(confFile)
	if cf != nil {
		return cf, nil
	}

	<-r.Client.ReadConfig(confFile, true)

	cf = r.LoadConfFile(confFile)
	if cf != nil {
		return cf, nil
	}

	if err := r.SnapshotService.Load(confFile); err == nil {
		return cf, nil
	}

	cf = r.LoadConfFile(confFile)
	if cf != nil {
		return cf, nil
	}

	return nil, fmt.Errorf("failed to Load conf file %s", confFile)
}

// PostConf posts the conf to the server with clientIps(blank/comma separated or all)
// returns crc and error.
func (r *Runner) PostConf(confFile, raw, clientIps string) (string, error) {
	return r.Client.PostConf(confFile, raw, clientIps)
}

// ListenerResults get the listener  results
func (r *Runner) ListenerResults(confFile, crc string) ([]ClientReportRspItem, error) {
	return r.Client.ListenerResults(confFile, crc)
}

// Register registers the change listener of conf file
func (r *Runner) Register(filename string, l ChangeListener) {
	if _, ok := r.listeners[filename]; !ok {
		r.listeners[filename] = make([]ChangeListener, 0)
	}

	r.listeners[filename] = append(r.listeners[filename], l)
	_, _ = r.ConfFile(filename)
}

// Unregister removes the register of the change listener of conf file
func (r *Runner) Unregister(filename string, listener ChangeListener) int {
	if _, ok := r.listeners[filename]; !ok {
		return 0
	}

	ls := make([]ChangeListener, 0, len(r.listeners[filename]))
	count := 0

	for _, l := range r.listeners[filename] {
		if l != listener {
			ls = append(ls, l)
		} else {
			count++
		}
	}

	r.listeners[filename] = ls

	return count
}

// UnregisterAll  removes all registers of the change listener of conf file
func (r *Runner) UnregisterAll(filename string) {
	delete(r.listeners, filename)
}

// triggerChange trigger the changes event
func (r *Runner) triggerChange(confFileName string, old, new FileRaw) []ClientReportItem {
	items := make([]ClientReportItem, 0)
	listeners, ok := r.listeners[confFileName]

	if !ok {
		return append(items, ClientReportItem{
			Msg:      "No listeners",
			ConfFile: confFileName,
			Time:     now.MakeNow().P})
	}

	for _, l := range listeners {
		msg, ok := l.OnChange(ConfFileChangeEvent{
			ConfFile:       confFileName,
			ConfFileFormat: PropertiesFmt,
			Old:            old.Content,
			Current:        new.Content,
			ChangedTime:    time.Now(),
		})

		items = append(items, ClientReportItem{
			Msg:      msg,
			Ok:       ok,
			ConfFile: confFileName,
			Crc:      new.Crc,
			Time:     now.MakeNow().P,
		})
	}

	return items
}
