// nolint gomnd
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/bingoohuang/gou/enc"
	"github.com/bingoohuang/gou/ran"
	"github.com/bingoohuang/typhon4g/base"

	"github.com/bingoohuang/properties"

	"github.com/bingoohuang/typhon4g"
	"github.com/sirupsen/logrus"
)

type myListener struct{}

// Make sure that myListener implements the interface typhon4g.ChangeListener
var _ base.ChangeListener = (*myListener)(nil)

// OnChange defines the callback when changes detected.
func (l myListener) OnChange(event base.ConfFileChangeEvent) (msg string, ok bool) {
	fmt.Printf("OnChange: %+v\n", event)
	// eat your own dog food here
	old, _ := properties.LoadString(event.Old)
	cur, _ := properties.LoadString(event.Current)

	properties.Diff(old, cur, func(e properties.DiffEvent) {
		fmt.Println(e)
	})

	return "my message:" + ran.String(10), true /*  true to means changed OK */
}

func main() {
	ty := typhon4g.LoadStart()
	var err error
	prop, err := ty.Properties("application.properties")
	if err != nil {
		logrus.Panic(err)
	}

	propContent, _ := properties.LoadString(prop.Raw())
	propContent.Set("Key", "Value")
	//newContent := propContent.String()
	//_, _ = ty.PostConf("hello.properties", newContent, "all")

	fmt.Println(prop.Doc.Map())

	//hello, err := ty.ConfFile("hello.json")
	//if err != nil {
	//	logrus.Panic(err)
	//}
	//fmt.Println("hello.json:", hello.Raw())

	var listener myListener
	prop.Register(&listener)

	var crc string

	enter := "Enter get/put/his(application.properties): "
	fmt.Print(enter)

	snr := bufio.NewScanner(os.Stdin)
	for ; snr.Scan(); fmt.Print(enter) {
		line := snr.Text()
		if len(line) == 0 {
			continue
		}

		fields := strings.Fields(line)
		fmt.Printf("Fields: %q\n", fields)

		cmd := fields[0]
		//for {
		//	time.Sleep(10 * time.Second)
		//	cmd := "put"
		switch cmd {
		case "put":
			content := "name=" + ran.String(10) + "\nage=" + ran.Num(3) + "\n"
			fmt.Println("putting:", content)
			crc, err = ty.PostConf("application.properties", content, "all")
			if err != nil {
				logrus.Panicf("error %v", err)
			}
			fmt.Println("crc:", crc)
		case "his":
			items, err := ty.ListenerResults("application.properties", crc)
			if err != nil {
				logrus.Panicf("error %v", err)
			}

			his := enc.JSONPretty(items)
			fmt.Println("items:\n", his)
		case "get":
			hello, err := ty.ConfFile("application.properties")
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(cmd, ":", hello.Raw())
			}
		default:
			fmt.Println("unknown cmd" + cmd)
			os.Exit(0)
		}
	}
}
