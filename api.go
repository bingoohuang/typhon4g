package typhon4g

import (
	"fmt"
	"runtime"
)

// LoadStart loads and starts the typhon client.
func LoadStart() (*Runner, error) {
	return LoadStartByFile("etc/typhon-context.properties")
}

// LoadStartByFile loads and starts the typhon client.
func LoadStartByFile(typhonContextFile string) (*Runner, error) {
	context, err := LoadContextFile(typhonContextFile)
	if err != nil {
		return nil, fmt.Errorf("unable to Load %s, %v", typhonContextFile, err)
	}

	return LoadStartByContext(context), nil
}

func LoadStartByContext(context *TyphonContext) *Runner {
	r := &Runner{
		C:               context,
		SnapshotService: &SnapshotService{C: context},
		MetaService:     &MetaService{C: context},
		ConfigService:   &ConfigService{C: context},
	}
	r.Start()
	runtime.SetFinalizer(r, func(r *Runner) { r.Stop() })
	return r
}
