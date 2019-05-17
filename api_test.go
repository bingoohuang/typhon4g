package typhon4g_test

import (
	"fmt"
	"github.com/bingoohuang/gou"
	"github.com/sirupsen/logrus"
	"testing"
	"time"
	"typhon4g"
)

type MyListener struct{}

// Make sure that MyListener implements the interface typhon4g.ConfFileChangeListener
var _ typhon4g.ConfFileChangeListener = (*MyListener)(nil)

func (l MyListener) OnChange(event typhon4g.ConfFileChangeEvent) (msg string, ok bool) {
	fmt.Println("OnChange", event)
	// eat your own dog food here
	return "your message", true /*  true to means changed OK */
}

func TestGetConf(t *testing.T) {
	prop, err := typhon4g.GetProperties("hello.properties")
	if err != nil {
		logrus.Panic(err)
	}
	fmt.Println("name:", prop.String("name"))
	fmt.Println("home:", prop.StringDefault("home", "中国"))
	fmt.Println("age:", prop.Int("age"))
	fmt.Println("adult", prop.Bool("adult"))

	hello, err := typhon4g.GetConfFile("hello.json")
	if err != nil {
		logrus.Panic(err)
	}
	fmt.Println("hello.json:", hello.Raw())

	var listener MyListener
	prop.Register(&listener)

	r := gou.RandomNum(3)
	crc, err := typhon4g.PostConf("hello.properties", "name=bingoo\nage="+r+"\n", "all")
	if err != nil {
		logrus.Panicf("error %v", err)
	}
	fmt.Println("crc:", crc)

	for i := 0; i < 1; i++ {
		fmt.Println("sleep 10 seconds")
		time.Sleep(time.Duration(10) * time.Second)
		fmt.Println(prop.String("name"))

		fmt.Println("hello.json:", hello.Raw())
	}

	items, err := typhon4g.GetListenerResults("hello.properties", crc)
	if err != nil {
		logrus.Panicf("error %v", err)
	}

	fmt.Println("items", items)
}
