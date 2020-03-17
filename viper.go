package typhon4g

import (
	"github.com/bingoohuang/properties"
	"github.com/bingoohuang/typhon4g/base"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// ViperListener defines the viper setter listener
type ViperListener struct{}

// Make sure that myListener implements the interface typhon4g.ChangeListener
var _ base.ChangeListener = (*ViperListener)(nil)

// OnChange defines the callback when changes detected.
func (l ViperListener) OnChange(event base.ConfFileChangeEvent) (msg string, ok bool) {
	old, _ := properties.LoadString(event.Old)
	cur, _ := properties.LoadString(event.Current)

	properties.Diff(old, cur, func(e properties.DiffEvent) {
		switch e.ChangeType {
		case properties.Added, properties.Modified:
			viper.Set(e.Key, e.RightValue)
		case properties.Removed:
			viper.Set(e.Key, nil)
		}
	})

	logrus.Debugf("onchanged %+v", event)

	return "got", true /*  true to means changed OK */
}
