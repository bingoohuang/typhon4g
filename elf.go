package typhon4g

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

// SecondsDuration converts i to duration in second unit.
func SecondsDuration(i int64) time.Duration {
	return time.Duration(i) * time.Second
}

// MillisDuration converts i to duration in millis unit
func MillisDuration(i int64) time.Duration {
	return time.Duration(i) * time.Millisecond
}

// MustMakeDirAll make sure the mkdir all success or panic
func MustMakeDirAll(dir string) string {
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		logrus.Panicf("failed to create dir %s, error %v", dir, err)
	}
	return dir
}

// Required make sure that the sth is required  or panic.
func Required(sth, name string) string {
	if sth != "" {
		return sth
	}

	logrus.Panic(name, " required")
	return ""
}
