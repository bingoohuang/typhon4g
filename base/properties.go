package base

import (
	"github.com/bingoohuang/properties"
	"github.com/sirupsen/logrus"
)

// PropertiesConfFile defines the properties format of conf file
type PropertiesConfFile struct {
	Conf
	Doc *properties.Doc
}

var _ Prop = (*PropertiesConfFile)(nil)

// NewPropertiesConfFile new a PropertiesConfFile file.
func NewPropertiesConfFile(confFile, raw string) *PropertiesConfFile {
	doc, err := properties.LoadString(raw)
	if err != nil {
		logrus.Warnf("LoadProperties %v", err)
		return nil
	}

	pcf := &PropertiesConfFile{
		Conf: Conf{raw: raw, confFile: confFile},
		Doc:  doc}

	pcf.updater = func(updated string) {
		if doc, err := properties.LoadString(updated); err != nil {
			logrus.Warnf("LoadProperties %v", err)
		} else {
			pcf.Doc = doc
		}
	}

	return pcf
}

// ConfFormat gets the format of conf file
func (p *PropertiesConfFile) ConfFormat() ConfFmt { return PropertiesFmt }

// Str get the string value of key specified by name.
func (p *PropertiesConfFile) Str(name string) string { return p.Doc.Str(name) }

// StrOr get the string value of key specified by name or defaultValue when value is empty or missed.
func (p *PropertiesConfFile) StrOr(name, defaultValue string) string {
	return p.Doc.StrOr(name, defaultValue)
}

// Bool get the bool value of key specified by name.
func (p *PropertiesConfFile) Bool(name string) bool { return p.Doc.Bool(name) }

// BoolOr get the bool value of key specified by name or defaultValue when value is empty or missed.
func (p *PropertiesConfFile) BoolOr(name string, defaultValue bool) bool {
	return p.Doc.BoolOr(name, defaultValue)
}

// Int get the int value of key specified by name.
func (p *PropertiesConfFile) Int(name string) int { return p.Doc.Int(name) }

// IntOr get the int value of key specified by name or defaultValue when value is empty or missed.
func (p *PropertiesConfFile) IntOr(name string, defaultValue int) int {
	return p.Doc.IntOr(name, defaultValue)
}

// Int32 get the int32 value of key specified by name.
func (p *PropertiesConfFile) Int32(name string) int32 { return int32(p.Doc.Int64(name)) }

// Int32Or get the int32 value of key specified by name or defaultValue when value is empty or missed.
func (p *PropertiesConfFile) Int32Or(name string, defaultValue int32) int32 {
	return int32(p.Doc.Int64Or(name, int64(defaultValue)))
}

// Int64 get the int64 value of key specified by name.
func (p *PropertiesConfFile) Int64(name string) int64 { return p.Doc.Int64(name) }

// Int64Or get the int64 value of key specified by name or defaultValue when value is empty or missed.
func (p *PropertiesConfFile) Int64Or(name string, defaultValue int64) int64 {
	return p.Doc.Int64Or(name, defaultValue)
}

// Float32 get the float32 value of key specified by name.
func (p *PropertiesConfFile) Float32(name string) float32 { return float32(p.Doc.Float64(name)) }

// Float32Or get the float32 value of key specified by name or defaultValue when value is empty or missed.
func (p *PropertiesConfFile) Float32Or(name string, defaultValue float32) float32 {
	return float32(p.Doc.Float64Or(name, float64(defaultValue)))
}

// Float64 get the float64 value of key specified by name.
func (p *PropertiesConfFile) Float64(name string) float64 { return p.Doc.Float64(name) }

// Float64Or get the float64 value of key specified by name or defaultValue when value is empty or missed.
func (p *PropertiesConfFile) Float64Or(name string, defaultValue float64) float64 {
	return p.Doc.Float64Or(name, defaultValue)
}
