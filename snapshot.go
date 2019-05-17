package typhon4g

import (
	"github.com/bingoohuang/gou"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path/filepath"
)

type SnapshotService struct {
	c *TyphonContext
}

func (s SnapshotService) init() {
	dir := s.c.SnapshotsDir
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			logrus.Warnf("failed to create snapshot dir %s, error %v", dir, err)
			return
		}
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		logrus.Warnf("failed to read snapshot dir %v", err)
		return
	}

	for _, f := range files {
		if !f.IsDir() {
			s.load(dir, f.Name())
		}
	}
}

func (s SnapshotService) load(dir, file string) {
	content, err := ioutil.ReadFile(filepath.Join(dir, file))
	if err != nil {
		logrus.Warnf("failed to read %s, error %v", file, err)
		return
	}

	raw := string(content)
	fc := &FileContent{
		AppID:    s.c.AppID,
		ConfFile: file,
		Content:  raw,
		Crc:      gou.Checksum(content),
	}

	s.c.recoverCache(fc)
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
