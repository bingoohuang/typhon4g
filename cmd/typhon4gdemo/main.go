// nolint gomnd
package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/spf13/viper"

	"github.com/bingoohuang/typhon4g"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	ty := typhon4g.LoadStart()
	ty.Register("application.properties", &typhon4g.ViperListener{Prefix: ""})
	ty.Register("helloyaml.yaml", &typhon4g.ViperListener{Prefix: "hello."})
}

func main() {
	for snr := bufio.NewScanner(os.Stdin); snr.Scan(); {
		for _, key := range viper.AllKeys() {
			fmt.Println(key, ":", viper.Get(key))
		}
	}
}
