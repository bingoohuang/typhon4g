package base

import (
	"reflect"
	"sync"
	"time"

	"github.com/bingoohuang/gonet"
)

type configClient interface {
	CreateMetaServers() []string
	// MetaGet gets the config servers address from the meta server.
	MetaGet(url string) ([]string, error)
	Polling(configServer string) error
	ReadConfig(confFile string, wait bool) <-chan bool

	// PostConf posts the conf to the server with clientIps(blank/comma separated or all)
	// returns crc and error.
	PostConf(confFile, raw, clientIps string) (string, error)

	ListenerResults(confFile, crc string) ([]ClientReportRspItem, error)
}

// ContextConfig defines the structure for the config.
type ContextConfig struct {
	// 配置中心服务器类型: apollo, typhon, 默认typhon
	ServerType string

	// 以下是apollo专有
	Cluster    string `default:"default"`
	DataCenter string
	LocalIP    string

	// AppID defines the global appID of the typhon client.
	AppID string `validate:"empty=false"`

	// MetaServers defines the meta servers.
	MetaServers string

	// ConfigServers defines the config servers.
	ConfigServers string

	// ConnectTimeoutMillis defines the http connect timeout in millis.
	ConnectTimeout time.Duration `default:"60s"`
	// PollingReadTimeoutMillis defines the read timeout in millis of config polling
	PollingReadTimeout time.Duration `default:"70s"`
	// RetryNetworkSleepSeconds defines the sleeping time in seconds before retry.
	RetryNetworkSleep time.Duration `default:"60s"`
	// ConfigRefreshIntervalSeconds defines the refresh loop interval of config service.
	ConfigRefreshInterval time.Duration `default:"5m"`
	// ConfigReadTimeoutMillis defines the read timeout in millis of config service.
	ConfigReadTimeout time.Duration `default:"5s"`
	// MetaRefreshIntervalSeconds defines the refresh loop interval of meta service.
	MetaRefreshInterval time.Duration `default:"5m"`
	// SnapshotsDir defines the snapshot directory of config.
	SnapshotsDir string

	RootPem   string
	ClientPem string
	ClientKey string
	PostAuth  string
}

// Context defines the context of typhon client.
type Context struct {
	ContextConfig

	Cache     map[string]*FileContent // file->content
	cacheLock sync.RWMutex

	Req     *gonet.ReqOption
	ReqPoll *gonet.ReqOption

	Client configClient

	MetaServersParsed   []string
	ConfigServersParsed []string
	addrLock            sync.Mutex
	FileRawChan         chan FileRawWait
}

// GetConfigServers gets the config servers.
func (c *Context) GetConfigServers() []string {
	c.addrLock.Lock()
	defer c.addrLock.Unlock()

	return c.ConfigServersParsed
}

// UpdateConfigServers updates the config server addresses.
func (c *Context) UpdateConfigServers(servers []string) bool {
	if reflect.DeepEqual(c.ConfigServersParsed, servers) {
		return false
	}

	c.addrLock.Lock()
	defer c.addrLock.Unlock()

	c.ConfigServersParsed = servers

	return true
}

// LoadConfFile loads the conf file by name confFile.
func (c *Context) LoadConfFile(confFile string) ConfFile {
	if fc := c.LoadConfCache(confFile); fc != nil {
		return fc.conf
	}

	return nil
}

// LoadConfCache load the config from the cache.
func (c *Context) LoadConfCache(confFile string) *FileContent {
	c.cacheLock.RLock()
	defer c.cacheLock.RUnlock()

	if fc, ok := c.Cache[confFile]; ok && fc != nil {
		return fc
	}

	c.Cache[confFile] = nil

	return nil
}

// WalkFileContents walks the cache.
func (c *Context) WalkFileContents(fn func(confFile string, fileContext *FileContent)) {
	c.cacheLock.RLock()
	defer c.cacheLock.RUnlock()

	for k, v := range c.Cache {
		fn(k, v)
	}
}
