package typhon4g

import (
	"bytes"
	"github.com/bingoohuang/gou"
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/thoas/go-funk"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

var runner *Runner

func init() {
	confFile := "etc/typhon-context.properties"
	context, err := LoadTyphonContext(confFile)
	if err != nil {
		logrus.Warnf("unable to load %s, %v", confFile, err)
		return
	}

	r := &Runner{
		context:  context,
		snapshot: &SnapshotService{c: context},
		meta:     &MetaService{c: context},
		config:   &ConfigService{c: context},
	}

	r.start()
	runner = r
}

func GetProperties(confFile string) *PropertiesConfFile {
	return GetConfFile(confFile).(*PropertiesConfFile)
}

func GetConfFile(confFile string) ConfFile {
	cf := runner.context.LoadConfFile(confFile)
	if cf == nil {
		runner.config.try(confFile)
	}

	return runner.context.LoadConfFile(confFile)
}

type TyphonContext struct {
	AppID       string `json:"appID"`
	MetaServers string `json:"metaServers"`

	MetaServerUrls   []string `json:"metaServerUrls"`
	ConfigServerUrls []string `json:"configServerAddr"`

	ConnectTimeoutMillis         int64 `json:"connectTimeoutMillis"`
	PollingReadTimeoutMillis     int64 `json:"pollingReadTimeoutMillis"`
	RetryNetworkSleepSeconds     int64 `json:"retryNetworkSleepSeconds"`
	ConfigRefreshIntervalSeconds int64 `json:"configRefreshIntervalSeconds"`
	ConfigReadTimeoutMillis      int64 `json:"configReadTimeoutMillis"`
	MetaRefreshIntervalSeconds   int64 `json:"metaRefreshIntervalSeconds"`

	SnapshotsDir string `json:"snapshotsDir"`

	Cache     map[string]*FileContent `json:"Cache"` // file->content
	cacheLock sync.RWMutex            `json:"-"`
}

func (t TyphonContext) LoadConfFile(confFile string) ConfFile {
	return t.getCache(confFile)
}

func (t TyphonContext) getCache(confFile string) ConfFile {
	t.cacheLock.RLock()
	defer t.cacheLock.RUnlock()

	if fc, ok := t.Cache[confFile]; ok {
		return fc.Conf
	}

	return nil
}

func (t TyphonContext) saveCaches(fcs []FileContent) {
	t.cacheLock.Lock()
	for _, fc := range fcs {
		if old, ok := t.Cache[fc.ConfFile]; ok {
			old.Conf.TriggerChange(fc.Content, time.Now())
		} else {
			fc.init()
			t.Cache[fc.ConfFile] = &fc
		}
	}
	t.cacheLock.Unlock()
}

func (t TyphonContext) recoverCache(fc *FileContent) {
	t.cacheLock.Lock()
	fc.init()
	t.Cache[fc.ConfFile] = fc
	t.cacheLock.Unlock()
}

func (t TyphonContext) iterateCache(fn func(confFile string, fileContext *FileContent)) {
	t.cacheLock.RLock()
	defer t.cacheLock.RUnlock()

	for k, v := range t.Cache {
		fn(k, v)
	}
}

func LoadTyphonContext(confFile string) (*TyphonContext, error) {
	if _, err := os.Stat(confFile); err != nil {
		return nil, err
	}

	f, err := ioutil.ReadFile(confFile)
	if nil != err {
		return nil, err
	}

	d, err := gou.LoadProperties(bytes.NewBuffer(f))
	if nil != err {
		return nil, err
	}

	sd := d.StringDefault("SnapshotsDir", "~/.typhon-client/snapshots")
	snapshotsDir, _ := homedir.Expand(sd)

	c := &TyphonContext{
		AppID:       d.String("appID"),
		MetaServers: d.StringDefault("metaServers", "http://127.0.0.1:11683"),

		ConnectTimeoutMillis:         d.IntDefault("connectTimeoutMillis", 1000),
		ConfigReadTimeoutMillis:      d.IntDefault("configReadTimeoutMillis", 5000),
		PollingReadTimeoutMillis:     d.IntDefault("pollingReadTimeoutMillis", 70000),
		RetryNetworkSleepSeconds:     d.IntDefault("retryNetworkSleepSeconds", 60),
		ConfigRefreshIntervalSeconds: d.IntDefault("configRefreshIntervalSeconds", 300),
		MetaRefreshIntervalSeconds:   d.IntDefault("metaRefreshIntervalSeconds", 300),

		SnapshotsDir: snapshotsDir,
		Cache:        make(map[string]*FileContent),
	}

	c.MetaServerUrls = funk.Map(gou.SplitN(c.MetaServers, ",", true, true),
		func(meta string) string { return meta + "/meta" }).([]string)

	return c, nil
}
