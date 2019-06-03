package typhon4g

import (
	"path/filepath"
	"strings"
	"time"
)

// ConfFileChangeListener defines the interface for change listener
type ConfFileChangeListener interface {
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

	// Register registers the change listener of conf file
	Register(ConfFileChangeListener)
	// Unregister removes the register of the change listener of conf file
	Unregister(ConfFileChangeListener) int
	// UnregisterAll  removes all registers of the change listener of conf file
	UnregisterAll()

	// TriggerChange trigger the changes event
	TriggerChange(old, new *FileContent, changedTime time.Time, triggerListeners bool) []ClientReportItem
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
