package typhon4g

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/bingoohuang/gou"
	"github.com/mitchellh/go-homedir"
	"github.com/thoas/go-funk"
)

// TyphonContext defines the context of typhon client.
type TyphonContext struct {
	// AppID defines the global appID of the typhon client.
	AppID string `json:"appID"`

	// MetaServers defines the meta servers.
	MetaServers []string `json:"metaServers"`
	// ConfigServers defines the config servers.
	ConfigServers []string `json:"configServers"`

	// ConnectTimeoutMillis defines the http connect timeout in millis.
	ConnectTimeoutMillis int64 `json:"connectTimeoutMillis"`
	// PollingReadTimeoutMillis defines the read timeout in millis of config polling
	PollingReadTimeoutMillis int64 `json:"pollingReadTimeoutMillis"`
	// RetryNetworkSleepSeconds defines the sleeping time in seconds before retry.
	RetryNetworkSleepSeconds int64 `json:"retryNetworkSleepSeconds"`
	// ConfigRefreshIntervalSeconds defines the refresh loop interval of config service.
	ConfigRefreshIntervalSeconds int64 `json:"configRefreshIntervalSeconds"`
	// ConfigReadTimeoutMillis defines the read timeout in millis of config service.
	ConfigReadTimeoutMillis int64 `json:"configReadTimeoutMillis"`
	// MetaRefreshIntervalSeconds defines the refresh loop interval of meta service.
	MetaRefreshIntervalSeconds int64 `json:"metaRefreshIntervalSeconds"`
	// SnapshotsDir defines the snapshot directory of config.
	SnapshotsDir string `json:"snapshotsDir"`

	cache     map[string]*FileContent // file->content
	cacheLock sync.RWMutex

	postAuth string
}

// LoadConfFile loads the conf file by name confFile.
func (c *TyphonContext) LoadConfFile(confFile string) ConfFile {
	c.cacheLock.RLock()
	defer c.cacheLock.RUnlock()

	if fc, ok := c.cache[confFile]; ok {
		return fc.Conf
	}

	return nil
}

// SaveFileContents saves the file contents to cache and snapshot.
func (c *TyphonContext) SaveFileContents(fcs []FileContent) *ClientReport {
	items := make([]ClientReportItem, 0)

	c.cacheLock.Lock()
	for _, fc := range fcs {
		if old, ok := c.cache[fc.ConfFile]; ok {
			subs := old.Conf.TriggerChange(*old, fc, time.Now())
			if subs != nil {
				items = append(items, subs...)
			}
		} else {
			fc.init()
			fc := fc
			c.cache[fc.ConfFile] = &fc
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

// RecoverFileContent recover the conf file from snapshot.
func (c *TyphonContext) RecoverFileContent(fc *FileContent) {
	c.cacheLock.Lock()
	defer c.cacheLock.Unlock()

	fc.init()
	c.cache[fc.ConfFile] = fc
}

// WalkFileContents walks the cache.
func (c *TyphonContext) WalkFileContents(fn func(confFile string, fileContext *FileContent)) {
	c.cacheLock.RLock()
	defer c.cacheLock.RUnlock()

	for k, v := range c.cache {
		fn(k, v)
	}
}

// LoadContextFile load the typhon context by file contextFile
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

// LoadContext loads the typhon context by reader.
func LoadContext(reader io.Reader) (*TyphonContext, error) {
	d, err := gou.LoadProperties(reader)
	if nil != err {
		return nil, err
	}

	sd := d.StringDefault("snapshotsDir", "~/.typhon-client/snapshots")
	snapshotsDir, _ := homedir.Expand(sd)

	c := &TyphonContext{
		AppID:                        Required(d.String("appID"), "appID"),
		postAuth:                     d.String("postAuth"),
		ConnectTimeoutMillis:         d.IntDefault("connectTimeoutMillis", 1000),
		ConfigReadTimeoutMillis:      d.IntDefault("configReadTimeoutMillis", 5000),
		PollingReadTimeoutMillis:     d.IntDefault("pollingReadTimeoutMillis", 70000),
		RetryNetworkSleepSeconds:     d.IntDefault("retryNetworkSleepSeconds", 60),
		ConfigRefreshIntervalSeconds: d.IntDefault("configRefreshIntervalSeconds", 300),
		MetaRefreshIntervalSeconds:   d.IntDefault("metaRefreshIntervalSeconds", 300),
	}

	c.cache = make(map[string]*FileContent)
	c.SnapshotsDir = MustMakeDirAll(filepath.Join(snapshotsDir, c.AppID))
	c.MetaServers = c.createMetaServers(d.StringDefault("metaServers", "http://127.0.0.1:11683"))

	return c, nil
}

func (c *TyphonContext) createMetaServers(metaServers string) []string {
	return funk.Map(gou.SplitN(metaServers, ",", true, true),
		func(meta string) string { return meta + "/meta" }).([]string)
}

func (c *TyphonContext) createConfigServers(configServers string) []string {
	return funk.Map(gou.SplitN(configServers, ",", true, true),
		func(url string) string { return url + "/client/config/" + c.AppID }).([]string)
}
