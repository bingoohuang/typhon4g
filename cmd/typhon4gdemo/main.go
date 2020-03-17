// nolint gomnd
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/spf13/viper"

	"github.com/bingoohuang/typhon4g"
)

func init() {
	typhon4g.LoadStart().Register("application.properties", &typhon4g.ViperListener{})
}

func main() {
	logrus.SetLevel(logrus.DebugLevel)

	enter := "Enter get(application.properties): "
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
		switch cmd {
		case "get":
			fmt.Println("viper get:", viper.GetString("name"))
		default:
			fmt.Println("unknown cmd" + cmd)
			os.Exit(0)
		}
	}
}
