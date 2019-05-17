package typhon4g

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/bingoohuang/gou"
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/thoas/go-funk"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var runner *Runner

func init() {
	r, err := CreateRunner("etc/typhon-context.properties")
	if err != nil {
		logrus.Warnf("unable to create typhon runner %v", err)
		return
	}

	r.Start()
	runner = r
}

func CreateRunner(contextFile string) (*Runner, error) {
	context, err := loadTyphonContext(contextFile)
	if err != nil {
		return nil, fmt.Errorf("unable to load %s, %v", contextFile, err)
	}

	return &Runner{
		context:  context,
		snapshot: &SnapshotService{c: context},
		meta:     &MetaService{c: context},
		config:   &ConfigService{c: context},
	}, nil
}

func GetProperties(confFile string) (*PropertiesConfFile, error) {
	return runner.GetProperties(confFile)
}

func GetConfFile(confFile string) (ConfFile, error) {
	return runner.GetConfFile(confFile)
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

func loadTyphonContext(confFile string) (*TyphonContext, error) {
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

	appID := d.String("appID")
	if appID == "" {
		return nil, errors.New("appID required")
	}

	sd := d.StringDefault("snapshotsDir", "~/.typhon-client/snapshots")
	snapshotsDir, _ := homedir.Expand(sd)

	c := &TyphonContext{
		AppID:       appID,
		MetaServers: d.StringDefault("metaServers", "http://127.0.0.1:11683"),

		ConnectTimeoutMillis:         d.IntDefault("connectTimeoutMillis", 1000),
		ConfigReadTimeoutMillis:      d.IntDefault("configReadTimeoutMillis", 5000),
		PollingReadTimeoutMillis:     d.IntDefault("pollingReadTimeoutMillis", 70000),
		RetryNetworkSleepSeconds:     d.IntDefault("retryNetworkSleepSeconds", 60),
		ConfigRefreshIntervalSeconds: d.IntDefault("configRefreshIntervalSeconds", 300),
		MetaRefreshIntervalSeconds:   d.IntDefault("metaRefreshIntervalSeconds", 300),

		SnapshotsDir: filepath.Join(snapshotsDir, appID),
		Cache:        make(map[string]*FileContent),
	}

	makeSnapshotDir(c.SnapshotsDir)

	c.MetaServerUrls = funk.Map(gou.SplitN(c.MetaServers, ",", true, true),
		func(meta string) string { return meta + "/meta" }).([]string)

	return c, nil
}

func makeSnapshotDir(dir string) {
	st, err := os.Stat(dir)
	if err == nil {
		if st.IsDir() {
			return
		} else {
			logrus.Panicf("make sure that snapshot dir %s is a directory", dir)
		}
	}

	if os.IsNotExist(err) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			logrus.Panicf("failed to create snapshot dir %s, error %v", dir, err)
		}
		return
	}

	logrus.Panicf("failed to stat snapshot dir %s, error %v", dir, err)
}
