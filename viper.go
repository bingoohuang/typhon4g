package typhon4g

import (
	"encoding/json"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/bingoohuang/properties"
	"github.com/bingoohuang/typhon4g/base"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

// ViperListener defines the viper setter listener
type ViperListener struct {
	Prefix string
}

// Make sure that myListener implements the interface typhon4g.ChangeListener
var _ base.ChangeListener = (*ViperListener)(nil)

// OnChange defines the callback when changes detected.
func (l ViperListener) OnChange(event base.ConfFileChangeEvent) (msg string, ok bool) {
	switch ext := strings.ToLower(filepath.Ext(event.ConfFile)); ext {
	case ".properties":
		l.propertiesDiff(event)
	case ".yaml", ".yml":
		if oldM, newM, ok := unmarshal(yaml.Unmarshal, event); ok {
			l.mapDiff(oldM, newM)
		}
	case ".toml", ".tml":
		if oldM, newM, ok := unmarshal(toml.Unmarshal, event); ok {
			l.mapDiff(oldM, newM)
		}
	case ".js", ".json":
		if oldM, newM, ok := unmarshal(json.Unmarshal, event); ok {
			l.mapDiff(oldM, newM)
		}
	default:
		logrus.Warnf("unsupported format for config file %s", event.ConfFile)
	}

	logrus.Debugf("onchanged %+v", event)

	return "got", true /*  true to means changed OK */
}

type unmarshaler func(p []byte, v interface{}) error

func unmarshal(f unmarshaler, event base.ConfFileChangeEvent) (a, b map[string]interface{}, ok bool) {
	oldM := make(map[string]interface{})
	if err := f([]byte(event.Old), &oldM); err != nil {
		logrus.Warnf("Unmarshal error: %v", err)
		return nil, nil, false
	}

	newM := make(map[string]interface{})
	if err := f([]byte(event.Current), &newM); err != nil {
		logrus.Warnf("Unmarshal error: %v", err)
		return nil, nil, false
	}

	return oldM, newM, true
}

func (l ViperListener) mapDiff(oldM map[string]interface{}, newM map[string]interface{}) {
	for _, item := range DiffMap(oldM, newM) {
		key := l.Prefix + item.Key.(string)

		switch item.ChangeType {
		case Added, Modified:
			viper.Set(key, item.RightValue)
		case Removed:
			_ = ViperUnset(key)
		}
	}
}

// ViperUnset unset the viper key.
func ViperUnset(key string) error {
	viper.Set(key, nil)

	return nil
}

func (l ViperListener) propertiesDiff(event base.ConfFileChangeEvent) {
	old, _ := properties.LoadString(event.Old)
	cur, _ := properties.LoadString(event.Current)

	properties.Diff(old, cur, func(e properties.DiffEvent) {
		key := l.Prefix + e.Key

		switch e.ChangeType {
		case properties.Added, properties.Modified:
			viper.Set(key, e.RightValue)
		case properties.Removed:
			_ = ViperUnset(key)
		}
	})
}

// ChangeType defines the type of chaging.
type ChangeType int

const (
	// Modified ...
	Modified ChangeType = iota
	// Added ...
	Added
	// Removed ...
	Removed
	// Same ...
	Same
)

// DiffItem defines ChangeEvent for properties diff
type DiffItem struct {
	ChangeType ChangeType
	Key        interface{}
	LeftValue  interface{}
	RightValue interface{}
}

// DiffMap compares two map.
func DiffMap(a, b interface{}) []DiffItem {
	am := reflect.ValueOf(a)
	bm := reflect.ValueOf(b)
	items := make([]DiffItem, 0)

	for _, ak := range am.MapKeys() {
		av, bv := am.MapIndex(ak), bm.MapIndex(ak)
		aki, avi := ak.Interface(), av.Interface()

		if !bv.IsValid() {
			items = append(items, DiffItem{ChangeType: Removed, Key: aki, LeftValue: avi, RightValue: nil})
			continue
		}

		if bvi := bv.Interface(); reflect.DeepEqual(avi, bvi) {
			items = append(items, DiffItem{ChangeType: Same, Key: aki, LeftValue: avi, RightValue: bvi})
		} else {
			items = append(items, DiffItem{ChangeType: Modified, Key: aki, LeftValue: avi, RightValue: bvi})
		}
	}

	for _, bk := range bm.MapKeys() {
		if av := am.MapIndex(bk); av.IsValid() {
			continue
		}

		bvi := bm.MapIndex(bk).Interface()

		items = append(items, DiffItem{ChangeType: Added, Key: bk.Interface(), LeftValue: nil, RightValue: bvi})
	}

	return items
}
