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
	ext := strings.ToLower(filepath.Ext(confFile))
	switch ext {
	case ".properties":
		return NewPropertiesConfFile(confFile, raw)
	default:
		return NewTxtConfFile(confFile, raw)
	}
}
