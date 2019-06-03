package typhon4g

type TxtConfFile struct {
	BaseConf
}

// ConfFormat gets the format of conf file
func (t *TxtConfFile) ConfFormat() ConfFmt {
	return TxtFmt
}

func NewTxtConfFile(confFile, raw string) *TxtConfFile {
	tcf := &TxtConfFile{
		BaseConf{
			confFile:  confFile,
			raw:       raw,
			listeners: make([]ConfFileChangeListener, 0),
		},
	}
	return tcf
}
