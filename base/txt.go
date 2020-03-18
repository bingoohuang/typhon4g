package base

// RawConfFile defines the txt format of conf file
type RawConfFile struct {
	Conf
}

// ConfFormat gets the format of conf file
func (t *RawConfFile) ConfFormat() ConfFmt { return TxtFmt }

// NewRawConfFile new a RawConfFile.
func NewRawConfFile(confFile, raw string) *RawConfFile {
	return &RawConfFile{Conf{confFile: confFile, raw: raw}}
}
