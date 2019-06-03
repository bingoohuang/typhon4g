package typhon4g

import (
	"bytes"
	"strconv"

	"github.com/bingoohuang/gou"
	"github.com/sirupsen/logrus"
)

// PropertiesConfFile defines the properties format of conf file
type PropertiesConfFile struct {
	BaseConf
	doc *gou.PropertiesDoc
}

var _ Prop = (*PropertiesConfFile)(nil)

// NewPropertiesConfFile new a PropertiesConfFile file.
func NewPropertiesConfFile(confFile, raw string) *PropertiesConfFile {
	doc, err := gou.LoadProperties(bytes.NewBufferString(raw))
	if err != nil {
		logrus.Warnf("LoadProperties %v", err)
		return nil
	}

	pcf := &PropertiesConfFile{
		BaseConf: BaseConf{raw: raw, confFile: confFile, listeners: make([]ConfFileChangeListener, 0)},
		doc:      doc}

	pcf.updater = func(updated string) {
		if doc, err := gou.LoadProperties(bytes.NewBufferString(updated)); err != nil {
			logrus.Warnf("LoadProperties %v", err)
		} else {
			pcf.doc = doc
		}
	}

	return pcf
}

// Map gets the map of conf file
func (p *PropertiesConfFile) Map() map[string]string {
	m := make(map[string]string)
	p.doc.Foreach(func(v, k string) bool {
		m[k] = v
		return true
	})

	return m
}

// ConfFormat gets the format of conf file
func (p *PropertiesConfFile) ConfFormat() ConfFmt {
	return PropertiesFmt
}

// Str get the string value of key specified by name.
func (p *PropertiesConfFile) Str(name string) string {
	v, _ := p.doc.Get(name)
	return v
}

// StrOr get the string value of key specified by name or defaultValue when value is empty or missed.
func (p *PropertiesConfFile) StrOr(name, defaultValue string) string {
	v := p.Str(name)
	if v == "" {
		return defaultValue
	}

	return v
}

// Bool get the bool value of key specified by name.
func (p *PropertiesConfFile) Bool(name string) bool {
	v := p.StrOr(name, "false")
	val, err := strconv.ParseBool(v)
	if err != nil {
		logrus.Warnf("parse bool fail for %s", v)
		return false
	}
	return val
}

// BoolOr get the bool value of key specified by name or defaultValue when value is empty or missed.
func (p *PropertiesConfFile) BoolOr(name string, defaultValue bool) bool {
	v := p.StrOr(name, "")
	if v == "" {
		return defaultValue
	}

	return p.Bool(name)
}

// Int get the int value of key specified by name.
func (p *PropertiesConfFile) Int(name string) int {
	v := p.StrOr(name, "0")
	val, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		logrus.Warnf("parse Int fail for %s", v)
		return 0
	}

	return int(val)
}

// IntOr get the int value of key specified by name or defaultValue when value is empty or missed.
func (p *PropertiesConfFile) IntOr(name string, defaultValue int) int {
	v := p.StrOr(name, "")
	if v == "" {
		return defaultValue
	}

	return p.Int(name)
}

// Int32 get the int32 value of key specified by name.
func (p *PropertiesConfFile) Int32(name string) int32 {
	v := p.StrOr(name, "0")
	val, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		logrus.Warnf("parse Int32 fail for %s", v)
		return 0
	}

	return int32(val)
}

// Int32Or get the int32 value of key specified by name or defaultValue when value is empty or missed.
func (p *PropertiesConfFile) Int32Or(name string, defaultValue int32) int32 {
	v := p.StrOr(name, "")
	if v == "" {
		return defaultValue
	}

	return p.Int32(name)
}

// Int64 get the int64 value of key specified by name.
func (p *PropertiesConfFile) Int64(name string) int64 {
	v := p.StrOr(name, "0")
	val, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		logrus.Warnf("parse Int64 fail for %s", v)
		return 0
	}

	return val
}

// Int64Or get the int64 value of key specified by name or defaultValue when value is empty or missed.
func (p *PropertiesConfFile) Int64Or(name string, defaultValue int64) int64 {
	v := p.StrOr(name, "")
	if v == "" {
		return defaultValue
	}

	return p.Int64(name)
}

// Float32 get the float32 value of key specified by name.
func (p *PropertiesConfFile) Float32(name string) float32 {
	v := p.StrOr(name, "0")
	val, err := strconv.ParseFloat(v, 32)
	if err != nil {
		logrus.Warnf("parse Float32 fail for %s", v)
		return 0
	}

	return float32(val)
}

// Float32Or get the float32 value of key specified by name or defaultValue when value is empty or missed.
func (p *PropertiesConfFile) Float32Or(name string, defaultValue float32) float32 {
	v := p.StrOr(name, "")
	if v == "" {
		return defaultValue
	}

	return p.Float32(name)
}

// Float64 get the float64 value of key specified by name.
func (p *PropertiesConfFile) Float64(name string) float64 {
	v := p.StrOr(name, "0")
	val, err := strconv.ParseFloat(v, 64)
	if err != nil {
		logrus.Warnf("parse Float64 fail for %s", v)
		return 0
	}
	return val
}

// Float64Or get the float64 value of key specified by name or defaultValue when value is empty or missed.
func (p *PropertiesConfFile) Float64Or(name string, defaultValue float64) float64 {
	v := p.StrOr(name, "")
	if v == "" {
		return defaultValue
	}

	return p.Float64(name)
}
