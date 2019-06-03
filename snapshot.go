package typhon4g

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/bingoohuang/gou"
	"github.com/sirupsen/logrus"
)

// SnapshotService defines the snapshot service of typhon client
type SnapshotService struct {
	C *TyphonContext
}

// Load loads the snapshot in file.
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

// Save saves the confFile and its content to snapshot.
func (s SnapshotService) Save(confFile, content string) {
	cf := filepath.Join(s.C.SnapshotsDir, confFile)
	err := ioutil.WriteFile(cf, []byte(content), 0644)
	if err != nil {
		logrus.Warnf("fail to write SnapshotService %s error %v", confFile, err)
	}
}

// LoadMeta loads meta from snapshot.
func (s SnapshotService) LoadMeta() string {
	confFile := filepath.Join(s.C.SnapshotsDir, "_meta")
	meta, err := ioutil.ReadFile(confFile)
	if err != nil {
		logrus.Warnf("fail to read SnapshotService %s error %v", confFile, err)
		return ""
	}

	return string(meta)
}

// SaveMeta saves the meta to snapshot.
func (s SnapshotService) SaveMeta(configServerUrls string) {
	s.Save("_meta", configServerUrls)
}

// SaveUpdates saves the updates to snapshot.
func (s SnapshotService) SaveUpdates(fcs []FileContent) {
	for _, fc := range fcs {
		s.Save(fc.ConfFile, fc.Content)
	}
}

const DeletedAt = ".deletedAt."

// SaveUpdates saves the updates to snapshot.
func (s SnapshotService) Clear(confFile string) error {
	from := filepath.Join(s.C.SnapshotsDir, confFile)
	to := filepath.Join(s.C.SnapshotsDir, confFile+DeletedAt+gou.FormatDateLayout(time.Now(), "yyyyMMddHHmmss"))

	return os.Rename(from, to)
}
