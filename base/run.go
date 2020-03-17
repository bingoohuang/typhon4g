package base

import (
	"context"
	"fmt"

	"github.com/bingoohuang/gou/str"
)

// Runner defines the typhon-client typhon service.
type Runner struct {
	C               *Context
	SnapshotService *SnapshotService
	ConfigService   *ConfigService
	MetaService     *MetaService
	PollingService  *PollingService
	cancelContext   context.Context
	cancelFunc      context.CancelFunc
}

// Start start the typhon.
func (r *Runner) Start() {
	r.initConfigServerUrls()

	r.MetaService.ConfigServersUpdater = func(addr []string) { r.SnapshotService.SaveMeta(addr) }
	r.MetaService.Try()

	r.cancelContext, r.cancelFunc = context.WithCancel(context.Background())

	go r.MetaService.Start(r.cancelContext)
	go r.ConfigService.Start(r.cancelContext)
	go r.PollingService.Start(r.cancelContext)
}

// Stop stops the runner.
func (r *Runner) Stop() {
	r.cancelFunc()
	close(r.C.FileRawChan)
}

func (r Runner) initConfigServerUrls() {
	if len(r.C.MetaServersParsed) > 0 {
		return
	}

	if cfs := r.SnapshotService.LoadMeta(); cfs != "" {
		r.C.UpdateConfigServers(str.SplitN(cfs, ",", true, true))
	}
}

// Properties gets the properties conf file.
func (r Runner) Properties(confFile string) (*PropertiesConfFile, error) {
	c, e := r.ConfFile(confFile)
	if c == nil {
		return nil, e
	}

	return c.(*PropertiesConfFile), e
}

// ConfFile gets the conf file.
func (r Runner) ConfFile(confFile string) (ConfFile, error) {
	cf := r.C.LoadConfFile(confFile)
	if cf != nil {
		return cf, nil
	}

	wait := r.C.Client.ReadConfig(confFile)
	<-wait

	cf = r.C.LoadConfFile(confFile)
	if cf != nil {
		return cf, nil
	}

	if err := r.SnapshotService.Load(confFile); err == nil {
		return cf, nil
	}

	cf = r.C.LoadConfFile(confFile)
	if cf != nil {
		return cf, nil
	}

	return nil, fmt.Errorf("failed to Load conf file %s", confFile)
}

// PostConf posts the conf to the server with clientIps(blank/comma separated or all)
// returns crc and error.
func (r Runner) PostConf(confFile, raw, clientIps string) (string, error) {
	return r.C.Client.PostConf(confFile, raw, clientIps)
}

// ListenerResults get the listener  results
func (r Runner) ListenerResults(confFile, crc string) ([]ClientReportRspItem, error) {
	return r.C.Client.ListenerResults(confFile, crc)
}
