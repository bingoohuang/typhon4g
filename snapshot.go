package typhon4g

import (
	"github.com/bingoohuang/gou"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"path/filepath"
)

type SnapshotService struct {
	c *TyphonContext
}

func (s SnapshotService) load(file string) error {
	content, err := ioutil.ReadFile(filepath.Join(s.c.SnapshotsDir, file))
	if err != nil {
		return err
	}

	raw := string(content)
	fc := &FileContent{
		AppID:    s.c.AppID,
		ConfFile: file,
		Content:  raw,
		Crc:      gou.Checksum(content),
	}

	s.c.recoverCache(fc)
	return nil
}

func (s SnapshotService) save(confFile, content string) {
	cf := filepath.Join(s.c.SnapshotsDir, confFile)
	err := ioutil.WriteFile(cf, []byte(content), 0644)
	if err != nil {
		logrus.Warnf("fail to write snapshot %s error %v", confFile, err)
	}
}

func (s SnapshotService) loadConfigServers() string {
	confFile := filepath.Join(s.c.SnapshotsDir, "_meta")
	meta, err := ioutil.ReadFile(confFile)
	if err != nil {
		logrus.Warnf("fail to read snapshot %s error %v", confFile, err)
		return ""
	}

	return string(meta)
}

func (s SnapshotService) saveConfigServers(configServerUrls string) {
	s.save("_meta", configServerUrls)
}

func (s SnapshotService) saveUpdates(fcs []FileContent) {
	for _, fc := range fcs {
		s.save(fc.ConfFile, fc.Content)
	}
}
