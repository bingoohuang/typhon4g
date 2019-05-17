package typhon4g

import (
	"path/filepath"
	"strings"
	"time"
)

type ConfFileChangeListener interface {
	OnChange(event ConfFileChangeEvent) (string, bool)
}

type ConfFile interface {
	Raw() string
	ConfFormat() ConfFmt
	Name() string

	Register(ConfFileChangeListener)
	Unregister(ConfFileChangeListener) int
	UnregisterAll()

	TriggerChange(old, new *FileContent, changedTime time.Time) []ClientReportItem
}

func NewConfFile(confFile, raw string) ConfFile {
	ext := filepath.Ext(confFile)
	ext = strings.ToLower(ext)
	switch ext {
	case ".properties":
		return MakePropertiesConfFile(confFile, raw)
	default:
		return MakeTxtConfFile(confFile, raw)
	}
}
