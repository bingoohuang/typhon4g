package typhon4g_test

import (
	"fmt"
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
	prop := typhon4g.GetProperties("hello.properties")
	fmt.Println("name:", prop.String("name"))
	fmt.Println("home:", prop.StringDefault("home", "中国"))
	fmt.Println("age:", prop.Int("age"))
	fmt.Println("adult", prop.Bool("adult"))

	hello := typhon4g.GetConfFile("hello.json")
	fmt.Println("hello.json:", hello.Raw())

	var listener MyListener
	prop.Register(&listener)

	for {
		fmt.Println("sleep 10 seconds")
		time.Sleep(time.Duration(10) * time.Second)
		fmt.Println(prop.String("name"))

		fmt.Println("hello.json:", hello.Raw())
	}
}
