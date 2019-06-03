package typhon4g

import (
	"io/ioutil"
	"path/filepath"

	"github.com/bingoohuang/gou"
	"github.com/sirupsen/logrus"
)

type SnapshotService struct {
	C *TyphonContext
}

func (s SnapshotService) Load(file string) (ConfFile, error) {
	content, err := ioutil.ReadFile(filepath.Join(s.C.SnapshotsDir, file))
	if err != nil {
		return nil, err
	}

	raw := string(content)
	fc := &FileContent{
		AppID:    s.C.AppID,
		ConfFile: file,
		Content:  raw,
		Crc:      gou.Checksum(content),
	}

	s.C.RecoverFileContent(fc)
	return fc.Conf, nil
}

func (s SnapshotService) Save(confFile, content string) {
	cf := filepath.Join(s.C.SnapshotsDir, confFile)
	err := ioutil.WriteFile(cf, []byte(content), 0644)
	if err != nil {
		logrus.Warnf("fail to write SnapshotService %s error %v", confFile, err)
	}
}

func (s SnapshotService) LoadMeta() string {
	confFile := filepath.Join(s.C.SnapshotsDir, "_meta")
	meta, err := ioutil.ReadFile(confFile)
	if err != nil {
		logrus.Warnf("fail to read SnapshotService %s error %v", confFile, err)
		return ""
	}

	return string(meta)
}

func (s SnapshotService) SaveMeta(configServerUrls string) {
	s.Save("_meta", configServerUrls)
}

func (s SnapshotService) saveUpdates(fcs []FileContent) {
	for _, fc := range fcs {
		s.Save(fc.ConfFile, fc.Content)
	}
}
