package typhon4g

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/bingoohuang/gor/validate"

	"github.com/bingoohuang/gonet"
	"github.com/bingoohuang/gor/defaults"
	"github.com/bingoohuang/gou/str"
	"github.com/bingoohuang/properties"
	"github.com/bingoohuang/typhon4g/apollo"
	"github.com/bingoohuang/typhon4g/base"
	"github.com/bingoohuang/typhon4g/typhon"
	homedir "github.com/mitchellh/go-homedir"
)

// LoadContextFile load the typhon context by file contextFile
func LoadContextFile(contextFile string) (*base.Context, error) {
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
func LoadContext(reader io.Reader) (*base.Context, error) {
	d, err := properties.Load(reader)
	if nil != err {
		return nil, err
	}

	sd := d.StrOr("snapshotsDir", "~/.typhon-client/snapshots")
	snapshotsDir, _ := homedir.Expand(sd)

	var config base.ContextConfig

	if err := d.Populate(&config, "key"); err != nil {
		return nil, err
	}

	if err := defaults.Set(&config); err != nil {
		return nil, err
	}

	if err := validate.Validate(&config); err != nil {
		return nil, err
	}

	c := &base.Context{ContextConfig: config}

	c.Cache = make(map[string]*base.FileContent)
	c.SnapshotsDir = base.MakeDirAll(filepath.Join(snapshotsDir, c.AppID))

	c.FileRawChan = make(chan base.FileRawWait)
	if c.ServerType == "apollo" {
		c.Client = apollo.MakeClient(c, c.FileRawChan)
	} else {
		c.Client = typhon.MakeClient(c, c.FileRawChan)
	}

	c.MetaServersParsed = c.Client.CreateMetaServers()
	c.ConfigServersParsed = str.SplitN(c.ConfigServers, ",", true, true)

	c.Req = gonet.NewReqOption()
	c.Req.TLSClientConfig = gonet.TLSConfigCreateClient(c.ClientKey, c.ClientPem, c.RootPem)
	c.Req.ConnectTimeout = c.ConnectTimeout
	c.Req.ReadWriteTimeout = c.ConfigReadTimeout

	c.ReqPoll = gonet.NewReqOption()
	c.ReqPoll.TLSClientConfig = gonet.TLSConfigCreateClient(c.ClientKey, c.ClientPem, c.RootPem)
	c.ReqPoll.ConnectTimeout = c.ConnectTimeout
	c.ReqPoll.ReadWriteTimeout = c.PollingReadTimeout

	return c, nil
}
