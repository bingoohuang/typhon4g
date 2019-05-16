package typhon4g_test

import (
	"fmt"
	"testing"
	"time"
	"typhon4g"
)

type MyListener struct {
}

var _ typhon4g.ConfFileChangeListener = (*MyListener)(nil)

func (l MyListener) OnChange(event typhon4g.ConfFileChangeEvent) (string, bool) {
	fmt.Errorf("%+v", event)
	return "", true
}

func TestGetConf(t *testing.T) {
	conf := typhon4g.GetProperties("hello.properties")
	fmt.Println(conf.String("name"))

	var listener MyListener
	conf.Register(&listener)

	for {
		fmt.Println("sleep 10 seconds")
		time.Sleep(time.Duration(10) * time.Second)
	}
}
