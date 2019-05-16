package typhon4g

type TxtConfFile struct {
	BaseConf
}

func (t TxtConfFile) ConfFormat() ConfFmt {
	return TXT
}

func MakeTxtConfFile(confFile, raw string) *TxtConfFile {
	tcf := &TxtConfFile{
		BaseConf{
			confFile:  confFile,
			raw:       raw,
			listeners: make([]ConfFileChangeListener, 0),
		},
	}
	return tcf
}
