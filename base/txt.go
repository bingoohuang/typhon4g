package base

// TxtConfFile defines the txt format of conf file
type TxtConfFile struct {
	Conf
}

// ConfFormat gets the format of conf file
func (t *TxtConfFile) ConfFormat() ConfFmt {
	return TxtFmt
}

// NewTxtConfFile new a TxtConfFile.
func NewTxtConfFile(confFile, raw string) *TxtConfFile {
	return &TxtConfFile{Conf{confFile: confFile, raw: raw,
		listeners: make([]ChangeListener, 0),
	}}
}
