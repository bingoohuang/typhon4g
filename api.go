package typhon4g

import (
	"fmt"
	"runtime"

	"github.com/bingoohuang/typhon4g/base"
)

// LoadStart loads and starts the typhon client.
func LoadStart() *base.Runner {
	r, e := LoadStartE()
	if e != nil {
		panic(e)
	}

	return r
}

// LoadStartE loads and starts the typhon client.
func LoadStartE() (*base.Runner, error) {
	return LoadStartByFileE("etc/typhon-context.properties")
}

// LoadStartByFile loads and starts the typhon client.
func LoadStartByFile(typhonContextFile string) *base.Runner {
	if r, e := LoadStartByFileE(typhonContextFile); e != nil {
		panic(e)
	} else {
		return r
	}
}

// LoadStartByFileE loads and starts the typhon client.
func LoadStartByFileE(typhonContextFile string) (*base.Runner, error) {
	context, err := LoadContextFile(typhonContextFile)
	if err != nil {
		return nil, fmt.Errorf("unable to Load %s, %v", typhonContextFile, err)
	}

	return LoadStartByContext(context), nil
}

// LoadStartByContext load the runner from the context.
func LoadStartByContext(context *base.Context) *base.Runner {
	r := &base.Runner{
		Context:         context,
		SnapshotService: &base.SnapshotService{Context: context},
		MetaService:     &base.MetaService{Context: context},
		ConfigService:   &base.ConfigService{Context: context},
		PollingService:  &base.PollingService{Context: context},
	}

	r.Start()
	runtime.SetFinalizer(r, func(r *base.Runner) { r.Stop() })

	return r
}
