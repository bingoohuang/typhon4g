package base

import (
	"path/filepath"
	"strings"
)

// ChangeListener defines the interface for change listener
type ChangeListener interface {
	OnChange(event ConfFileChangeEvent) (string, bool)
}

// ConfFile defines the interface of a typhon conf file.
type ConfFile interface {
	// Raw gets the raw content of conf file
	Raw() string
	// ConfFormat gets the format of conf file
	ConfFormat() ConfFmt
	// Name gets the name of conf file
	Name() string
	// UpdateRaw update the raw content of conf file
	UpdateRaw(updated string)
}

// NewConfFile creates a ConfFile interface by confFile and raw content.
func NewConfFile(confFile, raw string) ConfFile {
	ext := strings.ToLower(filepath.Ext(confFile))

	switch ext {
	case ".properties":
		return NewPropertiesConfFile(confFile, raw)
	default:
		return NewTxtConfFile(confFile, raw)
	}
}
