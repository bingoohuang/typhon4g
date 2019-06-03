package typhon4g

// TxtConfFile defines the txt format of conf file
type TxtConfFile struct {
	BaseConf
}

// ConfFormat gets the format of conf file
func (t *TxtConfFile) ConfFormat() ConfFmt {
	return TxtFmt
}

// NewTxtConfFile new a TxtConfFile.
func NewTxtConfFile(confFile, raw string) *TxtConfFile {
	return &TxtConfFile{BaseConf{confFile: confFile, raw: raw,
		listeners: make([]ConfFileChangeListener, 0),
	}}
}
