package typhon4g

import (
	"bytes"
	"fmt"
	"github.com/bingoohuang/gou"
	"github.com/mitchellh/go-homedir"
	"github.com/thoas/go-funk"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type TyphonContext struct {
	AppID string `json:"appID"`

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

func (c TyphonContext) LoadConfFile(confFile string) ConfFile {
	c.cacheLock.RLock()
	defer c.cacheLock.RUnlock()

	if fc, ok := c.Cache[confFile]; ok {
		return fc.Conf
	}

	return nil
}

func (c TyphonContext) SaveFileContents(fcs []FileContent) *ClientReport {
	items := make([]ClientReportItem, 0)

	c.cacheLock.Lock()
	for _, fc := range fcs {
		if old, ok := c.Cache[fc.ConfFile]; ok {
			subs := old.Conf.TriggerChange(old, &fc, time.Now())
			if subs != nil {
				items = append(items, subs...)
			}
		} else {
			fc.init()
			c.Cache[fc.ConfFile] = &fc
		}
	}
	c.cacheLock.Unlock()

	hostname, _ := os.Hostname()
	return &ClientReport{
		Host:  hostname,
		Pid:   fmt.Sprintf("%d", os.Getpid()),
		Bin:   os.Args[0],
		Items: items,
	}
}

func (c TyphonContext) RecoverFileContent(fc *FileContent) {
	c.cacheLock.Lock()
	defer c.cacheLock.Unlock()

	fc.init()
	c.Cache[fc.ConfFile] = fc
}

func (c TyphonContext) WalkFileContents(fn func(confFile string, fileContext *FileContent)) {
	c.cacheLock.RLock()
	defer c.cacheLock.RUnlock()

	for k, v := range c.Cache {
		fn(k, v)
	}
}

func LoadContextFile(contextFile string) (*TyphonContext, error) {
	if _, err := os.Stat(contextFile); err != nil {
		return nil, err
	}

	f, err := ioutil.ReadFile(contextFile)
	if nil != err {
		return nil, err
	}

	return LoadContext(bytes.NewBuffer(f))
}

func LoadContext(reader io.Reader) (*TyphonContext, error) {
	d, err := gou.LoadProperties(reader)
	if nil != err {
		return nil, err
	}

	sd := d.StringDefault("snapshotsDir", "~/.typhon-client/snapshots")
	snapshotsDir, _ := homedir.Expand(sd)

	c := &TyphonContext{

		AppID:    MustExists(d.String("appID"), "appID"),
		PostAuth: d.String("postAuth"),

		ConnectTimeoutMillis:         d.IntDefault("connectTimeoutMillis", 1000),
		ConfigReadTimeoutMillis:      d.IntDefault("configReadTimeoutMillis", 5000),
		PollingReadTimeoutMillis:     d.IntDefault("pollingReadTimeoutMillis", 70000),
		RetryNetworkSleepSeconds:     d.IntDefault("retryNetworkSleepSeconds", 60),
		ConfigRefreshIntervalSeconds: d.IntDefault("configRefreshIntervalSeconds", 300),
		MetaRefreshIntervalSeconds:   d.IntDefault("metaRefreshIntervalSeconds", 300),

		Cache: make(map[string]*FileContent),
	}

	c.SnapshotsDir = MustMakeDirAll(filepath.Join(snapshotsDir, c.AppID))
	c.MetaServerUrls = c.CreateMetaServerUrls(d.StringDefault("metaServers", "http://127.0.0.1:11683"))

	return c, nil
}

func (c TyphonContext) CreateMetaServerUrls(metaServers string) []string {
	return funk.Map(gou.SplitN(metaServers, ",", true, true),
		func(meta string) string { return meta + "/meta" }).([]string)
}

func (c TyphonContext) CreateConfigServerUrls(configServers string) []string {
	return funk.Map(gou.SplitN(configServers, ",", true, true),
		func(url string) string { return url + "/client/config/" + c.AppID }).([]string)
}
