package typhon4g

import (
	"bytes"
	"os"
	"time"

	"github.com/bingoohuang/gou"

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

type ChangeType int

const (
	ChangeTypeModify ChangeType = iota
	ChangeTypeNew
	ChangeTypeDeleted
)

//go:generate enumer -type=ChangeType -json

// DiffPropertiesChangeEvent defines ChangeEvent for properties diff
type DiffPropertiesChangeEvent struct {
	ChangeType ChangeType
	Key        string
	LeftValue  string
	RightValue string
}

func DiffProperties(l, r string, f func(DiffPropertiesChangeEvent)) error {
	ldoc, err := gou.LoadProperties(bytes.NewBufferString(l))
	if err != nil {
		return err
	}

	rdoc, err := gou.LoadProperties(bytes.NewBufferString(r))
	if err != nil {
		return err
	}

	lm := make(map[string]string)
	rm := make(map[string]string)
	ldoc.Foreach(func(v, k string) bool { lm[k] = v; return true })
	rdoc.Foreach(func(v, k string) bool { rm[k] = v; return true })

	lKeys := gou.MapKeysSorted(lm).([]string)
	for _, lk := range lKeys {
		lv := lm[lk]
		rv, ok := rm[lk]
		if ok {
			if lv != rv {
				f(DiffPropertiesChangeEvent{ChangeType: ChangeTypeModify, Key: lk, LeftValue: lv, RightValue: rv})
			}
			delete(rm, lk)
		} else {
			f(DiffPropertiesChangeEvent{ChangeType: ChangeTypeDeleted, Key: lk, LeftValue: lv, RightValue: ""})
		}
	}

	rKeys := gou.MapKeysSorted(rm).([]string)
	for _, rk := range rKeys {
		f(DiffPropertiesChangeEvent{ChangeType: ChangeTypeNew, Key: rk, LeftValue: "", RightValue: rm[rk]})
	}

	return nil
}
