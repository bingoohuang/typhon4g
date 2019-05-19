package typhon4g

import (
	"fmt"
	"github.com/bingoohuang/gou"
)

type Runner struct {
	C               *TyphonContext
	SnapshotService *SnapshotService
	ConfigService   *ConfigService
	MetaService     *MetaService
	PollingService  *PollingService
}

func (r Runner) Start() {
	r.initConfigServerUrls()

	r.MetaService.ConfigServersAddrUpdater = func(addr string) { r.SnapshotService.SaveMeta(addr) }
	r.MetaService.Try()

	r.ConfigService.UpdateFn = func(updates []FileContent) { r.SnapshotService.saveUpdates(updates) }
	r.ConfigService.Setting = *gou.GetDefaultSetting()
	r.ConfigService.Setting.ConnectTimeout = MillisDuration(r.C.ConnectTimeoutMillis)
	r.ConfigService.Setting.ReadWriteTimeout = MillisDuration(r.C.ConfigReadTimeoutMillis)

	r.PollingService = &PollingService{ConfigService: *r.ConfigService}
	r.PollingService.Setting.ReadWriteTimeout = MillisDuration(r.C.PollingReadTimeoutMillis)

	go r.MetaService.Start()
	go r.ConfigService.Start()
	go r.PollingService.Start()
}

func (r Runner) initConfigServerUrls() {
	if cfs := r.SnapshotService.LoadMeta(); cfs != "" {
		r.C.ConfigServerUrls = r.C.CreateConfigServerUrls(cfs)
	}
}

func (r Runner) Properties(confFile string) (*PropertiesConfFile, error) {
	c, e := r.ConfFile(confFile)
	return c.(*PropertiesConfFile), e
}

func (r Runner) ConfFile(confFile string) (ConfFile, error) {
	cf := r.C.LoadConfFile(confFile)
	if cf != nil {
		return cf, nil
	}

	_, cf = r.ConfigService.Try(confFile)
	if cf != nil {
		return cf, nil
	}

	if cf, err := r.SnapshotService.Load(confFile); err == nil {
		return cf, nil
	}

	return nil, fmt.Errorf("failed to Load conf file %s", confFile)
}

func (r Runner) PostConf(confFile, raw, clientIps string) (string, error) {
	return r.ConfigService.PostConf(confFile, raw, clientIps)
}

func (r Runner) ListenerResults(confFile, crc string) ([]ClientReportRspItem, error) {
	return r.ConfigService.ListenerResults(confFile, crc)
}
