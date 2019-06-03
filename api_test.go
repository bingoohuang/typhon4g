package typhon4g_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/bingoohuang/gou"
	"github.com/bingoohuang/typhon4g"
	"github.com/sirupsen/logrus"
)

var typhon *typhon4g.Runner

func init() {
	var err error
	if typhon, err = typhon4g.LoadStart(); err != nil {
		logrus.Panic(err)
	}
}

type MyListener struct{}

// Make sure that MyListener implements the interface typhon4g.ConfFileChangeListener
var _ typhon4g.ConfFileChangeListener = (*MyListener)(nil)

func (l MyListener) OnChange(event typhon4g.ConfFileChangeEvent) (msg string, ok bool) {
	fmt.Println("OnChange", event)
	// eat your own dog food here
	return "your message", true /*  true to means changed OK */
}

func TestGetConf(t *testing.T) {
	prop, err := typhon.Properties("hello.properties")
	if err != nil {
		logrus.Panic(err)
	}
	fmt.Println("name:", prop.Str("name"))
	fmt.Println("home:", prop.StrOr("home", "中国"))
	fmt.Println("age:", prop.Int("age"))
	fmt.Println("adult", prop.Bool("adult"))

	hello, err := typhon.ConfFile("hello.json")
	if err != nil {
		logrus.Panic(err)
	}
	fmt.Println("hello.json:", hello.Raw())

	var listener MyListener
	prop.Register(&listener)

	crc, err := typhon.PostConf("hello.properties", "name=bingoo\nage="+gou.RandomNum(3)+"\n", "all")
	if err != nil {
		logrus.Panicf("error %v", err)
	}
	fmt.Println("crc:", crc)

	for i := 0; i < 1; i++ {
		fmt.Println("sleep 10 seconds")
		time.Sleep(typhon4g.SecondsDuration(10))
		fmt.Println(prop.Str("name"))

		fmt.Println("hello.json:", hello.Raw())
	}

	items, err := typhon.ListenerResults("hello.properties", crc)
	if err != nil {
		logrus.Panicf("error %v", err)
	}

	fmt.Println("items", items)
}
