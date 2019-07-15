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

	"github.com/bingoohuang/gonet"

	"github.com/bingoohuang/properties"

	"github.com/bingoohuang/gou"
	homedir "github.com/mitchellh/go-homedir"
	funk "github.com/thoas/go-funk"
)

// Context defines the context of typhon client.
type Context struct {
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

	PostAuth string

	ReqOption *gonet.ReqOption

	RootPem   string
	ClientPem string
	ClientKey string
}

// LoadConfFile loads the conf file by name confFile.
func (c *Context) LoadConfFile(confFile string) ConfFile {
	if fc := c.LoadConfCache(confFile); fc != nil {
		return fc.conf
	}

	return nil
}

// ClearCache clears the conf cache by name confFile.
func (c *Context) ClearCache(confFile string) {
	c.cacheLock.RLock()
	defer c.cacheLock.RUnlock()

	if _, ok := c.cache[confFile]; ok {
		delete(c.cache, confFile)
	}
}
func (c *Context) LoadConfCache(confFile string) *FileContent {
	c.cacheLock.RLock()
	defer c.cacheLock.RUnlock()

	if fc, ok := c.cache[confFile]; ok {
		return fc
	}

	return nil
}

// SaveFileContents saves the file contents to cache and snapshot.
func (c *Context) SaveFileContents(fcs []FileContent, triggerListeners bool) *ClientReport {
	items := make([]ClientReportItem, 0)

	c.cacheLock.Lock()
	for _, fcItem := range fcs {
		fc := fcItem
		if old, ok := c.cache[fc.ConfFile]; !ok {
			fc.init()
			c.cache[fc.ConfFile] = &fc
		} else if old.Content != fc.Content {
			subs := old.conf.TriggerChange(old, &fc, time.Now(), triggerListeners)
			if subs != nil {
				items = append(items, subs...)
			}

			// caution: DO NOT directly replace to avoiding registered listeners losing.
			old.update(fc)
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
func (c *Context) RecoverFileContent(fc *FileContent) {
	c.cacheLock.Lock()
	defer c.cacheLock.Unlock()

	fc.init()
	c.cache[fc.ConfFile] = fc
}

// WalkFileContents walks the cache.
func (c *Context) WalkFileContents(fn func(confFile string, fileContext *FileContent)) {
	c.cacheLock.RLock()
	defer c.cacheLock.RUnlock()

	for k, v := range c.cache {
		fn(k, v)
	}
}

// LoadContextFile load the typhon context by file contextFile
func LoadContextFile(contextFile string) (*Context, error) {
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
func LoadContext(reader io.Reader) (*Context, error) {
	d, err := properties.Load(reader)
	if nil != err {
		return nil, err
	}

	sd := d.StrOr("snapshotsDir", "~/.typhon-client/snapshots")
	snapshotsDir, _ := homedir.Expand(sd)

	c := &Context{
		AppID:                        Required(d.Str("appID"), "appID"),
		ConnectTimeoutMillis:         d.Int64Or("connectTimeoutMillis", 1000),
		ConfigReadTimeoutMillis:      d.Int64Or("configReadTimeoutMillis", 5000),
		PollingReadTimeoutMillis:     d.Int64Or("pollingReadTimeoutMillis", 70000),
		RetryNetworkSleepSeconds:     d.Int64Or("retryNetworkSleepSeconds", 60),
		ConfigRefreshIntervalSeconds: d.Int64Or("configRefreshIntervalSeconds", 300),
		MetaRefreshIntervalSeconds:   d.Int64Or("metaRefreshIntervalSeconds", 300),

		PostAuth:  d.Str("postAuth"),
		RootPem:   d.Str("rootPem"),
		ClientPem: d.Str("clientPem"),
		ClientKey: d.Str("clientKey"),
	}

	c.cache = make(map[string]*FileContent)
	c.SnapshotsDir = MustMakeDirAll(filepath.Join(snapshotsDir, c.AppID))
	c.MetaServers = c.createMetaServers(d.StrOr("metaServers", "http://127.0.0.1:11683"))

	c.ReqOption = gonet.NewReqOption()
	c.ReqOption.TLSClientConfig = gonet.TLSConfigCreateClientMust(c.ClientKey, c.ClientPem, c.RootPem)

	return c, nil
}

func (c *Context) createMetaServers(metaServers string) []string {
	return funk.Map(gou.SplitN(metaServers, ",", true, true),
		func(meta string) string { return meta + "/meta" }).([]string)
}

func (c *Context) createConfigServers(configServers string) []string {
	return funk.Map(gou.SplitN(configServers, ",", true, true),
		func(url string) string { return url + "/client/config/" + c.AppID }).([]string)
}
