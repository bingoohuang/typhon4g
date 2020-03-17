package base

import (
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

// MakeDirAll make sure the mkdir all success or panic
func MakeDirAll(dir string) string {
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		logrus.Panicf("failed to create dir %s, error %v", dir, err)
	}

	return dir
}

// HTPPAddr try to prepends the  http:// to the address
func HTPPAddr(address string) string {
	addr := address
	if !(strings.HasPrefix(address, "http://") || strings.HasPrefix(address, "https://")) {
		addr = fmt.Sprintf("http://%s", address)
	}

	if strings.HasSuffix(addr, "/") {
		return addr[:len(addr)-1]
	}

	return addr
}
