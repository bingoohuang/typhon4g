package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/bingoohuang/gou"
	"github.com/bingoohuang/typhon4g"
	"github.com/sirupsen/logrus"
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
	fmt.Println("OnChange", event)
	// eat your own dog food here
	return "my message:" + gou.RandString(10), true /*  true to means changed OK */
}

func main() {
	var err error
	prop, err := ty.Properties("hello.properties")
	if err != nil {
		logrus.Panic(err)
	}
	fmt.Println(prop.Raw())

	hello, err := ty.ConfFile("hello.json")
	if err != nil {
		logrus.Panic(err)
	}
	fmt.Println("hello.json:", hello.Raw())

	var listener MyListener
	prop.Register(&listener)

	scanner := bufio.NewScanner(os.Stdin)
	var crc string

	for {
		fmt.Print("Enter put/his/ConfFile(default hello.properties): ")
		// Scans a line from Stdin(Console)
		scanner.Scan()

		// Holds the string that scanned
		cmd := scanner.Text()
		if cmd == "" {
			cmd = "hello.properties"
		}

		switch cmd {
		case "put":
			content := "name=" + gou.RandString(10) + "\nage=" + gou.RandomNum(3) + "\n"
			fmt.Println("puting:", content)
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
		default:
			hello, err := ty.ConfFile(cmd)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(cmd, ":", hello.Raw())
			}
		}
	}
}
