package typhon4g

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

var runner *Runner

func init() {
	r, err := CreateRunner("etc/typhon-context.properties")
	if err != nil {
		logrus.Warnf("unable to create default typhon runner %v", err)
		return
	}

	r.Start()
	runner = r
}

func CreateRunner(contextFile string) (*Runner, error) {
	context, err := LoadContextFile(contextFile)
	if err != nil {
		return nil, fmt.Errorf("unable to Load %s, %v", contextFile, err)
	}

	return &Runner{
		C:               context,
		SnapshotService: &SnapshotService{C: context},
		MetaService:     &MetaService{C: context},
		ConfigService:   &ConfigService{C: context},
	}, nil
}

func Properties(confFile string) (*PropertiesConfFile, error) {
	return runner.Properties(confFile)
}

func GetConfFile(confFile string) (ConfFile, error) {
	return runner.ConfFile(confFile)
}

func PostConf(confFile, raw, clientIps string) (string, error) {
	return runner.PostConf(confFile, raw, clientIps)
}

func ListenerResults(confFile, crc string) ([]ClientReportRspItem, error) {
	return runner.ListenerResults(confFile, crc)
}
