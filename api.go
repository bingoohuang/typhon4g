package typhon4g

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
)

var runner *Runner

func init() {
	r, err := CreateRunner("etc/typhon-context.properties")
	if err != nil {
		logrus.Warnf("unable to create typhon runner %v", err)
		return
	}

	r.Start()
	runner = r
}

func CreateRunner(contextFile string) (*Runner, error) {
	context, err := loadTyphonContext(contextFile)
	if err != nil {
		return nil, fmt.Errorf("unable to load %s, %v", contextFile, err)
	}

	return &Runner{
		context:  context,
		snapshot: &SnapshotService{c: context},
		meta:     &MetaService{c: context},
		config:   &ConfigService{c: context},
	}, nil
}

func GetProperties(confFile string) (*PropertiesConfFile, error) {
	return runner.GetProperties(confFile)
}

func GetConfFile(confFile string) (ConfFile, error) {
	return runner.GetConfFile(confFile)
}

func PostConf(confFile, raw, clientIps string) (string, error) {
	return runner.PostConf(confFile, raw, clientIps)
}

func GetListenerResults(confFile, crc string) ([]ClientReportRspItem, error) {
	return runner.GetListenerResults(confFile, crc)
}

func makeSnapshotDir(dir string) {
	st, err := os.Stat(dir)
	if err == nil {
		if st.IsDir() {
			return
		} else {
			logrus.Panicf("make sure that snapshot dir %s is a directory", dir)
		}
	}

	if os.IsNotExist(err) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			logrus.Panicf("failed to create snapshot dir %s, error %v", dir, err)
		}
		return
	}

	logrus.Panicf("failed to stat snapshot dir %s, error %v", dir, err)
}
