package typhon4g

// ConfFmt defines the conf file format.
type ConfFmt int

const (
	// PropertiesFmt means the conf file is in *.properties format.
	PropertiesFmt ConfFmt = iota
	// TxtFmt means the conf file is in *.txt format.
	TxtFmt
)

//go:generate enumer -type=ConfFmt -json
