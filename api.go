package typhon4g

import (
	"fmt"
	"runtime"
)

// MustLoadStart loads and starts the typhon client.
func MustLoadStart() *Runner {
	if r, e := LoadStart(); e != nil {
		panic(e)
	} else {
		return r
	}
}

// LoadStart loads and starts the typhon client.
func LoadStart() (*Runner, error) {
	return LoadStartByFile("etc/typhon-context.properties")
}

// MustLoadStartByFile loads and starts the typhon client.
func MustLoadStartByFile(typhonContextFile string) *Runner {
	if r, e := LoadStartByFile(typhonContextFile); e != nil {
		panic(e)
	} else {
		return r
	}
}

// LoadStartByFile loads and starts the typhon client.
func LoadStartByFile(typhonContextFile string) (*Runner, error) {
	context, err := LoadContextFile(typhonContextFile)
	if err != nil {
		return nil, fmt.Errorf("unable to Load %s, %v", typhonContextFile, err)
	}

	return LoadStartByContext(context), nil
}

func LoadStartByContext(context *Context) *Runner {
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
