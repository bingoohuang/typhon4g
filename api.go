package typhon4g

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

// LoadStart loads and starts the typhon client.
func LoadStart() (*Runner, error) {
	r, err := createRunner("etc/typhon-context.properties")
	if err != nil {
		logrus.Warnf("unable to create default typhon typhon %v", err)
		return nil, err
	}

	r.Start()
	return r, nil
}

func createRunner(contextFile string) (*Runner, error) {
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
