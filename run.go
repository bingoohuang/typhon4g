package typhon4g

import (
	"github.com/bingoohuang/gou"
	"time"
)

type Runner struct {
	context  *TyphonContext
	snapshot *SnapshotService
	config   *ConfigService
	meta     *MetaService
	polling  *PollingService
}

func (r Runner) Start() {
	r.snapshot.init()

	r.initConfigServerUrls()

	r.meta.configServerAddrUpdater = func(addr string) { r.snapshot.saveConfigServers(addr) }
	r.meta.try()

	r.config.updateFn = func(updates []FileContent) { r.snapshot.saveUpdates(updates) }
	r.config.setting = *gou.GetDefaultSetting()
	r.config.setting.ConnectTimeout = time.Duration(r.context.ConnectTimeoutMillis) * time.Millisecond
	r.config.setting.ReadWriteTimeout = time.Duration(r.context.ConfigReadTimeoutMillis) * time.Millisecond

	r.polling = &PollingService{ConfigService: *r.config}
	r.polling.setting.ReadWriteTimeout = time.Duration(r.context.PollingReadTimeoutMillis) * time.Millisecond

	go r.meta.start()
	go r.config.start()
	go r.polling.startPolling()
}

func (r Runner) initConfigServerUrls() {
	if configServers := r.snapshot.loadConfigServers(); configServers != "" {
		r.context.ConfigServerUrls = CreateConfigServerUrls(r.context.AppID, configServers)
	}
}

func (r Runner) GetProperties(confFile string) *PropertiesConfFile {
	return r.GetConfFile(confFile).(*PropertiesConfFile)
}

func (r Runner) GetConfFile(confFile string) ConfFile {
	cf := r.context.LoadConfFile(confFile)
	if cf == nil {
		r.config.try(confFile)
	}

	return r.context.LoadConfFile(confFile)
}
