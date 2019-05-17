package typhon4g

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/bingoohuang/gou"
	"github.com/mitchellh/go-homedir"
	"github.com/thoas/go-funk"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"
)

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

	PostAuth string `json:"postAuth"`
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

func (t TyphonContext) saveCaches(fcs []FileContent) *ClientReport {
	items := make([]ClientReportItem, 0)

	t.cacheLock.Lock()
	for _, fc := range fcs {
		if old, ok := t.Cache[fc.ConfFile]; ok {
			subs := old.Conf.TriggerChange(old, &fc, time.Now())
			if subs != nil {
				items = append(items, subs...)
			}
		} else {
			fc.init()
			t.Cache[fc.ConfFile] = &fc
		}
	}
	t.cacheLock.Unlock()

	if len(items) == 0 {
		return nil
	}

	hostname, _ := os.Hostname()
	return &ClientReport{
		Host:  hostname,
		Pid:   fmt.Sprintf("%d", os.Getpid()),
		Bin:   os.Args[0],
		Items: items,
	}
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
		PostAuth:    d.String("postAuth"),

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
