package main

import (
	"bufio"
	"fmt"
	"github.com/bingoohuang/gou"
	"github.com/bingoohuang/typhon4g"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

var ty *typhon4g.Runner

func init() {
	var err error
	if ty, err = typhon4g.LoadStart(); err != nil {
		logrus.Panic(err)
	}
}

type MyListener struct{}

// Make sure that MyListener implements the interface typhon4g.ConfFileChangeListener
var _ typhon4g.ConfFileChangeListener = (*MyListener)(nil)

func (l MyListener) OnChange(event typhon4g.ConfFileChangeEvent) (msg string, ok bool) {
	fmt.Printf("OnChange: %+v\n", event)
	// eat your own dog food here
	return "my message:" + gou.RandString(10), true /*  true to means changed OK */
}

func main() {
	var err error
	prop, err := ty.Properties("hello.properties")
	if err != nil {
		logrus.Panic(err)
	}
	fmt.Println(prop.Map())

	hello, err := ty.ConfFile("hello.json")
	if err != nil {
		logrus.Panic(err)
	}
	fmt.Println("hello.json:", hello.Raw())

	var listener MyListener
	prop.Register(&listener)

	var crc string

	enter := "Enter get/put/his(hello.properties): "
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
			content := "name=" + gou.RandString(10) + "\nage=" + gou.RandomNum(3) + "\n"
			fmt.Println("putting:", content)
			crc, err = ty.PostConf("hello.properties", content, "all")
			if err != nil {
				logrus.Panicf("error %v", err)
			}
			fmt.Println("crc:", crc)
		case "his":
			items, err := ty.ListenerResults("hello.properties", crc)
			if err != nil {
				logrus.Panicf("error %v", err)
			}

			his := gou.PrettyJsontify(items)
			fmt.Println("items:\n", his)
		case "get":
			hello, err := ty.ConfFile("hello.properties")
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
