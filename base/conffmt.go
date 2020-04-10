package base

// ConfFmt defines the conf file format.
type ConfFmt int

const (
	// PropertiesFmt means the conf file is in *.properties format.
	PropertiesFmt ConfFmt = iota
	// TxtFmt means the conf file is in *.txt format.
	TxtFmt
)

// https://github.com/alvaroloes/enumer
//go:generate enumer -type=ConfFmt -json
