package typhon4g

import (
	"bytes"
	"strconv"

	"github.com/bingoohuang/gou"
	"github.com/sirupsen/logrus"
)

type PropertiesConfFile struct {
	BaseConf
	Doc *gou.PropertiesDoc
}

var _ Prop = (*PropertiesConfFile)(nil)

func NewPropertiesConfFile(confFile, raw string) *PropertiesConfFile {
	doc, err := gou.LoadProperties(bytes.NewBufferString(raw))
	if err != nil {
		logrus.Warnf("LoadProperties %v", err)
		return nil
	}

	pcf := &PropertiesConfFile{
		BaseConf: BaseConf{
			raw:       raw,
			confFile:  confFile,
			listeners: make([]ConfFileChangeListener, 0),
		},
		Doc: doc,
	}

	pcf.updater = func(updated string) {
		if doc, err := gou.LoadProperties(bytes.NewBufferString(updated)); err != nil {
			logrus.Warnf("LoadProperties %v", err)
		} else {
			pcf.Doc = doc
		}
	}

	return pcf
}

// ConfFormat gets the format of conf file
func (p *PropertiesConfFile) ConfFormat() ConfFmt {
	return PropertiesFmt
}

func (p *PropertiesConfFile) Str(name string) string {
	value, _ := p.Doc.Get(name)
	return value
}

func (p *PropertiesConfFile) StrDefault(name, defaultValue string) string {
	value := p.Str(name)
	if value == "" {
		return defaultValue
	}

	return value
}

func (p *PropertiesConfFile) Bool(name string) bool {
	v := p.StrDefault(name, "false")
	val, err := strconv.ParseBool(v)
	if err != nil {
		logrus.Warnf("parse bool fail for %s", v)
		return false
	}
	return val
}

func (p *PropertiesConfFile) BoolDefault(name string, defaultValue bool) bool {
	v := p.StrDefault(name, "")
	if v == "" {
		return defaultValue
	}

	return p.Bool(name)
}

func (p *PropertiesConfFile) Int(name string) int {
	v := p.StrDefault(name, "0")
	val, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		logrus.Warnf("parse Int fail for %s", v)
		return 0
	}

	return int(val)
}

func (p *PropertiesConfFile) IntDefault(name string, defaultValue int) int {
	v := p.StrDefault(name, "")
	if v == "" {
		return defaultValue
	}

	return p.Int(name)
}

func (p *PropertiesConfFile) Int32(name string) int32 {
	v := p.StrDefault(name, "0")
	val, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		logrus.Warnf("parse Int32 fail for %s", v)
		return 0
	}

	return int32(val)
}

func (p *PropertiesConfFile) Int32Default(name string, defaultValue int32) int32 {
	v := p.StrDefault(name, "")
	if v == "" {
		return defaultValue
	}

	return p.Int32(name)
}

func (p *PropertiesConfFile) Int64(name string) int64 {
	v := p.StrDefault(name, "0")
	val, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		logrus.Warnf("parse Int64 fail for %s", v)
		return 0
	}

	return val
}

func (p *PropertiesConfFile) Int64Default(name string, defaultValue int64) int64 {
	v := p.StrDefault(name, "")
	if v == "" {
		return defaultValue
	}

	return p.Int64(name)
}

func (p *PropertiesConfFile) Float32(name string) float32 {
	v := p.StrDefault(name, "0")
	val, err := strconv.ParseFloat(v, 32)
	if err != nil {
		logrus.Warnf("parse Float32 fail for %s", v)
		return 0
	}

	return float32(val)
}

func (p *PropertiesConfFile) Float32Default(name string, defaultValue float32) float32 {
	v := p.StrDefault(name, "")
	if v == "" {
		return defaultValue
	}

	return p.Float32(name)
}

func (p *PropertiesConfFile) Float64(name string) float64 {
	v := p.StrDefault(name, "0")
	val, err := strconv.ParseFloat(v, 64)
	if err != nil {
		logrus.Warnf("parse Float64 fail for %s", v)
		return 0
	}
	return val
}

func (p *PropertiesConfFile) Float64Default(name string, defaultValue float64) float64 {
	v := p.StrDefault(name, "")
	if v == "" {
		return defaultValue
	}

	return p.Float64(name)
}
