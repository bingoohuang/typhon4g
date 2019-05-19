package typhon4g

import (
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

func SecondsDuration(i int64) time.Duration {
	return time.Duration(i) * time.Second
}

func MillisDuration(i int64) time.Duration {
	return time.Duration(i) * time.Millisecond
}

func MustMakeDirAll(dir string) string {
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		logrus.Panicf("failed to create dir %s, error %v", dir, err)
	}
	return dir
}

func MustExists(sth, name string) string {
	if sth != "" {
		return sth
	}

	logrus.Panic(name, " required")
	return ""
}
